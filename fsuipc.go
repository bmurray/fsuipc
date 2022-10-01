package fsuipc

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

var proc_executeCalclatorCode *syscall.LazyProc
var proc_init *syscall.LazyProc
var proc_start *syscall.LazyProc
var proc_isRunning *syscall.LazyProc
var proc_end *syscall.LazyProc

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
			path, err := os.Getwd()
			if err != nil {
				return nil, fmt.Errorf("cannot get cwd: %w", err)
			}
			dllPath = filepath.Join(path, "FSUIPC_WAPID.dll")
			if _, err := os.Stat(dllPath); os.IsNotExist(err) {
				return nil, fmt.Errorf("Cannot open file :%w", err)
			}
			fmt.Println("got secondary file")
		}
		log.Println("MOD", dllPath)
		mod := syscall.NewLazyDLL(dllPath)
		if err = mod.Load(); err != nil {
			return nil, err
		}
		log.Println("getting symbols")
		proc_executeCalclatorCode = mod.NewProc("fsuipcw_executeCalclatorCode")
		proc_start = mod.NewProc("fsuipcw_start")
		proc_isRunning = mod.NewProc("fsuipcw_isRunning")
		proc_end = mod.NewProc("fsuipcw_end")
		proc_init = mod.NewProc("fsuipcw_init")

	}
	log.Println("loaded", proc_start)

	args := []uintptr{uintptr(0)}
	r1, _, err := proc_init.Call(args...)
	if int32(r1) < 0 {
		return nil, fmt.Errorf("fsuipc_init error: %d %s", int32(r1), err)
	}
	log.Println("Initialized")
	args2 := []uintptr{}
	proc_start.Call(args2...)
	// if int32(r1) < 0 {
	// 	return nil, fmt.Errorf("fsuipc_start error: %d %s", int32(r1), err)
	// }

	return u, nil
}

func (u *FSUIPC) IsRunning() bool {
	r1, _, err := proc_isRunning.Call()
	if int32(r1) < 0 {
		log.Printf("fsuipc_isRunning error: %d %s", int32(r1), err)
		return false
	}
	log.Println("running", r1)
	return true
}
func (u *FSUIPC) End() {
	r1, _, err := proc_end.Call()
	if int32(r1) < 0 {
		log.Printf("fsuipc_end error: %d %s", int32(r1), err)
		return
	}
	//log.Println("stop", r1)
	return
}
func (u *FSUIPC) ExecuteCalclatorCode(code string) {
	args := []uintptr{
		uintptr(unsafe.Pointer(syscall.StringBytePtr(code))),
	}

	r1, r2, err := proc_executeCalclatorCode.Call(args...)
	log.Println("Execute", r1, r2, err)
	// if int32(r1) < 0 {

	// }
}
