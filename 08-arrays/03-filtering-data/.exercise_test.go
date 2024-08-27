package filteringdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var keys []string
var indices []int

func TestFilterData(t *testing.T) {
	{{if eq (index . "name") "even"}}
	keys = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}
	indices = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	assert.Equal(t, [10]string{"b", "d", "f", "h", "j", "", "", "", "", ""}, filterData(keys, indices))

	keys = []string{"a", "b"}
	indices = []int{0, 1}
	assert.Equal(t, [10]string{"b", "", "", "", "", "", "", "", "", ""}, filterData(keys, indices))

	keys = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}
	indices = []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
       	assert.Equal(t, [10]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}, filterData(keys, indices))

	keys = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}
	indices = []int{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
	assert.Equal(t, [10]string{"", "", "", "", "", "", "", "", "", ""}, filterData(keys, indices))

	keys = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}
	indices = []int{1}
	assert.Equal(t, [10]string{"", "", "", "", "", "", "", "", "", ""}, filterData(keys, indices))
	{{end}}

	{{if eq (index . "name") "odd"}}
	keys = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}
	indices = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	assert.Equal(t, [10]string{"a", "c", "e", "g", "i", "k", "", "", "", ""}, filterData(keys, indices))

	keys = []string{"a", "b"}
	indices = []int{0, 1}
	assert.Equal(t, [10]string{"a", "", "", "", "", "", "", "", "", ""}, filterData(keys, indices))

	keys = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}
	indices = []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	assert.Equal(t, [10]string{"", "", "", "", "", "", "", "", "", ""}, filterData(keys, indices))

	keys = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}
	indices = []int{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
	assert.Equal(t, [10]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}, filterData(keys, indices))

	keys = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}
	indices = []int{1}
	assert.Equal(t, [10]string{"", "", "", "", "", "", "", "", "", ""}, filterData(keys, indices))

	{{end}}

	{{if eq (index . "name") "greater"}}
	keys = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}
	indices = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	assert.Equal(t, [10]string{"a", "b", "c", "d", "e", "f", "", "", "", ""}, filterData(keys, indices))

	keys = []string{"a", "b"}
	indices = []int{0, 1}
	assert.Equal(t, [10]string{"a", "b", "", "", "", "", "", "", "", ""}, filterData(keys, indices))

	keys = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}
	indices = []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	assert.Equal(t, [10]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}, filterData(keys, indices))

	keys = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}
	indices = []int{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
	assert.Equal(t, [10]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}, filterData(keys, indices))

	keys = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}
	indices = []int{1}
	assert.Equal(t, [10]string{"", "", "", "", "", "", "", "", "", ""}, filterData(keys, indices))
	{{end}}

	{{if eq (index . "name") "smaller"}}
	keys = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}
	indices = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	assert.Equal(t, [10]string{"f", "g", "h", "i", "j", "k", "", "", "", ""}, filterData(keys, indices))

	keys = []string{"a", "b"}
	indices = []int{0, 1}
	assert.Equal(t, [10]string{"", "", "", "", "", "", "", "", "", ""}, filterData(keys, indices))

	keys = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}
	indices = []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	assert.Equal(t, [10]string{"", "", "", "", "", "", "", "", "", ""}, filterData(keys, indices))

	keys = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}
	indices = []int{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
	assert.Equal(t, [10]string{"", "", "", "", "", "", "", "", "", ""}, filterData(keys, indices))

	keys = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}
	indices = []int{1}
	assert.Equal(t, [10]string{"", "", "", "", "", "", "", "", "", ""}, filterData(keys, indices))
	{{end}}
}
