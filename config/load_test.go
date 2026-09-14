package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestOpenSave(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "app.conf")
	assert.NoError(t, os.WriteFile(filename, []byte("[Server]\nPort=1\nport=2"), 0600))
	cfg, err := Open(filename, Options{CaseSensitive: true})
	assert.NoError(t, err)
	values := struct {
		Port int `config:"Server.Port"`
	}{Port: 3}
	assert.NoError(t, cfg.Marshal(values))
	assert.NoError(t, cfg.Save())
	data, err := os.ReadFile(filename)
	assert.NoError(t, err)
	assert.Equal(t, "[Server]\nPort = 3\nport = 2\n", string(data))
}

func TestOpenErrors(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "app.conf")
	_, err := Open(filename, Options{})
	assert.ErrorIs(t, err, os.ErrNotExist)
	for _, tt := range []struct {
		name  string
		input string
		want  error
	}{
		{name: "duplicate key", input: "Port=1\nport=2", want: ErrDuplicateKey},
		{name: "size limit", input: strings.Repeat(" ", maxConfigSize+1), want: ErrConfigTooLarge},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, os.WriteFile(filename, []byte(tt.input), 0600))
			_, err := Open(filename, Options{})
			assert.ErrorIs(t, err, tt.want)
		})
	}
}
