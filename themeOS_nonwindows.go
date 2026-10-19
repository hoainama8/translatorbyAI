//go:build !windows

package main

func getWindowsTheme() ThemeStatus {
	return Unknown
}
