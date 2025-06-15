package config

import (
	"os"
	"path/filepath"
)

const (
	LoginURL    = "https://kitoakari-fc.com/slogin.php"
	CookieFile  = ".fukumimi_cookies"
)

func GetCookieFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return CookieFile
	}
	return filepath.Join(homeDir, CookieFile)
}