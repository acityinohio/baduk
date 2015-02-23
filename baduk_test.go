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
	err := b.Encode()
	if err != nil {
		t.Error("Error Encoding:", err)
	}
	expectEncode := "DWIYKgAQAAD__w=="
	if b.State != expectEncode {
		t.Error("Expected b.State: ", expectEncode, " got ", b.State)
	}
	//further spot check for encoding if verbose
	if testing.Verbose() {
		b.Init(4)
		b.SetB(0, 0)
		b.SetB(0, 1)
		b.SetW(1, 1)
		b.SetW(2, 2)
		b.Encode()
		t.Log(b.State) //generates "BGJiYmBgYGSAEQyAAAAA__8="
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
	str = "BGJiYmBgYGSAEQyAAAAA__8="
	err = b.Decode(str)
	if err != nil {
		t.Error("Error Decoding:", err)
	}
	if b.Size != 4 {
		t.Error("Expected size to be 4, got", b.Size)
	}
	if !(b.Grid[0][0].Black && b.Grid[0][1].Black && b.Grid[1][1].White && b.Grid[2][2].White) {
		t.Errorf("Expected different board, got" + b.PrettyString())
	}
	return
}
