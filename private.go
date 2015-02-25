package baduk

import "errors"

//Sets a Piece to either black or white
func (b *Board) set(x int, y int, isBlack bool) (err error) {
	if err = b.checkRange(x, y); err != nil {
		return err
	}
	if !b.Grid[y][x].Empty {
		err = errors.New("Piece is not empty")
		return
	}
	b.Grid[y][x].Black = isBlack
	b.Grid[y][x].White = !isBlack
	b.Grid[y][x].Empty = false
	//Check if setting a piece captures opponent's adjacent chains
	err = b.Grid[y][x].checkCapture(isBlack)
	if err != nil {
		//Reset to empty if there's an error
		b.setE(x, y)
		return
	}
	//If there are no liberties after checking opponent's chains,
	//Check whether opponent captures chains connected to this move
	if !b.Grid[y][x].hasLiberty() {
		err = b.Grid[y][x].checkCapture(!isBlack)
		if err != nil {
			b.setE(x, y)
			return
		}
	}
	return
}

//Checks surrounding stones to see if
//chains can be captured, empties
//chains without liberties
func (p *Piece) checkCapture(isBlack bool) (err error) {
	//Check chain liberties for each direction
	//Only check if chains are opposite of isBlack/not empty
	return
}

//Returns true if Piece has liberties
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

//Sets a Piece to empty
func (b *Board) setE(x, y int) {
	b.Grid[y][x].Black = false
	b.Grid[y][x].White = false
	b.Grid[y][x].Empty = true
	return
}

//Checks x,y against size
func (b *Board) checkRange(x, y int) error {
	switch {
	case x < 0 || x >= b.Size:
		return errors.New("x out of range")
	case y < 0 || y >= b.Size:
		return errors.New("y out of range")
	default:
		return nil
	}
}
