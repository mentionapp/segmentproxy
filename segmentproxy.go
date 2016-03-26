package segmentproxy

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/segmentio/analytics-go"
)

type Config struct {
	Segment *analytics.Client
	// Prefix is prepended to URLs in Register()
	Prefix string
	// EmailToID must be a function converting emails to ids
	EmailToID func(email string) (string, error)
}

type Mux interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

func Register(config *Config, mux Mux) {

	http.HandleFunc(config.Prefix+"/identify", func(w http.ResponseWriter, r *http.Request) {
		handle(&Identify{}, config, w, r)
	})

	http.HandleFunc(config.Prefix+"/group", func(w http.ResponseWriter, r *http.Request) {
		handle(&Group{}, config, w, r)
	})

	http.HandleFunc(config.Prefix+"/track", func(w http.ResponseWriter, r *http.Request) {
		handle(&Track{}, config, w, r)
	})
}

func handle(action Action, config *Config, w http.ResponseWriter, r *http.Request) {

	log.Printf("Handling request %v", r)

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed reading request body: %v", err)
		http.Error(w, "Failed reading request body", http.StatusBadRequest)
		return
	}

	if err := action.Unmarshal(buf); err != nil {
		log.Printf("Unmarshal failed: %v. json: %v", err, string(buf))
		http.Error(w, "Failed decoding body", http.StatusBadRequest)
		return
	}

	if email := action.GetEmail(); email != "" {
		id, err := config.EmailToID(email)
		if err != nil {
			log.Printf("EmailToID failed: %v", err)
			http.Error(w, "EmailToID failed", http.StatusInternalServerError)
			return
		}
		if id == "" {
			// We don't want the webhook client to retry
			http.Error(w, "No such user", http.StatusOK)
			return
		}
		action.SetUserID(id)
	}

	if err := action.Send(config.Segment); err != nil {
		msg := fmt.Sprintf("action.Send failed: %s", err)
		log.Print(err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	http.Error(w, "OK", http.StatusOK)
}
