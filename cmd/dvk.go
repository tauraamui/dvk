package main

import (
	"flag"
	"os"

	"github.com/tacusci/logging/v2"
	"github.com/tauraamui/dvk/pkg/module"
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
	flag.StringVar(&args.Module, "m", "", "define analysis method")
	flag.Parse()

	return args
}

func main() {
	args := resolveArgs()

	mods, err := module.LoadAllFromDir("modules")
	if err != nil {
		logging.Fatal(err.Error())
	}

	m := mods[args.Module]
	if m == nil {
		logging.Fatal("unable to find module of alias: %s", args.Module)
	}

	logsDir := os.DirFS(args.LogsRootDirPath)
	if err := m.ExecMain(logsDir); err != nil {
		logging.Fatal(err.Error())
	}
}
