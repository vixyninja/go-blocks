package version_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/vixyninja/go-blocks/version"
)

func withGlobals(t *testing.T, ver, commit, built string, fn func()) {
	t.Helper()
	oldV, oldC, oldB := version.Version, version.Commit, version.BuiltAt
	version.Version, version.Commit, version.BuiltAt = ver, commit, built
	t.Cleanup(func() { version.Version, version.Commit, version.BuiltAt = oldV, oldC, oldB })
	fn()
}

func TestGet_PopulatesInfo(t *testing.T) {
	built := time.Now().UTC().Format(time.RFC3339)
	withGlobals(t, "v1.2.3", "abcdef0", built, func() {
		info := version.Get()
		if info.Version != "v1.2.3" || info.Commit != "abcdef0" || info.BuiltAt != built {
			t.Fatalf("unexpected Info: %+v", info)
		}
		if info.GoVersion == "" || info.Compiler == "" || info.Platform == "" {
			t.Fatalf("runtime fields should not be empty: %+v", info)
		}
		// Time should parse successfully when BuiltAt is RFC3339
		if info.Time.IsZero() {
			t.Fatalf("expected parsed time, got zero time")
		}
	})
}

func TestGet_TimeZeroOnInvalidBuiltAt(t *testing.T) {
	withGlobals(t, "v", "c", "not-time", func() {
		info := version.Get()
		if !info.Time.IsZero() {
			t.Fatalf("expected zero time when BuiltAt invalid, got %v", info.Time)
		}
	})
}

func TestInfo_String_IncludesCoreFields(t *testing.T) {
	withGlobals(t, "1.0.0", "deadbeef", "2024-01-01T00:00:00Z", func() {
		s := version.Get().String()
		for _, want := range []string{"Version: 1.0.0", "Commit: deadbeef", "BuiltAt: 2024-01-01T00:00:00Z", "Go:", ") ", "/"} {
			if !strings.Contains(s, want) {
				t.Fatalf("String() missing %q in %q", want, s)
			}
		}
	})
}

func TestInfo_JSON_IsValidAndContainsFields(t *testing.T) {
	withGlobals(t, "1.0.1", "cafebabe", "2024-02-02T00:00:00Z", func() {
		js := version.Get().JSON()
		var got map[string]any
		if err := json.Unmarshal([]byte(js), &got); err != nil {
			t.Fatalf("invalid JSON: %v\n%s", err, js)
		}
		for _, k := range []string{"version", "commit", "builtAt", "goVersion", "compiler", "platform"} {
			if _, ok := got[k]; !ok {
				t.Fatalf("JSON missing key %q: %s", k, js)
			}
		}
	})
}

func TestIsDev(t *testing.T) {
	cases := []struct {
		name string
		ver  string
		want bool
	}{
		{"empty", "", true},
		{"dev_lower", "dev", true},
		{"dev_upper", "DEV", true},
		{"semver", "v1.0.0", false},
	}
	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			withGlobals(t, cs.ver, "c", "b", func() {
				if got := version.IsDev(); got != cs.want {
					t.Fatalf("IsDev()=%v want %v (ver=%q)", got, cs.want, cs.ver)
				}
			})
		})
	}
}

func TestReadRuntimeBuildInfo(t *testing.T) {
	bi, ok := version.ReadRuntimeBuildInfo()
	if !ok || bi == nil {
		t.Fatalf("expected runtime build info available")
	}
	if bi.MainPath == "" {
		t.Fatalf("Main.Path should not be empty")
	}
	if bi.Settings == nil {
		t.Fatalf("Settings should not be nil")
	}
}
