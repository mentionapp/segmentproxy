package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/arnaud-lb/segmentproxy"
	"github.com/segmentio/analytics-go"
)

func main() {

	var writeKey string

	flag.StringVar(&writeKey, "segment-write-key", "", "")
	flag.Parse()

	if writeKey == "" {
		flag.Usage()
		os.Exit(1)
	}

	conf := segmentproxy.Config{
		Segment: analytics.New(writeKey),
	}

	segmentproxy.Register(&conf, http.DefaultServeMux)

	log.Print("Listening on :8080 for incoming webhooks to forward to segment.com")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
