package register

// SAM (Synchronous Address Multiplexer) registers ($FFC0-$FFDF).
// The SAM chip controls memory mapping, display mode, and CPU clock speed.
// Each function uses a pair of addresses: clear (even) and set (odd).
const (
	// Display mode select bits (V0-V2)
	SAMV0Clear = 0xFFC0
	SAMV0Set   = 0xFFC1
	SAMV1Clear = 0xFFC2
	SAMV1Set   = 0xFFC3
	SAMV2Clear = 0xFFC4
	SAMV2Set   = 0xFFC5

	// Display offset address bits (F0-F6)
	SAMF0Clear = 0xFFC6
	SAMF0Set   = 0xFFC7
	SAMF1Clear = 0xFFC8
	SAMF1Set   = 0xFFC9
	SAMF2Clear = 0xFFCA
	SAMF2Set   = 0xFFCB
	SAMF3Clear = 0xFFCC
	SAMF3Set   = 0xFFCD
	SAMF4Clear = 0xFFCE
	SAMF4Set   = 0xFFCF
	SAMF5Clear = 0xFFD0
	SAMF5Set   = 0xFFD1
	SAMF6Clear = 0xFFD2
	SAMF6Set   = 0xFFD3

	// Page select bit (P1)
	SAMP1Clear = 0xFFD4
	SAMP1Set   = 0xFFD5

	// RAM size select bits (R0-R1)
	SAMR0Clear = 0xFFD6
	SAMR0Set   = 0xFFD7
	SAMR1Clear = 0xFFD8
	SAMR1Set   = 0xFFD9

	// Memory map type (M0-M1)
	SAMM0Clear = 0xFFDA
	SAMM0Set   = 0xFFDB
	SAMM1Clear = 0xFFDC
	SAMM1Set   = 0xFFDD

	// CPU rate select
	SAMRateClear = 0xFFDE // Normal speed (0.895 MHz)
	SAMRateSet   = 0xFFDF // Double speed (1.79 MHz)
)
