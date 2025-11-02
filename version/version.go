package version

import (
	"encoding/json"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

var (
	Version = "dev"
	Commit  = "none"
	BuiltAt = "unknown"
)

type Info struct {
	Version   string    `json:"version"`
	Commit    string    `json:"commit"`
	BuiltAt   string    `json:"builtAt"`
	GoVersion string    `json:"goVersion"`
	Compiler  string    `json:"compiler"`
	Platform  string    `json:"platform"`
	Time      time.Time `json:"-"`
}

func Get() Info {
	t := time.Time{}
	if tt, err := time.Parse(time.RFC3339, BuiltAt); err == nil {
		t = tt
	}
	return Info{
		Version:   Version,
		Commit:    Commit,
		BuiltAt:   BuiltAt,
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  runtime.GOOS + "/" + runtime.GOARCH,
		Time:      t,
	}
}

func (i Info) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Version: %s\nCommit: %s\nBuiltAt: %s\n", i.Version, i.Commit, i.BuiltAt)
	fmt.Fprintf(&b, "Go: %s (%s) %s\n", i.GoVersion, i.Compiler, i.Platform)
	return b.String()
}

func (i Info) JSON() string {
	b, _ := json.MarshalIndent(i, "", "  ")
	return string(b)
}

func IsDev() bool {
	return strings.EqualFold(Version, "dev") || Version == ""
}

type RuntimeBuildInfo struct {
	MainPath    string            `json:"mainPath"`
	MainVersion string            `json:"mainVersion"`
	Settings    map[string]string `json:"settings"`
}

func ReadRuntimeBuildInfo() (*RuntimeBuildInfo, bool) {
	bi, ok := debug.ReadBuildInfo()
	if !ok || bi == nil {
		return nil, false
	}
	out := &RuntimeBuildInfo{
		MainPath:    bi.Main.Path,
		MainVersion: bi.Main.Version,
		Settings:    map[string]string{},
	}
	for _, s := range bi.Settings {
		out.Settings[s.Key] = s.Value
	}
	return out, true
}
