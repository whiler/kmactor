package kmactor

import "github.com/go-vgo/robotgo"

const (
	DeviceKey   = 1
	DeviceMouse = 2

	KeyDown  = 1
	KeyPress = 2
	KeyUp    = 3

	MouseMove = 1
	MouseDown = 2
	MouseUp   = 3
)

var SupportedKeys = map[string]bool{
	"a": true, "b": true, "c": true, "d": true, "e": true, "f": true, "g": true, "h": true, "i": true, "j": true, "k": true, "l": true, "m": true, "n": true, "o": true, "p": true, "q": true, "r": true, "s": true, "t": true, "u": true, "v": true, "w": true, "x": true, "y": true, "z": true,
	"A": true, "B": true, "C": true, "D": true, "E": true, "F": true, "G": true, "H": true, "I": true, "J": true, "K": true, "L": true, "M": true, "N": true, "O": true, "P": true, "Q": true, "R": true, "S": true, "T": true, "U": true, "V": true, "W": true, "X": true, "Y": true, "Z": true,

	"~": true, "`": true,
	"!": true, "1": true,
	"@": true, "2": true,
	"#": true, "3": true,
	"$": true, "4": true,
	"%": true, "5": true,
	"^": true, "6": true,
	"&": true, "7": true,
	"*": true, "8": true,
	"(": true, "9": true,
	")": true, "0": true,
	"_": true, "-": true,
	"+": true, "=": true,
	"{": true, "[": true,
	"}": true, "]": true,
	"|": true, "\\": true,
	":": true, ";": true,
	"\"": true, "'": true,
	"<": true, ",": true,
	">": true, ".": true,
	"?": true, "/": true,

	"f1":  true,
	"f2":  true,
	"f3":  true,
	"f4":  true,
	"f5":  true,
	"f6":  true,
	"f7":  true,
	"f8":  true,
	"f9":  true,
	"f10": true,
	"f11": true,
	"f12": true,

	"esc":       true,
	"tab":       true,
	"capslock":  true,
	"shift":     true,
	"ctrl":      true,
	"alt":       true,
	"space":     true,
	"backspace": true,
	"enter":     true,

	"delete":   true,
	"end":      true,
	"home":     true,
	"pagedown": true,
	"pageup":   true,

	"up":    true,
	"down":  true,
	"left":  true,
	"right": true,
}

var SupportedMouses = map[string]bool{
	"left":   true,
	"center": true,
	"right":  true,
}

type Position struct {
	Width  int `json:"w"`
	Height int `json:"h"`
	Left   int `json:"l"`
	Top    int `json:"t"`
}

type Command struct {
	Device   int       `json:"d"`
	Action   int       `json:"a"`
	Key      string    `json:"k"`
	Position *Position `json:"p"`
}

func (self *Command) Reset() {
	self.Device = 0
	self.Action = 0
}

func Play(cmd *Command, width, height int) bool {
	var handled bool
	if cmd != nil {
		switch cmd.Device {
		case DeviceKey:
			switch cmd.Action {
			case KeyDown:
				if SupportedKeys[cmd.Key] {
					robotgo.KeyDown(cmd.Key)
					handled = true
				}
			case KeyPress:
				if SupportedKeys[cmd.Key] {
					robotgo.KeyPress(cmd.Key)
					handled = true
				}
			case KeyUp:
				if SupportedKeys[cmd.Key] {
					robotgo.KeyUp(cmd.Key)
					handled = true
				}
			}
		case DeviceMouse:
			switch cmd.Action {
			case MouseMove:
				if cmd.Position != nil && cmd.Position.Left >= 0 && width > 0 && cmd.Position.Top >= 0 && height > 0 {
					robotgo.Move(cmd.Position.Left*width/cmd.Position.Width, cmd.Position.Top*height/cmd.Position.Height)
					handled = true
				}
			case MouseDown:
				if SupportedMouses[cmd.Key] {
					robotgo.Toggle(cmd.Key, "down")
					handled = true
				}
			case MouseUp:
				if SupportedMouses[cmd.Key] {
					robotgo.Toggle(cmd.Key, "up")
					handled = true
				}
			}
		}
	}
	return handled
}

func GetScreenSize() (int, int) {
	return robotgo.GetScreenSize()
}
