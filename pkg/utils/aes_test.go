package utils

import (
	"testing"
)

func TestDecryptByAes(t *testing.T) {
	crpytoTxt, err := DecryptByAes("tE15VztAwmVUcfCjMNmKWw==")
	if err != nil {
		t.Fatal(err)
	}
	txt := "test123"
	if string(crpytoTxt) != txt {
		t.Error("Decrypt Failed")
	}

}

func TestEncryptByAes(t *testing.T) {
	txt := "test123"
	enc, err := EncryptByAes([]byte(txt)) //tE15VztAwmVUcfCjMNmKWw==
	if err != nil {
		t.Fatal(err)
	}
	t.Log(enc)
}
