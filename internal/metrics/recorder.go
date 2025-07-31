package metrics

import "time"

type Recorder interface {
	RecordHTTPRequest(method, endpoint, status string)
	RecordHTTPDuration(method, endpoint string, duration time.Duration)
	RecordRepositoryAnalysis(status string, duration time.Duration)
	RecordCacheHit()
	RecordCacheMiss()
}

type NoOpRecorder struct{}

func (n NoOpRecorder) RecordHTTPRequest(method, endpoint, status string)                  {}
func (n NoOpRecorder) RecordHTTPDuration(method, endpoint string, duration time.Duration) {}
func (n NoOpRecorder) RecordRepositoryAnalysis(status string, duration time.Duration)     {}
func (n NoOpRecorder) RecordCacheHit()                                                    {}
func (n NoOpRecorder) RecordCacheMiss()                                                   {}
