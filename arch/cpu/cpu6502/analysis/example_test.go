package analysis_test

import (
	"fmt"

	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/arch/cpu/analysis"
	"github.com/retroenv/retrogolib/arch/cpu/cpu6502"
	cpuanalysis "github.com/retroenv/retrogolib/arch/cpu/cpu6502/analysis"
)

func ExampleTransfer() {
	input := cpuanalysis.Input{
		Opcode:  cpu6502.LdaInst.Addressing[cpu6502.ImmediateAddressing],
		Variant: cpu6502.VariantNES6502,
	}
	input.Operand, _ = knownbits.NewConstant[uint16](8, 0x80)
	var result analysis.Result[cpuanalysis.Resource, uint8, cpuanalysis.Address]
	cpuanalysis.Transfer(input, &result)
	if !result.EffectsComplete {
		return
	}
	for _, write := range result.Writes {
		resource := write.View.Storage
		next, valid := analysis.Apply(8, input.Before[resource], write)
		if !valid {
			return
		}
		input.Before[resource] = next
	}
	value, known := input.Before[cpuanalysis.A].Constant(8)
	fmt.Printf("A=%02x known=%t\n", value, known)
	// Output: A=80 known=true
}
