package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"example.com/exam"
	"github.com/gocarina/gocsv"
)

type Pattern struct {
	FilePettern *regexp.Regexp
}

func (p *Pattern) String() string {
	if p.FilePettern == nil {
		p.FilePettern = regexp.MustCompile(`[\s\S]*`)
	}
	return p.FilePettern.String()
}

func (p *Pattern) Set(s string) error {
	pattern, err := regexp.Compile(s)
	if err != nil {
		return err
	}

	p.FilePettern = pattern
	return nil
}

func newWalkFunc(ignoreFileMap map[string]bool, pattern *regexp.Regexp, ignoreFile, resultFile *os.File) filepath.WalkFunc {
	return func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("unknown error: %w", err)
		}

		if info.IsDir() {
			return nil
		}

		if _, ok := ignoreFileMap[info.Name()]; ok {
			return nil
		}

		if !pattern.MatchString(info.Name()) {
			return nil
		}

		foundFile, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("unable to open file: %w", err)
		}

		results, err := exam.ParseAnswers(foundFile)
		if err != nil {
			return err
		}

		if err := gocsv.MarshalWithoutHeaders(results, resultFile); err != nil {
			return err
		}

		if ignoreFile != nil {
			if _, err := ignoreFile.WriteString(info.Name() + "\n"); err != nil {
				return err
			}
		}

		return nil
	}
}

func main() {

	pattern := Pattern{}
	flag.Var(&pattern, "r", "Regular expression indicating the name of the target file")

	var rawResultFilePath string
	flag.StringVar(&rawResultFilePath, "o", "results.csv", "Name of the file in which to record the tally results")

	var rawIgnoreFilePath string
	flag.StringVar(&rawIgnoreFilePath, "i", "ignore.txt", "Specify the destination for exporting the analyzed file name")

	var notAppendIgnoreFile bool
	flag.BoolVar(&notAppendIgnoreFile, "d", false, "Do not append analyzed files to the ignore file")

	flag.Parse()

	resultFilePath := filepath.Clean(rawResultFilePath)
	ignoreFilePath := filepath.Clean(rawIgnoreFilePath)

	args := flag.Args()

	sourcePath := filepath.Clean(args[0])

	resultFile, err := os.OpenFile(resultFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	ignoreFile, err := os.OpenFile(ignoreFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	ignoreFileMap := make(map[string]bool)
	ignoreFileMap[resultFile.Name()] = true
	ignoreFileMap[ignoreFile.Name()] = true
	ignoreFileScanner := bufio.NewScanner(ignoreFile)
	for ignoreFileScanner.Scan() {
		ignoreFileName := ignoreFileScanner.Text()
		ignoreFileMap[ignoreFileName] = true
	}

	var walkFunc func(path string, info fs.FileInfo, err error) error
	if notAppendIgnoreFile {
		walkFunc = newWalkFunc(ignoreFileMap, pattern.FilePettern, nil, resultFile)
	} else {
		walkFunc = newWalkFunc(ignoreFileMap, pattern.FilePettern, ignoreFile, resultFile)
	}

	err = filepath.Walk(sourcePath, walkFunc)
	if err != nil {
		panic(err)
	}

}
