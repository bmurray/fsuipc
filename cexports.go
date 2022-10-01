package fsuipc

import (
	"C"
	"log"
)

//export go_lvar
func go_lvar(id C.int, name *C.char) {
	log.Println("LVAR", id, name)
}

//export go_lvar_values
func go_lvar_values(name *C.char, value C.double) {
	n := C.GoString(name)
	if n == "dc6_082_obj" {
		log.Println("LVar values", n, float64(value))
	}
}

//export go_lvar_list
func go_lvar_list(id int, name *C.char) {
	n := C.GoString(name)
	if n == "dc6_082_obj" {
		log.Println("LVar name", n, id)
	}
}

//export go_test_func
func go_test_func() {
	log.Println("Test")
}
