package baduk

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"errors"
)

type Board struct {
	Size       int
	Grid       [][]Piece
	BlackScore int
	WhiteScore int
	State      string
}

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
	if size < 3 || size > 19 {
		err = errors.New("Size of board must be between 3 and 19")
		return
	}
	b.Size = size
	//Allocate the top-level slice
	b.Grid = make([][]Piece, size)
	for i := range b.Grid {
		//Allocate the intermediate slices
		b.Grid[i] = make([]Piece, size)
	}
	//Encode empty state into string
	if err = b.Encode(); err != nil {
		return err
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
func (b *Board) SetE(x, y int) (err error) {
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

//Encodes the Board state into a URL-safe string
func (b *Board) Encode() (err error) {
	var a bytes.Buffer
	//first byte of the buffer is size
	if err = a.WriteByte(byte(b.Size)); err != nil {
		return err
	}
	//use zlib to compress
	w := zlib.NewWriter(&a)
	for _, v := range b.Grid {
		for _, s := range v {
			switch {
			case s.Black:
				w.Write([]byte("b"))
			case s.White:
				w.Write([]byte("w"))
			default:
				w.Write([]byte("e"))
			}
		}
	}
	w.Close()
	b.State = base64.URLEncoding.EncodeToString(a.Bytes())
	return
}

//Initializes a Board from a URL-safe string
//encoded with Board.Encode
func (b *Board) Decode(str string) (err error) {
	data, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		return
	}
	//first byte of the data is size
	size := int(data[0])
	b.Init(size)
	//set up zlib reader
	rest := bytes.NewReader(data[1:])
	r, err := zlib.NewReader(rest)
	p := bufio.NewReader(r)
	if err != nil {
		return
	}
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			c, errr := p.ReadByte()
			if errr != nil {
				err = errr
				return
			}
			switch c {
			case []byte("b")[0]:
				err = b.SetB(x, y)
			case []byte("w")[0]:
				err = b.SetW(x, y)
			case []byte("e")[0]:
				err = b.SetE(x, y)
			default:
				err = errors.New("Piece not recognized during decode")
			}
			if err != nil {
				return
			}
		}
	}
	r.Close()
	return
}

//Scores the Board, and
//empties Pieces without liberties
func (b *Board) Score() {
}
