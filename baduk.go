//Package baduk implements a library for playing games
//of Baduk/Go. It's optimized for code simplicity,
//and does not include any AI support.
package baduk

import "errors"

//A Board represents information about the state
//of a Go game. Size represents the size of the board,
//Grid is the storage of Pieces, BlackScore and WhiteScore
//store the currently calculated scores of the Board, while
//State represents a compressed, base64-encoded string of
//the state of the board, suitable for use in URLs.
type Board struct {
	Size       int
	Grid       [][]Piece
	BlackScore int
	WhiteScore int
	State      string
}

//A Piece represents information about a piece on the
//Board. When both Black and White are false, the Piece
//is considered empty. If both Black and White are true,
//something is seriously wrong with the library.
type Piece struct {
	Black bool
	White bool
}

//Returns true if piece is empty
func (p *Piece) empty() bool {
	return !(p.Black || p.White)
}

//Returns err if piece is not empty
func (p *Piece) NotEmpty() error {
	if !p.empty() {
		return errors.New("Piece is not empty")
	}
	return nil
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
	return
}

//Sets a Piece to white on the Board
//x, y in range from 1 to Board.Size
func (b *Board) SetW(x, y int) (err error) {
	if err = checkRange(x, y, b.Size); err != nil {
		return err
	}
	if err = b.Grid[x][y].NotEmpty(); err != nil {
		return err
	}
	b.Grid[x][y].White = true
	//Calls Score to update Board
	b.Score()
	return
}

//Sets a Piece to black on the Board
//x, y in range from 1 to Board.Size
func (b *Board) SetB(x, y int) (err error) {
	if err = checkRange(x, y, b.Size); err != nil {
		return err
	}
	if err = b.Grid[x][y].NotEmpty(); err != nil {
		return err
	}
	b.Grid[x][y].Black = true
	//Calls Score to update Board
	b.Score()
	return
}

//Sets a Piece to empty on the Board
//x, y in range from 1 to Board.Size
//Used only by the Decode and Score
//function and not publicly scoped.
func (b *Board) setE(x, y int) (err error) {
	if err = checkRange(x, y, b.Size); err != nil {
		return err
	}
	b.Grid[x][y].Black = false
	b.Grid[x][y].White = false
	return
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

//Creates pretty string, suitable for use
//by fmt.Printf or any logging functions.
//Note that black and white circles are
//reversed compared to their unicode code points;
//this assumes your terminal has a dark background.
func (b *Board) PrettyString() (str string) {
	str = "\n"
	blk := "\u25cb"
	wht := "\u25cf"
	for x := 0; x < b.Size; x++ {
		for y := 0; y < b.Size; y++ {
			p := b.Grid[x][y]
			switch {
			case p.Black:
				str += blk
			case p.White:
				str += wht
			default:
				str += " "
			}
			if y != b.Size-1 {
				str += " - "
			}
		}
		str += "\n"
		if x != b.Size-1 {
			for y := 0; y < b.Size-1; y++ {
				str += "| - "
			}
			str += "|\n"
		}
	}
	str += "\n"
	return
}
