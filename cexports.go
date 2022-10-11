package fsuipc

import (
	"C"
	"log"
)
import "sync"

var mu sync.Mutex

var shared *FSUIPC

func SetShared(u *FSUIPC) {
	mu.Lock()
	defer mu.Unlock()
	shared = u
}
func Shared() *FSUIPC {
	mu.Lock()
	defer mu.Unlock()
	return shared
}

//export go_lvar
func go_lvar(id C.int, name *C.char) {
	log.Println("LVAR", id, name)
}

//export go_lvar_values
func go_lvar_values(name *C.char, value C.double) {
	mu.Lock()
	defer mu.Unlock()
	n := C.GoString(name)
	if shared != nil && shared.nameCallback != nil {
		shared.nameCallback(n, float64(value))
	}
	// if n == "dc6_082_obj" {
	// 	log.Println("LVar values", n, float64(value))
	// }
}

//export go_lvar_ids
func go_lvar_ids(id C.int, value C.double) {
	mu.Lock()
	defer mu.Unlock()
	// log.Println("LVar Id", id, value)
	if shared != nil && shared.idCallback != nil {
		shared.idCallback(int(id), float64(value))
	}
}

//export go_lvar_list
func go_lvar_list(id int, name *C.char) {
	// n := C.GoString(name)
	// if n == "dc6_082_obj" {
	// 	log.Println("LVar name", n, id)
	// }
}

//export go_test_func
func go_test_func() {
	log.Println("Test")
}
