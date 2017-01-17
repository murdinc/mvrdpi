package display

import (
	"bitbucket.org/gmcbay/i2c"
)

const (
	DISPLAYOFF          = 0xAE
	DISPLAYON           = 0xAF
	SETMULTIPLEX        = 0xA8
	SETREMAP            = 0xA1 // 0xA0
	DISPLAYALLON_RESUME = 0xA4
	NORMALDISPLAY       = 0xA6

	CHARGEPUMP         = 0x8D
	COLUMNADDR         = 0x21
	COMSCANDEC         = 0xC8
	COMSCANINC         = 0xC0
	EXTERNALVCC        = 0x1
	MEMORYMODE         = 0x20
	PAGEADDR           = 0x22
	SETCOMPINS         = 0xDA
	SETDISPLAYCLOCKDIV = 0xD5
	SETDISPLAYOFFSET   = 0xD3
	SETHIGHCOLUMN      = 0x10
	SETLOWCOLUMN       = 0x00
	SETPRECHARGE       = 0xD9
	SETSEGMENTREMAP    = 0xA1
	SETSTARTLINE       = 0x40
	SETVCOMDETECT      = 0xDB
	SWITCHCAPVCC       = 0x2
	SETCONTRAST        = 0x81

	// Scrolling constants
	ACTIVATE_SCROLL                      = 0x2F
	DEACTIVATE_SCROLL                    = 0x2E
	SET_VERTICAL_SCROLL_AREA             = 0xA3
	RIGHT_HORIZONTAL_SCROLL              = 0x26
	LEFT_HORIZONTAL_SCROLL               = 0x27
	VERTICAL_AND_RIGHT_HORIZONTAL_SCROLL = 0x29
	VERTICAL_AND_LEFT_HORIZONTAL_SCROLL  = 0x2A
)

type bitmap struct {
	cols        int
	rows        int
	bytesPerCol int
	data        []byte
}

func (b *bitmap) Init(cols int, rows int) {
	b.cols = cols
	b.rows = rows
	b.bytesPerCol = rows / 8
	b.data = make([]byte, cols*b.bytesPerCol)
}

func (b *bitmap) Clear() {
	for i, _ := range b.data {
		b.data[i] = 0
	}
}

func (b *bitmap) DrawPixel(x int, y int, on bool) {
	if x < 0 || x >= b.cols || y < 0 || y >= b.rows {
		return
	}
	memCol := x
	memRow := y / 8
	bitMask := 1 << (uint(y) % 8)
	offset := memRow + (b.rows / 8 * memCol)
	if on {
		b.data[offset] |= byte(bitMask)
	} else {
		b.data[offset] &= byte((0xff - bitMask))
	}
}

func (b *bitmap) ClearBlock(x0 int, y0 int, dx int, dy int) {
	xLength := x0 + dx
	yLength := y0 + dy
	for x := x0; x < xLength; x++ {
		for y := y0; y < yLength; y++ {
			b.DrawPixel(x, y, false)
		}
	}
}

type SSD1306 struct {
	bus    *i2c.I2CBus
	bmap   *bitmap
	addr   byte
	height int
	width  int
	pages  int
}

func NewDevice() *SSD1306 {
	return &SSD1306{}
}

func (l *SSD1306) Init(busNumber byte, addr byte, height int, width int) error {
	var err error
	l.bus, err = i2c.Bus(busNumber)
	l.bmap = &bitmap{}
	l.bmap.Init(width, height)
	l.addr = addr
	l.height = height
	l.width = width
	l.pages = height / 8
	return err
}

func (l *SSD1306) command(c int) {
	control := 0x00 // Co = 0, DC = 0
	l.bus.WriteByte(l.addr, byte(control), byte(c))
}

func (l *SSD1306) InitDevice() {

	l.command(DISPLAYOFF)
	l.command(SETDISPLAYCLOCKDIV)
	l.command(0x80)

	l.command(SETMULTIPLEX)
	l.command(0x3F)

	l.command(SETDISPLAYOFFSET)
	l.command(0x00)

	l.command(SETSTARTLINE)
	l.command(CHARGEPUMP)
	l.command(0x14)

	l.command(MEMORYMODE)
	l.command(0x00)

	l.command(SETREMAP)
	l.command(COMSCANDEC)
	//l.command(COMSCANINC)
	l.command(SETCOMPINS)
	l.command(0x12)

	l.command(SETPRECHARGE)
	l.command(0xF1)

	l.command(SETVCOMDETECT)
	l.command(0x40)

	l.command(DISPLAYALLON_RESUME)
	l.command(NORMALDISPLAY)

	l.command(SETCONTRAST)
	l.command(0xCF)

	l.Clear()
	l.Display()
	l.command(DISPLAYON)

}

func (s *SSD1306) Clear() {
	s.bmap.Clear()
	s.Display()
}

func (l *SSD1306) data(bytes []byte) {
	control := 0x40 // Co = 0, DC = 0
	l.bus.WriteByteBlock(l.addr, byte(control), bytes)
}

func (s *SSD1306) Display() {
	// s.displayBlock(0, 0, s.width, 0)

	s.command(COLUMNADDR)
	s.command(0)           // Column start address. (0 = reset)
	s.command(s.width - 1) // Column end address.
	s.command(PAGEADDR)
	s.command(0)           // Page start address. (0 = reset)
	s.command(s.pages - 1) // Page end address.
	length := len(s.bmap.data)
	for i := 0; i < length; i += 16 {
		s.data(s.bmap.data[i : i+16])
	}
}

func (s *SSD1306) DeactivateScroll() {
	s.command(DEACTIVATE_SCROLL)
}

func (s *SSD1306) ActivateScroll() {
	s.command(ACTIVATE_SCROLL)
}

func (s *SSD1306) GetPages() int {
	return s.pages
}

func (s *SSD1306) WriteData(d byte, pos int) {
	s.bmap.data[pos] = d
}

func (s *SSD1306) SetStartLine(pos int) {
	s.command(SETSTARTLINE | pos)
}

func (s *SSD1306) SetAndActiveScroll(speed int) {

}
