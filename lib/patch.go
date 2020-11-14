package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// type Patch struct {
// 	Prerequisites []string `json:"prerequisites"`
// 	ID            string   `json:"id"`
// 	Table         string   `json:"table"`
// 	Active        bool     `json:"active"`
// 	Description   string   `json:"description"`
// 	SQL           string   `json:"sql"`
// 	Path          string
// }

type Patch struct {
	Prerequisites []string `yaml:"prerequisites"` // ids of prerequisite patches
	ID            string   `yaml:"id"`
	Table         string   `yaml:"table"`
	Active        bool     `yaml:"active"`
	Description   string   `yaml:"description"`
	SQL           string   `yaml:"sql"`
	Path          string
}

func (p *Patch) SetActive(active bool) {
	p.Active = active
}

func (p *Patch) processPatch() {
	// prompt user
	fmt.Printf("Would you like to apply [%v] (Y/n)?\n", p.Path)

	var userIn string
	// this assumes user doesn't supply space in input
	if fmt.Scanln(&userIn); strings.ToLower(userIn) != "y" {
		logger.Printf("user skipped patch %v. User input: %v\n", p.Path, userIn)
		return
	}

	// sql tx
	ApplyPatchTx(p.Table, p.SQL)
}

func ParsePatches(patchesDir string) (map[string]*Patch, error) {
	patches := make(map[string]*Patch)

	err := filepath.Walk(patchesDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if err != nil {
			logger.Fatalf("error accessing a path %q: %v\n", path, err)
			return err
		}

		var patch Patch
		ParseFromYamlFile(path, &patch)
		patch.Path = path

		if patch.Active {
			if _, exists := patches[patch.ID]; exists {
				logger.Fatalf("error: found patches with duplicate id - %v\n", patch.ID)
				panic("duplicate patch id")
			}

			patches[patch.ID] = &patch
		}

		return err
	})

	if err != nil {
		logger.Fatalf("error walking the path %q: %v\n", patchesDir, err)
		return nil, err
	}

	return patches, nil
}

// this is opinionated -- if you got 1 bad patch in your patch folder it won't go through any of them
// one can make the subset of the patches that work into a new patch folder and apply that
func ProcessPatches(patches map[string]*Patch) {
	for id, patch := range patches {
		logger.Printf("checking prereq patch status of patch: %v\n", id)
		// check against prereqs
		for _, prereq_id := range patch.Prerequisites {
			logger.Printf("prerequisite found in patch metadata. prereq patch id: %v. checking its status...\n", prereq_id)
			if prereq_status := CheckPrereqStatus(prereq_id); !prereq_status {
				panic("unapplied prerequisite patch")
			}
		}
	}

	for id, patch := range patches {
		logger.Printf("processing patch: %v\n", id)
		patch.processPatch()
	}
}
