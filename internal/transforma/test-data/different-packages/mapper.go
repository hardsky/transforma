// +build transforma

package differentPackages

import (
	"github.com/hardsky/transforma/internal/transforma/test-data/different-packages/different"
)

func mapFooToBar(foo *Foo) *different.Bar {
	return &different.Bar{}
}

func mapBarToFoo(bar *different.Bar) *Foo {
	return &Foo{}
}
