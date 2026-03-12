package z80

// Undocumented DD prefix instructions - IXH/IXL operations

var DdIncIXH = &Instruction{Name: IncName, Unofficial: true, NoParamFunc: ddIncIXH}
var DdDecIXH = &Instruction{Name: DecName, Unofficial: true, NoParamFunc: ddDecIXH}
var DdIncIXL = &Instruction{Name: IncName, Unofficial: true, NoParamFunc: ddIncIXL}
var DdDecIXL = &Instruction{Name: DecName, Unofficial: true, NoParamFunc: ddDecIXL}

var DdLdIXHn = &Instruction{Name: LdName, Unofficial: true, ParamFunc: ddLdIXHn}
var DdLdIXLn = &Instruction{Name: LdName, Unofficial: true, ParamFunc: ddLdIXLn}

var DdLdBIXH = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: ddLdBIXH}
var DdLdBIXL = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: ddLdBIXL}
var DdLdCIXH = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: ddLdCIXH}
var DdLdCIXL = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: ddLdCIXL}
var DdLdDIXH = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: ddLdDIXH}
var DdLdDIXL = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: ddLdDIXL}
var DdLdEIXH = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: ddLdEIXH}
var DdLdEIXL = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: ddLdEIXL}
var DdLdAIXH = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: ddLdAIXH}
var DdLdAIXL = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: ddLdAIXL}

var DdLdIXHB = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: ddLdIXHB}
var DdLdIXHC = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: ddLdIXHC}
var DdLdIXHD = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: ddLdIXHD}
var DdLdIXHE = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: ddLdIXHE}
var DdLdIXHIXH = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: ddLdIXHIXH}
var DdLdIXHIXL = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: ddLdIXHIXL}
var DdLdIXHA = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: ddLdIXHA}

var DdLdIXLB = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: ddLdIXLB}
var DdLdIXLC = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: ddLdIXLC}
var DdLdIXLD = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: ddLdIXLD}
var DdLdIXLE = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: ddLdIXLE}
var DdLdIXLIXH = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: ddLdIXLIXH}
var DdLdIXLIXL = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: ddLdIXLIXL}
var DdLdIXLA = &Instruction{Name: LdName, Unofficial: true, NoParamFunc: ddLdIXLA}

var DdAddAIXH = &Instruction{Name: AddName, Unofficial: true, NoParamFunc: ddAddAIXH}
var DdAddAIXL = &Instruction{Name: AddName, Unofficial: true, NoParamFunc: ddAddAIXL}
var DdAdcAIXH = &Instruction{Name: AdcName, Unofficial: true, NoParamFunc: ddAdcAIXH}
var DdAdcAIXL = &Instruction{Name: AdcName, Unofficial: true, NoParamFunc: ddAdcAIXL}
var DdSubIXH = &Instruction{Name: SubName, Unofficial: true, NoParamFunc: ddSubIXH}
var DdSubIXL = &Instruction{Name: SubName, Unofficial: true, NoParamFunc: ddSubIXL}
var DdSbcAIXH = &Instruction{Name: SbcName, Unofficial: true, NoParamFunc: ddSbcAIXH}
var DdSbcAIXL = &Instruction{Name: SbcName, Unofficial: true, NoParamFunc: ddSbcAIXL}
var DdAndIXH = &Instruction{Name: AndName, Unofficial: true, NoParamFunc: ddAndIXH}
var DdAndIXL = &Instruction{Name: AndName, Unofficial: true, NoParamFunc: ddAndIXL}
var DdXorIXH = &Instruction{Name: XorName, Unofficial: true, NoParamFunc: ddXorIXH}
var DdXorIXL = &Instruction{Name: XorName, Unofficial: true, NoParamFunc: ddXorIXL}
var DdOrIXH = &Instruction{Name: OrName, Unofficial: true, NoParamFunc: ddOrIXH}
var DdOrIXL = &Instruction{Name: OrName, Unofficial: true, NoParamFunc: ddOrIXL}
var DdCpIXH = &Instruction{Name: CpName, Unofficial: true, NoParamFunc: ddCpIXH}
var DdCpIXL = &Instruction{Name: CpName, Unofficial: true, NoParamFunc: ddCpIXL}
