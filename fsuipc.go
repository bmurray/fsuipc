package fsuipc

// #include <stdlib.h>
import "C"
import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

type FSUIPC struct {
	x                           string
	proc_executeCalclatorCode   *syscall.LazyProc
	proc_init                   *syscall.LazyProc
	proc_start                  *syscall.LazyProc
	proc_isRunning              *syscall.LazyProc
	proc_end                    *syscall.LazyProc
	proc_getLvarList            *syscall.LazyProc
	proc_logLvars               *syscall.LazyProc
	proc_reload                 *syscall.LazyProc
	proc_getLvarIdFromName      *syscall.LazyProc
	proc_getLvarFromId          *syscall.LazyProc
	proc_getLvarFromName        *syscall.LazyProc
	proc_setLogLevel            *syscall.LazyProc
	proc_getLvarValues          *syscall.LazyProc
	proc_getLvarUpdateFrequency *syscall.LazyProc
	proc_setLvarUpdateFrequency *syscall.LazyProc

	// updates
	proc_registerUpdateCallback           *syscall.LazyProc
	proc_flagLvarForUpdateCallbackById    *syscall.LazyProc
	proc_flagLvarForUpdateCallbackByName  *syscall.LazyProc
	proc_registerLvarUpdateCallbackById   *syscall.LazyProc
	proc_registerLvarUpdateCallbackByName *syscall.LazyProc
}

func Double(d uintptr) float64 {
	s := (*float64)(unsafe.Pointer(&d))
	return *s
}

func New(name string) (*FSUIPC, error) {

	u := &FSUIPC{}

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
	u.proc_executeCalclatorCode = mod.NewProc("fsuipcw_executeCalclatorCode")
	u.proc_start = mod.NewProc("fsuipcw_start")
	u.proc_isRunning = mod.NewProc("fsuipcw_isRunning")
	u.proc_end = mod.NewProc("fsuipcw_end")
	u.proc_init = mod.NewProc("fsuipcw_init")

	u.proc_getLvarList = mod.NewProc("fsuipcw_getLvarList")
	u.proc_logLvars = mod.NewProc("fsuipcw_logLvars")
	u.proc_reload = mod.NewProc("fsuipcw_reload")
	u.proc_getLvarIdFromName = mod.NewProc("fsuipcw_getLvarIdFromName")
	u.proc_getLvarFromId = mod.NewProc("fsuipcw_getLvarFromId")
	u.proc_getLvarFromName = mod.NewProc("fsuipcw_getLvarFromName")
	u.proc_setLogLevel = mod.NewProc("fsuipcw_setLogLevel")
	u.proc_registerUpdateCallback = mod.NewProc("fsuipcw_registerUpdateCallback")
	u.proc_flagLvarForUpdateCallbackById = mod.NewProc("fsuipcw_flagLvarForUpdateCallbackById")
	u.proc_flagLvarForUpdateCallbackByName = mod.NewProc("fsuipcw_flagLvarForUpdateCallbackByName")
	u.proc_registerLvarUpdateCallbackById = mod.NewProc("fsuipcw_registerLvarUpdateCallbackById")
	u.proc_registerLvarUpdateCallbackByName = mod.NewProc("fsuipcw_registerLvarUpdateCallbackByName")
	u.proc_getLvarValues = mod.NewProc("fsuipcw_getLvarValues")
	u.proc_getLvarUpdateFrequency = mod.NewProc("fsuipcw_getLvarUpdateFrequency")
	u.proc_setLvarUpdateFrequency = mod.NewProc("fsuipcw_setLvarUpdateFrequency")

	log.Println("loaded", u.proc_start)
	u.Init()
	log.Println("Initialized")
	args2 := []uintptr{}
	u.proc_start.Call(args2...)
	// if int32(r1) < 0 {
	// 	return nil, fmt.Errorf("fsuipc_start error: %d %s", int32(r1), err)
	// }

	return u, nil
}

func (u *FSUIPC) IsRunning() bool {
	r1, _, err := u.proc_isRunning.Call()
	if int32(r1) < 0 {
		log.Printf("fsuipc_isRunning error: %d %s", int32(r1), err)
		return false
	}
	log.Println("running", r1)
	return true
}
func (u *FSUIPC) End() {
	r1, _, err := u.proc_end.Call()
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

	r1, r2, err := u.proc_executeCalclatorCode.Call(args...)
	if r1 > 0 {
		log.Println("Error: ExecuteCalclatorCode", r1, r2, err)
	}
}
func (u *FSUIPC) LogLVars() {
	args2 := []uintptr{}
	r1, r2, err := u.proc_start.Call(args2...)
	log.Println("LOG Lvars", r1, r2, err)

}
func (u *FSUIPC) Reload() {
	args2 := []uintptr{}
	r1, r2, err := u.proc_reload.Call(args2...)
	log.Println("LOG Reload", r1, r2, err)

}
func (u *FSUIPC) GetLvarIdFromName(name string) uintptr {
	args := []uintptr{
		uintptr(unsafe.Pointer(syscall.StringBytePtr(name))),
	}
	r1, _, _ := u.proc_getLvarIdFromName.Call(args...)
	//log.Println("LOG GetLvarIdFromName", r1, r2, err)

	return r1
}
func (u *FSUIPC) GetLvarFromName(name string) float64 {
	args := []uintptr{
		uintptr(unsafe.Pointer(syscall.StringBytePtr(name))),
	}
	r1, r2, err := u.proc_getLvarFromName.Call(args...)
	if r1 > 0 {
		log.Println("GetLvarFromName: ", err)
	}
	return Double(r2)
}

func (u *FSUIPC) GetLvarFromId(id uintptr) float64 {
	args := []uintptr{
		id,
	}
	r1, r2, err := u.proc_getLvarFromId.Call(args...)
	if r1 > 0 {
		log.Println("GetLvarFromId: ", err)
	}
	return Double(r2)
}
func (u *FSUIPC) SetLogLevel(level int) {
	args := []uintptr{
		uintptr(level),
	}
	u.proc_setLogLevel.Call(args...)
	//log.Println("LOG setLogLevel", r1, r2, err)
}
func (u *FSUIPC) GetLvarUpdateFrequency() int {
	args := []uintptr{}
	r1, _, _ := u.proc_getLvarUpdateFrequency.Call(args...)
	// log.Println("LOG GetLvarUpdateFrequency", r1, r2, err)
	return int(r1)
}
func (u *FSUIPC) SetLvarUpdateFrequency(freq int) {
	args := []uintptr{uintptr(freq)}
	u.proc_setLvarUpdateFrequency.Call(args...)
	// log.Println("LOG SetLvarUpdateFrequency", r1, r2, err)

}
