package main

import (
	"os/exec"
	"runtime"
	"strings"
)

// ThemeStatus thể hiện chế độ giao diện
type ThemeStatus string

const (
	Dark    ThemeStatus = "Dark"
	Light   ThemeStatus = "Light"
	Unknown ThemeStatus = "Unknown"
)

// GetSystemTheme trả về Theme hiện tại của hệ điều hành
func GetSystemTheme() ThemeStatus {
	switch runtime.GOOS {
	case "windows":
		return getWindowsTheme()
	case "darwin":
		return getDarwinTheme()
	case "linux":
		return getLinuxTheme()
	default:
		return Unknown
	}
}

// macOS: Gọi lệnh defaults read
func getDarwinTheme() ThemeStatus {
	cmd := exec.Command("defaults", "read", "-g", "AppleInterfaceStyle")
	out, err := cmd.Output()
	if err != nil {
		// Nếu key AppleInterfaceStyle không tồn tại, macOS đang ở Light Mode
		return Light
	}

	if strings.TrimSpace(string(out)) == "Dark" {
		return Dark
	}
	return Light
}

// Linux: Đọc qua gsettings (GNOME/Freedesktop)
func getLinuxTheme() ThemeStatus {
	// Kiểm tra chuẩn color-scheme mới (FreeDesktop / GNOME 42+)
	cmd := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "color-scheme")
	out, err := cmd.Output()
	if err == nil {
		val := strings.ToLower(string(out))
		if strings.Contains(val, "dark") {
			return Dark
		}
		if strings.Contains(val, "light") || strings.Contains(val, "default") {
			return Light
		}
	}

	// Fallback cho các bản GNOME cũ hơn (kiểm tra gtk-theme)
	cmdOld := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "gtk-theme")
	outOld, err := cmdOld.Output()
	if err == nil {
		if strings.Contains(strings.ToLower(string(outOld)), "dark") {
			return Dark
		}
	}

	return Unknown
}

// func main() {
// 	theme := GetSystemTheme()
// 	fmt.Printf("Hệ điều hành: %s\n", runtime.GOOS)
// 	fmt.Printf("Theme hệ thống hiện tại: %s\n", theme)
// }
