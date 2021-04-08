// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr API helper/wrapper for state and pub/sub
// Designed to return standard formatted errors (pkg/problem)
// ----------------------------------------------------------------------------

package dapr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/benc-uk/dapr-store/pkg/env"
	"github.com/benc-uk/dapr-store/pkg/problem"
	"github.com/gorilla/mux"
)

const (
	getStateURL      = "http://localhost:%d/v1.0/state/%s/%s"
	saveStateURL     = "http://localhost:%d/v1.0/state/%s"
	publishURL       = "http://localhost:%d/v1.0/publish/%s/%s"
	outputBindingURL = "http://localhost:%d/v1.0/bindings/%s"
	invokeURL        = "http://localhost:%d/v1.0/invoke/%s/method/%s"
)

type state struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type topic struct {
	PubSubName string `json:"pubsubname"`
	Topic      string `json:"topic"`
	Route      string `json:"route"`
}

type bindingOut struct {
	Metadata  map[string]string `json:"metadata"`
	Data      interface{}       `json:"data"`
	Operation string            `json:"operation"`
}

// Helper is our main struct
type Helper struct {
	Port        int
	ServiceName string
}

//
// NewHelper returns a new Dapr helper
//
func NewHelper(appName string) *Helper {
	// Fall back to default Dapr port of 3500
	daprPort := env.GetEnvInt("DAPR_HTTP_PORT", 3500)

	// Check for Dapr existence
	time.AfterFunc(time.Second*15, func() {
		daprResp, err := http.Get(fmt.Sprintf("http://localhost:%d/v1.0/healthz", daprPort))
		if err != nil || daprResp.StatusCode != 204 {
			log.Printf("### WARNING! Dapr process/sidecar NOT found for %s\n", appName)
		} else {
			log.Printf("### Dapr process/sidecar found for %s on port: %d", appName, daprPort)
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
func (h *Helper) PublishMessage(pubSubName string, topicName string, message interface{}) *problem.Problem {
	jsonPayload, err := json.Marshal(message)
	if err != nil {
		return problem.New("err://json-marshall", "Malformed JSON", 400, "Message could not be marshalled to JSON", h.ServiceName)
	}

	log.Printf("### Pub/sub helper, topic:%s payload:%+v\n", topicName, string(jsonPayload))

	daprURL := fmt.Sprintf(publishURL, h.Port, pubSubName, topicName)
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
		Metadata:  metadata,
		Data:      data,
		Operation: "create",
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

//
// RegisterTopicSubscriptions is a HTTP handler that lets Dapr know what topics we subscribe to
//
func (h *Helper) RegisterTopicSubscriptions(pubSubName string, topics []string, router *mux.Router) {
	router.HandleFunc("/dapr/subscribe", func(resp http.ResponseWriter, req *http.Request) {

		topicList := []topic{}
		for _, t := range topics {
			topicList = append(topicList, topic{
				PubSubName: pubSubName,
				Topic:      t,
				Route:      fmt.Sprintf("/receive/%s", t),
			})
		}
		json, _ := json.Marshal(topicList)

		resp.Header().Set("Content-Type", "application/json")
		_, _ = resp.Write(json)
	})
}

//
// RegisterTopicReceiver is a way to plug in a handler for receiving messages from a topic
//
func (h *Helper) RegisterTopicReceiver(topic string, router *mux.Router, handler func(body io.Reader) error) {
	route := fmt.Sprintf("/receive/%s", topic)
	router.HandleFunc(route, func(resp http.ResponseWriter, req *http.Request) {
		type cloudevent struct {
			ID   string      `json:"id"`
			Data interface{} `json:"data"`
		}
		event := &cloudevent{}
		var bodyBytes []byte
		bodyBytes, _ = ioutil.ReadAll(req.Body)

		// Basic validation checks
		err := json.Unmarshal(bodyBytes, &event)
		if err != nil {
			// Returning a non-200 will reschedule the received message
			problem.New("err://json-decode", "Event JSON decoding error", 500, err.Error(), h.ServiceName).Send(resp)
			return
		}
		log.Printf("### Received event from pub/sub topic %s %s", topic, event.ID)

		// This is a trick to reset the body reader
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// Pass the body to the handler
		// It would be really nice to pass the decoded data object/struct but we don't know the type
		err = handler(req.Body)
		if err != nil {
			// Returning a non-200 will reschedule the received message
			problem.New("err://handler-failed", "Message hander returned an error", 500, err.Error(), h.ServiceName).Send(resp)
			return
		}
	})
}

//
// InvokeGet calls another service with a HTTP GET
//
func (h *Helper) InvokeGet(service, method string) (*http.Response, error) {
	daprURL := fmt.Sprintf(invokeURL, h.Port, service, method)
	daprResp, err := http.Get(daprURL)
	return daprResp, err
}
