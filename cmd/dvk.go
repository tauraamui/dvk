package main

import (
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/tacusci/logging/v2"
	"github.com/tauraamui/dvk/pkg/module"
)

type args struct {
	LogsRootDirPath string   `short:"l" long:"logs" description:"Location of logs to analyise" default:"logs"`
	Module          string   `short:"m" long:"module" description:"Module of given alias to run"`
	Args            []string `short:"a" long:"arg" description:"Argument to parse to module's entry point"`
}

func resolveArgs() args {
	args := args{}
	flags.Parse(&args)

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
	if err := m.ExecMain(logsDir, args.Args); err != nil {
		logging.Fatal(err.Error())
	}
}
