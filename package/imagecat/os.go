package imagecat

import (
	"fmt"
	"runtime"
)

type SupportedOs int

const (
	Linux SupportedOs = iota
	Windows
)

var osMap = map[string]SupportedOs{
	"linux":   Linux,
	"windows": Windows,
}

func DetectOs() (SupportedOs, error) {
	runtimeOs := runtime.GOOS
	if os, ok := osMap[runtimeOs]; ok {
		return os, nil
	}
	return 0, fmt.Errorf("unsupported OS: %s", runtimeOs)
}
