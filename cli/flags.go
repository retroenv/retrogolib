// Package cli provides utilities for command-line interface applications.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Flag type constants for usage generation.
const (
	typeBool    = "bool"
	typeFloat64 = "float64"
	typeInt     = "int"
	typeInt64   = "int64"
	typeString  = "string"
	typeUint    = "uint"
	typeUint64  = "uint64"
)

// Struct tag constants.
const (
	tagArg        = "arg"
	tagDefault    = "default"
	tagEnv        = "env"
	tagFlag       = "flag"
	tagPositional = "positional"
	tagRequired   = "required"
	tagTrue       = "true"
	tagUsage      = "usage"
)

// Section groups flags for organized usage output.
type Section struct {
	Name  string
	Flags []FlagInfo
}

// FlagInfo contains metadata about a flag for usage generation.
type FlagInfo struct {
	Name     string
	Short    string // Short form (e.g., "v" for "-v")
	Usage    string
	Default  string
	Type     string
	Env      string
	Required bool
}

// PositionalInfo contains metadata about a positional argument.
type PositionalInfo struct {
	Name     string
	Usage    string
	Required bool
	Variadic bool // True for []string fields (consumes remaining args)
}

// requiredFlag tracks a required flag for validation.
type requiredFlag struct {
	name string
	ptr  any
}

// positionalArg tracks a positional argument for assignment.
type positionalArg struct {
	info PositionalInfo
	ptr  any
}

// FlagSet wraps flag.FlagSet with section-based usage generation.
type FlagSet struct {
	flags      *flag.FlagSet
	sections   []Section
	name       string
	required   []requiredFlag
	positional []positionalArg
}

// NewFlagSet creates a new FlagSet with the given program name.
// The usage line is auto-generated based on registered flags and positional arguments.
func NewFlagSet(name string) *FlagSet {
	fs := &FlagSet{
		flags: flag.NewFlagSet(name, flag.ContinueOnError),
		name:  name,
	}
	fs.flags.Usage = fs.showUsage
	return fs
}

// AddSection adds a named section with flags parsed from struct tags.
// The struct fields should have tags: `flag:"name" usage:"description"`.
// Short flags: `flag:"v,verbose"` creates both -v and -verbose.
// Optional tags: `default:"value"`, `env:"VAR_NAME"`, `required:"true"`.
func (f *FlagSet) AddSection(name string, opts any) {
	section := Section{Name: name}

	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		flagTag := field.Tag.Get(tagFlag)
		if flagTag == "" || flagTag == "-" {
			continue
		}

		fieldVal := v.Field(i)
		if !fieldVal.CanAddr() {
			continue
		}

		info := f.parseFlagField(field, fieldVal)
		if info == nil {
			continue
		}

		if info.Required {
			f.required = append(f.required, requiredFlag{name: info.Name, ptr: fieldVal.Addr().Interface()})
		}

		section.Flags = append(section.Flags, *info)
	}

	f.sections = append(f.sections, section)
}

// AddPositional registers positional arguments from struct tags.
// Fields should have tag: `arg:"positional" usage:"description"`.
// Optional: `required:"true"`. The last field can be []string for variadic args.
func (f *FlagSet) AddPositional(opts any) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		argTag := field.Tag.Get(tagArg)
		if argTag != tagPositional {
			continue
		}

		fieldVal := v.Field(i)
		if !fieldVal.CanAddr() {
			continue
		}

		info := PositionalInfo{
			Name:     strings.ToLower(field.Name),
			Usage:    field.Tag.Get(tagUsage),
			Required: field.Tag.Get(tagRequired) == tagTrue,
		}

		ptr := fieldVal.Addr().Interface()

		// Check if it's a variadic ([]string) argument.
		if _, ok := ptr.(*[]string); ok {
			info.Variadic = true
		} else if _, ok := ptr.(*string); !ok {
			// Only string and []string are supported for positional args.
			continue
		}

		f.positional = append(f.positional, positionalArg{info: info, ptr: ptr})
	}
}

