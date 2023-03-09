package exam

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Result struct {
	StudentNumber int     `csv:"studentNumber"`
	SubjectName   string  `csv:"subjectName"`
	Times         int     `csv:"times"`
	QuestionCount int     `csv:"questionCount"`
	CorrectCount  int     `csv:"correctCount"`
	CorrectRate   float64 `csv:"correctRate"`
}

var (
	fileNamePattern *regexp.Regexp = regexp.MustCompile(`^[^\s|_]+_\d+`)
)

func ParseFileName(fileName string) (string, int, error) {
	fileNameTokens := strings.Split(fileNamePattern.FindString(fileName), "_")
	if len(fileNameTokens) != 2 {
		return "", 0, errors.New("failed to parse " + fileName)
	}

	times, err := strconv.Atoi(fileNameTokens[1])
	if err != nil {
		return "", 0, fmt.Errorf("failed to parse times: %w", err)
	}

	return fileNameTokens[0], times, nil
}

func ParseAnswers(resultFile *os.File) ([]Result, error) {
	resultFileInfo, err := resultFile.Stat()
	if err != nil {
		return nil, err
	}

	subjectName, times, err := ParseFileName(resultFileInfo.Name())
	if err != nil {
		return nil, err
	}

	// ヘッダ行を読み飛ばす
	resultReader := csv.NewReader(resultFile)
	_, err = resultReader.Read()
	if err != nil {
		return nil, err
	}

	var results []Result
	for {
		row, err := resultReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		studentNumber, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, err
		}

		var questionCount, correctCount int
		for _, rawAnswer := range row[1:] {
			correct, err := strconv.ParseBool(rawAnswer)
			if err != nil {
				return nil, err
			}

			questionCount++
			if correct {
				correctCount++
			}
		}

		correctRate := float64(correctCount) / float64(questionCount)

		results = append(results, Result{
			StudentNumber: studentNumber,
			SubjectName:   subjectName,
			Times:         times,
			QuestionCount: questionCount,
			CorrectCount:  correctCount,
			CorrectRate:   correctRate,
		})
	}

	return results, nil
}
