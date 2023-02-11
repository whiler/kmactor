package kmactor

import (
	"github.com/go-vgo/robotgo"
)

const (
	MajorKey   = 1
	MajorMouse = 2

	KeyDown  = 1
	KeyPress = 2
	KeyUp    = 3

	MouseMove = 1
	MouseDown = 2
	MouseUp   = 3

	MouseKeyLeft   = 0
	MouseKeyCenter = 1
	MouseKeyRight  = 2
)

var width int
var height int

var keyboard = map[int]string{
	8:   "backspace",
	9:   "tab",
	13:  "enter",
	16:  "shift",
	17:  "ctrl",
	18:  "alt",
	20:  "capslock",
	27:  "esc",
	32:  "space",
	33:  "pageup",
	34:  "pagedown",
	35:  "end",
	36:  "home",
	37:  "left",
	38:  "up",
	39:  "right",
	40:  "down",
	46:  "delete",
	48:  "0",
	49:  "1",
	50:  "2",
	51:  "3",
	52:  "4",
	53:  "5",
	54:  "6",
	55:  "7",
	56:  "8",
	57:  "9",
	65:  "a",
	66:  "b",
	67:  "c",
	68:  "d",
	69:  "e",
	70:  "f",
	71:  "g",
	72:  "h",
	73:  "i",
	74:  "j",
	75:  "k",
	76:  "l",
	77:  "m",
	78:  "n",
	79:  "o",
	80:  "p",
	81:  "q",
	82:  "r",
	83:  "s",
	84:  "t",
	85:  "u",
	86:  "v",
	87:  "w",
	88:  "x",
	89:  "y",
	90:  "z",
	112: "f1",
	113: "f2",
	114: "f3",
	115: "f4",
	116: "f5",
	117: "f6",
	118: "f7",
	119: "f8",
	120: "f9",
	121: "f10",
	122: "f11",
	123: "f12",
	186: ";",
	187: "=",
	188: ",",
	189: "-",
	190: ".",
	191: "/",
	192: "`",
	219: "[",
	220: "\\",
	221: "]",
	222: "'",
}

var mouse = map[int]string{
	MouseKeyLeft:   "left",
	MouseKeyCenter: "center",
	MouseKeyRight:  "right",
}

type Command struct {
	Major    int       `json:"l"`
	Type     int       `json:"t"`
	Key      int       `json:"k"`
	Size     *Size     `json:"s"`
	Position *Position `json:"p"`
}

func (self *Command) Reset() {
	self.Major = 0
	self.Type = 0
	self.Key = 0
}

type Size struct {
	Width  int `json:"w"`
	Height int `json:"h"`
}

type Position struct {
	Left int `json:"l"`
	Top  int `json:"t"`
}

func Play(cmd *Command) bool {
	var key string
	var handled bool
	switch cmd.Major {
	case MajorKey:
		if key = keyboard[cmd.Key]; key != "" {
			switch cmd.Type {
			case KeyDown:
				robotgo.KeyDown(key)
				handled = true
			case KeyPress:
				robotgo.KeyPress(key)
				handled = true
			case KeyUp:
				robotgo.KeyUp(key)
				handled = true
			}
		}
	case MajorMouse:
		switch cmd.Type {
		case MouseMove:
			if cmd.Position != nil && cmd.Size != nil && cmd.Size.Width > 0 && cmd.Size.Height > 0 {
				robotgo.Move(cmd.Position.Left*width/cmd.Size.Width, cmd.Position.Top*height/cmd.Size.Height)
				handled = true
			}
		case MouseDown:
			if key = mouse[cmd.Key]; key != "" {
				robotgo.Toggle(key, "down")
				handled = true
			}
		case MouseUp:
			if key = mouse[cmd.Key]; key != "" {
				robotgo.Toggle(key, "up")
				handled = true
			}
		}
	}
	return handled
}

func Initialize() {
	width, height = robotgo.GetScreenSize()
}
