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
	switch cmd.Device {
	case DeviceKey:
		switch cmd.Action {
		case KeyDown:
			robotgo.KeyDown(cmd.Key)
			handled = true
		case KeyPress:
			robotgo.KeyPress(cmd.Key)
			handled = true
		case KeyUp:
			robotgo.KeyUp(cmd.Key)
			handled = true
		}
	case DeviceMouse:
		switch cmd.Action {
		case MouseMove:
			robotgo.Move(cmd.Position.Left*width/cmd.Position.Width, cmd.Position.Top*height/cmd.Position.Height)
			handled = true
		case MouseDown:
			robotgo.Toggle(cmd.Key, "down")
			handled = true
		case MouseUp:
			robotgo.Toggle(cmd.Key, "up")
			handled = true
		}
	}
	return handled
}

func GetScreenSize() (int, int) {
	return robotgo.GetScreenSize()
}
