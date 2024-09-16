package main

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

type Comment struct {
	Score int
	Text  string
	// but you can add anything you want
}

type Survey struct {
	mu *sync.Mutex
	S  map[string]passObj
}
type passObj struct {
	mu *sync.Mutex
	P  map[string]Comment
}

func NewSurvey() *Survey {
	tempComment := make(map[string]passObj)
	return &Survey{&sync.Mutex{}, tempComment}
}

func (s *Survey) AddFlight(flightName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.S[flightName]
	if exists {
		return errors.New("flight exists")
	} else {

		s.S[flightName] = passObj{&sync.Mutex{}, make(map[string]Comment)}

		return nil
	}
}

func (s *Survey) AddTicket(flightName, passengerName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, flighexists := s.S[flightName]

	if !flighexists {
		return errors.New("flight doesn't exists, unable to issue the ticket")
	}
	_, passengerExistInFlight := s.S[flightName].P[passengerName]
	if flighexists && passengerExistInFlight {
		return errors.New("duplicate ticket")

	}

	s.S[flightName].P[passengerName] = Comment{}

	return nil
}

func (s *Survey) AddComment(flightName, passengerName string, comment Comment) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, flighexists := s.S[flightName]
	if !flighexists {
		return errors.New("flight doesn't exists, unable to add the comment")
	}

	_, passengerExistInFlight := s.S[flightName].P[passengerName]
	if flighexists && !passengerExistInFlight {
		return errors.New("no ticket exists for this passenger")
	}

	if flighexists && passengerExistInFlight && (s.S[flightName].P[passengerName].Text != "") {
		return errors.New("this passenger has commented before")
	}

	if flighexists && passengerExistInFlight && (s.S[flightName].P[passengerName].Text == "") && ((comment.Score < 1) && (comment.Score > 10)) {
		return errors.New("incorrect score")
	}

	s.S[flightName].P[passengerName] = comment

	return nil
}

func (s *Survey) GetCommentsAverage(flightName string) (float64, error) {
	valObj, flightExists := s.S[flightName]
	if !flightExists {
		return 0, errors.New("flight doesn't exist. unable to return average point")
	}
	var tempFloat float64
	var commentNum int
	for _, commentInLoop := range valObj.P {
		if commentInLoop.Score != 0 {
			tempFloat = tempFloat + float64(commentInLoop.Score)
			commentNum++
		}
	}
	if commentNum == 0 {
		return 0, errors.New("no comment available for this flight")
	} else {
		tempFloat = tempFloat / float64(commentNum)
		return tempFloat, nil
	}

}

func (s *Survey) GetAllCommentsAverage() map[string]float64 {
	output := make(map[string]float64)
	for flightName, _ := range s.S {
		tempAverage, _ := s.GetCommentsAverage(flightName)
		if tempAverage != 0 {
			output[flightName] = tempAverage
		}
	}
	return output
}

func (s *Survey) GetComments(flightName string) ([]string, error) {
	var output []string

	_, flightexists := s.S[flightName]
	if !flightexists {
		return output, errors.New("flight doesn't exist. GetComment error")
	}

	for _, commentObj := range s.S[flightName].P {
		if commentObj.Text != "" {
			output = append(output, commentObj.Text)
		}
	}

	return output, nil
}

func (s *Survey) GetAllComments() map[string][]string {
	output := make(map[string][]string)
	for flightName, _ := range s.S {
		comments, _ := s.GetComments(flightName)
		output[flightName] = comments
	}
	return output
}

//////\

type healthMessage struct {
	Message string
}

type addFlightMessage struct {
	Name string
}

type ticketMessage struct {
	Flightname    string
	Passengername string
}

type commentMessage struct {
	Flightname    string
	Passengername string
	Score         int
	Text          string
}

type getMessgeWithAverage struct {
	Average map[string]float64
	Message string
}
type getMessgeWithAveragewithFlight struct {
	Average float64
	Message string
}

type getMessgeWithoutAverage struct {
	Message string
	Texts   map[string][]string
}
type getMessgeWithoutAveragewithFlight struct {
	Message string
	Texts   []string
}

type Server struct {
	portNumber     int
	instance       *echo.Echo
	surveyInstance *Survey
}

