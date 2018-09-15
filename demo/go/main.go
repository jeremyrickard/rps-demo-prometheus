// Copyright 2018, OpenCensus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	os "os"

	ocagent "contrib.go.opencensus.io/exporter/ocagent"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	"go.opencensus.io/trace"
)

func main() {
	// Register stats and trace exporters to export the collected data.
	serviceName := os.Getenv("SERVICE_NAME")
	if len(serviceName) == 0 {
		serviceName = "go-app"
	}

	agentHostName := os.Getenv("OCAGENT_TRACE_EXPORTER_ENDPOINT")
	if len(agentHostName) == 0 {
		agentHostName = "localhost"
	}

	exporter, err := ocagent.NewExporter(ocagent.WithInsecure(), ocagent.WithServiceName(serviceName), ocagent.WithAddress(agentHostName))
	if err != nil {
		log.Printf("Failed to create the agent exporter: %v", err)
	}

	trace.RegisterExporter(exporter)

	// Always trace for this demo. In a production application, you should
	// configure this to a trace.ProbabilitySampler set at the desired
	// probability.
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	client := &http.Client{Transport: &ochttp.Transport{Propagation: &tracecontext.HTTPFormat{}}}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "hello world")

		var jsonStr = []byte(`[ { "url": "http://blank.org", "arguments": [] } ]`)
		r, _ := http.NewRequest("POST", "http://lmolkova-oc-test.azurewebsites.net/api/forward", bytes.NewBuffer(jsonStr))
		r.Header.Set("Content-Type", "application/json")
		// Propagate the trace header info in the outgoing requests.
		r = r.WithContext(req.Context())
		resp, err := client.Do(r)
		if err != nil {
			log.Println(err)
		} else {
			// TODO: handle response
			resp.Body.Close()
		}
	})
	http.HandleFunc("/call_blank", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "hello world")

		r, _ := http.NewRequest("GET", "http://blank.org", nil)

		// Propagate the trace header info in the outgoing requests.
		r = r.WithContext(req.Context())
		resp, err := client.Do(r)
		if err != nil {
			log.Println(err)
		} else {
			// TODO: handle response
			resp.Body.Close()
		}
	})
	log.Fatal(http.ListenAndServe(":50030", &ochttp.Handler{Propagation: &tracecontext.HTTPFormat{}}))
}
