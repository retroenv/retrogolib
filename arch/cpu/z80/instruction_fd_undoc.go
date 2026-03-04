package z80

// Undocumented FD prefix instructions - IYH/IYL operations

var FdIncIYH = &Instruction{Name: IncName, Unofficial: true, NoParamFunc: fdIncIYH}
var FdDecIYH = &Instruction{Name: DecName, Unofficial: true, NoParamFunc: fdDecIYH}
var FdIncIYL = &Instruction{Name: IncName, Unofficial: true, NoParamFunc: fdIncIYL}
var FdDecIYL = &Instruction{Name: DecName, Unofficial: true, NoParamFunc: fdDecIYL}

var FdLdIYHn = &Instruction{Name: LdName, Unofficial: true, ParamFunc: fdLdIYHn}
var FdLdIYLn = &Instruction{Name: LdName, Unofficial: true, ParamFunc: fdLdIYLn}

var FdLdBIYH = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: fdLdBIYH}
var FdLdBIYL = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: fdLdBIYL}
var FdLdCIYH = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: fdLdCIYH}
var FdLdCIYL = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: fdLdCIYL}
var FdLdDIYH = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: fdLdDIYH}
var FdLdDIYL = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: fdLdDIYL}
var FdLdEIYH = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: fdLdEIYH}
var FdLdEIYL = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: fdLdEIYL}
var FdLdAIYH = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: fdLdAIYH}
var FdLdAIYL = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: fdLdAIYL}

var FdLdIYHB = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: fdLdIYHB}
var FdLdIYHC = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: fdLdIYHC}
var FdLdIYHD = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: fdLdIYHD}
var FdLdIYHE = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: fdLdIYHE}
var FdLdIYHIYH = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: fdLdIYHIYH}
var FdLdIYHIYL = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: fdLdIYHIYL}
var FdLdIYHA = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: fdLdIYHA}

var FdLdIYLB = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: fdLdIYLB}
var FdLdIYLC = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: fdLdIYLC}
var FdLdIYLD = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: fdLdIYLD}
var FdLdIYLE = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: fdLdIYLE}
var FdLdIYLIYH = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: fdLdIYLIYH}
var FdLdIYLIYL = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: fdLdIYLIYL}
var FdLdIYLA = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: fdLdIYLA}

var FdAddAIYH = &Instruction{Name: AddName, Unofficial: true, NoParamFunc: fdAddAIYH}
var FdAddAIYL = &Instruction{Name: AddName, Unofficial: true, NoParamFunc: fdAddAIYL}
var FdAdcAIYH = &Instruction{Name: AdcName, Unofficial: true, NoParamFunc: fdAdcAIYH}
var FdAdcAIYL = &Instruction{Name: AdcName, Unofficial: true, NoParamFunc: fdAdcAIYL}
var FdSubIYH = &Instruction{Name: SubName, Unofficial: true, NoParamFunc: fdSubIYH}
var FdSubIYL = &Instruction{Name: SubName, Unofficial: true, NoParamFunc: fdSubIYL}
var FdSbcAIYH = &Instruction{Name: SbcName, Unofficial: true, NoParamFunc: fdSbcAIYH}
var FdSbcAIYL = &Instruction{Name: SbcName, Unofficial: true, NoParamFunc: fdSbcAIYL}
var FdAndIYH = &Instruction{Name: AndName, Unofficial: true, NoParamFunc: fdAndIYH}
var FdAndIYL = &Instruction{Name: AndName, Unofficial: true, NoParamFunc: fdAndIYL}
var FdXorIYH = &Instruction{Name: XorName, Unofficial: true, NoParamFunc: fdXorIYH}
var FdXorIYL = &Instruction{Name: XorName, Unofficial: true, NoParamFunc: fdXorIYL}
var FdOrIYH = &Instruction{Name: OrName, Unofficial: true, NoParamFunc: fdOrIYH}
var FdOrIYL = &Instruction{Name: OrName, Unofficial: true, NoParamFunc: fdOrIYL}
var FdCpIYH = &Instruction{Name: CpName, Unofficial: true, NoParamFunc: fdCpIYH}
var FdCpIYL = &Instruction{Name: CpName, Unofficial: true, NoParamFunc: fdCpIYL}
