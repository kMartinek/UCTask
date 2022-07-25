/*
TO DO

- start server
- handle request
- recieve file
- process mp4 file
	- read initialization segment
	- write initialization segment to new file
- return new file

*/

package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/nats-io/nats.go"

	"os"
)

/*func getBoxSize(file byte[], position int) int{

}*/

func parseMP4(path string) string {
	dat, err := os.ReadFile(path)
	check(err)

	//check if result file exists
	//if true delete content or delete file

	//create new result file
	resultFile, err := os.Create("result.txt")
	check(err)

	i := 0
	var boxSize uint32
	//var boxType string
	for i < len(dat) {
		//read size of the box
		buf := bytes.NewReader(dat[i : i+4])
		err := binary.Read(buf, binary.BigEndian, &boxSize)
		check(err)

		resultFile.WriteString("Box Size: " + fmt.Sprintf("%d", boxSize) + "\n")
		fmt.Println("Box size: ", boxSize)
		//read type of the box

		dio := string(dat[i+4 : i+8])
		resultFile.WriteString("Box Type: " + dio + "\n")
		fmt.Println("boxType: ", dio)

		if dio == "ftyp" {
			resultFile.WriteString("ftyp data: " + string(dat[i+8:i+int(boxSize)]) + "\n")
			fmt.Println("Data:", string(dat[i+8:i+int(boxSize)]))
		}

		resultFile.WriteString("Box indexes: [" + fmt.Sprintf("%d", i) + ":" + fmt.Sprintf("%d", i+int(boxSize)) + "]\n\n")
		//read data
		i = i + int(boxSize)
		fmt.Println("Brojač: ", i)

	}
	resPath, err := filepath.Abs(filepath.Dir(resultFile.Name()))
	check(err)
	fmt.Println(resPath + "\\" + resultFile.Name())
	return resPath
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	parseMP4("C:\\Users\\Kac\\Downloads\\video.mp4")

	// Connect to a server

	// Connect to a server
	nc, err := nats.Connect("localhost:4222")
	check(err)

	nc.Subscribe("mp4InitSegment", func(m *nats.Msg) {
		fmt.Println("Dobio poruku")
		fmt.Println(string(m.Data))
		nc.Publish(m.Reply, []byte("parseMP4(string(m.Data))"))
	})

	runtime.Goexit()

}
