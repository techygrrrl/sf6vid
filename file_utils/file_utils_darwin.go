package file_utils

import "os/exec"

func OpenFile(filepath string) error {
	_, err := exec.Command("open", filepath).Output()
	return err
}
