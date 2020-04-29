// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr API helper/wrapper for state and pub/sub - returns standard formatted errors
// ----------------------------------------------------------------------------

package dapr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/benc-uk/dapr-store/pkg/env"
	"github.com/benc-uk/dapr-store/pkg/problem"
)

const (
	getStateURL      = "http://localhost:%d/v1.0/state/%s/%s"
	saveStateURL     = "http://localhost:%d/v1.0/state/%s"
	publishURL       = "http://localhost:%d/v1.0/publish/%s"
	outputBindingURL = "http://localhost:%d/v1.0/bindings/%s"
)

// DaprState is the payload for the Dapr state API
type state struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type bindingOut struct {
	Metadata map[string]string `json:"metadata"`
	Data     interface{}       `json:"data"`
}

// Helper is our main struct
type Helper struct {
	Port        int
	ServiceName string
}

// NewHelper returns a new Dapr helper
func NewHelper(appName string) *Helper {
	// Fall back to default Dapr port of 3500
	daprPort := env.GetEnvInt("DAPR_HTTP_PORT", 3500)

	// Check for Dapr existence
	time.AfterFunc(time.Second*15, func() {
		daprResp, err := http.Get(fmt.Sprintf("http://localhost:%d/v1.0/healthz", daprPort))
		if err != nil || daprResp.StatusCode != 200 {
			log.Println("### WARNING! Dapr process/sidecar NOT found")
		} else {
			log.Printf("### Dapr process/sidecar found on port: %d", daprPort)
		}
	})

	return &Helper{
		Port:        daprPort,
		ServiceName: appName,
	}
}

//
// GetState returns the state of given key
//
func (h *Helper) GetState(storeName, key string) ([]byte, *problem.Problem) {
	daprURL := fmt.Sprintf(getStateURL, h.Port, storeName, key)

	daprResp, err := http.Get(daprURL)
	if err != nil || (daprResp.StatusCode < 200 || daprResp.StatusCode > 299) {
		return nil, problem.NewAPIProblem(daprURL, "Dapr get state failed", h.ServiceName, daprResp, err)
	}

	defer daprResp.Body.Close()
	body, _ := ioutil.ReadAll(daprResp.Body)
	return body, nil
}

//
// SaveState stores value as serialized state into Dapr
//
func (h *Helper) SaveState(storeName, key string, value interface{}) *problem.Problem {
	daprPayload := state{
		Key:   key,
		Value: value,
	}

	jsonPayload, err := json.Marshal([]state{daprPayload})
	if err != nil {
		return problem.NewAPIProblem("err://json-marshall", "State JSON marshalling error", h.ServiceName, nil, err)
	}

	log.Printf("### State save helper, key:%s payload:%+v\n", key, string(jsonPayload))

	daprURL := fmt.Sprintf(saveStateURL, h.Port, storeName)
	daprResp, err := http.Post(daprURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil || (daprResp.StatusCode < 200 || daprResp.StatusCode > 299) {
		return problem.NewAPIProblem(daprURL, "Dapr save state failed", h.ServiceName, daprResp, err)
	}

	// All good
	return nil
}

//
// PublishMessage pushes a message to the given topic
//
func (h *Helper) PublishMessage(queueName string, message interface{}) *problem.Problem {
	jsonPayload, err := json.Marshal(message)
	if err != nil {
		return problem.New("err://json-marshall", "Malformed JSON", 400, "Message could not be marshalled to JSON", h.ServiceName)
	}

	daprURL := fmt.Sprintf(publishURL, h.Port, queueName)
	daprResp, err := http.Post(daprURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return problem.NewAPIProblem(daprURL, "Error publishing message", h.ServiceName, daprResp, err)
	}

	// All good
	return nil
}

//
// SendOutput sends some output to a binding
//
func (h *Helper) SendOutput(bindingName string, data interface{}, metadata map[string]string) *problem.Problem {
	daprPayload := bindingOut{
		Metadata: metadata,
		Data:     data,
	}

	jsonPayload, err := json.Marshal(daprPayload)
	if err != nil {
		return problem.NewAPIProblem("err://json-marshall", "State JSON marshalling error", h.ServiceName, nil, err)
	}

	daprURL := fmt.Sprintf(outputBindingURL, h.Port, bindingName)
	daprResp, err := http.Post(daprURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return problem.NewAPIProblem(daprURL, "Error sending output", h.ServiceName, daprResp, err)
	}

	// All good
	return nil
}
