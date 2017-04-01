package bind

import (
	"testing"

	"github.com/qxnw/lib4go/ut"
)

func TestBind1(t *testing.T) {
	binder := NewDefault("../;../script")
	ut.Expect(t, len(binder.Packages), 2)
	ut.Refute(t, len(binder.GlobalFunc), 0)
	ut.Refute(t, len(binder.Modules), 0)
	ut.Refute(t, len(binder.Types), 0)

}
func TestBind2(t *testing.T) {
	binder := NewDefault("")
	ut.Expect(t, len(binder.Packages), 0)
	ut.Refute(t, len(binder.GlobalFunc), 0)
	ut.Refute(t, len(binder.Modules), 0)
	ut.Refute(t, len(binder.Types), 0)

	binder = NewDefault(";a;")
	ut.Expect(t, len(binder.Packages), 1)
	ut.Refute(t, len(binder.GlobalFunc), 0)
	ut.Refute(t, len(binder.Modules), 0)
	ut.Refute(t, len(binder.Types), 0)

}
