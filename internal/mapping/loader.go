package mapping

import (
	config "db_migrate_server/internal/config"
	"db_migrate_server/internal/util"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Load(cfg config.Config) (*Root, error) {

	util.Info.Printf("mapping: Load start mappingPath=%s", cfg.MappingPath)

	entityFiles, err := listJSONFiles(cfg.MappingPath)
	if err != nil {
		return nil, fmt.Errorf("error reading json files : %w", err)
	}
	util.Info.Printf("mapping: found %d entity files", len(entityFiles))

	var entities []Entity
	for _, file := range entityFiles {
		util.Debug.Printf("mapping: loading entity file=%s", file)
		b, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("baca %s: %w", file, err)
		}

		var e Entity
		if err := json.Unmarshal(b, &e); err == nil && e.Entity != "" {
			entities = append(entities, e)
			continue
		}

		return nil, fmt.Errorf("format file %s tidak dikenali (bukan Entity)", file)
	}

	root := &Root{
		Version:  "",
		Entities: entities,
	}

	util.Info.Printf("mapping: Load success entities=%d", len(root.Entities))
	return root, nil
}

func listJSONFiles(path string) ([]string, error) {

	util.Info.Printf("mapping: listing json files path=%s", path)
	var out []string
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		name := strings.ToLower(d.Name())

		// abaikan default.json jika ada
		if strings.HasSuffix(name, ".json") && !strings.Contains(name, "default.json") {
			out = append(out, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(out)
	util.Debug.Printf("mapping: listJSONFiles result=%v", out)
	return out, nil
}
