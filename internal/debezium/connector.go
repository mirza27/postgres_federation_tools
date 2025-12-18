package debezium

import (
	"bytes"
	"db_migrate_server/internal/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type DebeziumConnector struct {
	Name            string
	ConnectorConfig map[string]string
	Config          *config.Config
}

func NewDebeziumConnector(name string, cfg *config.Config, snapshot bool) (*DebeziumConnector, error) {
	templatePath := filepath.Join("internal", "debezium", "connector.example.postgres.json")
	body, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("read connector template: %w", err)
	}

	var raw struct {
		Name   string            `json:"name"`
		Config map[string]string `json:"config"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("parse connector template: %w", err)
	}

	raw.Name = name
	raw.Config["database.hostname"] = cfg.SourceDatabaseHost
	raw.Config["database.port"] = strconv.Itoa(cfg.SourceDatabasePort)
	raw.Config["database.user"] = cfg.SourceDatabaseUser
	raw.Config["database.password"] = cfg.SourceDatabasePass
	raw.Config["database.dbname"] = cfg.SourceDatabaseName
	if snapshot {
		raw.Config["snapshot.mode"] = "initial"
	} else {
		raw.Config["snapshot.mode"] = "never"
	}

	return &DebeziumConnector{
		Name:            name,
		Config:          cfg,
		ConnectorConfig: raw.Config,
	}, nil
}

func (dc *DebeziumConnector) GetConnectorConfig() map[string]string {
	return dc.ConnectorConfig
}

// Create sends a POST to Debezium to create the connector.
func (dc *DebeziumConnector) Create() error {
	body, err := dc.toJSON()
	if err != nil {
		return err
	}

	baseUrl := fmt.Sprintf("%s:%d", dc.Config.DebeziumHost, dc.Config.DebeziumPort)
	url := strings.TrimRight(baseUrl, "/") + "/connectors"

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("connect error: %s (%s)", resp.Status, string(b))
	}
	return nil
}

func (dc *DebeziumConnector) toJSON() ([]byte, error) {
	payload := struct {
		Name   string            `json:"name"`
		Config map[string]string `json:"config"`
	}{
		Name:   dc.Name,
		Config: dc.ConnectorConfig,
	}
	return json.MarshalIndent(payload, "", "  ")
}

// Status calls Debezium to fetch connector status.
func (dc *DebeziumConnector) Status() (string, error) {

	baseUrl := fmt.Sprintf("%s:%d", dc.Config.DebeziumHost, dc.Config.DebeziumPort)
	url := strings.TrimRight(baseUrl, "/") + "/connectors/" + dc.Name + "/status"

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("status error: %s (%s)", resp.Status, string(b))
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
