package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"test-log-emitter/client"
	"test-log-emitter/config"
	"test-log-emitter/emitters"

	"gopkg.in/yaml.v2"

	flag "github.com/spf13/pflag"
)

func main() {
	var configFilePath *string = flag.String("config", "", "path to config file")
	flag.Parse()

	if *configFilePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	confContents, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	conf := new(config.Config)
	err = yaml.Unmarshal(confContents, conf)
	if err != nil {
		log.Fatal(err)
	}

	loggregatorClient, err := client.NewLoggregatorIngressClient(conf.Loggregator)
	if err != nil {
		log.Fatal(err)
	}
	emitter := emitters.NewSpikeEmitter(loggregatorClient)

	http.HandleFunc("/", ping)
	http.HandleFunc("/spike", http.HandlerFunc(emitter.EmitSpike("spike")))
	http.HandleFunc("/spike", NewSpikeHandler("spike"))
	http.HandleFunc("/spoke", emitSpoke(emitters.NewSpikeEmitter(loggregatorClient)))

	fmt.Printf("Starting cpu usage logger on port %d...", conf.ListenPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", conf.ListenPort), nil); err != nil {
		log.Fatal(err)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	message := "What do you want to emit today?\n"
	message = message + "* POST /spike - posts an app instance spike\n"

	if _, err := io.WriteString(w, message); err != nil {
		http.Error(w, fmt.Sprintf("Failed to resond to ping request: %v", err), http.StatusInternalServerError)
		return
	}
}

func emitSpoke(emitter *emitters.SpikeEmitter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Sorry, only POST methods are supported.", http.StatusMethodNotAllowed)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to read body: %v", err), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		reqMap := map[string]string{}
		if err := json.Unmarshal(body, &reqMap); err != nil {
			http.Error(w, fmt.Sprintf("Failed to unmarshal body: %v", err), http.StatusInternalServerError)
			return
		}

		spike, err := emitters.ParseSpike(reqMap)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse spike: %v", err), http.StatusInternalServerError)
			return
		}

		emitter.EmitSpoke(spike)
	}
}
