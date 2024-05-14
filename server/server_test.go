package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestServerIntegration(t *testing.T) {
	Convey("The server can recive ONE POST with payload", t, func() {
		sampleLogLine := &LogLine{
			Level:      "error",
			Message:    "Failed to connect to DB",
			ResourceID: "server-1234",
			Timestamp:  time.Now(),
			TraceID:    "abc-xyz-123",
			SpanID:     "span-456",
			Commit:     "5e5342f",
			Metadata:   Metadata{ParentResourceID: "server-0987"},
		}

		wr := bytes.Buffer{}

		Then(assert.NoError(t, json.NewEncoder(&wr).Encode(sampleLogLine)))
		req := httptest.NewRequest(http.MethodPost, "/", &wr)
		req.Header.Set("content-type", "application/json")
		w := httptest.NewRecorder()

		LogHandler(w, req)
		Convey("Then server respond with StatusAccepted", func() {
			So(w.Result().StatusCode, ShouldEqual, http.StatusAccepted)
		})

		Convey("The sent fields match the recived fields", func() {
			So(w.Body.String(), ShouldContainSubstring, sampleLogLine.ResourceID)
			So(w.Body.String(), ShouldContainSubstring, sampleLogLine.Message)
		})

		Reset(func() { w.Result().Body.Close() })
	})
}

func LogHandler(w *httptest.ResponseRecorder, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)

	var l LogLine
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.NewEncoder(w).Encode(l); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type (
	LogLine struct {
		Level      string    `json:"level"`
		Message    string    `json:"message"`
		ResourceID string    `json:"resource_id"`
		Timestamp  time.Time `json:"timestamp"`
		TraceID    string    `json:"trace_id"`
		SpanID     string    `json:"span_id"`
		Commit     string    `json:"commit"`
		Metadata   `json:"metadata"`
	}

	Metadata struct {
		ParentResourceID string `json:"parent_resource_id"`
	}
)

// ShouldPass a way to integrate `testify.assertion` with goconvey
func ShouldPass(actual any, expected ...any) string {
	if actual == true {
		return ""
	}
	return "suite test failed"
}

// Then rapper around So() for readability
func Then(assertion any) {
	So(assertion, ShouldPass)
}
