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
type DaprState struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type DaprBindingOut struct {
	Metadata map[string]string `json:"metadata"`
	Data     interface{}       `json:"data"`
}

// Helper is our man struct
type Helper struct {
	Port            int
	AppInstanceName string
	StateStoreName  string
	PubsubQueueName string
}

// NewHelper returns a new dapr helper
func NewHelper(daprPort int, appName, store, queue string) *Helper {
	storeName := "statestore"
	if store != "" {
		storeName = store
	}
	queueName := "queue"
	if store != "" {
		queueName = queue
	}

	return &Helper{
		Port:            daprPort,
		AppInstanceName: appName,
		StateStoreName:  storeName,
		PubsubQueueName: queueName,
	}
}

func BootstrapHelper(appName string) *Helper {
	daprStoreName := env.GetEnvString("DAPR_STORE_NAME", "statestore")
	daprTopicName := env.GetEnvString("DAPR_ORDERS_TOPIC", "orders-queue")

	daprPort := env.GetEnvInt("DAPR_HTTP_PORT", 0)
	if daprPort != 0 {
		log.Printf("### Dapr sidecar detected on port %v", daprPort)
	} else {
		log.Printf("### Dapr not detected (no DAPR_HTTP_PORT available), this is bad")
		log.Printf("### Exiting...")
		return nil
	}

	log.Printf("### Dapr store name: %s\n", daprStoreName)
	log.Printf("### Dapr orders topic: %s\n", daprTopicName)
	return NewHelper(daprPort, appName, daprStoreName, daprTopicName)
}

//
// GetState returns the state of given key
//
func (h *Helper) GetState(key string) ([]byte, *problem.Problem) {
	daprURL := fmt.Sprintf(getStateURL, h.Port, h.StateStoreName, key)

	daprResp, err := http.Get(daprURL)
	if err != nil || (daprResp.StatusCode < 200 || daprResp.StatusCode > 299) {
		return nil, problem.NewAPIProblem(daprURL, "Dapr get state failed", h.AppInstanceName, daprResp, err)
	}

	defer daprResp.Body.Close()
	body, _ := ioutil.ReadAll(daprResp.Body)
	return body, nil
}

//
// SaveState stores value as serialized state into Dapr
//
func (h *Helper) SaveState(key string, value interface{}) *problem.Problem {
	daprPayload := DaprState{
		Key:   key,
		Value: value,
	}

	jsonPayload, err := json.Marshal([]DaprState{daprPayload})
	if err != nil {
		return problem.NewAPIProblem("err://json-marshall", "State JSON marshalling error", h.AppInstanceName, nil, err)
	}

	log.Printf("### State save helper, key:%s payload:%+v\n", key, string(jsonPayload))

	daprURL := fmt.Sprintf(saveStateURL, h.Port, h.StateStoreName)
	daprResp, err := http.Post(daprURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil || (daprResp.StatusCode < 200 || daprResp.StatusCode > 299) {
		return problem.NewAPIProblem(daprURL, "Dapr save state failed", h.AppInstanceName, daprResp, err)
	}

	// All good
	return nil
}

//
// PublishMessage pushes a message to the given topic
//
func (h *Helper) PublishMessage(message interface{}) *problem.Problem {
	jsonPayload, err := json.Marshal(message)
	if err != nil {
		return problem.New("err://json-marshall", "Malformed JSON", 400, "Message could not be marshalled to JSON", h.AppInstanceName)
	}

	daprURL := fmt.Sprintf(publishURL, h.Port, h.PubsubQueueName)
	daprResp, err := http.Post(daprURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return problem.NewAPIProblem(daprURL, "Error publishing message", h.AppInstanceName, daprResp, err)
	}

	// All good
	return nil
}

//
// OutputBinding sends some output to a binding
//
func (h *Helper) SendOutput(bindingName string, data interface{}, metadata map[string]string) *problem.Problem {
	daprPayload := DaprBindingOut{
		Metadata: metadata,
		Data:     data,
	}

	jsonPayload, err := json.Marshal(daprPayload)
	if err != nil {
		return problem.NewAPIProblem("err://json-marshall", "State JSON marshalling error", h.AppInstanceName, nil, err)
	}

	daprURL := fmt.Sprintf(outputBindingURL, h.Port, bindingName)
	daprResp, err := http.Post(daprURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return problem.NewAPIProblem(daprURL, "Error sending output", h.AppInstanceName, daprResp, err)
	}

	// All good
	return nil
}

// orderNotifyPayload := `{
// 	"metadata": {
// 		"ContentType": "text/plain",
// 		"ContentEncoding": "UTF-8",
// 		"blobName": "order_` + order.ID + `.txt"
// 	},
// 	"data": "----------\nOrder title:` + order.Title + `\nOrder ID: ` + order.ID + `\nUser: ` + order.ForUser + `\nAmount: ` + fmt.Sprintf("%f", order.Amount) + `\n----------"
// }`

// // We silently consume errors here
// // - it's an optional component that might not be set up
// blobResp, err := http.Post(fmt.Sprintf("http://localhost:%d/v1.0/bindings/orders-notify", daprHelper.Port), "application/json", strings.NewReader(orderNotifyPayload))
// if err != nil || (blobResp.StatusCode != 200) {
// 	log.Printf("### Warning! Failed to output order notification file, %d %+v", blobResp.StatusCode, err)
// }
