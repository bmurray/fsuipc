package fsuipc

/*
#include <stdio.h>
#include <stdlib.h>
typedef void (*lvarListFunction) (int, const char *);
typedef void (*idCallbackFunction) (int[], double[]);
typedef void (*nameCallbackFunction) (const char *[], double[]);
typedef void (*lvarCallbackFunction) (const char *, double);
typedef void (*voidCallbackFunction) ();
extern void go_lvar(int id,  char * name);
extern void go_test_func();
extern void go_lvar_values(char * name, double value);
extern void go_lvar_list(int id, char * name);
void lvarList(int id, const char * name) {
	go_lvar_list(id, (char*)name);
}
void logger(const char * logString) {
	printf("Log: %s\n", logString);
	go_test_func();
}
void nameCallback(const char * lvarName[], double newValue[]) {
	for (int i = 0; lvarName[i] != NULL ; i++) {
		printf("Name callback %s = %f\n", lvarName[i], newValue[i]);
		go_lvar_values((char *)lvarName[i], newValue[i]);
	}
}
void idCallback(int id[], double newValue[]) {

	for (int i = 0; id[i] > 0 ; i++) {
		printf("ID callback %d = %f\n", id[i], newValue[i]);
	}
}
void lvarCallback(const char * name, double value) {
	go_lvar_values((char *)name, value);
}
void updateCallback() {
	printf("Update Callback\n");
}
*/
import "C"
import (
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

func (u *FSUIPC) Init() error {

	args := []uintptr{uintptr(C.logger)}
	r1, _, err := u.proc_init.Call(args...)
	if int32(r1) < 0 {
		return fmt.Errorf("fsuipc_init error: %d %s", int32(r1), err)
	}
	return nil
}
func (u *FSUIPC) GetLvarList() {

	log.Println("LVar List")
	f := C.lvarListFunction(C.lvarList)
	//f := LVar
	args := []uintptr{
		uintptr(unsafe.Pointer(f)),
	}
	r1, r2, err := u.proc_getLvarList.Call(args...)
	log.Println("Execute Lvar List", r1, r2, err)
}

func (u *FSUIPC) RegisterCallbacks() {

	u.RegisterUpdateCallback()
	u.RegisterLvarUpdateCallbackById()
	u.RegisterLvarUpdateCallbackByName()
}
func (u *FSUIPC) RegisterUpdateCallback() {
	f := C.voidCallbackFunction(C.updateCallback)
	args := []uintptr{uintptr(unsafe.Pointer(f))}
	r1, r2, err := u.proc_registerUpdateCallback.Call(args...)
	log.Println("Execute RegisterUpdateCallback", r1, r2, err)
}
func (u *FSUIPC) RegisterLvarUpdateCallbackById() {
	f := C.idCallbackFunction(C.idCallback)
	args := []uintptr{uintptr(unsafe.Pointer(f))}
	r1, r2, err := u.proc_registerLvarUpdateCallbackById.Call(args...)
	log.Println("Execute RegisterLvarUpdateCallbackById", r1, r2, err)
}
func (u *FSUIPC) RegisterLvarUpdateCallbackByName() {
	f := C.nameCallbackFunction(C.nameCallback)
	args := []uintptr{uintptr(unsafe.Pointer(f))}
	r1, r2, err := u.proc_registerLvarUpdateCallbackByName.Call(args...)
	log.Println("Execute RegisterLvarUpdateCallbackById", r1, r2, err)
}
func (u *FSUIPC) FlagLvarForUpdateCallbackByName(name string) {
	args := []uintptr{uintptr(unsafe.Pointer(syscall.StringBytePtr(name)))}
	r1, r2, err := u.proc_flagLvarForUpdateCallbackByName.Call(args...)
	log.Println("exec FlagLvarForUpdateCallbackByName", r1, r2, err)
}
func (u *FSUIPC) FlagLvarForUpdateCallbackById(id uintptr) {
	args := []uintptr{id}
	r1, r2, err := u.proc_flagLvarForUpdateCallbackById.Call(args...)
	log.Println("exec FlagLvarForUpdateCallbackById", r1, r2, err)
}

func (u *FSUIPC) GetLvarValues() {
	f := C.lvarCallbackFunction(C.lvarCallback)
	args := []uintptr{uintptr(unsafe.Pointer(f))}
	r1, r2, err := u.proc_getLvarValues.Call(args...)
	log.Println("Execute GetLvarValues", r1, r2, err)
}
