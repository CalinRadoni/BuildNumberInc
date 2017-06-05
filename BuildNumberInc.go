package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var newVersion int
var tokenFound bool

const verHi int = 1
const verLo int = 4

func genericUsage(out io.Writer) {
	fmt.Fprintf(out, "BuildNumberInc v.%d.%d, Copyright (c) 2016, 2017 Calin Radoni\n", verHi, verLo)
	fmt.Fprintf(out, "https://github.com/CalinRadoni/BuildNumberInc\n")
	fmt.Fprintf(out, "Released under the MIT License\n\n")
	fmt.Fprintf(out, "Usage: ./%s [-c] [-v] [-r] [-h] <fileName> <tokenName>\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(out, "Flags:\n")
	fmt.Fprintf(out, "  -c: search for a const, not for a #define\n")
	fmt.Fprintf(out, "  -v: verbose output\n")
	fmt.Fprintf(out, "  -r: do not add CR to line end\n")
	fmt.Fprintf(out, "  -h: help (this screen)\n")
	fmt.Fprintf(out, "\nExamples:\n")
	fmt.Fprintf(out, "  ./%s version.h SW_BUILD_NUMBER\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(out, "  ./%s -c -v version.h swBuildNumber\n\n", filepath.Base(os.Args[0]))
}

var flagsUsage = func() {
	genericUsage(os.Stderr)
}

func readAndProcessFile(fileName string, fileToken string, searchForConst bool) ([]string, error) {
	var matchString string
	var line string
	var content []string
	var err error
	var version int
	var posToken int    ///< number of capture group for `fileToken`
	var posTokenVal int ///< number of capture group for the value of `fileToken`

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if searchForConst {
		// ^(\s*)(const)(\s+\w+)(\s+\w+)*(\s+)(verBuild0)(\s)*(=)(\s)*([0-9]+)(.*)$
		matchString = "^(\\s*)(const)(\\s+\\w+)(\\s+\\w+)*(\\s+)("
		matchString += fileToken
		matchString += ")(\\s)*(=)(\\s)*([0-9]+)(.*)$"
		posToken = 6
		posTokenVal = 10
	} else {
		// ^(\s*)(#define)(\s+)(TOKEN)(\s+)([0-9]+)(.*)$
		matchString = "^(\\s*)(#define)(\\s+)("
		matchString += fileToken
		matchString += ")(\\s+)([0-9]+)(.*)$"
		posToken = 4
		posTokenVal = 6
	}

	regMatch := regexp.MustCompile(matchString)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
		if regMatch.MatchString(line) {
			regRes := regMatch.FindStringSubmatch(line)
			for i, val := range regRes {
				if i == 0 {
					line = ""
				} else {
					if i == posToken {
						if fileToken != val {
							err = errors.New("internal error ! Supplied file not changed")
							return nil, err
						}
					}
					if i == posTokenVal {
						version, err = strconv.Atoi(val)
						if err != nil {
							return nil, err
						}
						newVersion = version + 1
						line = line + strconv.Itoa(newVersion)
						tokenFound = true
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

func writeResultInFile(content []string, fileName string, addEndR bool) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	wr := bufio.NewWriter(file)
	for _, line := range content {
		if addEndR {
			fmt.Fprintf(wr, "%s\r\n", line)
		} else {
			fmt.Fprintf(wr, "%s\n", line)
		}
	}
	return wr.Flush()
}

func stringHaveSpaces(data string) bool {
	var pos int

	checkFunc := func(ch rune) bool {
		return unicode.IsSpace(ch)
	}
	pos = strings.IndexFunc(data, checkFunc)

	return (pos != -1)
}

func main() {
	var searchForConst bool
	var fileName string
	var fileToken string
	var flagVerbose bool
	var flagHelp bool
	var flagEndR bool

	flag.Usage = flagsUsage
	flag.BoolVar(&searchForConst, "c", false, "")
	flag.BoolVar(&flagVerbose, "v", false, "")
	flag.BoolVar(&flagHelp, "h", false, "")
	flag.BoolVar(&flagEndR, "r", false, "")
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		return
	}

	if flagHelp {
		genericUsage(os.Stdout)
		return
	}

	fileName = flag.Arg(0)

	fileToken = flag.Arg(1)
	if stringHaveSpaces(fileToken) {
		fmt.Fprintf(os.Stderr, "No whitespace allowed in fileToken !\n")
		return
	}

	newVersion = 0
	tokenFound = false

	content, err := readAndProcessFile(fileName, fileToken, searchForConst)
	if err != nil {
		log.Fatalf("Error in reading and processing: %v.\n", err)
	}

	if err = writeResultInFile(content, fileName, !flagEndR); err != nil {
		log.Fatalf("Write error: %v", err)
	}

	if flagVerbose {
		if tokenFound {
			fmt.Fprintf(os.Stdout, "BuildNumberInc: %s\\%s increased to %d\n", filepath.Base(fileName), fileToken, newVersion)
		} else {
			fmt.Fprintf(os.Stdout, "%s NOT found in %s !\n", fileToken, fileName)
		}
	} else {
		if !tokenFound {
			fmt.Fprintf(os.Stdout, "%s NOT found !\n", fileToken)
		}
	}
}
