package module

import (
	"errors"

	"github.com/d5/tengo/v2/stdlib"
	"github.com/tauraamui/tengox"
)

type Module struct {
	Alias string
	proc  *tengox.Compiled
}

func New(s []byte) (*Module, error) {
	script := tengox.NewScript(s)
	script.SetImports(stdlib.GetModuleMap("fmt"))

	proc, err := script.CompileRun()
	if err != nil {
		return nil, err
	}

	modAlias := proc.Get("MODULE_CMD_ALIAS").String()
	if len(modAlias) == 0 {
		return nil, errors.New("script missing alias value")
	}

	return &Module{
		Alias: modAlias,
		proc:  proc,
	}, nil
}

func (m *Module) ExecMain() error {
	_, err := m.proc.CallByName("main")
	return err
}
