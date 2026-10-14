package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestValidateReleaseDecision(t *testing.T) {
	cases := []struct {
		name string
		in   Release
		want bool
	}{
		{"complete", Release{"may", "Acme", "login"}, true},
		{"missing template", Release{"may", "Acme", ""}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if (ValidateRelease(tc.in) == nil) != tc.want {
				t.Fatalf("valid=%v, want %v", ValidateRelease(tc.in) == nil, tc.want)
			}
		})
	}
}

func TestPublishReleasePayloads(t *testing.T) {
	want := []map[string]string{{"name": "Acme"}, {"name": "sms-cutover-may-template", "body": "login"}}
	var got []map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("decode request: %v", err)
		}
		got = append(got, payload)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()
	c := &InfraiClient{base: server.URL, key: "test", http: server.Client()}
	if err := PublishRelease(c, Release{"may", "Acme", "login"}); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("payloads = %v, want %v", got, want)
	}
}
