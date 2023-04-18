package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/tarm/serial"
)

func main() {
	config := &serial.Config{
		Name:        "/dev/ttyUSB0",
		Baud:        9600,
		ReadTimeout: 1,
		Size:        8,
	}

	stream, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(stream)

	for {
		line, isPrefix, err := r.ReadLine()

		if err == nil && !isPrefix {
			decodedData, err := decode(line)
			if err != nil {
				log.Println(err)
				continue
			}

			text := fmt.Sprintf("%v,%v,%v\n", decodedData.Timestamp, decodedData.SensorID, decodedData.Value)

			f, err := os.OpenFile("data.csv",
				os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Println(err)
			}
			defer f.Close()
			if _, err := f.WriteString(text); err != nil {
				log.Println(err)
			}
		}
	}
}

type ProtocolMessage struct {
	Timestamp time.Time
	SensorID  string
	Value     string
}

func decode(data []byte) (*ProtocolMessage, error) {
	s := string(data)
	if !strings.HasPrefix(s, "###") || !strings.HasSuffix(s, "$$$") {
		return nil, errors.New("data does not comply to protocol")
	}

	result := strings.Split(strings.TrimSuffix(strings.TrimPrefix(s, "###"), "$$$"), ":")

	if len(result) < 2 {
		return nil, errors.New("data does not comply to protocol")
	}

	decodedData := ProtocolMessage{
		Timestamp: time.Now().UTC(),
		SensorID:  result[0],
		Value:     result[1],
	}

	return &decodedData, nil
}
