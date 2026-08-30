package cpu6809

type waitMode uint8

const (
	waitNone waitMode = iota
	waitSync
	waitCWAI
)

type interruptKind uint8

const (
	interruptNone interruptKind = iota
	interruptNMI
	interruptFIRQ
	interruptIRQ
)
