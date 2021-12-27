package module

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Table map[string]*Module

func LoadAllFromDir(dir string) (Table, error) {
	fs, err := fs.ReadDir(os.DirFS("."), dir)
	if err != nil {
		return nil, err
	}

	table := make(Table)
	for _, f := range fs {
		fc, err := ioutil.ReadFile(filepath.Join(".", dir, f.Name()))
		if err != nil {
			return nil, err
		}

		mod, err := New(fc)
		if err != nil {
			return nil, err
		}

		table[mod.Alias] = mod
	}

	return table, nil
}
