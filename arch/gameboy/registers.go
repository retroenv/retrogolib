package gameboy

// I/O Register addresses
const (
	// Joypad
	RegJOYP = 0xFF00 // Joypad (P1)

	// Serial Data Transfer
	RegSB = 0xFF01 // Serial transfer data (SB)
	RegSC = 0xFF02 // Serial transfer control (SC)

	// Timer
	RegDIV  = 0xFF04 // Divider Register (DIV)
	RegTIMA = 0xFF05 // Timer counter (TIMA)
	RegTMA  = 0xFF06 // Timer Modulo (TMA)
	RegTAC  = 0xFF07 // Timer Control (TAC)

	// Interrupts
	RegIF = 0xFF0F // Interrupt Flag (IF)
	RegIE = 0xFFFF // Interrupt Enable (IE)

	// Sound
	RegNR10 = 0xFF10 // Channel 1 Sweep register (NR10)
	RegNR11 = 0xFF11 // Channel 1 Sound length/Wave pattern duty (NR11)
	RegNR12 = 0xFF12 // Channel 1 Volume Envelope (NR12)
	RegNR13 = 0xFF13 // Channel 1 Frequency lo (NR13)
	RegNR14 = 0xFF14 // Channel 1 Frequency hi (NR14)

	RegNR21 = 0xFF16 // Channel 2 Sound Length/Wave Pattern Duty (NR21)
	RegNR22 = 0xFF17 // Channel 2 Volume Envelope (NR22)
	RegNR23 = 0xFF18 // Channel 2 Frequency lo data (NR23)
	RegNR24 = 0xFF19 // Channel 2 Frequency hi data (NR24)

	RegNR30 = 0xFF1A // Channel 3 Sound on/off (NR30)
	RegNR31 = 0xFF1B // Channel 3 Sound Length
	RegNR32 = 0xFF1C // Channel 3 Select output level (NR32)
	RegNR33 = 0xFF1D // Channel 3 Frequency's lower data (NR33)
	RegNR34 = 0xFF1E // Channel 3 Frequency's higher data (NR34)

	RegNR41 = 0xFF20 // Channel 4 Sound Length (NR41)
	RegNR42 = 0xFF21 // Channel 4 Volume Envelope (NR42)
	RegNR43 = 0xFF22 // Channel 4 Polynomial Counter (NR43)
	RegNR44 = 0xFF23 // Channel 4 Counter/consecutive; Initial (NR44)

	RegNR50 = 0xFF24 // Channel control / ON-OFF / Volume (NR50)
	RegNR51 = 0xFF25 // Selection of Sound output terminal (NR51)
	RegNR52 = 0xFF26 // Sound on/off (NR52)

	// Wave Pattern RAM
	WavePatternRAMStart = 0xFF30
	WavePatternRAMEnd   = 0xFF3F

	// LCD Display
	RegLCDC = 0xFF40 // LCD Control (LCDC)
	RegSTAT = 0xFF41 // LCDC Status (STAT)
	RegSCY  = 0xFF42 // Scroll Y (SCY)
	RegSCX  = 0xFF43 // Scroll X (SCX)
	RegLY   = 0xFF44 // LCDC Y-Coordinate (LY)
	RegLYC  = 0xFF45 // LY Compare (LYC)
	RegDMA  = 0xFF46 // DMA Transfer and Start Address (DMA)
	RegBGP  = 0xFF47 // BG Palette Data (BGP)
	RegOBP0 = 0xFF48 // Object Palette 0 Data (OBP0)
	RegOBP1 = 0xFF49 // Object Palette 1 Data (OBP1)
	RegWY   = 0xFF4A // Window Y Position (WY)
	RegWX   = 0xFF4B // Window X Position minus 7 (WX)

	// Boot ROM disable
	RegBootROMDisable = 0xFF50 // Boot ROM disable
)

// Interrupt flags
const (
	InterruptVBlank = 1 << 0 // V-Blank Interrupt
	InterruptLCDC   = 1 << 1 // LCDC Status Interrupt
	InterruptTimer  = 1 << 2 // Timer Interrupt
	InterruptSerial = 1 << 3 // Serial Interrupt
	InterruptJoypad = 1 << 4 // Joypad Interrupt
)

// Interrupt vectors
const (
	VectorVBlank = 0x0040 // V-Blank Interrupt vector
	VectorLCDC   = 0x0048 // LCDC Status Interrupt vector
	VectorTimer  = 0x0050 // Timer Interrupt vector
	VectorSerial = 0x0058 // Serial Interrupt vector
	VectorJoypad = 0x0060 // Joypad Interrupt vector
)

