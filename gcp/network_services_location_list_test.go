package gcp

import "testing"

func TestIsPolicyLocationID(t *testing.T) {
	cases := map[string]bool{
		"global":                  true,
		"us-central1":             true,
		"asia-east1":              true,
		"northamerica-northeast1": true,
		"us-central1-a":           false,
		"asia-east1-c":            false,
		"europe-west4-b":          false,
		"us":                      false,
		"":                        false,
	}
	for id, want := range cases {
		if got := isPolicyLocationID(id); got != want {
			t.Errorf("isPolicyLocationID(%q) = %v, want %v", id, got, want)
		}
	}
}
