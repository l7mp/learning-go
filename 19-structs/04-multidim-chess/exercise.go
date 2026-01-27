package multidimchess

// DO NOT REMOVE THIS COMMENT
//go:generate go run ../../exercises-cli.go -student-id=$STUDENT_ID generate

type Position struct {
	Coordinates []int
}

type Piece struct {
	Type     string // e.g. Queen, King, Bishop, etc.
	Color    string // White or Black
	Position Position
}

type Board struct {
	Pieces []Piece
}

// INSERT YOUR CODE HERE
