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

func TestAddPaginationParams_BothSet(t *testing.T) {
	params := map[string]string{}
	addPaginationParams(params, 10, 2)
	if params["limit"] != "10" {
		t.Errorf("params[limit] = %q, want %q", params["limit"], "10")
	}
	if params["page"] != "2" {
		t.Errorf("params[page] = %q, want %q", params["page"], "2")
	}
}

func TestAddPaginationParams_LimitZero(t *testing.T) {
	params := map[string]string{}
	addPaginationParams(params, 0, 2)
	if _, ok := params["limit"]; ok {
		t.Error("expected limit=0 to not be added to params")
	}
	if params["page"] != "2" {
		t.Errorf("params[page] = %q, want %q", params["page"], "2")
	}
}

func TestAddPaginationParams_PageZero(t *testing.T) {
	params := map[string]string{}
	addPaginationParams(params, 10, 0)
	if params["limit"] != "10" {
		t.Errorf("params[limit] = %q, want %q", params["limit"], "10")
	}
	if _, ok := params["page"]; ok {
		t.Error("expected page=0 to not be added to params")
	}
}
