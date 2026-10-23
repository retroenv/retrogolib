// Package knownbits describes fixed-width integers with known-zero and known-one bits.
//
// Value stores no width or identity. Callers supply the active width for each
// operation. Equal facts can describe different concrete values. The zero Value
// describes an unknown integer at any valid width. Join keeps facts common to
// both inputs; it does not combine constraints on one concrete value.
//
// Constructors clear bits outside the active width. Other operations reject
// noncanonical inputs. A false success result means an invalid width, fact,
// bit, or slice. A successful unknown result is valid. Callers must check success
// before they use a result. Result values are unspecified on failure.
//
// Arithmetic is modular. Carry and borrow are mathematical events. CPU adapters
// must map these events to their own flag conventions. This package has no CPU,
// control-flow, memory, or solver policy.
package knownbits
