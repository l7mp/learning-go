package basics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	g := newGame({{index .new "id"}}, "{{index .new "name"}}", {{index .new "price"}}, "{{index .new "genre"}}")

	assert.Equal(t, {{index . "new" "id"}},      g.id, "id")
	assert.Equal(t, "{{index . "new" "name"}}",  g.name, "name")
	assert.Equal(t, {{index . "new" "price"}},   g.price, "price")
	assert.Equal(t, "{{index . "new" "genre"}}", g.genre, "genre")
}

func TestString(t *testing.T) {
	g := newGame({{index . "string" "id"}}, "{{index . "string" "name"}}", {{index . "string" "price"}}, "{{index . "string" "genre"}}")

	assert.Equal(t, "{{index . "string" "item_result"}}", g.item.String(), "item string")
	assert.Equal(t, "{{index . "string" "game_result"}}", g.String(), "game string")
}

func TestList(t *testing.T) {
	games := newGameList()

	assert.Len(t, games, 3, "gamelist len")

 	assert.Equal(t, {{index . "list" 0 "id"}},      games[0].id, "game 0 id")
	assert.Equal(t, "{{index . "list" 0 "name"}}",  games[0].name, "game 0 name")
	assert.Equal(t, {{index . "list" 0 "price"}},   games[0].price, "game 0 price")
	assert.Equal(t, "{{index . "list" 0 "genre"}}", games[0].genre, "game 0 genre")

 	assert.Equal(t, {{index . "list" 1 "id"}},      games[1].id, "game 1 id")
	assert.Equal(t, "{{index . "list" 1 "name"}}",  games[1].name, "game 1 name")
	assert.Equal(t, {{index . "list" 1 "price"}},   games[1].price, "game 1 price")
	assert.Equal(t, "{{index . "list" 1 "genre"}}", games[1].genre, "game 1 genre")

 	assert.Equal(t, {{index . "list" 2 "id"}},      games[2].id, "game 2 id")
	assert.Equal(t, "{{index . "list" 2 "name"}}",  games[2].name, "game 2 name")
	assert.Equal(t, {{index . "list" 2 "price"}},   games[2].price, "game 2 price")
	assert.Equal(t, "{{index . "list" 2 "genre"}}", games[2].genre, "game 2 genre")
}

func TestById(t *testing.T) {
	g, err := queryById(newGameList(), {{index . "by_id" "id"}})
	assert.NoError(t, err, "no error")
 	assert.Equal(t, {{index . "by_id" "result" "id"}},      g.id, "game id")
	assert.Equal(t, "{{index . "by_id" "result" "name"}}",  g.name, "game name")
	assert.Equal(t, {{index . "by_id" "result" "price"}},   g.price, "game price")
	assert.Equal(t, "{{index . "by_id" "result" "genre"}}", g.genre, "game genre")

	g, err = queryById(newGameList(), 11)
	assert.EqualError(t, err, "no such game", "error")
}

func TestNameByPrice(t *testing.T) {
	names := listNameByPrice(newGameList(), {{index . "name_by_price" "price"}})

	assert.Len(t, names, {{index .name_by_price "result" "size"}}, "namelist len")
	if {{index .name_by_price "result" "size"}} > 0 {
	 	assert.Equal(t, "{{index .name_by_price "result" "list" 0}}", names[0], "names 1")
	}
	if {{index .name_by_price "result" "size"}} > 1 {
	 	assert.Equal(t, "{{index .name_by_price "result" "list" 1}}", names[1], "names 2")
	}
	if {{index .name_by_price "result" "size"}} > 2 {
	 	assert.Equal(t, "{{index .name_by_price "result" "list" 2}}", names[2], "names 3")
	}
	if {{index .name_by_price "result" "size"}} > 3 {
	 	assert.Equal(t, "{{index .name_by_price "result" "list" 3}}", names[3], "names 4")
	}
}

