// +build transforma

package anyNames

func mapperFooBar(f *Foo) *Bar {
	return &Bar{}
}

func mapperBarFoo(b *Bar) *Foo {
	return &Foo{}
}
