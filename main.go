package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

var selectedPort *enumerator.PortDetails
var buff bytes.Buffer

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
	}

	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}

	fmt.Println("Please indicate the number or name of the serial port you want to connect:")
	for i, p := range ports {
		if p.IsUSB {
			fmt.Printf("%d: %s\n", i, p.Name)
		}
	}

	fmt.Print("Input: ")
	var selection string
	_, err = fmt.Scanln(&selection)
	if err != nil {
		log.Fatal(err)
	}

	for i, p := range ports {
		if p.IsUSB && (selection == p.Name || selection == fmt.Sprintf("%d", i)) {
			selectedPort = p
			fmt.Println("Selected port:", p.Name)
			break
		}
	}

	port, err := serial.Open(selectedPort.Name, &serial.Mode{
		BaudRate: 9600,
	})

	if err != nil {
		log.Fatal(err)
	}

	go serialRead(port)

	s := <-sig
	port.Close()
	fmt.Println("Signal:", s)
}

func serialRead(port serial.Port) {
	for {
		read := make([]byte, 128)
		_, err := port.Read(read)
		if err != nil {
			log.Println(err)
		}
		read = bytes.Trim(read, "\x00")
		buff.Write(read)

		for {
			b, err := buff.ReadBytes('\n')
			if err == io.EOF {
				buff.Reset()
				buff.Write(b)
				break
			}
			fmt.Print(string(b))
		}
	}
}
