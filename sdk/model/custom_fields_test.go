package model

import (
	"encoding/json"
	"testing"
)

func TestIsCustomField(t *testing.T) {
	tests := []struct {
		key  string
		want bool
	}{
		{"custom_field_one", true},
		{"custom_field_9", true},
		{"custom_field_200", true},
		{"custom_plan_field_1", true},
		{"custom_plan_field_10", true},
		{"id", false},
		{"name", false},
		{"status", false},
		{"custom_field", false},
		{"custom_plan_field", false},
		{"custom_field_", true},
	}
	for _, tt := range tests {
		if got := IsCustomField(tt.key); got != tt.want {
			t.Errorf("IsCustomField(%q) = %v, want %v", tt.key, got, tt.want)
		}
	}
}

func TestExtractCustomFields(t *testing.T) {
	raw := map[string]json.RawMessage{
		"id":                  json.RawMessage(`"100"`),
		"name":                json.RawMessage(`"test"`),
		"custom_field_one":    json.RawMessage(`"value1"`),
		"custom_field_9":      json.RawMessage(`"value9"`),
		"custom_plan_field_1": json.RawMessage(`"plan1"`),
		"custom_field_empty":  json.RawMessage(`""`),
	}

	result := ExtractCustomFields(raw)

	if len(result) != 3 {
		t.Fatalf("expected 3 custom fields, got %d", len(result))
	}
	if result["custom_field_one"] != "value1" {
		t.Errorf("custom_field_one = %q, want %q", result["custom_field_one"], "value1")
	}
	if result["custom_field_9"] != "value9" {
		t.Errorf("custom_field_9 = %q, want %q", result["custom_field_9"], "value9")
	}
	if result["custom_plan_field_1"] != "plan1" {
		t.Errorf("custom_plan_field_1 = %q, want %q", result["custom_plan_field_1"], "plan1")
	}
}

func TestExtractCustomFields_Empty(t *testing.T) {
	raw := map[string]json.RawMessage{
		"id":   json.RawMessage(`"100"`),
		"name": json.RawMessage(`"test"`),
	}

	result := ExtractCustomFields(raw)
	if result != nil {
		t.Errorf("expected nil for no custom fields, got %v", result)
	}
}

func TestMergeCustomFields(t *testing.T) {
	params := map[string]string{
		"workspace_id": "1",
		"name":         "test",
	}
	customFields := map[string]string{
		"custom_field_one": "val1",
		"custom_field_9":   "val9",
		"":                 "skip",
	}
	MergeCustomFields(params, customFields)

	if params["custom_field_one"] != "val1" {
		t.Errorf("custom_field_one = %q, want %q", params["custom_field_one"], "val1")
	}
	if params["custom_field_9"] != "val9" {
		t.Errorf("custom_field_9 = %q, want %q", params["custom_field_9"], "val9")
	}
	if params["workspace_id"] != "1" {
		t.Errorf("workspace_id should not be overwritten")
	}
}

func TestMergeCustomFields_Nil(t *testing.T) {
	params := map[string]string{"workspace_id": "1"}
	MergeCustomFields(params, nil)
	if len(params) != 1 {
		t.Errorf("expected 1 param, got %d", len(params))
	}
}

func TestStory_UnmarshalJSON_CustomFields(t *testing.T) {
	data := []byte(`{"id":"100","name":"Test","custom_field_one":"cf1","custom_field_9":"cf9","custom_plan_field_1":"pf1"}`)
	var s Story
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if s.ID != "100" {
		t.Errorf("ID = %q, want %q", s.ID, "100")
	}
	if s.Name != "Test" {
		t.Errorf("Name = %q, want %q", s.Name, "Test")
	}
	if len(s.CustomFields) != 3 {
		t.Fatalf("expected 3 custom fields, got %d", len(s.CustomFields))
	}
	if s.CustomFields["custom_field_one"] != "cf1" {
		t.Errorf("custom_field_one = %q, want %q", s.CustomFields["custom_field_one"], "cf1")
	}
	if s.CustomFields["custom_field_9"] != "cf9" {
		t.Errorf("custom_field_9 = %q, want %q", s.CustomFields["custom_field_9"], "cf9")
	}
	if s.CustomFields["custom_plan_field_1"] != "pf1" {
		t.Errorf("custom_plan_field_1 = %q, want %q", s.CustomFields["custom_plan_field_1"], "pf1")
	}
}

