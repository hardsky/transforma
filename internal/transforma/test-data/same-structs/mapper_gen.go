//+build !transforma

package sameStructs

func mapAtoB(a *A) *B {
	res := &B{}
	res.x = a.x
	return res

}
