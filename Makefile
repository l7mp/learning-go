# exercise subdirs
EXERCISE_DIRS=\
	01-getting-started \
	02-types \
	03-variables \
	04-type-system \
	06-switch \
	07-loops \
	08-arrays \
	12-funcs \
	15-pointers \
	19-structs \
	20-goroutines \
	21-channels \
	22-context \
	23-misc \
	24-generics

# LAB subdirs
LAB_DIRS=\
	99-labs

STUDENT_ID_FILE=STUDENT_ID

# files to be generated
.PHONY: init check generate test clean realclean

# check if STUDENT_ID is set
check:
	export STUDENT_ID=$(STUDENT_ID)
	@if [ ! -s "$(STUDENT_ID_FILE)" -o ! -r "$(STUDENT_ID_FILE)" ]; then \
		echo "ERROR: '$(STUDENT_ID_FILE)' is not readable or empty"; exit 1; \
	fi

# init the placeholder files
init:
	@for dir in $(EXERCISE_DIRS); do \
		find $$dir/* -type d -exec sh -c 'echo \# PLEASE RUN make generate > "$$0/README.md"' {} \; ; \
		find $$dir/* -type d -exec sh -c 'echo // PLEASE RUN make generate > "$$0/exercise_test.go"' {} \; ; \
	done

# run go generate
generate:
	export STUDENT_ID=$(STUDENT_ID)
	go generate ./...

# run the tests
test:
	export STUDENT_ID=$(STUDENT_ID)
	go test ./... -v -count 1

# generate reports
report:
	export STUDENT_ID=$(STUDENT_ID)
	if  [ ! -f report_tmp ]; then \
		go test ./... -v -count 1 -parallel 1 > report_tmp; \
  	fi
	go run exercises-cli.go -verbose report

# generate the hmtl report
report-html:
	export STUDENT_ID=$(STUDENT_ID)
	go run exercises-cli.go -report-dir=$(REPORT_DIR) -verbose html > report.html

# clean up generated files
clean:
	for dir in $(EXERCISE_DIRS); do \
		find $$dir -name "README.md" -type f -print | xargs rm -f;  \
		find $$dir -name "exercise_test.go" -type f -print | xargs rm -f;  \
	done

# also wipe student id
realclean: clean
	echo "PLEASE SET STUDENT ID" > $(STUDENT_ID_FILE)

