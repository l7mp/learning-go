package basics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	g := newGame({{index .new "id"}}, "{{index .new "name"}}", {{index .new "price"}}, "{{index .new "genre"}}")

	assert.Equal(t, g.id, {{index . "new" "id"}}, "id")
	assert.Equal(t, g.name, "{{index . "new" "name"}}", "name")
	assert.Equal(t, g.price, {{index . "new" "price"}}, "price")
	assert.Equal(t, g.genre, "{{index . "new" "genre"}}", "genre")
}

func TestString(t *testing.T) {
	g := newGame({{index . "string" "id"}}, "{{index . "string" "name"}}", {{index . "string" "price"}}, "{{index . "string" "genre"}}")

	assert.Equal(t, g.item.String(), "{{index . "string" "item_result"}}", "item string")
	assert.Equal(t, g.String(), "{{index . "string" "game_result"}}", "game string")
}

func TestList(t *testing.T) {
	games := newGameList()

	assert.Len(t, games, 3, "gamelist len")

 	assert.Equal(t, games[0].id, {{index . "list" 0 "id"}}, "game 0 id")
	assert.Equal(t, games[0].name, "{{index . "list" 0 "name"}}", "game 0 name")
	assert.Equal(t, games[0].price, {{index . "list" 0 "price"}}, "game 0 price")
	assert.Equal(t, games[0].genre, "{{index . "list" 0 "genre"}}", "game 0 genre")

 	assert.Equal(t, games[1].id, {{index . "list" 1 "id"}}, "game 1 id")
	assert.Equal(t, games[1].name, "{{index . "list" 1 "name"}}", "game 1 name")
	assert.Equal(t, games[1].price, {{index . "list" 1 "price"}}, "game 1 price")
	assert.Equal(t, games[1].genre, "{{index . "list" 1 "genre"}}", "game 1 genre")

 	assert.Equal(t, games[2].id, {{index . "list" 2 "id"}}, "game 2 id")
	assert.Equal(t, games[2].name, "{{index . "list" 2 "name"}}", "game 2 name")
	assert.Equal(t, games[2].price, {{index . "list" 2 "price"}}, "game 2 price")
	assert.Equal(t, games[2].genre, "{{index . "list" 2 "genre"}}", "game 2 genre")
}

func TestById(t *testing.T) {
	g, err := queryById(newGameList(), {{index . "by_id" "id"}})
	assert.NoError(t, err, "no error")
 	assert.Equal(t, g.id, {{index . "by_id" "result" "id"}}, "game id")
	assert.Equal(t, g.name, "{{index . "by_id" "result" "name"}}", "game name")
	assert.Equal(t, g.price, {{index . "by_id" "result" "price"}}, "game price")
	assert.Equal(t, g.genre, "{{index . "by_id" "result" "genre"}}", "game genre")

	g, err = queryById(newGameList(), 11)
	assert.EqualError(t, err, "No such game", "error")
}

func TestNameByPrice(t *testing.T) {
	games := listNameByPrice(newGameList(), {{index . "name_by_price" "price"}})

 	assert.Equal(t, games, "{{index . "name_by_price" "result"}}", "names by prices")
}

