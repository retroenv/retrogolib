package register

// RIOT (6532) registers ($0280-$0297).
// The RIOT provides 128 bytes of RAM, two 8-bit I/O ports, and a programmable timer.
const (
	SWCHA  = 0x0280 // Port A: joystick directions (read/write)
	SWACNT = 0x0281 // Port A DDR (data direction register)
	SWCHB  = 0x0282 // Port B: console switches (read)
	SWBCNT = 0x0283 // Port B DDR
	INTIM  = 0x0284 // Timer output (read)
	INSTAT = 0x0285 // Timer interrupt status (read)
	TIM1T  = 0x0294 // Set 1-clock interval timer (write)
	TIM8T  = 0x0295 // Set 8-clock interval timer (write)
	TIM64T = 0x0296 // Set 64-clock interval timer (write)
	T1024T = 0x0297 // Set 1024-clock interval timer (write)
)

// Timer interval prescaler values (in CPU clock cycles per tick).
const (
	TimerInterval1    = 1
	TimerInterval8    = 8
	TimerInterval64   = 64
	TimerInterval1024 = 1024
)

// Console switch bits in SWCHB.
const (
	SwitchReset  = 0x01 // Reset button (active low)
	SwitchSelect = 0x02 // Select button (active low)
	SwitchBW     = 0x08 // B/W-Color switch (0=B/W, 1=Color)
	SwitchP0Diff = 0x40 // Player 0 difficulty (0=B/expert, 1=A/novice)
	SwitchP1Diff = 0x80 // Player 1 difficulty (0=B/expert, 1=A/novice)
)

// Joystick direction bits in SWCHA.
// Each player uses 4 bits (active low).
const (
	Joy0Right = 0x80 // Player 0 right
	Joy0Left  = 0x40 // Player 0 left
	Joy0Down  = 0x20 // Player 0 down
	Joy0Up    = 0x10 // Player 0 up
	Joy1Right = 0x08 // Player 1 right
	Joy1Left  = 0x04 // Player 1 left
	Joy1Down  = 0x02 // Player 1 down
	Joy1Up    = 0x01 // Player 1 up
)

// RIOTNames maps RIOT register addresses to their names.
var RIOTNames = map[uint16]string{
	SWCHA:  "SWCHA",
	SWACNT: "SWACNT",
	SWCHB:  "SWCHB",
	SWBCNT: "SWBCNT",
	INTIM:  "INTIM",
	INSTAT: "INSTAT",
	TIM1T:  "TIM1T",
	TIM8T:  "TIM8T",
	TIM64T: "TIM64T",
	T1024T: "T1024T",
}
