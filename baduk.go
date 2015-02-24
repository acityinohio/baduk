//Package baduk implements a library for playing games
//of Baduk/Go. It's optimized for code simplicity,
//and doesn't include any AI support.
package baduk

import "errors"

//A Board represents information about the state
//of a Go game. Size represents the size of the board,
//Grid is the storage of Pieces.
type Board struct {
	Size int
	Grid [][]Piece
}

//A Piece represents information about a piece on the
//Board. Contains pointers to adjacent pieces. If it's
//a border, the pointer is nil.
type Piece struct {
	Black bool
	White bool
	Empty bool
	Up    *Piece //y-1
	Down  *Piece //y+1
	Left  *Piece //x-1
	Right *Piece //x+1
}

//Initializes an empty Board
func (b *Board) Init(size int) (err error) {
	if size < 4 || size > 19 {
		err = errors.New("Size of board must be between 4 and 19")
		return
	}
	b.Size = size
	//Allocate the top-level slice
	b.Grid = make([][]Piece, size)
	for i := range b.Grid {
		//Allocate the intermediate slices
		b.Grid[i] = make([]Piece, size)
	}
	//Set Pieces to Empty, connect them via pointers
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			up, down, left, right := true, true, true, true
			b.Grid[y][x].Empty = true
			//If border, don't connect
			if y == 0 {
				up = false
			}
			if y == size-1 {
				down = false
			}
			if x == 0 {
				left = false
			}
			if x == size-1 {
				right = false
			}
			if up {
				b.Grid[y][x].Up = &b.Grid[y-1][x]
			}
			if down {
				b.Grid[y][x].Down = &b.Grid[y+1][x]
			}
			if left {
				b.Grid[y][x].Left = &b.Grid[y][x-1]
			}
			if right {
				b.Grid[y][x].Right = &b.Grid[y][x+1]
			}
		}
	}
	return
}

//Sets a Piece to white on the Board
//x, y in range from 1 to Board.Size
func (b *Board) SetW(x, y int) (err error) {
	if err = checkRange(x, y, b.Size); err != nil {
		return err
	}
	if !b.Grid[y][x].Empty {
		err = errors.New("Piece is not empty")
		return
	}
	b.Grid[y][x].White = true
	b.Grid[y][x].Black = false
	b.Grid[y][x].Empty = false
	return
}

//Sets a Piece to black on the Board
//x, y in range from 1 to Board.Size
func (b *Board) SetB(x, y int) (err error) {
	if err = checkRange(x, y, b.Size); err != nil {
		return err
	}
	if !b.Grid[y][x].Empty {
		err = errors.New("Piece is not empty")
		return
	}
	b.Grid[y][x].White = false
	b.Grid[y][x].Black = true
	b.Grid[y][x].Empty = false
	return
}

//Creates pretty string, suitable for use
//by fmt.Printf or any logging functions.
//Note that black and white circles are
//reversed compared to their unicode code points;
//this assumes your terminal has a dark background.
func (b *Board) PrettyString() (str string) {
	str = "\n"
	blk := "\u25cb"
	wht := "\u25cf"
	for y := 0; y < b.Size; y++ {
		for x := 0; x < b.Size; x++ {
			p := b.Grid[y][x]
			switch {
			case p.Black:
				str += blk
			case p.White:
				str += wht
			default:
				str += " "
			}
			if x != b.Size-1 {
				str += " - "
			}
		}
		str += "\n"
		if y != b.Size-1 {
			for x := 0; x < b.Size-1; x++ {
				str += "| - "
			}
			str += "|\n"
		}
	}
	str += "\n"
	return
}
