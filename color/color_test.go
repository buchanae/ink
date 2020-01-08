package color

import "testing"

func TestHex(t *testing.T) {
	var tests = []struct {
		hex    int
		expect RGBA
	}{
		{0xff0000, RGBA{1, 0, 0, 1}},
		{0x00ff00, RGBA{0, 1, 0, 1}},
		{0x0000ff, RGBA{0, 0, 1, 1}},
		{0xff, RGBA{0, 0, 1, 1}},
		{0xff000000, RGBA{0, 0, 0, 1}},
		{0xffffff, RGBA{1, 1, 1, 1}},
		{0xffff00, RGBA{1, 1, 0, 1}},
	}

	for i, test := range tests {
		v := Hex(test.hex)
		if v != test.expect {
			t.Errorf("%d: given %x wanted %v, got %v",
				i, test.hex, test.expect, v)
		}
	}
}

func TestHexString(t *testing.T) {
	var tests = []struct {
		hex    string
		expect RGBA
	}{
		{"#ff0000", RGBA{1, 0, 0, 1}},
		{"#00ff00", RGBA{0, 1, 0, 1}},
		{"#0000ff", RGBA{0, 0, 1, 1}},
		{"#ffffff", RGBA{1, 1, 1, 1}},
		{"#ffff00", RGBA{1, 1, 0, 1}},
		// invalid
		{"#ff", RGBA{0, 0, 0, 0}},
		{"#ff000000", RGBA{0, 0, 0, 0}},
	}

	for i, test := range tests {
		v := HexString(test.hex)
		if v != test.expect {
			t.Errorf("%d: given %s wanted %v, got %v",
				i, test.hex, test.expect, v)
		}
	}
}
