package transactionlog

import (
	"context"
	"os"
	"testing"
	"time"
)

const logFile = "/tmp/translog_test.log"

func TestFileTransactionLogger(t *testing.T) {
	// make sure we start with an empty logfile
	file, err := os.OpenFile(logFile, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)
	file.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger1, err := NewFileTransactionLogger(logFile)
	if err != nil {
		t.Errorf("NewFileTransactionLogger 1 error: %s", err.Error())
	}

	logger1.Run(ctx)

	// record errors
	go func() {
		errCh := logger1.Err()
		for {
			select {
			case err := <-errCh:
				t.Errorf("runner error signal: %s", err.Error())
			case <-ctx.Done():
				return
			}
		}
	}()

	logger1.Write("a", "1")
	logger1.Write("b", "2")
	logger1.Write("a", "3")

	// sleep a bit so that the logger goroutine has enough time to pick up the write events
	time.Sleep(100 * time.Millisecond)

	// create another logger to read back the log
	logger2, err := NewFileTransactionLogger(logFile)
	if err != nil {
		t.Errorf("NewFileTransactionLogger 2 error: %s", err.Error())
	}

	events, err := logger2.ReadEvents()
	if err != nil {
		t.Errorf("read-events: error %s", err.Error())
	}

	if len(events) != 3 {
		t.Errorf("read-events: wrong number of records %d (expected 3)", len(events))
	}

	if events[0].Type != "a" || events[0].Value != "1" {
		t.Errorf("read-events: error on record 1: key: %s, value: %s", events[0].Type, events[0].Value)
	}

	if events[1].Type != "b" || events[1].Value != "2" {
		t.Errorf("read-events: error on record 2: key: %s, value: %s", events[1].Type, events[1].Value)
	}

	if events[2].Type != "a" || events[2].Value != "3" {
		t.Errorf("read-events: error on record 3: key: %s, value: %s", events[2].Type, events[2].Value)
	}

	// write a new record
	logger1.Write("c", "4")

	// sleep a bit so that the logger goroutine has enough time to pick up the write events
	time.Sleep(100 * time.Millisecond)

	// and read again
	logger3, err := NewFileTransactionLogger(logFile)
	if err != nil {
		t.Errorf("NewFileTransactionLogger 3 error: %s", err.Error())
	}

	events, err = logger3.ReadEvents()
	if err != nil {
		t.Errorf("read-events: error %s", err.Error())
	}

	if len(events) != 4 {
		t.Errorf("read-events: wrong number of records %d (expected 4)", len(events))
	}

	if events[0].Type != "a" || events[0].Value != "1" {
		t.Errorf("read-events: error on record 1: key: %s, value: %s", events[0].Type, events[0].Value)
	}

	if events[1].Type != "b" || events[1].Value != "2" {
		t.Errorf("read-events: error on record 2: key: %s, value: %s", events[1].Type, events[1].Value)
	}

	if events[2].Type != "a" || events[2].Value != "3" {
		t.Errorf("read-events: error on record 3: key: %s, value: %s", events[2].Type, events[2].Value)
	}

	if events[3].Type != "c" || events[3].Value != "4" {
		t.Errorf("read-events: error on record 4: key: %s, value: %s", events[3].Type, events[3].Value)
	}
}
