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
	24-generics \
	25-testing

# LAB subdirs
LAB_DIRS=\
	99-labs

STUDENT_ID_FILE=STUDENT_ID

# the exercises: directories whose README.md and exercise_test.go are generated from the
# templates by hashing the student id. "make generate" writes them and "make clean"
# removes them, so the two stay symmetric. Add new exercises here.
EXERCISES=\
	01-getting-started/01-hello-world \
	02-types/01-booleans \
	02-types/02-numbers \
	02-types/03-strings \
	02-types/04-printf \
	03-variables/01-repaint \
	03-variables/02-path-split \
	04-type-system/01-construct-duration \
	04-type-system/02-secret-protocol-header \
	04-type-system/03-read-secret-register \
	06-switch/01-richter-scale \
	06-switch/02-grades \
	07-loops/01-factorial-sum-abs \
	07-loops/02-digits \
	08-arrays/01-message_queue \
	08-arrays/02-n-arithmetic \
	08-arrays/03-filtering-data \
	12-funcs/01-fibonacci \
	12-funcs/02-callbacks \
	12-funcs/03-closures \
	15-pointers/01-basic \
	15-pointers/02-new \
	19-structs/01-basics \
	19-structs/02-interfaces-with-structs \
	19-structs/03-struct-embedding \
	19-structs/04-multidim-chess \
	20-goroutines/01-concurrent-primes \
	20-goroutines/02-concurrent-word-count \
	20-goroutines/03-sleep-sort \
	21-channels/01-pipeline \
	21-channels/02-channel-multiplex \
	21-channels/03-channel-broadcast \
	22-context/02-nested \
	22-context/05-threadpool \
	23-misc/01-scanning \
	23-misc/02-map-as-sets \
	24-generics/01-sorting \
	24-generics/02-functional \
	25-testing/01-testify


# files to be generated
.PHONY: init check generate test ci-test lab-test clean realclean

# check if STUDENT_ID is set
check:
	export STUDENT_ID=$(STUDENT_ID)
	@if [ ! -s "$(STUDENT_ID_FILE)" -o ! -r "$(STUDENT_ID_FILE)" ]; then \
		echo "ERROR: '$(STUDENT_ID_FILE)' is not readable or empty"; exit 1; \
	fi

# init the placeholder files
init:
	@for dir in $(EXERCISES); do \
		echo '# PLEASE RUN make generate' > $$dir/README.md; \
		echo '// PLEASE RUN make generate' > $$dir/exercise_test.go; \
	done

# generate the README and the test for each exercise
generate:
	@for dir in $(EXERCISES); do \
		(cd $$dir && go run $(CURDIR)/exercises-cli.go generate) || exit 1; \
	done

# run the tests
test:
	export STUDENT_ID=$(STUDENT_ID)
	go test ./... -v -count 1

# LAB narrows the lab tests to a single lab, e.g. "LAB=02 make lab-test" while you are
# working through lab 02. Only the applications you have already started are built, so the
# later labs do not have to be solved yet. Labs 02 to 06 have tests.
LAB ?=

LAB_RUN=$(if $(LAB),-run TestLab$(LAB))

# run the exercise tests plus the lab tests. The lab tests build the container images from
# the code in 99-labs/code, bring up a Kubernetes cluster and deploy into it, so they need
# a container engine; see 99-labs/ci/README.md.
ci-test: test lab-test

# run only the lab tests
lab-test:
	@if ! docker info >/dev/null 2>&1; then \
		echo "ERROR: no usable container engine."; \
		echo "       The lab tests build images and run Kubernetes in containers."; \
		echo "       See 99-labs/ci/README.md."; \
		exit 1; \
	fi
	cd 99-labs/ci && LABS_SOURCE=$(CURDIR) go test ./ $(LAB_RUN) -v -count 1 -timeout 60m

# clean up generated files
clean:
	@for dir in $(EXERCISES); do \
		rm -fv $$dir/README.md $$dir/exercise_test.go; \
	done

# also wipe student id
realclean: clean
	echo "PLEASE SET STUDENT ID" > $(STUDENT_ID_FILE)
