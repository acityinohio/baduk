package baduk

import "errors"

//Sets a Piece to empty on the Board
//x, y in range from 1 to Board.Size
func (b *Board) setE(x, y int) (err error) {
	if err = checkRange(x, y, b.Size); err != nil {
		return err
	}
	b.Grid[y][x].Black = false
	b.Grid[y][x].White = false
	b.Grid[y][x].Empty = true
	return
}

//Returns true if piece has liberties
func (p *Piece) hasLiberty() bool {
	if p.Up != nil && p.Up.Empty {
		return true
	}
	if p.Down != nil && p.Down.Empty {
		return true
	}
	if p.Left != nil && p.Left.Empty {
		return true
	}
	if p.Right != nil && p.Right.Empty {
		return true
	}
	return false
}

//Checks x,y against size
func checkRange(x, y, size int) error {
	switch {
	case x < 0 || x >= size:
		return errors.New("x out of range")
	case y < 0 || y >= size:
		return errors.New("y out of range")
	default:
		return nil
	}
}
