package anyNames

type Foo struct {
	Name string
}

type Bar struct {
	OtherName string `trf:"Name"`
}
