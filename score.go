package baduk

import (
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
	empty := make(chan Piece)
	var wg sync.WaitGroup
	wg.Add(1)
	go b.countStones(blk, wht, empty, &wg)
	//Do the counting
	wg.Add(1)
	go sumChan(blk, &black, &wg)
	wg.Add(1)
	go sumChan(wht, &white, &wg)
	//mirror checkChain() from scoring.go, but focused on assembling empty chains
	wg.Add(1)
	go checkEmptyChains(blk, wht, empty, &wg)
	wg.Wait()
	return
}

//Counts stones, sends counts or Pieces down channels
func (b *Board) countStones(blk chan int, wht chan int, empty chan Piece, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, j := range b.Grid {
		for _, i := range j {
			if i.Black {
				blk <- 1
			}
			if i.White {
				wht <- 1
			}
			if i.Empty {
				empty <- i
			}
		}
	}
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

//Assembles chains of Empty pieces
func checkEmptyChains(blk chan int, wht chan int, empty chan Piece, wg *sync.WaitGroup) {
	defer wg.Done()
	chains := make([]map[Piece]bool, 0)
	checked := make(map[Piece]bool)
	for lib := range empty {
		if checked[lib] {
			continue
		}
		//set up new WaitGroup to signal end of recursion
		var ewg sync.WaitGroup
		ewg.Add(1)
		//Initialize chain channel
		chainChan := make(chan map[Piece]bool, 1)
		go func() { chainChan <- make(map[Piece]bool) }()
		//crawls the empties, recursively
		go emptyCrawler(lib, chainChan, &ewg)
		ewg.Wait()
		chain := <-chainChan
		//adds chain to "checked"
		for i, v := range chain {
			checked[i] = v
		}
		//adds chain to chains
		chains = append(chains, chain)
	}
	wg.Add(1)
	scoreEmptyChains(blk, wht, chains, wg)
	return
}

func emptyCrawler(p Piece, chainChan chan map[Piece]bool, ewg *sync.WaitGroup) {
	defer ewg.Done()
	chain := <-chainChan
	chain[p] = true
	chainChan <- chain
	if p.Up != nil && p.Up.Empty && !chain[*p.Up] {
		ewg.Add(1)
		go emptyCrawler(*p.Up, chainChan, ewg)
	}
	if p.Down != nil && p.Down.Empty && !chain[*p.Down] {
		ewg.Add(1)
		go emptyCrawler(*p.Down, chainChan, ewg)
	}
	if p.Left != nil && p.Left.Empty && !chain[*p.Left] {
		ewg.Add(1)
		go emptyCrawler(*p.Left, chainChan, ewg)
	}
	if p.Right != nil && p.Right.Empty && !chain[*p.Right] {
		ewg.Add(1)
		go emptyCrawler(*p.Right, chainChan, ewg)
	}
	return
}

//Checks empty area for scoring by assembling
//empty chains then checking for border encapsulations
func scoreEmptyChains(blk chan int, wht chan int, chains []map[Piece]bool, wg *sync.WaitGroup) {
	defer wg.Done()
	//Investigate empty chain's borders
	for _, chain := range chains {
		bBord, wBord := false, false
		for lib := range chain {
			bBord = bBord || lib.hasBlackBorder()
			wBord = wBord || lib.hasWhiteBorder()
			if bBord && wBord {
				break
			}
		}
		if bBord && wBord {
			continue
		} else if bBord {
			blk <- len(chain)
		} else if wBord {
			wht <- len(chain)
		}
	}
	close(blk)
	close(wht)
}

//Returns true if piece has adjacent Black piece
func (p *Piece) hasBlackBorder() bool {
	if p.Up != nil && p.Up.Black {
		return true
	} else if p.Down != nil && p.Down.Black {
		return true
	} else if p.Left != nil && p.Left.Black {
		return true
	} else if p.Right != nil && p.Right.Black {
		return true
	} else {
		return false
	}
}

//Returns true if piece has adjacent White piece
func (p *Piece) hasWhiteBorder() bool {
	if p.Up != nil && p.Up.White {
		return true
	} else if p.Down != nil && p.Down.White {
		return true
	} else if p.Left != nil && p.Left.White {
		return true
	} else if p.Right != nil && p.Right.White {
		return true
	} else {
		return false
	}
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
