package bui

import "C"
import (
	"path/filepath"
	"runtime"
)

func FindNodeDLL() string {
	if runtime.GOARCH == "amd64" {
		p, _ := filepath.Abs("miniblink_x64.dll")
		return p
	} else {
		p, _ := filepath.Abs("node.dll")
		return p
	}
}

func FindMbDLL() string {
	if runtime.GOARCH == "amd64" {
		p, _ := filepath.Abs("mb_x64.dll")
		return p
	} else {
		p, _ := filepath.Abs("mb.dll")
		return p
	}
}
