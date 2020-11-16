package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v4"
)

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
	// check this whether this patch is applied already
	if patchStatus := GetPatchStatus(p.ID); patchStatus != pgx.ErrNoRows {
		logger.Printf("[%v] is already applied! skipping..\n", p.Path)
		return
	}
	// prompt user
	fmt.Printf("Would you like to apply [%v] (Y/n)?\n", p.Path)

	var userIn string
	// this assumes user doesn't supply space in input
	if fmt.Scanln(&userIn); strings.ToLower(userIn) != "y" {
		logger.Printf("user skipped patch %v. User input: %v\n", p.Path, userIn)
		return
	}
	logger.Printf("Applying patch at [%v]\n", p.Path)
	// sql tx
	err := ApplyPatch(p.Table, p.SQL)
	// record patch in the status table
	if err == nil {
		SetPatchStatus(p.ID)
	}
}

func ParsePatches(patchesDir string) (map[string]*Patch, error) {
	patches := make(map[string]*Patch)

	_, err := os.Stat(patchesDir)
	if os.IsNotExist(err) {
		logger.Panicf("patch path [%v] doesn't exist\n", patchesDir)
	}

	err = filepath.Walk(patchesDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if err != nil {
			logger.Panicf("error accessing a path %q: %v\n", path, err)
		}

		var patch Patch
		ParseFromYamlFile(path, &patch)
		patch.Path = path

		if patch.Active {
			if _, exists := patches[patch.ID]; exists {
				logger.Panicf("error: found patches with duplicate id - %v\n", patch.ID)
			}

			patches[patch.ID] = &patch
		}

		return err
	})

	if err != nil {
		logger.Panicf("error walking the path %q: %v\n", patchesDir, err)
	}

	return patches, nil
}

// this is opinionated -- if you got 1 bad patch in your patch folder it won't go through any of them
// one can make the subset of the patches that work into a new patch folder and apply that
func ProcessPatches(patches map[string]*Patch) {
	for id, patch := range patches {
		logger.Printf("checking prereq patch status of patch: %v\n", id)
		// check against prereqs
		for _, prereqID := range patch.Prerequisites {
			logger.Printf("prerequisite found in patch metadata. prereq patch id: %v.\n", prereqID)
			logger.Println("checking its status...")
			if prereqStatus := GetPatchStatus(prereqID); prereqStatus == pgx.ErrNoRows {
				panic("unapplied prerequisite patch")
			}
		}
		logger.Printf("prerequisites of [%v] are all applied!\n", patch.Path)
	}

	for id, patch := range patches {
		logger.Printf("processing patch: [%v] id: %v\n", patch.Path, id)
		patch.processPatch()
	}
}
