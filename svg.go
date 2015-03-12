package baduk

import "strconv"

//Creates a string representing an SVG
//view of the board, suitable for use
//inline in web templates. Wrap in an external
//div with particular width/height in CSS
//to control size. Yay for resolution
//independence!
func (b *Board) PrettySVG() (svg string) {
	svg += `<svg id="board" width="100%" height="100%" 
		<symbol id="blackstone" viewBox="0 0 120 120">
			<circle cx=60 cy=60 r=45 fill="#000000" />
			<circle cx=80 cy=80 r=10 fill="#ffffff" />
			<circle cx=40 cy=40 r=5 fill="#999999" />
		</symbol>
		<symbol id="whitestone" viewBox="0 0 120 120">
			<circle cx=60 cy=60 r=45 fill="#ffffff" stroke="#000000" stroke-width="3"/>
			<circle cx=80 cy=40 r=10 fill="#aaaaaa" />
			<circle cx=40 cy=80 r=5 fill="#dddddd" />
		</symbol>
		<symbol id="grid" viewBox="0 0 1000 1000">
		`
	base := b.Size * 2
	scale := 1000 / b.Size
	begin := 1000 / base
	end := (base - 1) * begin
	//Make grid
	for i := 1; i < b.Size*2; i += 2 {
		svg += "<line x1=" + strconv.Itoa(begin) + " y1=" + strconv.Itoa(i*begin/base) +
			" x2=" + strconv.Itoa(end) + " y2=" + strconv.Itoa(i*begin/base) +
			" stroke=\"black\" stroke-width=\"10\" />\n"
		svg += "<line x1=" + strconv.Itoa(i*begin/base) + " y1=" + strconv.Itoa(begin) +
			" x2=" + strconv.Itoa(i*begin/base) + " y2=" + strconv.Itoa(end) +
			" stroke=\"black\" stroke-width=\"10\" />\n"
	}
	//Place pieces
	for y := 0; y < b.Size; y++ {
		for x := 0; x < b.Size; x++ {
			if b.Grid[y][x].Black {
				svg += "<use x=" + strconv.Itoa(x*scale) + " y=" + strconv.Itoa(y*scale) +
					" width=" + strconv.Itoa(scale) + " height=" + strconv.Itoa(scale) +
					"xlink:href=\"#blackstone\" />\n"
			} else if b.Grid[y][x].White {
				svg += "<use x=" + strconv.Itoa(x*scale) + " y=" + strconv.Itoa(y*scale) +
					" width=" + strconv.Itoa(scale) + " height=" + strconv.Itoa(scale) +
					"xlink:href=\"#whitestone\" />\n"
			} else {
				continue
			}
		}
	}
	svg += `</symbol>
		<use x="10%" y="10%" width="80%" height="80%" xlink:href="#grid" />
		</svg>
		`
	return
}
