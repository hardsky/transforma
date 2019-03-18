package skipField

type Foo struct {
	SameName string `trf:"-"`
	OtheSame string
	Mapped   string
}

type Bar struct {
	SameName string
	OtheSame string `trf:"-"`
	Mapped   string
}
