package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const portNoSample = 1235

var addrSample string = fmt.Sprintf("http://127.0.0.1:%d/", portNoSample)
var ctjSample string = `application/json`

func TestSample(t *testing.T) {
	server := NewServer(portNoSample)
	go server.Start()
	time.Sleep(10000 * time.Millisecond)

	{
		resp, err := http.DefaultClient.Post(addrSample+"flights", ctjSample, bytes.NewBufferString(`{ "Name" : "f9" }`))
		assert.NoError(t, err)
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"Message" : "OK"}`, string(body))
		assert.Equal(t, resp.StatusCode, 201)
	}

	{
		resp, err := http.DefaultClient.Post(addrSample+"tickets", ctjSample,
			bytes.NewBufferString(`{ "FlightName" : "f9", "PassengerName" : "p1"}`))
		assert.NoError(t, err)
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"Message" : "OK"}`, string(body))
		assert.Equal(t, resp.StatusCode, 201)
	}
	{
		resp, err := http.DefaultClient.Post(addrSample+"tickets", ctjSample,
			bytes.NewBufferString(`{ "FlightName" : "f9", "PassengerName" : "p2"}`))
		assert.NoError(t, err)
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"Message" : "OK"}`, string(body))
		assert.Equal(t, resp.StatusCode, 201)
	}
	{
		resp, err := http.DefaultClient.Post(addrSample+"comments", ctjSample,
			bytes.NewBufferString(`{ "FlightName" : "f9", "PassengerName" : "p1", "Score" : 8, "Text" : "random comment" }`))
		assert.NoError(t, err)
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"Message" : "OK"}`, string(body))
		assert.Equal(t, resp.StatusCode, 201)
	}

	{
		resp, err := http.DefaultClient.Get(addrSample + "comments/f9?average=true")
		assert.NoError(t, err)
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"Message" : "OK", "Average" : 8 }`, string(body))
		assert.Equal(t, resp.StatusCode, 200)
	}

	{
		resp, err := http.DefaultClient.Get(addrSample + "comments/f9?average=false")
		assert.NoError(t, err)
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"Message" : "OK", "Texts" : [ "random comment" ]}`, string(body))
		assert.Equal(t, resp.StatusCode, 200)
	}

}
