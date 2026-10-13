package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	TH32CS_SNAPPROCESS = 0x00000002
)

type PROCESSENTRY32 struct {
	DwSize              uint32      // indica el tamaño de la estructura
	CntUsage            uint32      // antiguamente contaba las referencias a procesos. Hoy en dia ya no se usa
	Th32ProcessID       uint32      // el PID del proceso
	Th32DefaultHeapID   uintptr     // uintptr maneja automáticamente 32 o 64 bits según tu Windows
	Th32ModuleID        uint32      // antiguamente servia para identificar el modulo principal del proceso. Hoy en dia esta obsoleto
	CntThreads          uint32      // la cantidad de hilos
	Th32ParentProcessID uint32      // el PID del proceso padre. Te dice que programa abrio a este
	PcPriClassBase      int32       // Le indica al planificador de windos cuanta atencion de la CPU debe recibir este programa a diferencia de los demas
	DwFlags             uint32      // reservado para uso interno de Windows (flags de control)
	SzExeFile           [260]uint16 // array de 260 caracteres para guardar el nombre del .exe en UTF-16
}

func main() {
	// Cargamos el archivo DLL de windows (kernel32.dll) con la funcion MustLoadDLL,
	// si la DLL no existe, o no se puede cargar hace un panic y detiene el programa.
	kernel32 := syscall.MustLoadDLL("kernel32.dll")

	// Buscamos la función "CreateToolhelp32Snapshot" dentro de la DLL que ya cargamos,
	// esta funcion que buscamos es la encargada de darnos los procesos.
	// CreateToolhelp32Snapshot pide la bandera de procesos 0x00000002 y un ID 0.
	createSnapshotProc := kernel32.MustFindProc("CreateToolhelp32Snapshot")

	// El método .Call() es el puente final donde Go le entrega el control al sistema operativo y espera una respuesta.
	// r1: El resultado principal de la función, r2: Un segundo valor de retorno secundario que algunas funciones de Windows usan para datos extra (no la vamos a usar),
	// err: El error del sistema operativo (syscall.Errno).
	r1, _, err := createSnapshotProc.Call(uintptr(TH32CS_SNAPPROCESS), uintptr(0))

	if r1 == uintptr(syscall.InvalidHandle) {
		panic(fmt.Sprintf("Error al tomar la instantánea: %v", err))
	}

	// Asegurarnos de cerrar el handle de la foto al terminar (buena práctica).
	defer syscall.CloseHandle(syscall.Handle(r1))

	// Creamos una instancia de nuestro molde (struct).
	var pe32 PROCESSENTRY32

	// Importante en Windows para que no tire error indicar el tamaño ezacto del struct en bytes.
	pe32.DwSize = uint32(unsafe.Sizeof(pe32))

	// Buscamos la funcion "Process32First", este lee el primer proceso para guardarlo en nuestro struct.
	process32FirstProc := kernel32.MustFindProc("Process32FirstW")

	// En Go, para pasarle la dirección de memoria de una variable a una función de bajo nivel, usamos el paquete unsafe con la función unsafe.Pointer(&pe32).
	// success devuelve TRUE o 1 si encontro el primer proceso y relleno pe32. Si devuelve FALSE o 0 significa que la lista esta vacia o hubo un error.
	success, _, err := process32FirstProc.Call(r1, uintptr(unsafe.Pointer(&pe32)))

	// Validamos si Process32First encontró el primer proceso con éxito
	if success == 0 {
		fmt.Println("No se encontraron procesos o falló la lectura.")
		return
	}

	// Buscamos la funcion "Process32Next" para ir pasando al siguiente proceso en la lista hasta que no queden más.
	process32NextProc := kernel32.MustFindProc("Process32NextW")

	// Armamos un bucle for para recorrer todos los procesos
	for {

		// Buscamos el nombre del proceso, y lo convertimos a string
		nombreProceso := syscall.UTF16ToString(pe32.SzExeFile[:])

		// Imprimimos el ID del proceso junto al nombre del mismo
		fmt.Printf("ID: %d \t Nombre del Proceso: %s\n", pe32.Th32ProcessID, nombreProceso)

		// Y luego llamamos a Process32Next para saltar al siguiente
		ret, _, _ := process32NextProc.Call(r1, uintptr(unsafe.Pointer(&pe32)))

		// Si ret es 0, significa que llegamos al final de la lista y el bucle termina
		if ret == 0 {
			break
		}
	}
}
