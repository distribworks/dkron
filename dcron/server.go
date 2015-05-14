package dcron

import (
	"encoding/json"
	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
	"github.com/tylerb/graceful"
	"net/http"
	"time"
)

func ServerInit() {
	loadConfig()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)

	middle := interpose.New()
	middle.UseHandler(router)

	srv := &graceful.Server{
		Timeout: 1 * time.Second,
		Server:  &http.Server{Addr: ":8080", Handler: middle},
	}

	log.Infoln("Running HTTP server on 8080")

	certFile := config.GetString("certFile")
	keyFile := config.GetString("keyFile")
	if certFile != "" && keyFile != "" {
		srv.ListenAndServeTLS(certFile, keyFile)
	} else {
		srv.ListenAndServe()
	}
	log.Debug("Exiting")
}

func Index(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(jobs)
}