func TestStory_MarshalJSON_CustomFields(t *testing.T) {
	s := Story{
		ID:   "100",
		Name: "Test",
		CustomFields: map[string]string{
			"custom_field_one": "cf1",
		},
	}
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal output error: %v", err)
	}
	if raw["id"] != "100" {
		t.Errorf("id = %q, want %q", raw["id"], "100")
	}
	if raw["custom_field_one"] != "cf1" {
		t.Errorf("custom_field_one = %q, want %q", raw["custom_field_one"], "cf1")
	}
}

func TestStory_MarshalJSON_NoCustomFields(t *testing.T) {
	s := Story{ID: "100", Name: "Test"}
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	for k := range raw {
		if IsCustomField(k) {
			t.Errorf("unexpected custom field %q in output", k)
		}
	}
}

func TestBug_UnmarshalJSON_CustomFields(t *testing.T) {
	data := []byte(`{"id":"500","title":"Bug","custom_field_1":"v1","custom_field_50":"v50"}`)
	var b Bug
	if err := json.Unmarshal(data, &b); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if b.ID != "500" {
		t.Errorf("ID = %q, want %q", b.ID, "500")
	}
	if len(b.CustomFields) != 2 {
		t.Fatalf("expected 2 custom fields, got %d", len(b.CustomFields))
	}
	if b.CustomFields["custom_field_1"] != "v1" {
		t.Errorf("custom_field_1 = %q, want %q", b.CustomFields["custom_field_1"], "v1")
	}
}

func TestBug_MarshalJSON_CustomFields(t *testing.T) {
	b := Bug{
		ID:    "500",
		Title: "Bug",
		CustomFields: map[string]string{
			"custom_field_1": "v1",
		},
	}
	data, err := json.Marshal(b)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal output error: %v", err)
	}
	if raw["custom_field_1"] != "v1" {
		t.Errorf("custom_field_1 = %q, want %q", raw["custom_field_1"], "v1")
	}
}

func TestTask_UnmarshalJSON_CustomFields(t *testing.T) {
	data := []byte(`{"id":"300","name":"Task","custom_field_one":"t1","custom_plan_field_5":"p5"}`)
	var task Task
	if err := json.Unmarshal(data, &task); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if task.ID != "300" {
		t.Errorf("ID = %q, want %q", task.ID, "300")
	}
	if len(task.CustomFields) != 2 {
		t.Fatalf("expected 2 custom fields, got %d", len(task.CustomFields))
	}
	if task.CustomFields["custom_field_one"] != "t1" {
		t.Errorf("custom_field_one = %q, want %q", task.CustomFields["custom_field_one"], "t1")
	}
	if task.CustomFields["custom_plan_field_5"] != "p5" {
		t.Errorf("custom_plan_field_5 = %q, want %q", task.CustomFields["custom_plan_field_5"], "p5")
	}
}

func TestTask_MarshalJSON_CustomFields(t *testing.T) {
	task := Task{
		ID:   "300",
		Name: "Task",
		CustomFields: map[string]string{
			"custom_field_one": "t1",
		},
	}
	data, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal output error: %v", err)
	}
	if raw["custom_field_one"] != "t1" {
		t.Errorf("custom_field_one = %q, want %q", raw["custom_field_one"], "t1")
	}
}

func TestStory_RoundTrip(t *testing.T) {
	original := `{"id":"100","name":"Test","status":"open","custom_field_one":"cf1","custom_field_9":"cf9"}`
	var s Story
	if err := json.Unmarshal([]byte(original), &s); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var s2 Story
	if err := json.Unmarshal(data, &s2); err != nil {
		t.Fatalf("second Unmarshal error: %v", err)
	}

	if s2.ID != s.ID || s2.Name != s.Name || s2.Status != s.Status {
		t.Errorf("standard fields mismatch after round-trip")
	}
	if len(s2.CustomFields) != len(s.CustomFields) {
		t.Fatalf("custom fields count mismatch: got %d, want %d", len(s2.CustomFields), len(s.CustomFields))
	}
	for k, v := range s.CustomFields {
		if s2.CustomFields[k] != v {
			t.Errorf("custom field %q = %q, want %q", k, s2.CustomFields[k], v)
		}
	}
}
