package testfuncs

import "os/exec"

func ExecBashScript(scriptPath string) (string, error) {
	out, err := exec.Command("sh", "-c", "chmod +x "+scriptPath).Output()
	if err != nil {
		return string(out), err
	}
	out, err = exec.Command("sh", "-c", scriptPath).Output()
	if err != nil {
		return string(out), err
	}

	return string(out), nil
}

func ExecBashScriptFromSpecificPath(dirPath string, scriptPath string) (string, error) {
	cmd := exec.Command("sh", "-c", "chmod +x "+scriptPath)
	cmd.Dir = dirPath
	out, err := cmd.Output()
	if err != nil {
		return string(out), err
	}
	cmd = exec.Command("sh", "-c", scriptPath)
	cmd.Dir = dirPath
	out, err = cmd.Output()
	if err != nil {
		return string(out), err
	}

	return string(out), nil
}
