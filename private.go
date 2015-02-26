package baduk

import (
	"errors"
	"log"
)

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
	b.Grid[y][x].checkCapture(isBlack)
	//If there are no liberties after checking opponent's chains,
	//check whether opponent captures chains connected to this move
	if !b.Grid[y][x].hasLiberty() {
		b.Grid[y][x].checkChain(!isBlack)
	}
	return
}

//Checks surrounding stones to see if
//chains can be captured, empties
//chains without liberties
func (p *Piece) checkCapture(checkWhite bool) {
	//Check chain liberties for each direction
	//Only check if piece isn't border, matches checkWhite
	//Or if it doesn't have any liberties
	if p.Up != nil && !p.Up.Empty && p.Up.White == checkWhite && !p.Up.hasLiberty() {
		p.Up.checkChain(checkWhite) //check chain liberties
	}
	if p.Down != nil && !p.Down.Empty && p.Down.White == checkWhite && !p.Down.hasLiberty() {
		p.Down.checkChain(checkWhite) //check chain liberties
	}
	if p.Left != nil && !p.Left.Empty && p.Left.White == checkWhite && !p.Left.hasLiberty() {
		p.Left.checkChain(checkWhite) //check chain liberties
	}
	if p.Right != nil && !p.Right.Empty && p.Right.White == checkWhite && !p.Right.hasLiberty() {
		p.Left.checkChain(checkWhite) //check chain liberties
	}
	return
}

//If chain has no liberties, Empty it
func (p *Piece) checkChain(checkWhite bool) {
	//Create map of Piece addresses representing chain
	//Needs to be buffered to prevent deadlock
	chainChan := make(chan map[*Piece]bool, 1)
	libChan := make(chan bool)
	done := make(chan bool)
	go func() { chainChan <- make(map[*Piece]bool) }()
	go crawler(p, checkWhite, chainChan, libChan, done)
	for {
		select {
		//If no liberties found, "empty" all pieces in chain
		case <-done:
			emptyChain(<-chainChan)
			return
		//if liberty found, return
		case <-libChan:
			return
		default:
		}
	}
}

//Check up, down, left, right pieces (If not nil && same color as checkWhite)
//if any adjacent checkWhite piece has liberties, quit without capturing
//Otherwise, "travel" to it, add piece to map, check its directions for liberties
//Keep traveling, but don't go to already traveled places in map
func crawler(pi *Piece, checkWhite bool, chainChan chan map[*Piece]bool, libChan chan bool, done chan bool) {
	/*go func() {
		for {
			select {
			case <-done:
				return
			default:
				still <- true
			}
		}
	}()*/
	chain := <-chainChan
	chain[pi] = true
	chainChan <- chain
	if pi.hasLiberty() {
		libChan <- true
		return
	}
	if pi.Up != nil && pi.Up.White == checkWhite && !chain[pi.Up] {
		log.Print("Going up", chain)
		go crawler(pi.Up, checkWhite, chainChan, libChan, done)
	}
	if pi.Down != nil && pi.Down.White == checkWhite && !chain[pi.Down] {
		log.Print("Going down", chain)
		go crawler(pi.Down, checkWhite, chainChan, libChan, done)
	}
	if pi.Left != nil && pi.Left.White == checkWhite && !chain[pi.Left] {
		log.Print("Going left", chain)
		go crawler(pi.Left, checkWhite, chainChan, libChan, done)
	}
	if pi.Right != nil && pi.Right.White == checkWhite && !chain[pi.Right] {
		log.Print("Going right", chain)
		go crawler(pi.Right, checkWhite, chainChan, libChan, done)
	}
	done <- true
	return
}

//Empties chain represented by map[*Piece]bool
func emptyChain(chain map[*Piece]bool) {
	for p := range chain {
		p.Black = false
		p.White = false
		p.Empty = true
	}
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
