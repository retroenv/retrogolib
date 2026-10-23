package analysis_test

import (
	"fmt"

	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/arch/cpu/analysis"
	"github.com/retroenv/retrogolib/arch/cpu/z80"
	cpuanalysis "github.com/retroenv/retrogolib/arch/cpu/z80/analysis"
)

func ExampleTransfer() {
	input := cpuanalysis.Input{
		Opcode: z80.OpcodeInfo{
			Opcode: 0x06,
			Size:   2,
		},
	}
	input.Before[cpuanalysis.BC], _ = knownbits.NewConstant[uint16](16, 0x9934)
	input.Operand, _ = knownbits.NewConstant[uint16](8, 0x12)
	var result analysis.Result[cpuanalysis.Resource, uint16, cpuanalysis.Address]
	cpuanalysis.Transfer(input, &result)
	if !result.EffectsComplete {
		return
	}
	for _, write := range result.Writes {
		resource := write.View.Storage
		width := cpuanalysis.StorageWidth(resource)
		next, valid := analysis.Apply(width, input.Before[resource], write)
		if !valid {
			return
		}
		input.Before[resource] = next
	}
	value, known := input.Before[cpuanalysis.BC].Constant(16)
	fmt.Printf("BC=%04x known=%t\n", value, known)
	// Output: BC=1234 known=true
}
