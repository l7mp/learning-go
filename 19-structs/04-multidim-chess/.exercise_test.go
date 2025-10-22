package multidimchess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeepCopyFunctions(t *testing.T) {
	{{if eq (index . "name") "32-dimensional" -}}
	origPos := Position{Coordinates: []int{3, 5}}
	newPos := copyPosition(origPos)
	assert.Equal(t, 32, len(newPos.Coordinates))
	assert.Equal(t, []int{3, 5}, newPos.Coordinates[:2])
	for _, v := range newPos.Coordinates[2:] {
		assert.Equal(t, 0, v)
	}
	newPos.Coordinates[0] = 99
	assert.NotEqual(t, origPos.Coordinates[0], newPos.Coordinates[0])

	p := Piece{Type: "Queen", Color: "White", Position: origPos}
	cp := copyPiece(p)
	assert.Equal(t, p.Type, cp.Type)
	assert.Equal(t, p.Color, cp.Color)
	cp.Position.Coordinates[0] = 99
	assert.NotEqual(t, p.Position.Coordinates[0], cp.Position.Coordinates[0])

	board := Board{
		Pieces: []Piece{
			{Type: "King", Color: "Black", Position: Position{Coordinates: []int{0, 0}}},
			{Type: "Pawn", Color: "White", Position: Position{Coordinates: []int{1, 2}}},
		},
	}
	cb := copyBoard(board)
	assert.Equal(t, len(board.Pieces), len(cb.Pieces))
	cb.Pieces[0].Position.Coordinates[0] = 99
	assert.NotEqual(t, board.Pieces[0].Position.Coordinates[0], cb.Pieces[0].Position.Coordinates[0])
	{{- end -}}

	{{if eq (index . "name") "2-dimensional" -}}
	origPos := Position{Coordinates: []int{1, 2, 3, 4, 5}}
	newPos := copyPosition(origPos)
	assert.Equal(t, 2, len(newPos.Coordinates))
	assert.Equal(t, []int{1, 2}, newPos.Coordinates)
	newPos.Coordinates[0] = 99
	assert.NotEqual(t, origPos.Coordinates[0], newPos.Coordinates[0])

	p := Piece{Type: "Bishop", Color: "Black", Position: origPos}
	cp := copyPiece(p)
	assert.Equal(t, p.Type, cp.Type)
	assert.Equal(t, p.Color, cp.Color)
	assert.Equal(t, 2, len(cp.Position.Coordinates))
    cp.Position.Coordinates[0] = 99
	assert.NotEqual(t, p.Position.Coordinates[0], cp.Position.Coordinates[0])

	board := Board{
		Pieces: []Piece{
			{Type: "Rook", Color: "White", Position: Position{Coordinates: []int{1, 2, 3}}},
			{Type: "Knight", Color: "Black", Position: Position{Coordinates: []int{4, 5, 6}}},
		},
	}
	cb := copyBoard(board)
	assert.Equal(t, len(board.Pieces), len(cb.Pieces))
	assert.Equal(t, 2, len(cb.Pieces[0].Position.Coordinates))
	cb.Pieces[0].Position.Coordinates[0] = 99
	assert.NotEqual(t, board.Pieces[0].Position.Coordinates[0], cb.Pieces[0].Position.Coordinates[0])
	{{- end -}}

	{{if eq (index . "name") "N-dimensional" -}}
	origPos := Position{Coordinates: []int{1, 2, 3, 4, 5}}
	newPos := copyPosition(origPos)
	assert.Equal(t, origPos.Coordinates, newPos.Coordinates)
	newPos.Coordinates[0] = 99
	assert.NotEqual(t, origPos.Coordinates[0], newPos.Coordinates[0])

	p := Piece{Type: "Queen", Color: "White", Position: origPos}
	cp := copyPiece(p)
	assert.Equal(t, p.Type, cp.Type)
	assert.Equal(t, p.Color, cp.Color)
	cp.Position.Coordinates[0] = 99
	assert.NotEqual(t, p.Position.Coordinates[0], cp.Position.Coordinates[0])

	board := Board{
		Pieces: []Piece{
			{Type: "King", Color: "White", Position: Position{Coordinates: []int{1, 1}}},
			{Type: "Queen", Color: "Black", Position: Position{Coordinates: []int{2, 2}}},
		},
	}
	cb := copyBoard(board)
	assert.Equal(t, len(board.Pieces), len(cb.Pieces))
	cb.Pieces[0].Position.Coordinates[0] = 99
	assert.NotEqual(t, board.Pieces[0].Position.Coordinates[0], cb.Pieces[0].Position.Coordinates[0])
	{{- end -}}
	{{- "\n"}}
}
