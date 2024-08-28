package structsinterfaces

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func test{{index . "interface" "func1" "name"}}(e {{index . "interface" "name"}}) {{index . "interface" "func1" "retval"}} {
	return e.{{index . "interface" "func1" "name"}}()
}

func test{{index . "interface" "func2" "name"}}(e {{index . "interface" "name"}}) {{index . "interface" "func2" "retval"}} {
	return e.{{index . "interface" "func2" "name"}}()
}

func TestStructsInterfaces(t *testing.T) {
	s1 := New{{index . "struct1" "name"}}({{index . "struct1" "params"}})

	assert.Equal(t, s1.{{index . "interface" "func1" "name"}}(), test{{index . "interface" "func1" "name"}}(s1))
	assert.Equal(t, s1.{{index . "interface" "func2" "name"}}(), test{{index . "interface" "func2" "name"}}(s1))

	s2 := New{{index . "struct2" "name"}}({{index . "struct2" "params"}})
	assert.Equal(t, s2.{{index . "interface" "func1" "name"}}(), test{{index . "interface" "func1" "name"}}(s2))
	assert.Equal(t, s2.{{index . "interface" "func2" "name"}}(), test{{index . "interface" "func2" "name"}}(s2))
}
