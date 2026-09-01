package cartridge

// ControlBytes returns the 2 control bytes of the iNES header based on the cartridge configuration.
func ControlBytes(battery, mirror byte, mapper uint16, hasTrainer bool) (byte, byte) {
	mapper8 := byte(mapper)

	var control1, control2 byte
	if battery&1 != 0 {
		control1 |= batteryFlag
	}

	// The iNES header cannot encode mapper-controlled single-screen modes.
	switch MirrorMode(mirror) {
	case MirrorVertical:
		control1 |= verticalMirroringFlag
	case Mirror4:
		control1 |= fourScreenFlag
	}

	control1 |= mergeNibbles(mapper8, control1)
	control2 |= mergeNibbles(highNibble(mapper8), control2)

	if hasTrainer {
		control1 |= trainerFlag
	}
	return control1, control2
}

func highNibble(b byte) byte {
	return b >> 4
}

func mergeNibbles(highNibble byte, lowNibble byte) byte {
	return highNibble<<4 | lowNibble
}
