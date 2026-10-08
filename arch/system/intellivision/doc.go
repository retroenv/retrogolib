// Package intellivision defines the base Intellivision memory map.
//
// The CP1610 has 65,536 word addresses. A device can supply fewer than 16 data
// bits at each address. This package gives the primary device ranges. Cartridge
// ranges depend on the cartridge. Some device addresses also have mirrors, and
// access to STIC and graphics memory depends on the display phase.
package intellivision
