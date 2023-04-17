package testfuncs

import (
	"fmt"
	"os"
	"path/filepath"
)

func PullPackagedChart(repoName string, repoURL string, artifact string, version string, targetPath string) error {
	localPath := filepath.Join(targetPath, fmt.Sprintf("%s-%s.tgz", artifact, version))
	if _, err := os.Stat(localPath); err == nil {
		fmt.Printf("helm pull is skipped, the file '%s' already exists", localPath)
		return nil
	}
	scriptContent := "" +
		fmt.Sprintf("helm repo add %s %s", repoName, repoURL) + "\n" +
		fmt.Sprintf("helm pull %s/%s --version %s", repoName, artifact, version)

	path, err := CreateTemporalFile(scriptContent)
	if err != nil {
		return err
	}
	defer os.Remove(path)
	_, err = ExecBashScriptFromSpecificPath(targetPath, path)

	return err
}
