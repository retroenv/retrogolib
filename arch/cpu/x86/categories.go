package x86

import "github.com/retroenv/retrogolib/set"

// BranchingInstructions contains all instructions that can change control flow.
// Used by disassemblers and static analysis tools to identify branch points.
var BranchingInstructions = set.NewFromSlice([]string{
	CallName,
	IntName,
	IntoName,
	IretName,
	JbName,
	JbeName,
	JlName,
	JleName,
	JmpName,
	JnbName,
	JnbeName,
	JnlName,
	JnleName,
	JnoName,
	JnpName,
	JnsName,
	JnzName,
	JoName,
	JpName,
	JsName,
	JzName,
	RetName,
	RetfName,
})

// NotExecutingFollowingOpcodeInstructions contains instructions that never
// execute the following opcode (unconditional jumps, returns, halts).
// Used by disassemblers to identify basic block boundaries.
var NotExecutingFollowingOpcodeInstructions = set.NewFromSlice([]string{
	HltName,
	IretName,
	JmpName,
	RetName,
	RetfName,
})

// MemoryReadInstructions contains instructions that read from memory.
// Used for memory access analysis and data flow tracking.
var MemoryReadInstructions = set.NewFromSlice([]string{
	AdcName,
	AddName,
	AndName,
	BoundName,
	BsfName,
	BsrName,
	BtName,
	BtcName,
	BtrName,
	BtsName,
	CmpName,
	CmpsbName,
	CmpswName,
	CmpxchgName,
	DecName,
	DivName,
	IdivName,
	ImulName,
	IncName,
	IretName,
	LeaName,
	LodsbName,
	LodswName,
	MovName,
	MovsbName,
	MovswName,
	MovsxName,
	MovzxName,
	MulName,
	OrName,
	PopName,
	PopaName,
	RclName,
	RcrName,
	RetName,
	RetfName,
	RolName,
	RorName,
	SarName,
	SbbName,
	ScasbName,
	ScaswName,
	ShldName,
	ShlName,
	ShrdName,
	ShrName,
	SubName,
	TestName,
	XaddName,
	XchgName,
	XlatName,
	XorName,
})

// MemoryWriteInstructions contains instructions that write to memory.
// Used for memory access analysis and side effect tracking.
var MemoryWriteInstructions = set.NewFromSlice([]string{
	BtcName,
	BtrName,
	BtsName,
	CallName,
	CmpxchgName,
	DecName,
	EnterName,
	IncName,
	IntName,
	IntoName,
	MovName,
	MovsbName,
	MovswName,
	PushName,
	PushaName,
	RclName,
	RcrName,
	RolName,
	RorName,
	SarName,
	ShldName,
	ShlName,
	ShrdName,
	ShrName,
	StosbName,
	StoswName,
	XaddName,
	XchgName,
})

// MemoryReadWriteInstructions contains instructions that both read and write memory.
// These instructions perform atomic read-modify-write operations.
var MemoryReadWriteInstructions = set.NewFromSlice([]string{
	AdcName,
	AddName,
	AndName,
	BtcName,
	BtrName,
	BtsName,
	CmpxchgName,
	DecName,
	IncName,
	MovsbName,
	MovswName,
	OrName,
	RclName,
	RcrName,
	RolName,
	RorName,
	SarName,
	SbbName,
	ShldName,
	ShlName,
	ShrdName,
	ShrName,
	SubName,
	XaddName,
	XchgName,
	XorName,
})

// ConditionalJumpInstructions contains all conditional jump instructions.
// Used for control flow analysis and branch prediction analysis.
var ConditionalJumpInstructions = set.NewFromSlice([]string{
	JbName,
	JbeName,
	JlName,
	JleName,
	JnbName,
	JnbeName,
	JnlName,
	JnleName,
	JnoName,
	JnpName,
	JnsName,
	JnzName,
	JoName,
	JpName,
	JsName,
	JzName,
})

// StringInstructions contains all string manipulation instructions.
// These operate on memory pointed to by SI/DI registers.
var StringInstructions = set.NewFromSlice([]string{
	CmpsbName,
	CmpswName,
	InsbName,
	InswName,
	LodsbName,
	LodswName,
	MovsbName,
	MovswName,
	OutsbName,
	OutswName,
	ScasbName,
	ScaswName,
	StosbName,
	StoswName,
})

// RepeatableInstructions contains instructions that can be prefixed with REP/REPZ/REPNZ.
var RepeatableInstructions = set.NewFromSlice([]string{
	CmpsbName,
	CmpswName,
	InsbName,
	InswName,
	LodsbName,
	LodswName,
	MovsbName,
	MovswName,
	OutsbName,
	OutswName,
	ScasbName,
	ScaswName,
	StosbName,
	StoswName,
})

// PrivilegedInstructions contains instructions that require privilege level 0.
// These are typically only usable in kernel mode or real mode.
var PrivilegedInstructions = set.NewFromSlice([]string{
	CliName,
	HltName,
	InvdName,
	LmswName,
	StiName,
	WbinvdName,
})

// FlagModifyingInstructions contains instructions that explicitly modify CPU flags.
var FlagModifyingInstructions = set.NewFromSlice([]string{
	ClcName,
	CldName,
	CliName,
	CmcName,
	StcName,
	StdName,
	StiName,
})

// PortIOInstructions contains instructions that perform port I/O operations.
var PortIOInstructions = set.NewFromSlice([]string{
	InName,
	InsbName,
	InswName,
	OutName,
	OutsbName,
	OutswName,
})
