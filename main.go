package main

import (
	"net/http"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"racemap.com/tz-service/handlers"
	"racemap.com/tz-service/logger"
	"racemap.com/tz-service/timezones"
)

func main() {
	// init log
	log := logger.InitLogger()
	log.Info("Init Logger Instance")

	logMiddleware := logger.BuildMiddleware(log)

	log.Info("Begin to load Timezone Database")
	tzService, err := timezones.InitTimezoneService()

	if err != nil {
		panic(err)
	}
	defer tzService.Close()
	log.Info("Init Timezone Database")

	// build handlers for routes
	tzHandler := handlers.TimezoneHandler(tzService, log)
	statusHandler := handlers.StatusHandler()

	r := mux.NewRouter()
	r.HandleFunc("/api", tzHandler)
	r.HandleFunc("/status", statusHandler)

	// add middlewares
	r.Use(logMiddleware)

	http.Handle("/", r)

	var port = "8080"

	log.Info("Start HTTP Server on Port " + port)
	log.Fatal(http.ListenAndServe(":"+port, gorillaHandlers.CORS()(r)))

	log.Info("Stopped Server")
}
