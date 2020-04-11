// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service
// ----------------------------------------------------------------------------

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/benc-uk/dapr-store/common"
	"github.com/benc-uk/go-starter/pkg/envhelper"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload" // Autoloads .env file if it exists
)

// API type is a wrap of he common API with local version
type API struct {
	*common.API
}

var (
	healthy        = true               // Simple health flag
	version        = "0.0.1"            // App version number, set at build time with -ldflags "-X 'main.version=1.2.3'"
	buildInfo      = "No build details" // Build details, set at build time with -ldflags "-X 'main.buildInfo=Foo bar'"
	serviceName    = "products"
	daprPort       int
	daprStateStore string
	products       map[string]common.Product
)

//
// Main entry point, will start HTTP service
//
func main() {
	log.SetOutput(os.Stdout) // Personal preference on log output
	log.Printf("### Dapr Store: %v v%v starting...", serviceName, version)

	// Port to listen on, change the default as you see fit
	serverPort := envhelper.GetEnvInt("PORT", 9002)
	daprStateStore = envhelper.GetEnvString("DAPR_STORE_NAME", "statestore")

	daprPort = envhelper.GetEnvInt("DAPR_HTTP_PORT", 0)
	if daprPort != 0 {
		log.Printf("### Dapr sidecar detected on port %v", daprPort)
	} else {
		log.Printf("### Dapr not detected (no DAPR_HTTP_PORT available), this is bad")
		log.Printf("### Exiting...")
		os.Exit(1)
	}

	// Use gorilla/mux for routing
	router := mux.NewRouter()

	// Add middleware for logging and CORS
	router.Use(corsMiddleware)
	router.Use(loggingMiddleware)

	// Wrapper type with anonymous inner field
	api := API{
		&common.API{
			ServiceName: serviceName,
			Healthy:     healthy,
			Version:     version,
			BuildInfo:   buildInfo,
		}}

	api.AddCommonRoutes(router)
	api.addRoutes(router)

	// Load all state into memory as Dapr state API is really basic
	time.AfterFunc(3*time.Second, func() {
		http.Get(fmt.Sprintf("http://localhost:%d/reload", serverPort))
		log.Printf("### Loaded %d product data items from state into memory\n", len(products))
	})

	// Start server
	log.Printf("### Dapr state store: %v\n", daprStateStore)
	log.Printf("### Server listening on %v\n", serverPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", serverPort), router)
	if err != nil {
		panic(err.Error())
	}
}

//
// Change CORS settings here
//
func corsMiddleware(handler http.Handler) http.Handler {
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	corsOrigins := handlers.AllowedOrigins([]string{"*"})
	return handlers.CORS(corsOrigins, corsMethods)(handler)
}

//
// Change request logging here
//
func loggingMiddleware(handler http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, handler)
}

//
// Return the product catalog
//
// func loadState() {
// 	fmt.Println("LOOOOOOOOOOOOOAAAAAAAAAADDDDDDDDDDDDDDDDDDD")

// 	// Fake HTTP context
// 	r := httptest.NewRecorder()

// 	// Load index which is just an array of keys/ids
// 	data, err := common.GetState(r, daprPort, daprStateStore, serviceName, "index")
// 	if err != nil {
// 		fmt.Println("### Error loading product index", err.Error())
// 		return
// 	}
// 	productIds := []int{}
// 	err = json.Unmarshal(data, &productIds)
// 	if err != nil {
// 		fmt.Println("### Error decoding product index", err.Error())
// 		return
// 	}

// 	// Load each object, and push into array
// 	products = make(map[string]common.Product)
// 	for _, id := range productIds {
// 		p := common.Product{}
// 		data, err := common.GetState(r, daprPort, daprStateStore, serviceName, strconv.Itoa(id))
// 		if err != nil {
// 			fmt.Printf("### Error loading product '%v' %+v\n", id, err.Error())
// 			continue
// 		}
// 		err = json.Unmarshal(data, &p)
// 		if err != nil {
// 			fmt.Printf("### Error decoding product '%v' %+v\n", id, err.Error())
// 			continue
// 		}
// 		products[strconv.Itoa(id)] = p
// 	}

// 	fmt.Println(products)
// }
