//go:build ignore

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/l7mp/learning-go/internals/lib"
)

func Usage() {
	fmt.Fprintf(os.Stderr, "exercises-cli <generate|check|report|html>\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

var (
	studentId = flag.String("student-id", "", "student id; optional, default is to read from file STUDENT_ID")
	verbose   = flag.Bool("verbose", false, "verbose mode; default is off")
	reportDir = flag.String("report-dir", "", "path to directory where the report files are stored")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// log.SetFlags(0)
	log.SetPrefix("exercises-cli: ")
	flag.Usage = Usage
	flag.Parse()

	//if len(flag.Args()) != 1 {
	//	flag.Usage()
	//	os.Exit(1)
	//}

	id, err := lib.GetStudentId(studentId)
	if err != nil {
		log.Fatalf("Cannot find student id: %s", err)
	}

	if *verbose {
		log.Printf("Using student id %q\n", id)
	}

	switch flag.Arg(0) {
	case "generate":
		if err := lib.Generate(id, *verbose); err != nil {
			log.Fatalf("Cannot generate exercise for student id %q: %s", id, err)
		}
	case "check":
		fmt.Println("Would check stuff")
	case "report":
		if err := lib.CreateStudentReport(id, *verbose); err != nil {
			log.Fatalf("Cannot create report for student id %q: %s", id, err)
		}
	case "html":
		fmt.Println(reportDir)
		if err := lib.CreateTeacherReport(*reportDir, *verbose); err != nil {
			log.Fatalf("Cannot create teacher report %s", err)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
}
