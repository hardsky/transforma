//+build !transforma

package oneFile

type A struct {
	x int
}

type B struct {
	x int
}

func mapAtoB(a *A) *B {
	res := &B{}
	res.x = a.x
	return res

}
