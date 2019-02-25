package metrics

import (
	"net/http"
	"time"
)

// var Metrics = make(chan time.Duration)
// var Avg = make(chan time.Duration)
// var GetAvg = make(chan bool)

func NewMetric(window int) *Metric {
	return &Metric{
		Measurement: make(chan time.Duration),
		Avg:         make(chan time.Duration),
		GetAvg:      make(chan bool),
		Window:      window,
	}
}

type Metric struct {
	Measurement chan time.Duration
	Avg         chan time.Duration
	GetAvg      chan bool
	Window      int
}

func (m *Metric) TimeIt(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		f.ServeHTTP(w, r)
		toSend := time.Since(start)
		m.Measurement <- toSend
	})
}

func (m *Metric) Listen() {
	var timings []time.Duration
	for {
		select {
		case timing := <-m.Measurement:
			timings = append(timings, timing)
			if len(timings) > m.Window {
				timings = timings[1:(m.Window + 1)]
			}
		case <-m.GetAvg:
			if len(timings) < 1 {
				m.Avg <- 0
				continue
			}
			var sum time.Duration
			for _, t := range timings {
				sum += t
			}
			mean := int64(sum) / int64(len(timings))
			m.Avg <- time.Duration(mean)
		}
	}
}
