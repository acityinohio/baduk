package baduk

import (
	"bufio"
	"bytes"
	"compress/flate"
	"encoding/base64"
	"errors"
)

//Encodes the Board state into a compressed,
//base64-encoded URL-safe string
func (b *Board) Encode() (err error) {
	var a bytes.Buffer
	//first byte of the buffer is size
	if err = a.WriteByte(byte(b.Size)); err != nil {
		return err
	}
	//use flate to compress
	dict := []byte{2, 1, 0}
	w, err := flate.NewWriterDict(&a, flate.BestCompression, dict)
	for _, v := range b.Grid {
		for _, s := range v {
			switch {
			case s.Black:
				w.Write(dict[0:1])
			case s.White:
				w.Write(dict[1:2])
			default:
				w.Write(dict[2:3])
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
	err = b.Init(size)
	if err != nil {
		return
	}
	//set up flate reader with dict
	dict := []byte{2, 1, 0}
	r := flate.NewReaderDict(bytes.NewReader(data[1:]), dict)
	p := bufio.NewReader(r)
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			c, errr := p.ReadByte()
			if errr != nil {
				err = errr
				return
			}
			switch c {
			case dict[0]:
				err = b.SetB(x, y)
			case dict[1]:
				err = b.SetW(x, y)
			case dict[2]:
				err = b.setE(x, y)
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
