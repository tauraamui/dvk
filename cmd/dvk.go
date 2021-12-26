package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/tacusci/logging/v2"
)

type Args struct {
	LogsRootDirPath string
	Module          string
}

func resolveArgs() Args {
	args := Args{
		LogsRootDirPath: "logs",
	}

	flag.StringVar(&args.LogsRootDirPath, "ldir", "logs", "location of logs to analyise")
	flag.StringVar(&args.Module, "method", "", "define analysis method")
	flag.Parse()

	return args
}

type moduleCmds map[string]*tengo.Compiled

func loadModules(modulesDir string) moduleCmds {
	c := moduleCmds{}

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

		modAlias := proc.Get("MODULE_CMD_ALIAS").String()

		if len(modAlias) > 0 {
			// TODO (tauraamui): guard against multiple module alias collisions
			c[modAlias] = proc
		}
		fmt.Printf("MOD ALIAS: %s\n", modAlias)
	}

	return c
}

func main() {
	args := resolveArgs()
	fmt.Printf("ARGS: %v\n", args)

	loadModules("modules")
}
