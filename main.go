package main

import (
	"fmt"
	"net/http"

	timezone "github.com/evanoberholster/timezoneLookup"
	"github.com/gorilla/mux"
	"racemap.com/tz-service/handlers"
	"racemap.com/tz-service/logger"
)

func main() {
	// init log
	log := logger.InitLogger()
	log.Info("Init Logger Instance")

	logMiddleware := logger.BuildMiddleware(log)

	tzService, err := timezone.LoadTimezones(timezone.Config{
		DatabaseType: "memory",
		DatabaseName: "assets/timezone",
		Snappy:       true,
	})

	if err != nil {
		fmt.Println(err)
	}
	log.Info("Init Timezone Database")

	// build handlers for routes
	tzHandler := handlers.TimezoneHandler(tzService)
	statusHandler := handlers.StatusHandler()

	r := mux.NewRouter()
	r.HandleFunc("/api", tzHandler)
	r.HandleFunc("/status", statusHandler)

	// add middlewares
	r.Use(logMiddleware)

	http.Handle("/", r)

	var port = "8000"

	log.Info("Start HTTP Server on Port " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))

	log.Info("Stopped Server")
}
