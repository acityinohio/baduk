package baduk

import (
	"errors"
	"sync"
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
//chains can be captured
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
		p.Right.checkChain(checkWhite) //check chain liberties
	}
	return
}

//If chain has no liberties, Empty it
func (p *Piece) checkChain(checkWhite bool) {
	chainChan := make(chan map[*Piece]bool, 1)
	libChan := make(chan bool)
	done := make(chan bool)
	var wg sync.WaitGroup
	//Initialize channel containing map representing chain
	go func() { chainChan <- make(map[*Piece]bool) }()
	//Crawl chain, use WaitGroup to detect when recursion finishes
	wg.Add(1)
	go crawler(p, checkWhite, chainChan, libChan, &wg)
	//Send done signal if recursion is finished
	go func() {
		wg.Wait()
		done <- true
	}()
	select {
	//If no liberties found, "empty" all pieces in chain
	case <-done:
		emptyChain(<-chainChan)
		return
	//if liberty found, return
	case <-libChan:
		return
	}
}

//Check up, down, left, right pieces (If not nil && same color as checkWhite)
//if any adjacent checkWhite piece has liberties, quit without capturing
//Otherwise, crawl to it, check its directions for liberties
//Recursively crawl, but don't go to already crawled places in chainChan
func crawler(pi *Piece, checkWhite bool, chainChan chan map[*Piece]bool, libChan chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	chain := <-chainChan
	chain[pi] = true
	chainChan <- chain
	if pi.hasLiberty() {
		libChan <- true
		return
	}
	if pi.Up != nil && pi.Up.White == checkWhite && !chain[pi.Up] {
		wg.Add(1)
		go crawler(pi.Up, checkWhite, chainChan, libChan, wg)
	}
	if pi.Down != nil && pi.Down.White == checkWhite && !chain[pi.Down] {
		wg.Add(1)
		go crawler(pi.Down, checkWhite, chainChan, libChan, wg)
	}
	if pi.Left != nil && pi.Left.White == checkWhite && !chain[pi.Left] {
		wg.Add(1)
		go crawler(pi.Left, checkWhite, chainChan, libChan, wg)
	}
	if pi.Right != nil && pi.Right.White == checkWhite && !chain[pi.Right] {
		wg.Add(1)
		go crawler(pi.Right, checkWhite, chainChan, libChan, wg)
	}
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
