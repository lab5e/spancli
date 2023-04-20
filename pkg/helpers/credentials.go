package helpers

import (
	"os"
	"path"
)

func credFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(home, ".spanrc"), nil

}

// RemoveCredentials removes the local credentials file
func RemoveCredentials() {
	fileName, err := credFile()
	if err != nil {
		return
	}
	os.Remove(fileName)
}

// ReadCredentials reads credentials from a local file. If no file is found the returned string is empty
func ReadCredentials() string {
	fileName, err := credFile()
	if err != nil {
		return ""
	}
	buf, err := os.ReadFile(fileName)
	if err != nil {
		return ""
	}
	return string(buf)
}

// WriteCredentials writes (JWT) credentials to a local file
func WriteCredentials(creds string) error {
	fileName, err := credFile()
	if err != nil {
		return err
	}
	buf := []byte(creds)
	if err := os.WriteFile(fileName, buf, 0700); err != nil {
		return err
	}
	return nil
}