// LCDC (LCD Control) register bits
const (
	LCDCBGDisplay     = 1 << 0 // BG Display (0=Off, 1=On)
	LCDCOBJDisplay    = 1 << 1 // OBJ (Sprite) Display (0=Off, 1=On)
	LCDCOBJSize       = 1 << 2 // OBJ (Sprite) Size (0=8x8, 1=8x16)
	LCDCBGTileMap     = 1 << 3 // BG Tile Map Display Select (0=9800-9BFF, 1=9C00-9FFF)
	LCDCBGWindowTiles = 1 << 4 // BG & Window Tile Data Select (0=8800-97FF, 1=8000-8FFF)
	LCDCWindowDisplay = 1 << 5 // Window Display (0=Off, 1=On)
	LCDCWindowTileMap = 1 << 6 // Window Tile Map Display Select (0=9800-9BFF, 1=9C00-9FFF)
	LCDCDisplayEnable = 1 << 7 // LCD Display Enable (0=Off, 1=On)
)

// STAT (LCDC Status) register bits
const (
	STATMode0       = 0 << 0 // Mode Flag (Bit 1-0)
	STATMode1       = 1 << 0 // 00: HBlank
	STATMode2       = 2 << 0 // 01: VBlank
	STATMode3       = 3 << 0 // 10: Searching Sprites Atts
	STATModeMask    = 3 << 0 // 11: Transferring Data to LCD Driver
	STATCoincidence = 1 << 2 // Coincidence Flag (0:LYC<>LY, 1:LYC=LY)
	STATHBlankInt   = 1 << 3 // Mode 0 HBlank Interrupt (1=Enable)
	STATVBlankInt   = 1 << 4 // Mode 1 VBlank Interrupt (1=Enable)
	STATOAMInt      = 1 << 5 // Mode 2 OAM Interrupt (1=Enable)
	STATLYCInt      = 1 << 6 // LYC=LY Coincidence Interrupt (1=Enable)
)

// Timer Control (TAC) register bits
const (
	TACClockSelect = 3 << 0 // Input Clock Select
	TACClock4096   = 0 << 0 // 4096 Hz
	TACClock262144 = 1 << 0 // 262144 Hz
	TACClock65536  = 2 << 0 // 65536 Hz
	TACClock16384  = 3 << 0 // 16384 Hz
	TACTimerStop   = 0 << 2 // Timer Stop (0=Stop, 1=Start)
	TACTimerStart  = 1 << 2
)

// Joypad (P1) register bits
const (
	JOYPRight  = 1 << 0 // Right or A
	JOYPLeft   = 1 << 1 // Left or B
	JOYPUp     = 1 << 2 // Up or Select
	JOYPDown   = 1 << 3 // Down or Start
	JOYPAction = 1 << 4 // Action buttons (A, B, Select, Start)
	JOYPDir    = 1 << 5 // Direction buttons (Right, Left, Up, Down)
)

// Serial Control (SC) register bits
const (
	SCShiftClock    = 1 << 0 // Shift Clock (0=External Clock, 1=Internal Clock)
	SCClockSpeed    = 1 << 1 // Clock Speed (0=Normal, 1=Fast) - CGB Mode Only
	SCTransferStart = 1 << 7 // Transfer Start Flag (0=No transfer, 1=Start)
)

