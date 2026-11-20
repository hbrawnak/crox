package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hbrawnak/crox/internal/adapter/httpapi"
)

func TestHealthz(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantAllow  string // checked only when set
	}{
		{"get health", http.MethodGet, "/healthz", http.StatusOK, ""},
		{"head health", http.MethodHead, "/healthz", http.StatusOK, ""},
		{"post health", http.MethodPost, "/healthz", http.StatusMethodNotAllowed, "GET, HEAD"},
		{"unknown path", http.MethodGet, "/nope", http.StatusNotFound, ""},
	}

	h := httpapi.NewHandler()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			resp := rec.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tc.wantStatus {
				t.Fatalf("%s %s: status = %d, want %d", tc.method, tc.path, resp.StatusCode, tc.wantStatus)
			}

			if tc.wantAllow != "" {
				if got := resp.Header.Get("Allow"); got != tc.wantAllow {
					t.Errorf("%s %s: Allow = %q, want %q", tc.method, tc.path, got, tc.wantAllow)
				}
			}

			if tc.wantStatus != http.StatusOK {
				return
			}

			if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
				t.Errorf("%s %s: Content-Type = %q, want application/json", tc.method, tc.path, ct)
			}

			if tc.method == http.MethodHead {
				return
			}

			var body struct {
				Status string `json:"status"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("%s %s: decode body: %v", tc.method, tc.path, err)
			}
			if body.Status != "ok" {
				t.Errorf("%s %s: status field = %q, want %q", tc.method, tc.path, body.Status, "ok")
			}
		})
	}
}