// Parse parses command-line arguments and assigns positional arguments.
func (f *FlagSet) Parse(args []string) ([]string, error) {
	if err := f.flags.Parse(args); err != nil {
		return nil, fmt.Errorf("parsing flags: %w", err)
	}
	if err := f.validateRequired(); err != nil {
		return nil, err
	}

	remaining := f.flags.Args()
	remaining, err := f.assignPositional(remaining)
	if err != nil {
		return nil, err
	}
	return remaining, nil
}

// ShowUsage prints the usage message with sections.
func (f *FlagSet) ShowUsage() {
	f.showUsage()
}

func (f *FlagSet) parseFlagField(field reflect.StructField, fieldVal reflect.Value) *FlagInfo {
	flagTag := field.Tag.Get(tagFlag)

	// Parse short,long format (e.g., "v,verbose").
	var short, long string
	if idx := strings.Index(flagTag, ","); idx != -1 {
		short = flagTag[:idx]
		long = flagTag[idx+1:]
	} else {
		long = flagTag
	}

	info := &FlagInfo{
		Name:     long,
		Short:    short,
		Usage:    field.Tag.Get(tagUsage),
		Default:  field.Tag.Get(tagDefault),
		Env:      field.Tag.Get(tagEnv),
		Required: field.Tag.Get(tagRequired) == tagTrue,
	}

	// Check environment variable before using default.
	if info.Env != "" {
		if envVal := os.Getenv(info.Env); envVal != "" {
			info.Default = envVal
		}
	}

	ptr := fieldVal.Addr().Interface()
	if !f.registerFlag(info, ptr) {
		return nil
	}
	return info
}

func (f *FlagSet) registerFlag(info *FlagInfo, ptr any) bool {
	// Collect names to register (long and optionally short).
	names := []string{info.Name}
	if info.Short != "" {
		names = append(names, info.Short)
	}

	switch p := ptr.(type) {
	case *string:
		f.registerString(info, p, names)
	case *bool:
		f.registerBool(info, p, names)
	case *int:
		f.registerInt(info, p, names)
	case *int64:
		f.registerInt64(info, p, names)
	case *uint:
		f.registerUint(info, p, names)
	case *uint64:
		f.registerUint64(info, p, names)
	case *float64:
		f.registerFloat64(info, p, names)
	default:
		return false
	}
	return true
}

func (f *FlagSet) registerString(info *FlagInfo, p *string, names []string) {
	info.Type = typeString
	defaultVal := info.Default
	if defaultVal == "" {
		defaultVal = *p
		info.Default = defaultVal
	}
	for _, name := range names {
		f.flags.StringVar(p, name, defaultVal, info.Usage)
	}
}

func (f *FlagSet) registerBool(info *FlagInfo, p *bool, names []string) {
	info.Type = typeBool
	def := info.Default == tagTrue
	for _, name := range names {
		f.flags.BoolVar(p, name, def, info.Usage)
	}
}

func (f *FlagSet) registerInt(info *FlagInfo, p *int, names []string) {
	info.Type = typeInt
	def, _ := strconv.Atoi(info.Default)
	for _, name := range names {
		f.flags.IntVar(p, name, def, info.Usage)
	}
}

func (f *FlagSet) registerInt64(info *FlagInfo, p *int64, names []string) {
	info.Type = typeInt64
	def, _ := strconv.ParseInt(info.Default, 10, 64)
	for _, name := range names {
		f.flags.Int64Var(p, name, def, info.Usage)
	}
}

func (f *FlagSet) registerUint(info *FlagInfo, p *uint, names []string) {
	info.Type = typeUint
	def, _ := strconv.ParseUint(info.Default, 10, 64)
	for _, name := range names {
		f.flags.UintVar(p, name, uint(def), info.Usage)
	}
}

func (f *FlagSet) registerUint64(info *FlagInfo, p *uint64, names []string) {
	info.Type = typeUint64
	def, _ := strconv.ParseUint(info.Default, 10, 64)
	for _, name := range names {
		f.flags.Uint64Var(p, name, def, info.Usage)
	}
}

func (f *FlagSet) registerFloat64(info *FlagInfo, p *float64, names []string) {
	info.Type = typeFloat64
	def, _ := strconv.ParseFloat(info.Default, 64)
	for _, name := range names {
		f.flags.Float64Var(p, name, def, info.Usage)
	}
}

