package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

//
// GetState gets state
//
func GetState(resp http.ResponseWriter, port int, store string, service string, key string) (data []byte, err error) {
	daprURL := fmt.Sprintf("http://localhost:%d/v1.0/state/%s/%s", port, store, key)
	//fmt.Println(daprURL)

	daprResp, err := http.Get(daprURL)
	// fmt.Printf("### STATE GET DEBUG RESP %+v\n", daprResp)
	// fmt.Printf("### STATE GET DEBUG ERR %+v\n", err)
	if err != nil || (daprResp.StatusCode < 200 || daprResp.StatusCode > 299) {
		SendDaprProblem(daprURL, resp, daprResp, err, service)
		return nil, errors.New("Failed to get state object from Dapr")
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
		Problem{"json-error", "State JSON marshalling error", 500, err.Error(), service}.HttpSend(resp)
		return
	}
	//fmt.Printf("### STATE SAVE DEBUG PAYLOAD %+v\n", string(jsonPayload))

	daprURL := fmt.Sprintf("http://localhost:%d/v1.0/state/%s", port, store)
	//fmt.Printf("### STATE SAVE DEBUG URL %+v\n", daprURL)
	daprResp, err := http.Post(daprURL, "application/json", bytes.NewBuffer(jsonPayload))
	// fmt.Printf("### STATE SAVE DEBUG RESP %+v\n", daprResp)
	// fmt.Printf("### STATE SAVE DEBUG ERR %+v\n", err)

	if err != nil || (daprResp.StatusCode < 200 || daprResp.StatusCode > 299) {
		SendDaprProblem(daprURL, resp, daprResp, err, service)
		return err
	}
	return nil
}
