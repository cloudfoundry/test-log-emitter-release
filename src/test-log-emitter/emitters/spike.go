package emitters

import (
	"fmt"
	"time"

	"code.cloudfoundry.org/go-loggregator"
)

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
