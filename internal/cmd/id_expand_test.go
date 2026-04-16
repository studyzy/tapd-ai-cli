package cmd

import "testing"

func TestExpandShortID_ShortNumeric(t *testing.T) {
	got := expandShortID("28841", "61692131")
	want := "1161692131000028841"
	if got != want {
		t.Errorf("expandShortID(\"28841\", \"61692131\") = %q, want %q", got, want)
	}
}

func TestExpandShortID_SingleDigit(t *testing.T) {
	got := expandShortID("1", "61692131")
	want := "1161692131000000001"
	if got != want {
		t.Errorf("expandShortID(\"1\", \"61692131\") = %q, want %q", got, want)
	}
}

func TestExpandShortID_NineDigits(t *testing.T) {
	got := expandShortID("123456789", "61692131")
	want := "1161692131123456789"
	if got != want {
		t.Errorf("expandShortID(\"123456789\", \"61692131\") = %q, want %q", got, want)
	}
}

func TestExpandShortID_AlreadyLong(t *testing.T) {
	longID := "1161692131001028841"
	got := expandShortID(longID, "61692131")
	if got != longID {
		t.Errorf("expandShortID(%q) = %q, want unchanged", longID, got)
	}
}

func TestExpandShortID_TenDigits(t *testing.T) {
	id := "1234567890"
	got := expandShortID(id, "61692131")
	if got != id {
		t.Errorf("expandShortID(%q) = %q, want unchanged", id, got)
	}
}

func TestExpandShortID_NonNumeric(t *testing.T) {
	got := expandShortID("abc123", "61692131")
	if got != "abc123" {
		t.Errorf("expandShortID(\"abc123\") = %q, want \"abc123\"", got)
	}
}

func TestExpandShortID_Empty(t *testing.T) {
	got := expandShortID("", "61692131")
	if got != "" {
		t.Errorf("expandShortID(\"\") = %q, want \"\"", got)
	}
}

func TestExpandShortID_EmptyWorkspaceID(t *testing.T) {
	got := expandShortID("28841", "")
	if got != "28841" {
		t.Errorf("expandShortID(\"28841\", \"\") = %q, want \"28841\"", got)
	}
}
