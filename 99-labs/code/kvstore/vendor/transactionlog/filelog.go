package transactionlog

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
)

// FileTransactionLogger is an implementation of the TransactionLogger interface that stores the
// log records in a file.
type FileTransactionLogger struct {
	// events is a channel for sending events to the log.
	events chan Event
	// errors is a channel for receiving errors from the log.
	errors chan error
	// lastSequence is the last used event sequence number.
	lastSequence uint64
	// file is the file backing the transaction log.
	file *os.File
}

// NewFileTransactionLogger returns a new transaction logger implementation.
func NewFileTransactionLogger(filename string) (TransactionLogger, error) {
	log.Printf("Opening/creating transaction log in file %s", filename)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	return &FileTransactionLogger{file: file}, nil
}

// Write stores a key-value pair in the transaction log.
func (l *FileTransactionLogger) Write(key, value string) {
	log.Printf("Writing record: key=%s, value=%s", key, value)
	l.events <- Event{Type: key, Value: value}
}

// Err returns the channel on which the transaction log writer signals errors.
func (l *FileTransactionLogger) Err() chan error {
	return l.errors
}

// Run starts the event recording goroutine of the transaction log.
func (l *FileTransactionLogger) Run(ctx context.Context) {
	l.events = make(chan Event, 16) // bidirectional channel
	l.errors = make(chan error, 1)  // bidirectional channel

	go func() { // Writes are served in the background
		defer close(l.events) // Close channels when the goroutine ends
		defer close(l.errors)

		for {
			select {
			case e := <-l.events: // Retrieve the next Event
				l.lastSequence++       // Increment sequence number
				_, err := fmt.Fprintf( // Write the event to the log
					l.file, "%d\t%s\t%s\n",
					l.lastSequence, e.Type, e.Value)
				if err != nil {
					l.errors <- err // Write error to the err channel
					return          // Stop the goroutine on error
				}
			case <-ctx.Done():
				log.Printf("Closing transaction log")
				l.file.Close()
				return
			}
		}
	}()
}

// ReadEvents synchronously replays all events stored in the transaction log or returns an error.
func (l *FileTransactionLogger) ReadEvents() ([]Event, error) {
	scanner := bufio.NewScanner(l.file) // Create a Scanner for l.file
	events := []Event{}                 // Create empty event list.

	for scanner.Scan() {
		line := scanner.Text()
		var e Event

		_, err := fmt.Sscanf(line, "%d\t%s\t%s", &e.Sequence, &e.Type, &e.Value)
		if err != nil {
			return []Event{}, fmt.Errorf("input parse error: %w", err)
		}
		// Sanity check! Are the sequence numbers in increasing order?
		if l.lastSequence >= e.Sequence {
			return []Event{}, fmt.Errorf("transaction numbers out of sequence")
		}
		l.lastSequence = e.Sequence
		events = append(events, e)

		log.Printf("Read-events: New event %#v", e)

	}

	if err := scanner.Err(); err != nil {
		return []Event{}, fmt.Errorf("transaction log read failure: %w", err)
	}

	return events, nil
}
