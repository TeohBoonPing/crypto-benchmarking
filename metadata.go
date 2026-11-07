package main

import "runtime"

type RuntimeMetadata struct {
	GoVersion    string
	OS           string
	Architecture string
	CPUCount     int
	GOMAXPROCS   int
}

func GetRuntimeMetadata() RuntimeMetadata {
	return RuntimeMetadata{
		GoVersion:    runtime.Version(),
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		CPUCount:     runtime.NumCPU(),
		GOMAXPROCS:   runtime.GOMAXPROCS(0),
	}
}
