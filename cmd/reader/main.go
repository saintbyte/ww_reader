package main

import (
	"fmt"
	"github.com/adrianmo/go-nmea"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"log/slog"
	"strings"
)

func normalizeLine(s string) string {
	s = strings.Trim(s, "\r")
	return s
}
func readFromChan(port io.ReadWriteCloser, channel chan string) {
	var line string
	var bufStr string
	line = ""
	for {
		buf := make([]byte, 32)
		n, err := port.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from serial port: ", err)
			}
		} else {
			buf = buf[:n]
			//fmt.Println("Rx: ", hex.EncodeToString(buf))
			if strings.Contains(string(buf), "\n") {
				bufStr = string(buf)
				arr := strings.Split(bufStr, "\n")
				line = line + arr[0]
				channel <- normalizeLine(line)
				if len(arr) > 2 {
					for i := 1; i < len(arr)-1; i++ {
						channel <- normalizeLine(arr[i])
					}
				}
				line = arr[len(arr)-1]
			} else {
				line = line + string(buf)
			}

		}
	}
}
func parseString(seq string) {
	fmt.Println("received", seq)
	s, err := nmea.Parse(seq)
	if err != nil {
		slog.Error("Nmea parse error:", err)
		return
	}
	if s.DataType() == nmea.TypeXDR {

	}
}

func defineCustomNMEAParsers() {
	err = nmea.RegisterParser("XYZ", func(s nmea.BaseSentence) (nmea.Sentence, error) {
		// This example uses the package builtin parsing helpers
		// you can implement your own parsing logic also
		p := nmea.NewParser(s)
		return LILType{
			BaseSentence: s,
			Time:         p.Time(0, "time"),
			Label:        p.String(1, "label"),
			Counter:      p.Int64(2, "counter"),
			Value:        p.Float64(3, "value"),
		}, p.Err()
	})
}

func main() {
	// stty -F /dev/ttyUSB0 115200 cs8 -cstopb -parenb
	defineCustomNMEAParsers()
	var err error
	options := serial.OpenOptions{
		PortName:              "/dev/ttyUSB0",
		BaudRate:              9600,
		DataBits:              8,
		StopBits:              1,
		MinimumReadSize:       0,
		InterCharacterTimeout: 100,
	}

	port, err := serial.Open(options)
	if err != nil {
		slog.Error("Open port error:", err)
		return
	}
	defer port.Close()

	outputChan := make(chan string)
	go readFromChan(port, outputChan)
	for {
		select {
		case msg1 := <-outputChan:
			parseString(msg1)
		}
	}
}
