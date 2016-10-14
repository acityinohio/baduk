package baduk

import "testing"

func TestInit(t *testing.T) {
	var b Board
	err := b.Init(3)
	if err.Error() != "Size of board must be between 4 and 19" {
		t.Error("Expected error for small board, got", err)
	}
	err = b.Init(20)
	if err.Error() != "Size of board must be between 4 and 19" {
		t.Error("Expected error for large board, got", err)
	}
	err = b.Init(13)
	if err != nil {
		t.Error("Error Initializing:", err)
	}
	return
}

func TestEncode(t *testing.T) {
	var b Board
	b.Init(13)
	expectEncode := "DWIYKgAQAAD__w=="
	realEncode, err := b.Encode()
	if err != nil {
		t.Error("Error encoding empty 13x13 board:", err)
	}
	if realEncode != expectEncode {
		t.Error("Expected Encoded string", expectEncode, " got ", realEncode)
	}
	//further spot check for encoding if verbose
	if testing.Verbose() {
		b.Init(4)
		b.SetB(0, 0)
		b.SetB(0, 1)
		b.SetW(1, 1)
		b.SetW(2, 2)
		enc, _ := b.Encode()
		t.Log(enc) //generates "BGJiYGAAUkAAJhgAAQAA__8="
	}
}

func TestDecode(t *testing.T) {
	var b Board
	str := "DWIYKgAQAAD__w=="
	err := b.Decode(str)
	if err != nil {
		t.Error("Error Decoding:", err)
	}
	if b.Size != 13 {
		t.Error("Expected size to be 13, got", b.Size)
	}
	//Check as generated above
	str = "BGJiYGAAUkAAJhgAAQAA__8="
	err = b.Decode(str)
	if err != nil {
		t.Error("Error Decoding:", err)
	}
	if b.Size != 4 {
		t.Error("Expected size to be 4, got", b.Size)
	}
	if !(b.Grid[0][0].Black && b.Grid[1][0].Black && b.Grid[1][1].White && b.Grid[2][2].White) {
		t.Errorf("Expected different board, got" + b.PrettyString())
	}
	return
}

func TestHasLiberty(t *testing.T) {
	var b Board
	err := b.Decode("BGJiYGAAUkAAJhgAAQAA__8=")
	if err != nil {
		t.Error("Error Decoding:", err)
	}
	b.SetB(1, 0)
	if b.Grid[0][0].hasLiberty() {
		t.Error("Expected false, got true with piece at 0,0", b.Grid[0][0])
		t.Logf(b.PrettyString())
	}
	if !b.Grid[1][1].hasLiberty() {
		t.Error("Expected false, got true with piece at 1,1", b.Grid[1][1])
		t.Logf(b.PrettyString())
	}
	return
}

func TestCheckCapture(t *testing.T) {
	var b Board
	b.Decode("BGJiYmBgYmRgYGQAEwyAAAAA__8=")
	//Attempt top corner black chain capture
	b.SetW(2, 0)
	expectStr := "BGJgYIQiCA0IAAD__w=="
	realStr, err := b.Encode()
	if err != nil {
		t.Error("Error encoding:", err)
	}
	if realStr != expectStr {
		t.Error("Expected top board, got bottom board")
		var c Board
		c.Decode(expectStr)
		t.Logf(c.PrettyString())
		t.Logf(b.PrettyString())
	}
	//Reset board
	b.Decode("BGJiYmBgYmRgYGQAEwyAAAAA__8=")
	//Attempt suicide on bottom left
	//Should return to prior state
	b.SetB(0, 3)
	expectStr = "BGJiYmBgYmRgYGQAEwyAAAAA__8="
	realStr, err = b.Encode()
	if err != nil {
		t.Error("Error encoding:", err)
	}
	if realStr != expectStr {
		t.Error("Expected top board, got bottom board")
		var c Board
		c.Decode(expectStr)
		t.Logf(expectStr + "\n")
		t.Logf(realStr + "\n")
	}
	//Reset board with new setup
	b.Decode("BATAAQEAAACCoPD_6KgtDbE9AAD__w==")
	//Attempt middle board chain captures
	b.SetB(2, 3)
	expectStr = "BGJiYmBgYmBiYABRDEwMgAAAAP__"
	realStr, err = b.Encode()
	if err != nil {
		t.Error("Error encoding:", err)
	}
	if realStr != expectStr {
		t.Error("Expected top board, got bottom board")
		var c Board
		c.Decode(expectStr)
		t.Logf(c.PrettyString())
		t.Logf(b.PrettyString())
	}
}

func TestScore(t *testing.T) {
	var b Board
	b.Init(4)
	//check empty board
	black, white := b.Score()
	if black != 0 || white != 0 {
		t.Error("For empty board, expected black: 0, white: 0, got black:", black, ", white:", white)
	}
	//check early game
	b.SetB(0, 0)
	b.SetW(2, 1)
	black, white = b.Score()
	if black != 1 || white != 1 {
		t.Error("Expected black: 1, white: 1, got black:", black, ", white:", white)
		t.Logf(b.PrettyString())
	}
	//check later game
	b.Decode("BGJiYmBgYmRiYGQAEwyAAAAA__8=")
	black, white = b.Score()
	if black != 8 || white != 6 {
		t.Error("For this board, expected black: 8, white: 6, got black:", black, ", white:", white)
		t.Errorf(b.PrettyString())
	}
}
