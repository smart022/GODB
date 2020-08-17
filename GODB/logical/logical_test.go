package logical

import (
	//"os"
	"testing"
)

func TestValRef(t *testing.T) {

	var cot ValueRef
	tt := NewTrueValRef("", 0)
	cot = &tt
	t.Log(cot.address())
}
