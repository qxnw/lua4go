package bind

import (
	"github.com/arsgo/lib4go/draw"
	"github.com/arsgo/lib4go/script"
	"github.com/qxnw/lib4go/images"
	lua "github.com/yuin/gopher-lua"
)

func getImageDrawTypeBinder() *TypeBinder {
	return &TypeBinder{
		Name: "draw",
		NewFunc: map[string]lua.LGFunction{
			"new": typeNewDraw,
		},
		Methods: map[string]lua.LGFunction{
			"draw_image": typeDrawImage,
			"draw_font":  typeDrawFont,
			"save":       typeDrawSave,
		},
	}
}

// Constructor
func typeNewDraw(L *lua.LState) int {
	with := L.CheckInt(1)
	height := L.CheckInt(2)
	count := L.GetTop()
	ud := L.NewUserData()
	var err error
	if count > 2 {
		path := L.CheckString(3)
		ud.Value, err = images.NewImageFromFile(with, height, path)
		L.SetMetatable(ud, L.GetTypeMetatable("draw"))
		L.Push(ud)
		L.Push(script.New(L, err))
		return 2
	}
	ud.Value = draw.NewDraw(with, height)
	L.SetMetatable(ud, L.GetTypeMetatable("draw"))
	L.Push(ud)
	return 1

}

// Checks whether the first lua argument is a *LUserData with *Person and returns this *Person.
func checkDraw(L *lua.LState) *images.Image {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*images.Image); ok {
		return v
	}
	L.RaiseError("bad argument  (draw expected, got %s)", ud.Type().String())
	return nil
}
func typeDrawImage(L *lua.LState) int {
	p := checkDraw(L)
	path := L.CheckString(2)
	sx := L.CheckInt(3)
	sy := L.CheckInt(4)
	ex := L.CheckInt(5)
	ey := L.CheckInt(6)
	if L.GetTop() > 6 {
		w := L.CheckInt(7)
		h := L.CheckInt(8)
		err := p.DrawImageWithScale(path, sx, sy, ex, ey, w, h)
		return pushValues(L, err)
	}
	err := p.DrawImage(path, sx, sy, ex, ey)
	return pushValues(L, err)
}

func typeDrawFont(L *lua.LState) int {
	p := checkDraw(L)
	path := L.CheckString(2)
	text := L.CheckString(3)
	color := L.CheckString(4)
	font := L.CheckInt(5)
	sx := L.CheckInt(6)
	sy := L.CheckInt(7)
	err := p.DrawFont(path, text, color, float64(font), sx, sy)
	return pushValues(L, err)
}
func typeDrawSave(L *lua.LState) int {
	p := checkDraw(L)
	path := L.CheckString(2)
	err := p.Save(path)
	return pushValues(L, err)
}
