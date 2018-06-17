package enum_test

import (
	"GoroutinePool/enum"
	"fmt"
	"reflect"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Color2 int32

const (
	ColorNone  Color2 = 0
	ColorRed   Color2 = 1
	ColorGreen Color2 = 2
	ColorBlue  Color2 = 3
)

func (c Color2) String() string {
	switch c {
	case ColorNone:
		return "None"
	case ColorRed:
		return "Red"
	case ColorGreen:
		return "Green"
	case ColorBlue:
		return "Blue"
	default:
		return "Unknown"
	}
}

func (c *Color2) Parse(cs string) error {
	*c = ColorNone // Default unless overridden
	switch cs {
	case "Red":
		*c = ColorRed
	case "Green":
		*c = ColorGreen
	case "Blue":
		*c = ColorBlue
	default:
		return fmt.Errorf("couldn't parse %q into a Color")
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var EColor = Color(0).None() // Helper variable used by consuming code (improves cross-package consumption)
type Color uint8             // I want Color enum variables to be unsigned 8-bit values

// Define Color Symbols:
func (Color) None() Color  { return Color(0) }
func (Color) Red() Color   { return Color(1) }
func (Color) Green() Color { return Color(2) }
func (Color) Blue() Color  { return Color(3) }

// String returns the enum value's symbol or a string with an integer value if value has no matching symbol
func (c Color) String() string {
	return enum.StringInt(c, reflect.TypeOf(c))
}

// Parse sets c if s matches a symbol or is a number which can be parsed.
func (c *Color) Parse(s string) error {
	enumVal, err := enum.ParseInt(reflect.TypeOf(c), s, true, false)
	if enumVal != nil {
		*c = enumVal.(Color)
	}
	return err
}

func ExampleGetSymbols() {
	// Display all the enum's symbols and each symbol's value
	enum.GetSymbols(reflect.TypeOf(EColor), func(enumSymbolName string, enumSymbolValue interface{}) (stop bool) {
		fmt.Println(enumSymbolName, enumSymbolValue)
		return false
	})

	c := EColor.Red()
	fmt.Println(c)

	err := c.Parse("Bluex")
	fmt.Println(c, err)

	err = c.Parse("Blue")
	fmt.Println(c, err)

	c = Color(123)
	fmt.Println(c)

	err = c.Parse("0x15")
	fmt.Println(c, err)

	switch c {
	case EColor.Red():
		fmt.Println(c)
	case EColor.Blue():
		fmt.Println(c)
	case EColor.Green():
		fmt.Println(c)
	case Color(21):
		fmt.Println("No symbol: 21")
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var ESASProtocol = SASProtocol("").None()

type SASProtocol string

func (SASProtocol) None() SASProtocol         { return SASProtocol("") }
func (SASProtocol) Https() SASProtocol        { return SASProtocol("https") }
func (SASProtocol) HttpsAndHttp() SASProtocol { return SASProtocol("https,http") }

func (p SASProtocol) String() string {
	return enum.String(p, reflect.TypeOf(p))
}

func (p *SASProtocol) Parse(s string) error {
	v, err := enum.Parse(reflect.TypeOf(p), s, false)
	if err == nil {
		*p = v.(SASProtocol)
	}
	return err
}

func ExampleStringEnum() {
	p := ESASProtocol.HttpsAndHttp()
	fmt.Println(p)

	err := p.Parse("foo")
	fmt.Println(p, err)

	err = p.Parse("Https")
	fmt.Println(p, err)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type protection uint32

func (protection) None() protection    { return protection(0) }
func (protection) Read() protection    { return protection(1) }
func (protection) Write() protection   { return protection(2) }
func (protection) Execute() protection { return protection(4) }

func (p protection) String() string {
	return enum.StringUintFlags(uint64(p), reflect.TypeOf(p), 16)
}

func (p *protection) Parse(s string) error {
	v, err := enum.ParseUintFlags(reflect.TypeOf(p), s, true)
	if err == nil {
		*p = protection(v)
	}
	return err
}

func ExampleUintFlags() {
	EProtection := protection(0)
	pr := protection(EProtection.Write() | EProtection.Read())
	fmt.Println(pr)
	err := pr.Parse("read, execute, 0x1001")
	fmt.Println(pr, err)
}
