# 💻 Native Windows Process Manager in Go

Una herramienta de administración de procesos para Windows desarrollada en **Go (Golang)** puro, diseñada desde una perspectiva de primeros principios para interactuar directamente con la API nativa del sistema operativo (`kernel32.dll`) sin dependencias de terceros.

---

## 🛠️ Arquitectura y Funcionamiento Interno

El proyecto evita abstracciones de alto nivel para comunicarse directamente con el kernel de Windows a través de llamadas al sistema (*syscalls*):

1. **Captura de Estado (*Snapshot*):** Utiliza `CreateToolhelp32Snapshot` para generar una instantánea global de todos los procesos activos en el sistema.
2. **Enumeración de Procesos:** Recorre la instantánea de manera secuencial mediante las funciones de la API Wide-char (`Process32FirstW` y `Process32NextW`) para garantizar la correcta decodificación de caracteres Unicode (UTF-16).
3. **Mapeo de Memoria a Bajo Nivel:** Replica la estructura nativa de C `PROCESSENTRY32` en Go, calculando el alineamiento exacto de bytes mediante `unsafe.Sizeof` y punteros seguros con `unsafe.Pointer`.
4. **Gestión de Handles y Recursos:** Implementa un ciclo de vida estricto para los manejadores del sistema operativo, liberando memoria de forma segura utilizando `defer syscall.CloseHandle`.

---

## 🚀 Características Principales

* **Cero Dependencias Externas:** Construido utilizando exclusivamente la biblioteca estándar de Go y el paquete `syscall`.
* **Soporte Unicode Completo:** Decodificación precisa de nombres de ejecutables (`SzExeFile`) evitando problemas de codificación ANSI.
* **Extracción de Metadatos:** Lectura exacta de Identificadores de Proceso (`PID`) y nombres de ejecutables en tiempo de ejecución.
* **Arquitectura Modular:** Preparado para soportar tanto visualización de procesos como terminación interactiva a demanda.

---

## 📦 Instalación y Uso

### Prerrequisitos
* Sistema operativo **Windows** (x86 / x64).
* Entorno de desarrollo **Go** instalado.

### Clonar el repositorio
```bash
git clone [https://github.com/TU_USUARIO/windows-process-manager.git](https://github.com/TU_USUARIO/windows-process-manager.git)
cd windows-process-manager
```

### Ejecutar la Aplicacion
go run main.go

### 🧠 Conceptos de Sistemas Aplicados
- **Carga Dinámica de DLLs**: Interacción en tiempo de ejecución con librerías nativas del sistema usando syscall.MustLoadDLL y MustFindProc.
- **Aritmetica de Punteros**: Uso del paquete unsafe para la interoperabilidad directa entre estructuras de Go y estructuras nativas de la API de Windows en C.
- **Manejo de Errores a Nivel de SO**: Validación de códigos de retorno de la API Win32 y control de procesos protegidos del sistema.