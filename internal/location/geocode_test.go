package location

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDetectLocation_Success(t *testing.T) {
	// Mock IP geolocation service.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"success","lat":51.5000,"lon":-0.1000}`))
	}))
	defer srv.Close()

	oldBase := geoBaseURL
	geoBaseURL = srv.URL
	defer func() { geoBaseURL = oldBase }()

	loc, err := DetectLocation()
	if err != nil {
		t.Fatalf("DetectLocation() err = %v, want nil", err)
	}
	if loc.Lat != 51.5 || loc.Lon != -0.1 {
		t.Errorf("DetectLocation() = %+v, want Lat=51.5 Lon=-0.1", loc)
	}
}

func TestDetectLocation_ErrorOnFailureStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"fail","message":"test error"}`))
	}))
	defer srv.Close()

	oldBase := geoBaseURL
	geoBaseURL = srv.URL
	defer func() { geoBaseURL = oldBase }()

	_, err := DetectLocation()
	if err == nil {
		t.Fatal("DetectLocation() err = nil, want error for status=fail")
	}
}

