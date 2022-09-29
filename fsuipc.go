package fsuipc

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

var proc_executeCalclatorCode *syscall.LazyProc

type FSUIPC struct {
	x string
}

func New(name string) (*FSUIPC, error) {

	u := &FSUIPC{}
	if proc_executeCalclatorCode == nil {
		exePath, err := os.Executable()
		if err != nil {
			return nil, err
		}
		dllPath := filepath.Join(filepath.Dir(exePath), "FSUIPC_WAPID.dll")
		if _, err := os.Stat(dllPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("Cannot open file :%w", err)
		}
		mod := syscall.NewLazyDLL(dllPath)
		if err = mod.Load(); err != nil {
			return nil, err
		}
		proc_executeCalclatorCode = mod.NewProc("executeCalclatorCode")
	}

	////

	return u, nil
}
