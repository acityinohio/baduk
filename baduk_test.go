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
		t.Logf(b.PrettyString())
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
	t.Logf(b.PrettyString())
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
	b.Decode("BGJiYmAAUgyMDGCCARAAAP__")
	t.Logf(b.PrettyString())
	//b.SetB(0, 3)
	b.SetW(2, 0)
	t.Logf(b.PrettyString())
}
