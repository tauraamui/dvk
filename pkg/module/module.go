package module

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"sync"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/tacusci/logging/v2"
	"github.com/tauraamui/dvk/stdlibx"
	"github.com/tauraamui/tengox"
)

type Module struct {
	ProcOptions
	proc *tengox.Compiled
}

func New(s []byte) (*Module, error) {
	script := tengox.NewScript(s)
	modMap := tengo.NewModuleMap()
	modMap.AddMap(stdlib.GetModuleMap("fmt", "times", "text"))
	mapStdlibModules(modMap)

	script.SetImports(modMap)

	proc, err := script.CompileRun()
	if err != nil {
		return nil, err
	}

	opts, err := extractOptsFromProc(proc)
	if err != nil {
		return nil, err
	}

	return &Module{
		opts,
		proc,
	}, nil
}

func mapStdlibModules(m *tengo.ModuleMap) {
	for n, s := range stdlibx.Modules {
		s = strings.Replace(s, ";;;;", "`", -1)
		m.AddSourceModule(n, []byte(s))
	}
}

func extractOptsFromProc(proc *tengox.Compiled) (ProcOptions, error) {
	var opts ProcOptions
	modAlias := proc.Get("DVK_MODULE_CMD_ALIAS")
	if modAlias.IsUndefined() {
		return opts, errors.New("script missing alias value")
	}

	opts.CmdAlias = modAlias.String()
	opts.SeekMin = 0
	opts.SeekMax = 100000

	if seekMin := proc.Get("DVK_SEEK_MIN"); !seekMin.IsUndefined() {
		opts.SeekMin = seekMin.Int()
	}

	if seekMax := proc.Get("DVK_SEEK_MAX"); !seekMax.IsUndefined() {
		opts.SeekMax = seekMax.Int()
	}

	return opts, nil
}

type ProcOptions struct {
	CmdAlias         string
	SeekMin, SeekMax int
}

func (m *Module) ExecMain(logsDir fs.FS, args []string, optArgs map[string]string) error {
	logs, err := loadDirLogs(logsDir, m.SeekMin, m.SeekMax)
	if err != nil {
		return err
	}

	modArgs, err := resolveModArgs(logs, args, m.proc)
	if err != nil {
		return err
	}

	if err := loadOptionalArgs(optArgs, m.proc); err != nil {
		return err
	}

	_, err = m.proc.CallByName("main", modArgs...)
	return err
}

func resolveModArgs(logs []interface{}, args []string, proc *tengox.Compiled) ([]interface{}, error) {
	if mainFunc, ok := proc.Get("main").Object().(*tengo.CompiledFunction); ok {
		if 1+len(args) != mainFunc.NumParameters {
			return nil, fmt.Errorf("expected %d additional arguemnts", mainFunc.NumParameters-1)
		}
	}

	modArgs := []interface{}{logs}
	for _, a := range args {
		modArgs = append(modArgs, a)
	}

	return modArgs, nil
}

func loadOptionalArgs(optArgs map[string]string, proc *tengox.Compiled) error {
	for k, v := range optArgs {
		if err := proc.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

func loadDirLogs(logsDir fs.FS, min, max int) ([]interface{}, error) {
	entries, err := fs.ReadDir(logsDir, ".")
	if err != nil {
		return nil, err
	}

	logs := []interface{}{}
	readEntryFiles(logsDir, entries, &logs)

	if min < 0 {
		return nil, fmt.Errorf("seek min %d must be greater than 0", min)
	}

	if max < 0 {
		return nil, fmt.Errorf("seek max %d must be greater than 0", max)
	}

	if min == max {
		return nil, fmt.Errorf("seek min %d and seek max %d must not equal", min, max)
	}

	if min > max {
		return nil, fmt.Errorf("seek min %d cannot be more than seek max %d", min, max)
	}

	if min < len(logs) {
		if max <= len(logs) {
			logs = logs[min:max]
		}
	}

	return logs, nil
}

func readEntryFiles(fsys fs.FS, entries []fs.DirEntry, dest *[]interface{}) {
	lines := make(chan string)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go readFile(&wg, fsys, entries, lines)

	for l := range lines {
		*dest = append(*dest, l)
	}
}

func readFile(wg *sync.WaitGroup, fsys fs.FS, entries []fs.DirEntry, lines chan string) {
	defer wg.Done()
	readingLinesWG := sync.WaitGroup{}
	for _, e := range entries {
		if e.IsDir() {
			return
		}

		f, err := fsys.Open(e.Name())
		if err != nil {
			logging.Error("unable to open %s: %s", e.Name(), err.Error())
			return
		}

		readingLinesWG.Add(1)
		go fileLines(&readingLinesWG, f, lines)
	}
	readingLinesWG.Wait()
	close(lines)
}

func fileLines(wg *sync.WaitGroup, fd fs.File, lines chan<- string) {
	defer wg.Done()

	rr := bufio.NewReader(fd)
	count := 0
	for {
		count++
		line, isPrefix, err := rr.ReadLine()
		if err != nil {
			if err == io.EOF {
				return
			}
		}

		if isPrefix {
			logging.Error("line %d is too long: %v", count, err)
		}

		lines <- string(line)
	}
}
