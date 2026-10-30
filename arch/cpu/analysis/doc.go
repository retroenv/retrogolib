// Package analysis defines local machine-transfer records and logical register views.
//
// Targets own resource IDs and address types. Storage ID zero is valid. Bit offsets
// do not depend on memory byte order. Callers own state, memory policy, aliases,
// joins, caches, and the slices in Result. A transfer must read its input snapshot
// before it publishes writes. Apply all invalidations before installing new facts.
//
// EffectsComplete describes the target's declared data-resource scope and local
// instruction memory effects. It does not describe a callee or interrupt handler.
// If it is false, use a conservative barrier. Otherwise invalidate every write,
// including unknown writes, before using facts. An unknown memory address or zero
// byte count requires a conservative memory update. A known value never removes
// a volatile event. Memory events describe semantic accesses, not bus cycles.
//
// Result slices belong to the caller. Reset and the next transfer reuse their
// storage. Copy a slice when a result must remain independent of later transfers.
package analysis
