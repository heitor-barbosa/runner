package runner

import "testing"

func TestParseResponseSuccess(t *testing.T) {
	response, err := parseResponse([]byte(`{"success":true,"data":"SIGNATURE"}`))
	if err != nil {
		t.Fatalf("parseResponse returned error: %v", err)
	}
	if !response.Success || response.Data != "SIGNATURE" {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestParseResponseStructuredError(t *testing.T) {
	response, err := parseResponse([]byte(`{"success":false,"errorCode":"INVALID","errorMessage":"bad input"}`))
	if err != nil {
		t.Fatalf("parseResponse returned error: %v", err)
	}
	if response.Success || response.ErrorCode != "INVALID" || response.ErrorMessage != "bad input" {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestParseResponseRejectsInvalidJSON(t *testing.T) {
	if _, err := parseResponse([]byte(`not-json`)); err == nil {
		t.Fatal("expected invalid JSON to be rejected")
	}
}
