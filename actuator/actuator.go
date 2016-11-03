package actuator

import (
	"net/http"
)

func InitializeActuator() {
	loadGitInfo()

	http.HandleFunc("/info", infoHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/metrics", metricsHandler)
	http.HandleFunc("/env", envHandler)

}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(info()))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(health()))
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(metrics()))
}

func envHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(env()))
}