func (f *FlagSet) validateRequired() error {
	var missing []string
	for _, req := range f.required {
		if isZeroValue(req.ptr) {
			missing = append(missing, req.name)
		}
	}
	if len(missing) > 0 {
		return errors.New("missing required flag(s): " + strings.Join(missing, ", "))
	}
	return nil
}

func (f *FlagSet) assignPositional(args []string) ([]string, error) {
	var missing []string

	for i, pos := range f.positional {
		if pos.info.Variadic {
			// Variadic argument consumes all remaining args.
			if p, ok := pos.ptr.(*[]string); ok {
				*p = args
				args = nil
			}
			break
		}

		if len(args) == 0 {
			// No more arguments available.
			if pos.info.Required {
				missing = append(missing, pos.info.Name)
			}
			continue
		}

		// Assign the next argument.
		if p, ok := pos.ptr.(*string); ok {
			*p = args[0]
			args = args[1:]
		}

		// Check remaining required positional args.
		if len(args) == 0 {
			for j := i + 1; j < len(f.positional); j++ {
				if f.positional[j].info.Required && !f.positional[j].info.Variadic {
					missing = append(missing, f.positional[j].info.Name)
				}
			}
		}
	}

	if len(missing) > 0 {
		return nil, errors.New("missing required argument(s): " + strings.Join(missing, ", "))
	}
	return args, nil
}

func (f *FlagSet) showUsage() {
	fmt.Println(f.buildUsageLine())
	fmt.Println()

	for _, section := range f.sections {
		fmt.Printf("%s:\n", section.Name)
		for _, fl := range section.Flags {
			f.printFlag(fl)
		}
		fmt.Println()
	}

	if len(f.positional) > 0 {
		fmt.Println("Positional arguments:")
		for _, pos := range f.positional {
			f.printPositional(pos.info)
		}
		fmt.Println()
	}
}

func (f *FlagSet) buildUsageLine() string {
	var sb strings.Builder
	sb.WriteString("usage: ")
	sb.WriteString(f.name)

	if len(f.sections) > 0 {
		sb.WriteString(" [options]")
	}
	for _, pos := range f.positional {
		sb.WriteByte(' ')
		sb.WriteString(formatPositionalUsage(pos.info))
	}
	return sb.String()
}

func formatPositionalUsage(pos PositionalInfo) string {
	name := pos.Name
	if pos.Variadic {
		name += "..."
	}
	if pos.Required {
		return name
	}
	return "[" + name + "]"
}

func (f *FlagSet) printPositional(pos PositionalInfo) {
	name := pos.Name
	if pos.Variadic {
		name += "..."
	}
	fmt.Printf("  %s\n", name)

	usage := pos.Usage
	if pos.Required {
		usage += " (required)"
	}
	fmt.Printf("    \t%s\n", usage)
}

func (f *FlagSet) printFlag(fl FlagInfo) {
	// Build flag name display (e.g., "-v, -verbose" or just "-verbose").
	var flagDisplay string
	if fl.Short != "" {
		flagDisplay = fmt.Sprintf("-%s, -%s", fl.Short, fl.Name)
	} else {
		flagDisplay = "-" + fl.Name
	}

	// Bool flags don't show type indicator.
	if fl.Type == typeBool {
		fmt.Printf("  %s\n", flagDisplay)
	} else {
		fmt.Printf("  %s %s\n", flagDisplay, fl.Type)
	}

	usage := fl.Usage
	if fl.Required {
		usage += " (required)"
	}
	if fl.Env != "" {
		usage += " [env: " + fl.Env + "]"
	}
	if fl.Default != "" && fl.Type != typeBool && !strings.Contains(fl.Usage, "(default:") {
		usage += " (default: " + fl.Default + ")"
	}
	fmt.Printf("    \t%s\n", usage)
}

func isZeroValue(ptr any) bool {
	switch p := ptr.(type) {
	case *string:
		return *p == ""
	case *bool:
		return !*p
	case *int:
		return *p == 0
	case *int64:
		return *p == 0
	case *uint:
		return *p == 0
	case *uint64:
		return *p == 0
	case *float64:
		return *p == 0
	}
	return false
}
