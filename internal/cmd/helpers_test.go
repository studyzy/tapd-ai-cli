package cmd

import (
	"testing"
)

func TestAddOptionalParam_Empty(t *testing.T) {
	params := map[string]string{}
	addOptionalParam(params, "status", "")
	if _, ok := params["status"]; ok {
		t.Error("expected empty value to not be added to params")
	}
}

func TestAddOptionalParam_NonEmpty(t *testing.T) {
	params := map[string]string{}
	addOptionalParam(params, "status", "open")
	if params["status"] != "open" {
		t.Errorf("params[status] = %q, want %q", params["status"], "open")
	}
}
