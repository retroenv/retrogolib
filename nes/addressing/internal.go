package addressing

// internal types
type (
	// AbsoluteX defines absolute addressing using the X register
	AbsoluteX uint16
	// AbsoluteY defines absolute addressing using the Y register
	AbsoluteY uint16
	// IndirectX defines indirect addressing using the X register
	IndirectX uint16
	// IndirectY defines indirect addressing using the Y register
	IndirectY uint16
	// ZeroPageX defines zeropage addressing using the X register
	ZeroPageX uint8
	// ZeroPageY defines zeropage addressing using the Y register
	ZeroPageY uint8
)
