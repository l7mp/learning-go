# exercise subdirs
EXERCISE_DIRS=\
	01-getting-started \
	02-types \
	03-variables \
	04-type-system \
	06-switch \
	08-arrays \
	19-structs \
	22-goroutines

STUDENT_ID_FILE=STUDENT_ID

# files to be generated
.PHONY: check generate test clean realclean

# check if STUDENT_ID is set
check:
	export STUDENT_ID=$(STUDENT_ID)
	@if [ ! -s "$(STUDENT_ID_FILE)" -o ! -r "$(STUDENT_ID_FILE)" ]; then \
		echo "ERROR: '$(STUDENT_ID_FILE)' is not readable or has zero content"; exit 1; \
	fi

# run go generate
generate:
	export STUDENT_ID=$(STUDENT_ID)
	go generate ./...

# run the tests
test:
	export STUDENT_ID=$(STUDENT_ID)
	go test ./... -v

# clean up generated files
clean:
	for dir in $(EXERCISE_DIRS); do \
		find $$dir -name "README.md" -type f -print | xargs rm -f;  \
		find $$dir -name "exercise_test.go" -type f -print | xargs rm -f;  \
	done

# also wipe student id
realclean: clean
	echo "PLEASE SET STUDENT ID" > $(STUDENT_ID_FILE)
