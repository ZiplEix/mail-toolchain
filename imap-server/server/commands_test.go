package server

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseUIDFetch_ValidInputs(t *testing.T) {
	tests := []struct {
		rangeStr  string
		fieldsRaw string
		expect    *UIDFetchRequest
		name      string
	}{
		{
			"1", "(UID FLAGS)",
			&UIDFetchRequest{1, 1, []string{"UID", "FLAGS"}},
			"Simple UID",
		},
		{
			"1:*", "(UID RFC822.SIZE FLAGS)",
			&UIDFetchRequest{1, -1, []string{"UID", "RFC822.SIZE", "FLAGS"}},
			"UID range with wildcard",
		},
		{
			"2:4", "(UID BODY.PEEK[HEADER.FIELDS (From To Subject)])",
			&UIDFetchRequest{2, 4, []string{"UID", "BODY.PEEK[HEADER.FIELDS (From To Subject)]"}},
			"With HEADER.FIELDS block",
		},
		{
			"*", "(BODY[])",
			&UIDFetchRequest{-1, -1, []string{"BODY[]"}},
			"Wildcard range only",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseUIDFetch(test.rangeStr, test.fieldsRaw)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if !reflect.DeepEqual(got, test.expect) {
				t.Errorf("got %+v, want %+v", got, test.expect)
			}
		})
	}
}

func TestParseUIDFetch_InvalidInputs(t *testing.T) {
	tests := []struct {
		rangeStr  string
		fieldsRaw string
		errMsg    string
		name      string
	}{
		{"a", "(UID)", "invalid range", "Non-integer range"},
		{"1:abc", "(UID)", "invalid range end", "Non-integer end"},
		{"1:3", "(BODY.PEEK[HEADER.FIELDS (From To]", "unterminated HEADER.FIELDS", "Unclosed header fields"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := parseUIDFetch(test.rangeStr, test.fieldsRaw)
			if err == nil {
				t.Errorf("expected error but got nil")
				return
			}
			if !strings.Contains(err.Error(), test.errMsg) {
				t.Errorf("expected error to contain '%s', got '%v'", test.errMsg, err)
			}
		})
	}
}

func TestExtractHeaders(t *testing.T) {
	raw := "Subject: Hello\r\nFrom: test@example.com\r\nTo: you@example.com\r\n\r\nBody here"
	wanted := []string{"From", "To"}
	expected := "From: test@example.com\r\nTo: you@example.com\r\n\r\n"
	result := extractHeaders(raw, wanted)
	if result != expected {
		t.Errorf("extractHeaders failed, got: %q, want: %q", result, expected)
	}
}

func TestParseUIDFetch_ComplexHeaderFields(t *testing.T) {
	rangeStr := "1:3"
	fieldsRaw := "(UID RFC822.SIZE FLAGS BODY.PEEK[HEADER.FIELDS (From To Cc Bcc Subject Date Message-ID Priority X-Priority References Newsgroups In-Reply-To Content-Type Reply-To)])"

	expected := &UIDFetchRequest{
		RangeStart: 1,
		RangeEnd:   3,
		Fields: []string{
			"UID",
			"RFC822.SIZE",
			"FLAGS",
			"BODY.PEEK[HEADER.FIELDS (From To Cc Bcc Subject Date Message-ID Priority X-Priority References Newsgroups In-Reply-To Content-Type Reply-To)]",
		},
	}

	got, err := parseUIDFetch(rangeStr, fieldsRaw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %+v, want %+v", got, expected)
	}
}
