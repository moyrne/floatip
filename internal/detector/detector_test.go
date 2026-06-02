package detector

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDetectPublicIP_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("1.2.3.4\n"))
	}))
	defer ts.Close()

	ip, err := DetectPublicIP(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ip != "1.2.3.4" {
		t.Fatalf("expected 1.2.3.4, got %s", ip)
	}
}

func TestDetectPublicIP_Non200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := DetectPublicIP(ts.URL)
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestDetectPublicIP_EmptyBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("  \n"))
	}))
	defer ts.Close()

	_, err := DetectPublicIP(ts.URL)
	if err == nil {
		t.Fatal("expected error for empty body")
	}
}

func TestDetectPublicIP_Unreachable(t *testing.T) {
	_, err := DetectPublicIP("http://192.0.2.1:9")
	if err == nil {
		t.Fatal("expected error for unreachable service")
	}
}

func TestDetectPublicIP_TrimsWhitespace(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("  203.0.113.5  \n"))
	}))
	defer ts.Close()

	ip, err := DetectPublicIP(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ip != "203.0.113.5" {
		t.Fatalf("expected 203.0.113.5, got %s", ip)
	}
}

func TestDetectPublicIP_EmptyURL(t *testing.T) {
	_, err := DetectPublicIP("")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}
