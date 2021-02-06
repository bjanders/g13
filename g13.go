package g13

import (
	"bytes"
	"image/color"

	"github.com/google/gousb"
)

// USB vendor and product IDs
const (
	USBVendorLogitech = 0x046d
	USBProductG13     = 0xc21c

	USBEndpointG13Key = 1
	USBEndpointG13LCD = 2
)

// G13 keys
const (
	KeyG1 = iota
	KeyG2
	KeyG3
	KeyG4
	KeyG5
	KeyG6
	KeyG7
	KeyG8
	KeyG9
	KeyG10
	KeyG11
	KeyG12
	KeyG13
	KeyG14
	KeyG15
	KeyG16
	KeyG17
	KeyG18
	KeyG19
	KeyG20
	KeyG21
	KeyG22
	Undef1
	LightState
	KeyBD
	KeyL1
	KeyL2
	KeyL3
	KeyL4
	KeyM1
	KeyM2
	KeyM3
	KeyMR
	KeyLeft
	KeyRight
	KeyStick
	Undef2
	Light
	Light2
	Undef3
)

var Keys = []string{
	"G1",
	"G2",
	"G3",
	"G4",
	"G5",
	"G6",
	"G7",
	"G8",
	"G9",
	"G10",
	"G11",
	"G12",
	"G13",
	"G14",
	"G15",
	"G16",
	"G17",
	"G18",
	"G19",
	"G20",
	"G21",
	"G22",
	"Undef1",
	"LightState",
	"BD",
	"L1",
	"L2",
	"L3",
	"L4",
	"M1",
	"M2",
	"M3",
	"MR",
	"Left",
	"Right",
	"Stick",
	"Undef2",
	"Light",
	"Light2",
	"Undef3",
}

// G13 is the G13
type G13 struct {
	ctx        *gousb.Context
	device     *gousb.Device
	intf       *gousb.Interface
	intfDone   func()
	inEndpoint *gousb.InEndpoint
	state      [8]byte
	color      color.RGBA
}

type KeyState struct {
	Key  int
	Down bool
}

type StickState struct {
	X int
	Y int
}

// NewG13 creates a new G13
func NewG13() (*G13, error) {
	var err error

	g13 := G13{}
	g13.ctx = gousb.NewContext()
	g13.device, err = g13.ctx.OpenDeviceWithVIDPID(USBVendorLogitech, USBProductG13)
	if err != nil {
		return nil, err
	}
	g13.device.SetAutoDetach(true)
	g13.intf, g13.intfDone, err = g13.device.DefaultInterface()
	if err != nil {
		return nil, err
	}
	g13.inEndpoint, err = g13.intf.InEndpoint(USBEndpointG13Key)
	g13.color = color.RGBA{255, 255, 255, 255}
	g13.SetColor(g13.color)

	return &g13, nil
}

func (g13 *G13) WatchKeys() (chan KeyState, chan StickState, chan bool) {
	keyCh := make(chan KeyState)
	stickCh := make(chan StickState)
	backlightCh := make(chan bool)
	go readKeys(g13, keyCh, stickCh, backlightCh)
	return keyCh, stickCh, backlightCh
}

func (g13 *G13) SetModeLEDs(leds byte) {
	var data = [5]byte{5, leds}

	g13.device.Control(gousb.ControlOut|gousb.ControlClass|gousb.ControlInterface, 0x09,
		0x305, 0x00, data[:])

}

func (g13 *G13) SetMode(b, v byte) {
	var data = [5]byte{5}
	data[b] = v

	g13.device.Control(gousb.ControlOut|gousb.ControlClass|gousb.ControlInterface, 0x09,
		0x305, 0x00, data[:])

}

func (g13 *G13) SetColor(c color.Color) error {

	var rgba = color.RGBAModel.Convert(c).(color.RGBA)
	var data = [5]byte{5, rgba.R, rgba.G, rgba.B, rgba.A}

	g13.device.Control(gousb.ControlOut|gousb.ControlClass|gousb.ControlInterface, 0x09,
		0x307, 0x00, data[:])
	g13.color = rgba
	return nil
}

func (g13 *G13) Color() color.Color {
	return g13.color
}

func readKeys(g13 *G13, keyCh chan KeyState, stickCh chan StickState, backlightCh chan bool) error {
	var data [8]byte

	stream, err := g13.inEndpoint.NewStream(8, 1)
	if err != nil {
		return err
	}
	defer stream.Close()

	for {
		_, err := stream.Read(data[:])
		if err != nil {
			return err
		}
		if !bytes.Equal(data[1:3], g13.state[1:3]) {
			select {
			case stickCh <- StickState{int(data[1]), int(data[2])}:
			default:
			}
		}
		for b := 0; b < 5; b++ {
			for i := 0; i < 8; i++ {
				d := data[b+3] & (1 << i)
				s := g13.state[b+3] & (1 << i)
				if d != s {
					switch key := b*8 + i; key {
					case Undef3:
					case Light2:
					case LightState:
						select {
						case backlightCh <- d != 0:
						default:
						}
					default:
						select {
						case keyCh <- KeyState{key, d != 0}:
						default:
						}
					}
				}
			}
		}
		copy(g13.state[:], data[:])
	}
}
