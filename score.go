package baduk

import (
	"log"
	"strconv"
	"sync"
)

//Scores the Board. Scoring follows this algorithm:
//A color's score = total pieces + total empty pieces completely
//enclosed by said color. If empty pieces are enclosed by both colors,
//then empty territory is contested and not added to either score.
//This method is consistent with a simple version of Chinese area scoring.
func (b *Board) Score() (black, white int) {
	//Edge case: return 0,0 if entire board is empty
	if b.isEmpty() {
		return
	}
	//Counts black/white/empty Pieces, sends down channels
	blk, wht := make(chan int, 1), make(chan int, 1)
	empty := make(chan *Piece)
	var wg sync.WaitGroup
	wg.Add(1)
	go b.countStones(blk, wht, empty, &wg)
	//Do the counting
	wg.Add(1)
	go sumChan(blk, &black, &wg)
	wg.Add(1)
	go sumChan(wht, &white, &wg)
	wg.Add(1)
	go checkEmptyArea(empty, &wg)
	wg.Wait()
	return
}

//Counts stones, sends counts or Pieces down channels
func (b *Board) countStones(blk chan int, wht chan int, empty chan *Piece, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, j := range b.Grid {
		for _, i := range j {
			if i.Black {
				blk <- 1
			} else if i.White {
				wht <- 1
			} else {
				empty <- &i
			}
		}
	}
	close(blk)
	close(wht)
	close(empty)
	return
}

//Takes channel, sums it, writes to total
func sumChan(pipe chan int, total *int, wg *sync.WaitGroup) {
	defer wg.Done()
	sum := 0
	for i := range pipe {
		sum += i
	}
	*total = sum
	return
}

//Figures out how to score empty area
//Currently dumby function to prevent deadlock
func checkEmptyArea(empty chan *Piece, wg *sync.WaitGroup) {
	defer wg.Done()
	empties := make(map[*Piece]bool)
	for p := range empty {
		empties[p] = true
	}
	log.Print(empties)
}

//Returns true if Board is empty
func (b *Board) isEmpty() bool {
	for _, j := range b.Grid {
		for _, i := range j {
			if !i.Empty {
				return false
			}
		}
	}
	return true
}

//Makes a string suitable for Sprintf output that
//declares the winner
func (b *Board) ScorePretty() (str string) {
	black, white := b.Score()
	switch {
	case black > white:
		str = "Black wins, by " + strconv.Itoa(black-white) + "\n"
	case white > black:
		str = "White wins, by " + strconv.Itoa(white-black) + "\n"
	case black == white:
		str = "Tie game."
	}
	return
}
