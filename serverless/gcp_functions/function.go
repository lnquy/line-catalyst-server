// Package p contains an HTTP Cloud Function.
package p

import (
	"log"
	"net/http"
	"os"
)

// Wake sends a request to Heroku's Catalyst dyno to make sure it's not going inactive.
func Wake(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	resp, err := http.Get(os.Getenv("CATALYST_URL"))
	if err != nil {
		log.Printf("Failed to make request to Catalyst: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Catalyst response with not OK status: %s", resp.Status)
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}

	log.Println("Catalyst is active")
	w.Write([]byte("Catalyst is active"))
}
