//+build !transforma

package anyNames

func mapperFooBar(f *Foo) *Bar {
	res := &Bar{}
	res.OtherName = f.Name
	return res

}

func mapperBarFoo(b *Bar) *Foo {
	res := &Foo{}
	res.Name = b.OtherName
	return res

}
