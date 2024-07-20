package file_utils

import (
	"os/exec"
)

func OpenFile(filepath string) error {
	_, err := exec.Command("xdg-open", filepath).Output()
	return err
}
