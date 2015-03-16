package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	follow = flag.Bool("f", false, "Follow the files being tailed.")
)

func main() {
	flag.Parse()
	var filename = flag.Arg(0)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	linesToKeep := 10
	currentLines := make([]string, 0, linesToKeep)
	buf := make([]byte, 512)

	for {
		n, err := file.Read(buf)
		if n > 0 {
			s := string(buf[:n])
			lines := strings.Split(s, "\n")

			for _, line := range lines {
				if len(currentLines) == linesToKeep {
					currentLines = currentLines[1:]
				}

				currentLines = append(currentLines, line)
			}
		}

		if err != nil {
			break
		}
	}

	for _, line := range currentLines {
		fmt.Println(line)
	}
}
