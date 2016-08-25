package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func readAndProcessFile(fileName string) ([]string, error) {
	var line string
	var content []string
	var err error
	var version int

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	regMatch := regexp.MustCompile("^(\\s*)(#define)(\\s+)(SW_VER_BUILD)(\\s+)([0-9]+)(.*)$")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
		if regMatch.MatchString(line) {
			regRes := regMatch.FindStringSubmatch(line)
			for i, val := range regRes {
				if i == 0 {
					line = ""
				} else {
					if i == 6 {
						version, err = strconv.Atoi(val)
						if err != nil {
							return nil, err
						}
						line = line + strconv.Itoa(version+1)
					} else {
						line = line + val
					}
				}
			}
		}
		content = append(content, line)
	}

	return content, scanner.Err()
}

func writeResultInFile(content []string, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	wr := bufio.NewWriter(file)
	for _, line := range content {
		fmt.Fprintln(wr, line)
	}
	return wr.Flush()
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Source file name not specified !")
	}

	fileName := os.Args[1]

	content, err := readAndProcessFile(fileName)
	if err != nil {
		log.Fatalf("Read error: %v", err)
	}

	if err = writeResultInFile(content, fileName); err != nil {
		log.Fatalf("Write error: %v", err)
	}
}
