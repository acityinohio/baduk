package baduk

import "strconv"

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
	return
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
