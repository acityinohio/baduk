package baduk

import "testing"

func TestInit(t *testing.T) {
	var b Board
	err := b.Init(2)
	if err.Error() != "Size of board must be between 3 and 19" {
		t.Error("Expected error for small board, got", err)
	}
	err = b.Init(20)
	if err.Error() != "Size of board must be between 3 and 19" {
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
	expectEncode := "DXicSh0qABAAAP__JWRCrg=="
	if b.State != expectEncode {
		t.Error("Expected b.State: ", expectEncode, " got ", b.State)
	}
}

func TestDecode(t *testing.T) {
	var b Board
	str := "DXicSh0qABAAAP__JWRCrg=="
	err := b.Decode(str)
	if err != nil {
		t.Error("Error Decoding:", err)
	}
	if b.Size != 13 {
		t.Error("Expected size to be 13, got", b.Size)
	}
	return
}
