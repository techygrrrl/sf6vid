package file_utils

import "os/exec"

func OpenFile(filepath string) error {
	_, err := exec.Command("PowerShell", "-Command", "Start-Process", filepath).Output()
	if err != nil {
		panic(err)
	}
	return err
}
