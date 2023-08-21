// A generic package for implementing transaction logs.
package transactionlog

import "context"

// Event holds a record written to, or read from the transaction log.
type Event struct {
	// Sequence is a unique record ID, in monotonically increasing order.
	Sequence uint64
	// Type is the type of the event being recorded.
	Type string
	// Value is the new value put into the transaction log.
	Value string
}

// TransactionLogger is an interface for transaction implementations.
type TransactionLogger interface {
	// Write stores a key-value pair in the transaction log.
	Write(key, value string)
	// Run starts the event recording goroutine of the transaction log. Can be stopped with
	// cancelling the context.
	Run(context.Context)
	// ReadEvents synchronously replays all events stored in the transaction log or returns an
	// error.
	ReadEvents() ([]Event, error)
	// Err returns the channel on which the transaction log writer signals errors.
	Err() chan error
}
