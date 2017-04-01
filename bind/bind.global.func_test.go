package bind

import (
	"testing"

	"github.com/qxnw/lib4go/ut"
	"github.com/yuin/gopher-lua"
)

type userData struct {
	Name string
	ID   int
}

func TestBindGlobalGetParams(t *testing.T) {
	s := lua.NewState()
	values := globalGetParams(s)
	ut.Expect(t, len(values), 0)

	pushValues(s, "a")
	values = globalGetParams(s)
	ut.Expect(t, len(values), 1)
	ut.Expect(t, values[0].(string), "a")

	pushValues(s, 11)
	values = globalGetParams(s)
	ut.Expect(t, len(values), 1)
	ut.Expect(t, values[0].(string), "11")

	pushValues(s, "11", 2)
	values = globalGetParams(s)
	ut.Expect(t, len(values), 2)
	ut.Expect(t, values[0].(string), "11")
	ut.Expect(t, values[1].(string), "2")

	ud := &userData{Name: "colin", ID: 10}
	pushValues(s, ud)
	values = globalGetParams(s)
	ut.Expect(t, len(values), 1)
	ut.Expect(t, values[0], ud)

}
