package settings

import (
	"path/filepath"
	"time"

	"github.com/astronetes/sdk-go/internal/fsys"
	v1 "github.com/astronetes/sdk-go/k8s/operator/api/v1"
	"gopkg.in/yaml.v3"
)

type Settings struct {
	Cloud       v1.Provider           `json:"cloud,omitempty"`
	Controllers map[string]Controller `json:"controllers,omitempty"`
	Monitoring  Monitoring            `json:"monitoring,omitempty"`
}

type Controller struct {
	Timeout       time.Duration         `json:"timeout,omitempty"`
	MaxConditions int32                 `json:"max_conditions,omitempty"`
	Phases        []ReconciliationPhase `json:"phases,omitempty"`
}

type ReconciliationPhase struct {
	Name    v1.PhaseCode    `json:"name,omitempty"`
	Backoff []time.Duration `json:"backoff,omitempty"`
}

type Monitoring struct {
	Address  string `json:"adress,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func LoadFromConfigFile(path string) (Settings, error) {
	settings := Settings{}
	dirPath, filename := filepath.Split(path)
	buf, err := fsys.GetFileContent(dirPath, filename)
	if err != nil {
		return settings, err
	}
	err = yaml.Unmarshal(buf, &settings)
	return settings, err
}

type ReconciliationPhaseSettings struct {
}
