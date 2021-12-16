package main

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/tacusci/logging/v2"
)

const modulesDir = "modules"

func main() {
	fs, err := fs.ReadDir(os.DirFS("."), modulesDir)
	if err != nil {
		logging.Fatal(err.Error())
	}

	for _, f := range fs {
		fc, err := ioutil.ReadFile(filepath.Join(".", modulesDir, f.Name()))
		if err != nil {
			logging.Fatal(err.Error())
		}
		script := tengo.NewScript(fc)
		script.SetImports(stdlib.GetModuleMap("fmt"))
		proc, err := script.Compile()
		if err != nil {
			logging.Fatal(err.Error())
		}

		proc.Run()
	}
}
