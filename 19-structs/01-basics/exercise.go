package warmup

//go:generate go run ../../exercises-cli.go -student-id=$STUDENT_ID generate

// INSERT YOUR STRUCTS HERE

// newGame returns a new game struct.
func newGame(id int, name string, price int, genre string) game {
	// INSERT YOUR CODE HERE
}

// String stringifies an item.
func (i item) String() string {
	// INSERT YOUR CODE HERE
}

// String stringifies a game.
func (g game) String() string {
	// INSERT YOUR CODE HERE
}

// newGameList creates a game store.
func newGameList() []game {
	// INSERT YOUR CODE HERE
}

// queryById returns the game in the specified store with the given id or returns a "No such game" error.
func queryById(games []game, id int) (game, error) {
	// INSERT YOUR CODE HERE
}

// listNameByPrice returns the name of the game(s) with price equal or smaller than a given price.
func listNameByPrice(games []game, price int) []string {
	// INSERT YOUR CODE HERE
}
