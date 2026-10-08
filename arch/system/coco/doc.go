// Package coco provides constants and definitions for the TRS-80 Color Computer (CoCo) system.
//
// The CoCo uses a Motorola 6809E processor with a SAM (Synchronous Address Multiplexer)
// chip for memory management and a PIA (MC6821 Peripheral Interface Adapter) pair for
// I/O handling. The system has 4-64 KB of RAM depending on model.
//
// Memory map (64 KB address space):
//
//	$0000-$7FFF  RAM (32 KB standard, 64 KB extended)
//	$8000-$9FFF  Extended BASIC ROM (8 KB)
//	$A000-$BFFF  Color BASIC ROM (8 KB)
//	$C000-$FEFF  Cartridge ROM space (16 KB minus I/O)
//	$FF00-$FF03  PIA 0 (keyboard, joystick, cassette)
//	$FF04-$FF1F  PIA 0 mirrors
//	$FF20-$FF23  PIA 1 (VDG control, serial, sound)
//	$FF24-$FF3F  PIA 1 mirrors
//	$FF40-$FF5F  Floppy disk controller
//	$FF90-$FFBF  GIME registers (CoCo 3 only)
//	$FFC0-$FFDF  SAM registers
//	$FFF0-$FFFF  Interrupt vectors
package coco
