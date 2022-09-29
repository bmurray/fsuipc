package fsuipc

import (
	"syscall"
)

var proc_executeCalclatorCode *syscall.LazyProc

type FSUIPC struct {
	x string
}

func New(name string) (*FSUIPC, error) {

	u := &FSUIPC{}
	if proc_executeCalclatorCode == nil {

	}

	////

	return u, nil
}
