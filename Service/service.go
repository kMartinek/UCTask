package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

func extractNameFromPath(path string) string {
	pathArray := strings.Split(path, "\\")
	nameWithExtension := pathArray[len(pathArray)-1]
	nameWithoutExtension := strings.Split(nameWithExtension, ".")[0]
	return nameWithoutExtension
}

func getBoxInfo(data []byte, startIndex int) (uint32, string, error) {
	var err error = nil
	var boxSize uint32
	var boxType string

	buf := bytes.NewReader(data[startIndex : startIndex+4])
	readErr := binary.Read(buf, binary.BigEndian, &boxSize)
	if readErr != nil {
		err = readErr
	}

	boxType = string(data[startIndex+4 : startIndex+8])

	return boxSize, boxType, err
}

func extractMP4Init(path string) (string, error) {

	//Open file
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	i := 0
	for i < len(data) {

		//Get box size and type
		boxSize, boxType, err := getBoxInfo(data, i)
		if err != nil {
			return "", err
		}

		if boxType == "ftyp" {

			//Check if next box is "moov"
			nextBoxSize, nextBoxType, err := getBoxInfo(data, i+int(boxSize))
			if err != nil {
				return "", err
			}

			if nextBoxType == "moov" {

				//Create new result file
				fileName := extractNameFromPath(path)
				resultFileName := fileName + "_init_" + fmt.Sprintf("%d", time.Now().Unix()) + ".mp4"
				resultFile, err := os.Create(resultFileName)
				if err != nil {
					return "", err
				}

				//Write init segment to File
				resultFile.Write(data[i : i+int(boxSize)+int(nextBoxSize)])

				//Return resultFile path
				resPath, err := filepath.Abs(filepath.Dir(resultFile.Name()))
				if err != nil {
					return "", err
				}
				return (resPath + "\\" + resultFileName), nil
			}
		}
		//Move on to the next box
		i = i + int(boxSize)
	}
	//Whole file analyzed and segment not found
	return "Init segment not found", nil
}

func main() {

	// Connect to a server
	const server = "localhost:4222"

	nc, err := nats.Connect(server)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to:", server)

	//Handle request
	nc.Subscribe("mp4InitSegment", func(m *nats.Msg) {

		fmt.Println("Request recieved...")
		fmt.Println("Filepath received: ", string(m.Data))

		response, err := extractMP4Init(string(m.Data))
		if err != nil {
			nc.Publish(m.Reply, []byte(err.Error()))
		} else {
			nc.Publish(m.Reply, []byte(response))
		}
	})

	runtime.Goexit()
}
