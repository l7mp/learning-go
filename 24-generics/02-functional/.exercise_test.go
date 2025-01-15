package functional

import (
{{ if or (eq .name "filter") (eq .name "reduce") }}
	"strconv"
{{ end }}
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGenericFunction(t *testing.T) {
	{{ if eq (index . "name") "fmap"}}
	valuesInt := []int{2,5,8,88,12}
	valuesStr := []string{"alm","bel","but","balg"}
	mappedInt := fmap(valuesInt, func(i int) int {return 3 * i - 1})
	mappedStr := fmap(valuesStr, func(s string) string { return s + "a" })
	assert.Equal(t, []int{5,14,23,263,35},mappedInt)
	assert.Equal(t,[]string{"alma","bela","buta","balga"},mappedStr)
	{{end}}
	{{ if eq (index . "name") "reduce"}}
	valuesInt := []int{2,5,8,88,12}
	valuesStr := []string{"22","3","4","56"}
	reducedInt := reduce(valuesInt, 1, func(a ,b int) int {return a * b})
	reducedStr := reduce(valuesStr, 0,func(a int, b string) int {
		parsed, _ := strconv.Atoi(b)
		return parsed + a
	})
	assert.Equal(t, 84480,reducedInt)
	assert.Equal(t, 85,reducedStr)
	{{end}}
	{{ if eq (index . "name") "filter"}}
	valuesInt := []int{2,5,8,88,12}
	valuesStr := []string{"22","3","4","56","alma","korte","banan"}
	filteredInt := filter(valuesInt, func(i int) bool { return i%2 == 0 })
	filteredStr := filter(valuesStr, func(s string) bool {
		_,err := strconv.Atoi(s)
		return err != nil
	})
	assert.Equal(t, []int{2,8,88,12},filteredInt)
	assert.Equal(t, []string{"alma","korte","banan"},filteredStr)
	{{end}}
	{{ if eq (index . "name") "flatten"}}
	valuesInt := [][]int{[]int{2,5,8,88,12},[]int{3,4}}
	valuesStr := [][]string{[]string{"Boldog","karacsonyt"},[]string{"te","mocskos","allat"}}
	assert.Equal(t, []int{2,5,8,88,12,3,4},flatten(valuesInt))
	assert.Equal(t, []string{"Boldog","karacsonyt","te","mocskos","allat"},flatten(valuesStr))
	{{end}}
}
