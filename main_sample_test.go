package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlightAddSample(t *testing.T) {
	var err error
	s := NewSurvey()
	err = s.AddFlight("my-flight")
	assert.NoError(t, err)
	err = s.AddFlight("my-flight")
	assert.Error(t, err)
}

func TestTicketAddSample(t *testing.T) {
	var err error
	s := NewSurvey()

	_ = s.AddFlight("abc")

	err = s.AddTicket("abc", "p2")
	assert.NoError(t, err)
}

func TestSample(t *testing.T) {
	var err error
	s := NewSurvey()

	err = s.AddFlight("my-flight")
	assert.NoError(t, err)
	err = s.AddFlight("my-flight")
	assert.Error(t, err)

	err = s.AddTicket("my-flight", "p1")
	assert.NoError(t, err)
	err = s.AddTicket("my-flight", "p2")
	assert.NoError(t, err)

	err = s.AddComment("my-flight", "p1", Comment{Score: 9, Text: "good"})
	assert.NoError(t, err)

	err = s.AddComment("other flight", "p2", Comment{Score: 8, Text: "ok"})
	assert.Error(t, err)

	assert.EqualValues(t, 9, s.GetAllCommentsAverage()["my-flight"])

	comments, err := s.GetComments("my-flight")
	assert.NoError(t, err)
	assert.Equal(t, []string{"good"}, comments)
	assert.Equal(t, []string{"good"}, s.GetAllComments()["my-flight"])
}
