package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/joho/godotenv/autoload" // Autoloads .env file if it exists
)

//
// Main entry point, will start HTTP service
//
func main() {
	r := mux.NewRouter()

	a := api{
		message: "hello",
	}
	r.PathPrefix("/api/set").HandlerFunc(a.setAPI)
	r.PathPrefix("/api/get").HandlerFunc(a.getAPI)

	log.Printf("### Server listening on %v\n", 8000)
	err := http.ListenAndServe(fmt.Sprintf(":%d", 8000), r)
	if err != nil {
		panic(err.Error())
	}
}

type api struct {
	message string
}

func (a *api) setAPI(resp http.ResponseWriter, req *http.Request) {
	fmt.Println(a)
	a.message = "I have been changed"
	fmt.Println(a)
}

func (a *api) getAPI(resp http.ResponseWriter, req *http.Request) {
	fmt.Println(a)
	resp.Write([]byte("Message is " + a.message))
}
