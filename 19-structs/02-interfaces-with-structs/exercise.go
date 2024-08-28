package structsinterfaces

// DO NOT REMOVE THIS COMMENT
//go:generate go run ../../exercises-cli.go -student-id=$STUDENT_ID generate



// {{index . "interface" "name"}} has two functions: {{index . "interface" "func1" "name"}} and {{index . "interface" "func2" "name"}}
// INSERT YOUR CODE HERE

// INSERT YOUR STRUCTS HERE

// New{{index . "struct1"}} returns a new {{index . "struct1"}} struct.
func New{{index . "struct1" "name"}}({{index . "struct1" "new"}}) {{index . "struct1" "name"}} {
	// INSERT YOUR CODE HERE
}

// New{{index . "struct2"}} returns a new {{index . "struct2"}} struct.
func New{{index . "struct2" "name"}}({{index . "struct2" "new"}}) {{index . "struct2" "name"}} {
	// INSERT YOUR CODE HERE
}

// IMPLEMENT INTERFACE FUNCTIONS FOR YOUR STRUCTS HERE
