package audio

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestFormatValidation(t *testing.T) {
	valid := Format{
		SampleRate: 48000, Channels: 2, Samples: 256, Format: FormatS16,
	}
	assert.NoError(t, valid.Validate())
	assert.Equal(t, 1024, valid.BufferSize())
	for _, change := range []func(*Format){
		func(f *Format) { f.SampleRate = 0 },
		func(f *Format) { f.SampleRate = -1 },
		func(f *Format) { f.Channels = 0 },
		func(f *Format) { f.Channels = 257 },
		func(f *Format) { f.Samples = 0 },
		func(f *Format) { f.Samples = -1 },
		func(f *Format) { f.Samples = int(^uint(0) >> 1) },
		func(f *Format) { f.Format = 0xffff },
	} {
		format := valid
		change(&format)
		assert.Error(t, format.Validate())
		assert.Equal(t, 0, format.BufferSize())
	}
}
