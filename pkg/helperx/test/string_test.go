package helperx_test

import (
	"reflect"
	"testing"

	helperx "github.com/vixyninja/go-blocks/pkg/helperx"
)

func TestIsEmpty_IsNotEmpty(t *testing.T) {
	cases := []struct {
		in       string
		empty    bool
		notEmpty bool
	}{
		{"", true, false},
		{" ", false, true},
		{"abc", false, true},
	}
	for _, cs := range cases {
		if got := helperx.IsEmpty(cs.in); got != cs.empty {
			t.Fatalf("IsEmpty(%q)=%v want %v", cs.in, got, cs.empty)
		}
		if got := helperx.IsNotEmpty(cs.in); got != cs.notEmpty {
			t.Fatalf("IsNotEmpty(%q)=%v want %v", cs.in, got, cs.notEmpty)
		}
	}
}

func TestConcatStrings(t *testing.T) {
	if got, want := helperx.ConcatStrings(), ""; got != want {
		t.Fatalf("ConcatStrings()=%q want %q", got, want)
	}
	if got, want := helperx.ConcatStrings("a"), "a"; got != want {
		t.Fatalf("ConcatStrings(a)=%q want %q", got, want)
	}
	if got, want := helperx.ConcatStrings("a", "b", "c"), "abc"; got != want {
		t.Fatalf("ConcatStrings(a,b,c)=%q want %q", got, want)
	}
}

func TestConcatStringsWithSeparator(t *testing.T) {
	if got, want := helperx.ConcatStringsWithSeparator(","), ""; got != want {
		t.Fatalf("ConcatStringsWithSeparator(',')=%q want %q", got, want)
	}
	if got, want := helperx.ConcatStringsWithSeparator(",", "a"), "a"; got != want {
		t.Fatalf("ConcatStringsWithSeparator(',', 'a')=%q want %q", got, want)
	}
	if got, want := helperx.ConcatStringsWithSeparator(",", "a", "b", "c"), "a,b,c"; got != want {
		t.Fatalf("ConcatStringsWithSeparator(',', a,b,c)=%q want %q", got, want)
	}
}

func TestStringToUint64(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		want    uint64
		wantErr bool
	}{
		{"empty", "", 0, true},
		{"invalid", "abc", 0, true},
		{"ok", "42", 42, false},
	}
	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			got, err := helperx.StringToUint64(cs.in)
			if (err != nil) != cs.wantErr {
				t.Fatalf("err=%v wantErr=%v", err, cs.wantErr)
			}
			if got != cs.want {
				t.Fatalf("got=%v want=%v", got, cs.want)
			}
		})
	}
}

func TestStringToInt64(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		want    int64
		wantErr bool
	}{
		{"empty", "", 0, true},
		{"invalid", "x", 0, true},
		{"ok_negative", "-7", -7, false},
		{"ok_zero", "0", 0, false},
	}
	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			got, err := helperx.StringToInt64(cs.in)
			if (err != nil) != cs.wantErr {
				t.Fatalf("err=%v wantErr=%v", err, cs.wantErr)
			}
			if got != cs.want {
				t.Fatalf("got=%v want=%v", got, cs.want)
			}
		})
	}
}

func TestStringToInt(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		want    int
		wantErr bool
	}{
		{"empty", "", 0, true},
		{"invalid", "a1", 0, true},
		{"ok_trim", "  15 ", 15, false},
		{"ok_plain", "27", 27, false},
	}
	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			got, err := helperx.StringToInt(cs.in)
			if (err != nil) != cs.wantErr {
				t.Fatalf("err=%v wantErr=%v", err, cs.wantErr)
			}
			if got != cs.want {
				t.Fatalf("got=%v want=%v", got, cs.want)
			}
		})
	}
}

func TestBytesStringConversions(t *testing.T) {
	s := "hello"
	b := helperx.StringToBytes(s)
	if !reflect.DeepEqual(b, []byte{'h', 'e', 'l', 'l', 'o'}) {
		t.Fatalf("StringToBytes mismatch: %v", b)
	}
	if got := helperx.BytesToString(b); got != s {
		t.Fatalf("BytesToString(%v)=%q want %q", b, got, s)
	}
}

func TestReverseString(t *testing.T) {
	if got, want := helperx.ReverseString("abc"), "cba"; got != want {
		t.Fatalf("ReverseString('abc')=%q want %q", got, want)
	}
	// Unicode aware
	if got, want := helperx.ReverseString("áb😊"), "😊bá"; got != want {
		t.Fatalf("ReverseString('áb😊')=%q want %q", got, want)
	}
}