func NewServer(port int) *Server {

	e := echo.New()
	f := NewSurvey()
	return &Server{
		portNumber:     port,
		instance:       e,
		surveyInstance: f,
	}
}

func (s *Server) HealthCheck(c echo.Context) error {
	c.JSON(200, healthMessage{Message: "OK"})
	return nil
}

func (s *Server) addFlight(c echo.Context) error {
	newMes := new(addFlightMessage)
	if err := c.Bind(newMes); err != nil {
		c.JSON(http.StatusBadRequest, healthMessage{Message: "bad request"})
		return err
	}
	err := s.surveyInstance.AddFlight(newMes.Name)

	if err != nil {
		c.JSON(http.StatusBadRequest, healthMessage{Message: "bad request"})

	}
	c.JSON(201, healthMessage{Message: "OK"})
	return nil
}

func (s *Server) addTicket(c echo.Context) error {

	newMes := new(ticketMessage)

	if err := c.Bind(newMes); err != nil {
		c.JSON(http.StatusBadRequest, healthMessage{Message: "bad request"})
		return err
	}

	err := s.surveyInstance.AddTicket(newMes.Flightname, newMes.Passengername)
	if err != nil {
		c.JSON(http.StatusBadRequest, healthMessage{Message: "bad request"})

	}
	c.JSON(201, healthMessage{Message: "OK"})
	return nil

}

func (s *Server) addComment(c echo.Context) error {

	newMes := new(commentMessage)

	if err := c.Bind(newMes); err != nil {
		c.JSON(http.StatusBadRequest, healthMessage{Message: "bad request"})
		return err
	}

	err := s.surveyInstance.AddComment(newMes.Flightname, newMes.Passengername, Comment{Text: newMes.Text, Score: newMes.Score})
	if err != nil {
		c.JSON(http.StatusBadRequest, healthMessage{Message: "bad request"})

	}
	c.JSON(201, healthMessage{Message: "OK"})
	return nil
}

func (s *Server) getComment(c echo.Context) error {
	average := c.QueryParam("average")

	if average == "true" {
		e := c.JSON(200, getMessgeWithAverage{Message: "OK", Average: s.surveyInstance.GetAllCommentsAverage()})
		if e != nil {
			return fmt.Errorf("error in sending json data of average with this error: %w\n", e)
		}

	} else {

		e := c.JSON(200, getMessgeWithoutAverage{Message: "OK", Texts: s.surveyInstance.GetAllComments()})
		if e != nil {
			return fmt.Errorf("error in sending json data without average with this error: %w\n", e)
		}

	}
	return nil
}

func (s *Server) getCommentwithFlight(c echo.Context) error {
	average := c.QueryParam("average")
	FlightName := c.Param("flightname")

	if average == "true" {

		ave, err := s.surveyInstance.GetCommentsAverage(FlightName)
		if err != nil {
			return fmt.Errorf("coudnt to retrive average for flight %s. this is error: %w\n", FlightName, err)
		}

		e := c.JSON(200, getMessgeWithAveragewithFlight{Message: "OK", Average: ave})

		if e != nil {
			return fmt.Errorf("error in sending json data of average with this error: %w\n", e)
		}

	} else {

		ave, err := s.surveyInstance.GetComments(FlightName)
		if err != nil {
			return fmt.Errorf("coudnt to retrive cooments for flight %s. this is error: %w\n", FlightName, err)
		}

		e := c.JSON(200, getMessgeWithoutAveragewithFlight{Message: "OK", Texts: ave})
		if e != nil {
			return fmt.Errorf("error in sending json data of comment  with this error: %w\n", e)
		}

	}
	return nil
}

// blocking // I must review the blocking and non blocking here aziCom
func (server *Server) Start() {
	listenAddress := fmt.Sprintf(":%d", server.portNumber)
	server.instance.GET("/", server.HealthCheck)
	server.instance.POST("/flights", server.addFlight)
	server.instance.POST("/tickets", server.addTicket)
	server.instance.POST("/comments", server.addComment)
	server.instance.GET("/comments", server.getComment)
	server.instance.GET("/comments/:flightname", server.getCommentwithFlight)
	server.instance.Start(listenAddress)

}

// /aziCom: better to pass the server or survey to the handler as argument receiver?
func main() {
	e := NewServer(8080)
	e.Start()

}
