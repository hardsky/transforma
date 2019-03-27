package transforma

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestGenerateAnyNames(t *testing.T) {
	testDataDir := "./test-data"
	gotFile := "mapper_gen.go"
	wantFile := "mapper-want.txt"

	cases := []string{
		"any-names",
		"one-file",
		"same-structs",
		"skip-field",
		"str-to-int",
	}

	for _, c := range cases {
		pkgPath := filepath.Join(testDataDir, c)
		// t.Log("package:", pkgPath)

		err := Generate("./" + pkgPath)
		if err != nil {
			t.Error(err)
			continue
		}

		gotFileName := filepath.Join(pkgPath, gotFile)
		// t.Log("got file:", gotFileName)

		got, err := ioutil.ReadFile(gotFileName)
		if err != nil {
			t.Error(err)
			continue

		}

		wantFileName := filepath.Join(pkgPath, wantFile)
		// t.Log("want file:", wantFileName)

		want, err := ioutil.ReadFile(wantFileName)
		if err != nil {
			t.Error(err)
			continue
		}

		if string(want) != string(got) {
			t.Errorf("error with %q test", c)
		}
	}
}
