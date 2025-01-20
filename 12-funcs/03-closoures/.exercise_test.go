package closures

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProxy(t *testing.T) {
	{{ if eq (index . "name") "auth" }}
	authOne := proxy("admin","admin", func()string{return "Hello"})
	authTwo := proxy("admin","badmin", func()string{return "World!"})
	c1,cerr1 := authOne("admin","admin")
	c2,cerr2 := authTwo("admin","badmin")
	b1,berr1 := authOne("admin","badmin")
	b2, berr2 := authTwo("admin","admin")
	assert.Equal(t, "Hello", c1)
	assert.Equal(t, "World!", c2)
	assert.Nil(t, cerr1)
	assert.Nil(t, cerr2)
	assert.Equal(t, "", b1)
	assert.Equal(t, "", b2)
	assert.NotNil(t, berr1)
	assert.NotNil(t, berr2)
	{{ end }}
	{{ if eq (index . "name") "limiter" }}
	limitedOne := proxy(2, func() int { return 2 })
	limitedTwo := proxy(3, func() int { return 3 })
	for i := 0; i < 2; i++ {
		ret, err := limitedOne()
		assert.Nil(t, err)
		assert.Equal(t, 2, ret)
	}
	ret, err := limitedOne()
	assert.NotNil(t, err)
	assert.Equal(t, 0, ret)
	for i := 0; i < 3; i++ {
		ret, err = limitedTwo()
		assert.Nil(t, err)
		assert.Equal(t, 3, ret)
	}
	ret, err = limitedTwo()
	assert.NotNil(t, err)
	assert.Equal(t, 0, ret)
	{{ end }}
	{{ if eq (index . "name") "load-balancer" }}
	balancedOne := proxy(func()int{return 1},func() int { return 2 })
	balancedTwo := proxy(func()int{return 1},func() int { return 2 })
	assert.Equal(t, 1, balancedOne())
	assert.Equal(t, 1, balancedTwo())
	assert.Equal(t, 2, balancedOne())
	assert.Equal(t, 2, balancedTwo())
	{{ end }}
	{{ if eq (index . "name") "on-off" }}
	onOffOne := proxy(func(inp string) int{return len(inp)})
	onOffTwo := proxy(func(inp string) int{return len(inp)})
	one, err1 := onOffOne("alma")
	two, err2 := onOffTwo("ban")
	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Equal(t, 4, one)
	assert.Equal(t, 3, two)
	one, err1 = onOffOne("alma")
	two, err2 = onOffTwo("ban")
	assert.NotNil(t, err1)
	assert.NotNil(t, err2)
	assert.Equal(t, 0, one)
	assert.Equal(t, 0, two)
	{{ end }}
}