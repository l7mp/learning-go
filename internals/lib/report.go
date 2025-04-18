package lib

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type reportCard struct {
	StudentID      string                   `json:"student_id"`
	ExerciseScores map[string]exerciseStats `json:"exercise_scores"`
	Summary        string                   `json:"report"`
}
type exerciseStats struct {
	Score           float64         `json:"score"`
	TestCaseSuccess map[string]bool `json:"testCaseSuccess"`
}

// getTestCaseName takes a line and returns the name of the test case if the line is a go test output
func getTestCaseName(line string) string {
	re := regexp.MustCompile(`^=== RUN\s+(\w+)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 {
		testName := matches[1]
		return testName
	}
	return ""
}

// newPrefilledReport takes a slice of string and returns a reportCard that has keys matching the names of the tested exercises (e.g 01-hello-world)
func newPrefilledReport(id string, lines []string) reportCard {
	ret := reportCard{StudentID: id, ExerciseScores: make(map[string]exerciseStats)}
	re := regexp.MustCompile(`^ok\s+\S+/([^/\s]+)/([^/\s]+)\s+\d+\.\d+s`)
	for _, line := range lines {
		if match := re.FindStringSubmatch(line); match != nil {
			matchString := match[1] + "/" + match[2]
			ret.ExerciseScores[matchString] = exerciseStats{}
		}
	}
	return ret
}

func generateReportSummary(id string, failed int, passed int) string {
	if failed != 0 {
		return fmt.Sprintf("%s didn't pass all the homework tests, %d test cases have passed out of %d test cases", id, passed, passed+failed)
	} else {
		return fmt.Sprintf("%s did pass all the homework tests, %d/%d test cases passed", id, passed, passed)
	}
}

// extractExerciseScores takes in a string slice (test output line by line) and returns the exerciseStats representation alongside the total of failed and successful test cases
func extractExerciseScores(lines []string, verbose bool) (exStats []exerciseStats, failTotal int, succTotal int) {
	exStats = make([]exerciseStats, 0)
	exStat := exerciseStats{Score: 0.0, TestCaseSuccess: make(map[string]bool)}
	var testCaseName string
	failTotal = 0
	succTotal = 0
	var succ, fail float64
	if verbose {
		fmt.Println("Parsing test outputs")
	}
	for _, line := range lines {
		if strings.Contains(line, "RUN") {
			testCaseName = getTestCaseName(line)
		} else if strings.Contains(line, "PASS:") {
			succTotal++
			succ += 1.0
			exStat.TestCaseSuccess[testCaseName] = true
		} else if strings.Contains(line, "FAIL:") {
			failTotal++
			fail += 1.0
			exStat.TestCaseSuccess[testCaseName] = false
		} else if strings.Contains(line, "github.com/") {
			exStat.Score = succ / (fail + succ) * 100
			exStats = append(exStats, exStat)
			exStat = exerciseStats{Score: 0, TestCaseSuccess: make(map[string]bool)}
			fail, succ = 0.0, 0.0
		}
	}
	if verbose {
		fmt.Println("Finished parsing test outputs")
	}
	return
}

// CreateStudentReport takes the output of the tests and generates a report card for the student
func CreateStudentReport(id string, verbose bool) error {
	f, err := os.Open(ReportFile)
	if err != nil {
		return fmt.Errorf("error opening %q: %w", ReportFile, err)
	}
	if verbose {
		defer func() {
			if err != nil {
				fmt.Printf("the program exited with error %s", err.Error())
			} else {
				fmt.Printf("the program finished, successfully")
			}
		}()
	}
	var lines []string
	scanner := bufio.NewScanner(f)
	if verbose {
		fmt.Println("Reading test outputs...")
	}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if verbose {
		fmt.Println("Read test outputs")
	}
	report := newPrefilledReport(id, lines)
	exStats, failTotal, succTotal := extractExerciseScores(lines, verbose)
	index := 0
	for k, _ := range report.ExerciseScores {
		if math.IsNaN(exStats[index].Score) {
			delete(report.ExerciseScores, k)
			continue
		}
		report.ExerciseScores[k] = exStats[index]
		index++
	}
	report.Summary = generateReportSummary(id, failTotal, succTotal)
	outName := fmt.Sprintf("/%s-report.json", report.StudentID)
	output, err := os.Create(ReportOutDir + outName)
	if err != nil {
		return err
	}
	defer output.Close()
	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		return err
	}
	return nil
}

// createVisualReport takes a report card as argument and templates the html accordingly
func createVisualReport(rp []reportCard, file *os.File) error {
	temp, err := template.New("report_template.html").ParseFiles("internals/static/report_template.html")
	if err != nil {
		return err
	}
	err = temp.Execute(file, rp)
	return nil
}

// CreateTeacherReport templates the html file generating templated reports (uses reports defined in reportDir)
func CreateTeacherReport(reportDir string, verbose bool) error {
	dir, err := os.ReadDir(reportDir)
	if err != nil {
		return err
	}
	if verbose {
		fmt.Println("Reading student reports...")
	}
	rpCards := make([]reportCard, len(dir))
	for i, file := range dir {
		if !strings.Contains(file.Name(), "json") {
			continue
		}
		filePath := filepath.Join(reportDir, file.Name())
		report, err := os.Open(filePath)
		if err != nil {
			if verbose {
				fmt.Printf("error opening %q for reading: %s", file.Name(), err)
			}
			continue
		}
		var rep reportCard
		if err = json.NewDecoder(report).Decode(&rep); err != nil {
			fmt.Println(err.Error())
			if verbose {
				fmt.Printf("error encoding %q: %s", file.Name(), err)
			}
		}
		rpCards[i] = rep
		report.Close()
	}
	if verbose {
		fmt.Println("Finished reading student reports")
	}
	if err = createVisualReport(rpCards, os.Stdout); err != nil {
		if verbose {
			fmt.Printf("error creating visual report: %s", err)
		}
		return err
	}
	return nil
}
