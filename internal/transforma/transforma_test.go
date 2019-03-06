package transforma

import "testing"

func TestGenerate(t *testing.T) {
	err := Generate("./test-data/same-structs")
	if err != nil {
		t.Error(err)
	}
}
