package structembedding

// DO NOT REMOVE THIS COMMENT
//go:generate go run ../../exercises-cli.go -student-id=$STUDENT_ID generate

// INSERT YOUR CODE HERE

// Author represents information about the book's author
type Author struct {
	// TODO: Define the Author struct fields
}

// Book represents information about a book
type Book struct {
	// TODO: Define the Book struct fields, embedding the Author struct
}

// Article represents information about a article
type Article struct {
	// TODO: Define the Article struct fields, embedding the Author struct
}

// ParseBook parses the given JSON data into a Book struct
func ParseBook(jsonData []byte) (Book, error)

// ParseArticle parses the given JSON data into a Article struct
func ParseArticle(jsonData []byte) (Article, error)
