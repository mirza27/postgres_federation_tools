package mapping

import (
	config "db_migrate_server/internal/config"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Load(cfg config.Config) (*Root, error) {

	conf, err := setDefaultConfig(cfg.DefaultConfigPath)
	if err != nil {
		return nil, fmt.Errorf("gagal memuat default config dari %s: %w", cfg.DefaultConfigPath, err)
	}

	entityFiles, err := listJSONFiles(cfg.MappingPath)
	if err != nil {
		return nil, fmt.Errorf("error reading json files : %w", err)
	}

	var entities []Entity
	for _, file := range entityFiles {
		b, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("baca %s: %w", file, err)
		}

		var e Entity
		if err := json.Unmarshal(b, &e); err == nil && e.Entity != "" {
			entities = append(entities, e)
			continue
		}

		return nil, fmt.Errorf("format file %s tidak dikenali (bukan Entity, []Entity, atau Root berisi entities)", file)
	}

	root := &Root{
		Version:  conf.Version,
		Sources:  conf.Sources,
		Target:   conf.Target,
		Engine:   conf.Engine,
		Entities: entities,
	}

	return root, nil
}

func setDefaultConfig(path string) (*Root, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("default config error: %w", err)
	}

	// Struct penampung fleksibel agar bisa menangkap variasi field
	var tmp struct {
		// Field Root standar
		Version  string   `json:"version"`
		Sources  []Source `json:"sources"`
		Target   Target   `json:"target"`
		Topics   *Topics  `json:"topics,omitempty"`
		Engine   Engine   `json:"engine"`
		Entities []Entity `json:"entities"`

		// Bentuk singkat (opsional) yang mungkin ada di default.json
		Name          string `json:"name"`
		Type          string `json:"type"`
		DefaultSchema string `json:"default_schema"`
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return nil, fmt.Errorf("gagal parse default.json: %w", err)
	}

	// Susun Root dari tmp
	res := &Root{
		Version:  tmp.Version, // kosong jika tidak ada di file — itu OK
		Sources:  tmp.Sources,
		Target:   tmp.Target,
		Topics:   tmp.Topics,
		Engine:   tmp.Engine,
		Entities: tmp.Entities,
	}

	// Jika "sources" kosong tapi ada triple (name/type/default_schema) di top-level, konversi jadi satu Source
	if len(res.Sources) == 0 && tmp.Name != "" && tmp.Type != "" {
		res.Sources = []Source{{
			Name:          tmp.Name,
			Type:          tmp.Type,
			DefaultSchema: tmp.DefaultSchema,
		}}
	}

	// Pastikan Entities tidak nil
	if res.Entities == nil {
		res.Entities = []Entity{}
	}

	return res, nil
}

func listJSONFiles(path string) ([]string, error) {

	var out []string
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		name := strings.ToLower(d.Name())

		// if not default.json
		if strings.HasSuffix(name, ".json") && !strings.Contains(name, "default.json") {
			out = append(out, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(out)
	return out, nil
}
