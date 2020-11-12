package lib

import (
	"log"
	"os"
	"path/filepath"
)

type Patch struct {
	Prerequisites []string `json:"prerequisites"`
	ID            string   `json:"id"`
	Table         string   `json:"table"`
	Active        bool     `json:"active"`
	Description   string   `json:"description"`
	SQL           string   `json:"sql"`
}

func (p *Patch) setActive(active bool) {
	p.Active = active
}

func ParsePatches(patchesDir string) ([]*Patch, error) {
	var result []*Patch

	err := filepath.Walk(patchesDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if err != nil {
			log.Fatalf("error accessing a path %q: %v\n", path, err)
			return err
		}

		var patch Patch
		ParseFromJsonFile(path, &patch)
		result = append(result, &patch)
		return err
	})

	if err != nil {
		log.Fatalf("error walking the path %q: %v\n", patchesDir, err)
		return nil, err
	}

	return result, nil
}
