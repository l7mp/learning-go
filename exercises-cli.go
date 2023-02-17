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
	fmt.Fprintf(os.Stderr, "exercises-cli <generate|check>\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

var (
	flagId  = flag.String("student-id", "", "student id; optional, default is to read from file STUDENT_ID")
	verbose = flag.Bool("verbose", false, "verbose mode; default is off")
)

func main() {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetFlags(0)
	log.SetPrefix("exercises-cli: ")
	flag.Usage = Usage
	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	id, err := lib.GetStudentId(flagId)
	if err != nil {
		log.Fatalf("Cannot find student id: %s", err)
	}

	if *verbose {
		log.Printf("Using student id %q\n", id)
	}

	switch flag.Arg(0) {
	case "generate":
		if err := lib.Generate(id, *verbose); err != nil {
			log.Fatalf("Cannot generate exercise for student id %s: %s", id, err)
		}
	case "check":
		fmt.Println("Would check stuff")
	default:
		flag.Usage()
		os.Exit(1)
	}
}
