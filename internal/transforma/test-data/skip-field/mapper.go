// +build transforma

package skipField

func mapperFooBar(f *Foo) *Bar {
	return &Bar{}
}

func mapperBarFoo(b *Bar) *Foo {
	return &Foo{}
}
