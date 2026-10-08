package intellivision

// AddressSpaceWords is the number of CP1610 word addresses.
const AddressSpaceWords = 0x10000

// OpenBusValue is the value read when no device drives the data bus.
const OpenBusValue = 0xFFFF

// Primary device address ranges on the base Intellivision.
const (
	// STICStart is the first STIC register address.
	STICStart = 0x0000
	// STICEnd is the last primary STIC register address.
	STICEnd = 0x003F

	// ScratchpadRAMStart is the first scratchpad RAM address.
	ScratchpadRAMStart = 0x0100
	// ScratchpadRAMEnd is the last scratchpad RAM address.
	ScratchpadRAMEnd = 0x01EF

	// PSGStart is the first sound generator register address.
	PSGStart = 0x01F0
	// PSGEnd is the last sound generator register address.
	PSGEnd = 0x01FF

	// SystemRAMStart is the first system RAM address.
	SystemRAMStart = 0x0200
	// SystemRAMEnd is the last system RAM address.
	SystemRAMEnd = 0x035F

	// ExecutiveROMStart is the first Executive ROM address.
	ExecutiveROMStart = 0x1000
	// ExecutiveROMEnd is the last Executive ROM address.
	ExecutiveROMEnd = 0x1FFF

	// GraphicsROMStart is the first graphics ROM address.
	GraphicsROMStart = 0x3000
	// GraphicsROMEnd is the last graphics ROM address.
	GraphicsROMEnd = 0x37FF

	// GraphicsRAMStart is the first primary graphics RAM address.
	GraphicsRAMStart = 0x3800
	// GraphicsRAMEnd is the last primary graphics RAM address.
	GraphicsRAMEnd = 0x39FF
)

// Executive ROM entry addresses.
const (
	// ResetAddress is the address where the CP1610 starts after reset.
	ResetAddress = 0x1000
	// InterruptAddress is the entry address for the STIC interrupt.
	InterruptAddress = 0x1004
)
