package chip8

// Instructions maps instruction names to their definitions.
var Instructions = map[string]*Instruction{
	AddName:  AddInst,
	AndName:  AndInst,
	CallName: CallInst,
	ClsName:  ClsInst,
	DrwName:  DrwInst,
	JpName:   JpInst,
	LdName:   LdInst,
	OrName:   OrInst,
	RetName:  RetInst,
	RndName:  RndInst,
	SeName:   SeInst,
	ShlName:  ShlInst,
	ShrName:  ShrInst,
	SkpName:  SkpInst,
	SknpName: SknpInst,
	SneName:  SneInst,
	SubName:  SubInst,
	SubnName: SubnInst,
	XorName:  XorInst,
}
