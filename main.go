package main

type Comment struct {
	// do not modify or remove these fields
	Score int
	Text  string
	// but you can add anything you want
}

type Survey struct {
}

func NewSurvey() *Survey {
	return &Survey{}
}

func (s *Survey) AddFlight(flightName string) error {
	return nil
}

func (s *Survey) AddTicket(flightName, passengerName string) error {
	return nil
}

func (s *Survey) AddComment(flightName, passengerName string, comment Comment) error {
	return nil
}

func (s *Survey) GetCommentsAverage(flightName string) (float64, error) {
	return 0, nil
}

func (s *Survey) GetAllCommentsAverage() map[string]float64 {
	return nil
}

func (s *Survey) GetComments(flightName string) ([]string, error) {
	return nil, nil
}

func (s *Survey) GetAllComments() map[string][]string {
	return nil
}
