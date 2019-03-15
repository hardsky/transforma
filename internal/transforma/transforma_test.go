package transforma

import "testing"

func TestGenerateAnyNames(t *testing.T) {
	err := Generate("./test-data/any-names")
	if err != nil {
		t.Error(err)
	}
}
