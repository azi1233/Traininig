package main

import (
	"errors"
	"sync"
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
	_, exists := s.S[flightName]
	if exists {
		return errors.New("flight exists")
	} else {
		s.mu.Lock()
		s.S[flightName] = passObj{&sync.Mutex{}, make(map[string]Comment)}
		s.mu.Unlock()
		return nil
	}
}

func (s *Survey) AddTicket(flightName, passengerName string) error {

	_, flighexists := s.S[flightName]

	if !flighexists {
		return errors.New("flight doesn't exists, unable to issue the ticket")
	}
	_, passengerExistInFlight := s.S[flightName].P[passengerName]
	if flighexists && passengerExistInFlight {
		return errors.New("duplicate ticket")

	}

	s.S[flightName].mu.Lock()
	s.S[flightName].P[passengerName] = Comment{}
	s.S[flightName].mu.Unlock()

	return nil
}

func (s *Survey) AddComment(flightName, passengerName string, comment Comment) error {

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
	s.S[flightName].mu.Lock()
	s.S[flightName].P[passengerName] = comment
	s.S[flightName].mu.Unlock()

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
