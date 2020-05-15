package bui

import "runtime"

func FindDLL() string {
	if runtime.GOARCH == "amd64" {
		return "ui_x64.dll"
	}
	return "ui.dll"
}