// Sound register bits
const (
	// NR10 - Channel 1 Sweep
	NR10SweepShift = 7 << 0 // Number of sweep shift (n: 0-7)
	NR10SweepUp    = 0 << 3 // Sweep Increase/Decrease (0=Addition, 1=Subtraction)
	NR10SweepDown  = 1 << 3
	NR10SweepTime  = 7 << 4 // Sweep Time (0-7)

	// NR11/NR21 - Sound Length/Wave Pattern
	NR11WavePattern = 3 << 6 // Wave Pattern Duty (0-3)
	NR11SoundLength = 63     // Sound length data (t1: 0-63)

	// NR12/NR22 - Volume Envelope
	NR12EnvSweep     = 7 << 0  // Number of envelope sweep (n: 0-7)
	NR12EnvDirection = 1 << 3  // Envelope Direction (0=Decrease, 1=Increase)
	NR12InitVolume   = 15 << 4 // Initial Volume of envelope (0-15)

	// NR14/NR24 - Frequency Hi
	NR14Initial   = 1 << 7 // Initial (1=Restart Sound)
	NR14Counter   = 1 << 6 // Counter/consecutive selection
	NR14Frequency = 7 << 0 // Frequency's higher 3 bits

	// NR30 - Channel 3 Sound on/off
	NR30SoundEnable = 1 << 7 // Sound Channel 3 off (0=Stop, 1=Playback)

	// NR32 - Channel 3 Select output level
	NR32OutputLevel = 3 << 5 // Select output level (0-3)

	// NR43 - Channel 4 Polynomial Counter
	NR43DividingRatio  = 7 << 0  // Dividing Ratio of Frequencies (r)
	NR43CounterStep    = 1 << 3  // Counter Step/Width (0=15 bits, 1=7 bits)
	NR43ShiftClockFreq = 15 << 4 // Shift Clock Frequency (s)

	// NR50 - Channel control / ON-OFF / Volume
	NR50SO1Volume = 7 << 0 // SO1 output level (volume) (0-7)
	NR50VinSO1    = 1 << 3 // Vin→SO1 ON/OFF
	NR50SO2Volume = 7 << 4 // SO2 output level (volume) (0-7)
	NR50VinSO2    = 1 << 7 // Vin→SO2 ON/OFF

	// NR51 - Selection of Sound output terminal
	NR51Sound1SO1 = 1 << 0 // Sound 1 → SO1 ON/OFF
	NR51Sound2SO1 = 1 << 1 // Sound 2 → SO1 ON/OFF
	NR51Sound3SO1 = 1 << 2 // Sound 3 → SO1 ON/OFF
	NR51Sound4SO1 = 1 << 3 // Sound 4 → SO1 ON/OFF
	NR51Sound1SO2 = 1 << 4 // Sound 1 → SO2 ON/OFF
	NR51Sound2SO2 = 1 << 5 // Sound 2 → SO2 ON/OFF
	NR51Sound3SO2 = 1 << 6 // Sound 3 → SO2 ON/OFF
	NR51Sound4SO2 = 1 << 7 // Sound 4 → SO2 ON/OFF

	// NR52 - Sound on/off
	NR52Sound1ON   = 1 << 0 // Sound 1 ON flag (Read Only)
	NR52Sound2ON   = 1 << 1 // Sound 2 ON flag (Read Only)
	NR52Sound3ON   = 1 << 2 // Sound 3 ON flag (Read Only)
	NR52Sound4ON   = 1 << 3 // Sound 4 ON flag (Read Only)
	NR52AllSoundON = 1 << 7 // All sound on/off (0=Stop, 1=Operate)
)

// OAMEntry represents an Object Attribute Memory structure.
type OAMEntry struct {
	Y      uint8 // Y Position
	X      uint8 // X Position
	TileID uint8 // Tile/Pattern Number
	Attr   uint8 // Attributes/Flags
}

// OAM Attribute bits
const (
	OAMAttrPalette  = 1 << 4 // Palette number (0=OBP0, 1=OBP1)
	OAMAttrXFlip    = 1 << 5 // X flip (0=Normal, 1=Horizontally mirrored)
	OAMAttrYFlip    = 1 << 6 // Y flip (0=Normal, 1=Vertically mirrored)
	OAMAttrPriority = 1 << 7 // BG and Window over OBJ (0=No, 1=BG and Window colors 1-3 over the OBJ)
)

// Screen dimensions
const (
	ScreenWidth  = 160 // Screen width in pixels
	ScreenHeight = 144 // Screen height in pixels
	TileSize     = 8   // Size of a tile in pixels
	TilesPerRow  = 32  // Number of tiles per row in tilemap
	TilesPerCol  = 32  // Number of tiles per column in tilemap
)

// Timing constants
const (
	CPUFrequency   = 4194304 // CPU frequency in Hz (4.194304 MHz)
	CyclesPerFrame = 70224   // CPU cycles per frame
	LinesPerFrame  = 154     // Total lines per frame (144 visible + 10 vblank)
	CyclesPerLine  = 456     // CPU cycles per line
	VBlankLines    = 10      // Number of VBlank lines
)

// LCD mode timings (in CPU cycles)
const (
	OAMSearchCycles     = 80  // Mode 2: OAM Search
	PixelTransferCycles = 172 // Mode 3: Pixel Transfer (minimum)
	HBlankCycles        = 204 // Mode 0: HBlank (minimum)
)

// Palette colors (for DMG)
const (
	ColorWhite     = 0 // Color 0 (lightest)
	ColorLightGray = 1 // Color 1
	ColorDarkGray  = 2 // Color 2
	ColorBlack     = 3 // Color 3 (darkest)
)
