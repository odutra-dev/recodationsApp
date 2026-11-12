package presenter

import "time"

// ResponseDoc provides a concrete non-generic struct for OpenAPI/Swagger 2.0 generation.
type ResponseDoc struct {
	Data         any                 `json:"data"`
	Message      string              `json:"message"`
	TraceID      string              `json:"traceId"`
	Success      bool                `json:"success"`
	Timestamp    time.Time           `json:"timestamp"`
	DetailsError *string             `json:"detailsError,omitempty"`
	Path         *string             `json:"path,omitempty"`
	APIVersion   uint16              `json:"apiVersion"`
	Errors       []string            `json:"errors,omitempty"`
	ErrorValids  map[string][]string `json:"errorValids,omitempty"`
}
