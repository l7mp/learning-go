package sorting

import (
	"math"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSorting(t *testing.T) {
	{{ if eq (index . "name") "string-float64"}}
	arrTypeOne := []string{"banan","alma","64","cseresznye","jojo","toki","wa","ugoki","desu"}
	arrTypeTwo := []float64{1.1,-0.1,66,-22,0,math.Pi,6.2834, math.E}
	sortedTwo = sortSlice(arrTypeTwo)
	assert.Equal(t,[]string{"64","alma","banan","cseresznye","desu","jojo","toki","ugoki","wa"}, sortSlice(arrTypeOne))
	assert.Equal(t, []float64{-22,-0.1,0,1.1, math.E, math.Pi,6.2834,66}, sortSlice(arrTypeTwo))
	{{end}}
	{{ if eq (index . "name") "int64-float64"}}
	arrTypeOne := []int64{3,-1,22,-11,8080,5432,123,53,179,-66}
	arrTypeTwo := []float64{1.1,-0.1,66,-22,0,math.Pi,6.2834, math.E}
	assert.Equal(t, []int64{-66,-11,-1,3,22,53,123,179,5432,8080},sortSlice(arrTypeOne))
	assert.Equal(t, []float64{-22,-0.1,0,1.1, math.E, math.Pi,6.2834,66}, sortSlice(arrTypeTwo))
	{{end}}
	{{ if eq (index . "name") "uint8-int32"}}
	arrTypeOne := []uint8{1,2,7,5,4,6,3,9,8}
	arrTypeTwo := []int32{3,-1,22,-11,8080,5432,123,53,179,-66}
	assert.Equal(t, []uint8{1,2,3,4,5,6,7,8,9}, sortSlice(arrTypeOne))
	assert.Equal(t, []int32{-66,-11,-1,3,22,53,123,179,5432,8080},sortSlice(arrTypeTwo))
	{{end}}
	{{ if eq (index . "name") "string-uint8"}}
	arrTypeOne := []string{"banan","alma","64","cseresznye","jojo","toki","wa","ugoki","desu"}
	arrTypeTwo := []uint8{1,2,7,5,4,6,3,9,8}
	assert.Equal(t,[]string{"64","alma","banan","cseresznye","desu","jojo","toki","ugoki","wa"}, sortSlice(arrTypeOne))
	assert.Equal(t, []uint8{1,2,3,4,5,6,7,8,9}, sortSlice(arrTypeTwo))
	{{end}}
	{{ if eq (index . "name") "int32-uint64"}}
	arrTypeOne := []int32{3,-1,22,-11,8080,5432,123,53,179,-66}
	arrTypeTwo := []uint64{11111,12,57,25,94,66,31,19,8888}
	assert.Equal(t, []int32{-66,-11,-1,3,22,53,123,179,5432,8080},sortSlice(arrTypeOne))
	assert.Equal(t, []uint64{12,19,25,31,57,66,94,8888,11111},sortSlice(arrTypeTwo))
	{{end}}
}
