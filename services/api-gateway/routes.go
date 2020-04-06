package main

//
// Basic REST API microservice, template/reference code
// Ben Coleman, July 2019, v1
//

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//
// Call Orders service
//
func (h API) ordersAPI(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "application/json")

	url := fmt.Sprintf("http://localhost:%v/v1.0/invoke/orders/method/sss", daprPort)
	fmt.Println(url)
	r, err := http.Get(url)
	if err != nil {
		fmt.Println("ERROR")
		resp.Write([]byte("{error: \"BAD\"}"))
		resp.WriteHeader(r.StatusCode)
		return
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(body))
	resp.WriteHeader(r.StatusCode)
	resp.Write(body)
}
