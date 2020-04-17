package state

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/benc-uk/dapr-store/pkg/problem"
)

// DaprState is the payload for the Dapr state API
type DaprState struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

//
// GetState gets state
//
func GetState(resp http.ResponseWriter, port int, store string, service string, key string) (data []byte, err error) {
	daprURL := fmt.Sprintf("http://localhost:%d/v1.0/state/%s/%s", port, store, key)

	daprResp, err := http.Get(daprURL)
	if err != nil || (daprResp.StatusCode < 200 || daprResp.StatusCode > 299) {
		problem.Send("Dapr get state failed", daprURL, resp, daprResp, err, service)
		return nil, errors.New("Failed to call Dapr get state API")
	}

	defer daprResp.Body.Close()
	body, _ := ioutil.ReadAll(daprResp.Body)
	return body, nil
}

//
// SaveState gets state
//
func SaveState(resp http.ResponseWriter, port int, store string, service string, key string, value interface{}) (err error) {
	daprPayload := DaprState{
		Key:   key,
		Value: value,
	}

	jsonPayload, err := json.Marshal([]DaprState{daprPayload})
	if err != nil {
		problem.Send("State JSON marshalling error", "err://json-marshall", resp, nil, err, service)
		return errors.New("Failed to call Dapr save state API")
	}

	log.Printf("### State save helper, key:%s payload:%+v\n", key, string(jsonPayload))

	daprURL := fmt.Sprintf("http://localhost:%d/v1.0/state/%s", port, store)
	daprResp, err := http.Post(daprURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil || (daprResp.StatusCode < 200 || daprResp.StatusCode > 299) {
		problem.Send("Dapr save state failed", daprURL, resp, daprResp, err, service)
		return errors.New("Failed to call Dapr save state API")
	}
	return nil
}
