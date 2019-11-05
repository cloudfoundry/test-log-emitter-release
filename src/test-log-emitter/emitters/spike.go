package emitters

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"code.cloudfoundry.org/go-loggregator"
)

// type GaugeValue struct {
// 	Name  string  `json:"name"`
// 	Value float64 `json:"value"`
// 	Unit  string  `json:"unit"`
// }

// type GaugeMetric struct {
// 	SourceId   string
// 	InstanceId string
// 	Tags       map[string]string
// 	Values     []GaugeValue
// }

type Spike struct {
	SourceId          string
	InstanceId        string
	ProcessInstanceId string
	Start             time.Time
	End               time.Time
}

func ParseSpike(m map[string]string) (*Spike, error) {
	spike := Spike{
		SourceId:          m["source_id"],
		InstanceId:        m["instance_id"],
		ProcessInstanceId: m["process_instance_id"],
	}
	var err error

	if spike.Start, err = time.Parse(time.RFC3339, m["spike_start"]); err != nil {
		return nil, fmt.Errorf("Failed to parse spike_start: %v", err)
	}

	if spike.End, err = time.Parse(time.RFC3339, m["spike_end"]); err != nil {
		return nil, fmt.Errorf("Failed to parse spike_end: %v", err)
	}

	return &spike, nil
}

type SpikeEmitter struct {
	client *loggregator.IngressClient
}

func NewSpikeEmitter(client *loggregator.IngressClient) *SpikeEmitter {
	return &SpikeEmitter{client: client}
}

func (e SpikeEmitter) Emit(spike *Spike) {
	tags := map[string]string{
		"process_instance_id": spike.ProcessInstanceId,
	}

	e.client.EmitGauge(
		loggregator.WithGaugeSourceInfo(spike.SourceId, spike.InstanceId),
		loggregator.WithGaugeValue("spike_start", float64(spike.Start.Unix()), "seconds"),
		loggregator.WithGaugeValue("spike_end", float64(spike.End.Unix()), "seconds"),
		loggregator.WithEnvelopeTags(tags),
	)
}

func (e SpikeEmitter) EmitSpoke(spike *Spike) {
	tags := map[string]string{
		"process_instance_id": spike.ProcessInstanceId,
	}

	e.client.EmitGauge(
		loggregator.WithGaugeSourceInfo(spike.SourceId, spike.InstanceId),
		loggregator.WithGaugeValue("spoke_start", float64(spike.Start.Unix()), "seconds"),
		loggregator.WithGaugeValue("spoke_end", float64(spike.End.Unix()), "seconds"),
		loggregator.WithEnvelopeTags(tags),
	)
}

func (e SpikeEmitter) EmitSpike() http.HandlerFunc {
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

		spike, err := ParseSpike(reqMap)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse spike: %v", err), http.StatusInternalServerError)
			return
		}

		e.Emit(spike)
	}
}

func (e SpikeEmitter) EmitBadSpike() http.HandlerFunc {
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

		spike, err := ParseSpike(reqMap)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse spike: %v", err), http.StatusInternalServerError)
			return
		}

		e.EmitSpoke(spike)
	}
}
