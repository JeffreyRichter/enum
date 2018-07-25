package enum_test

import (
	"fmt"
	"github.com/JeffreyRichter/enum/enum"
	"log"
	"reflect"
	"strconv"
)

// A ColorIdiomaticEnum value is a signed 16-bit integer
// NOTE: Defining a type improves compile-time type-safety
type ColorIdiomaticEnum int16

// NOTE: const ensures that integer values can't change while program runs
const (
	// Define ColorIdiomaticEnum's "symbols" and their values
	// NOTE: It's typical to use a common prefix to help discoverability in code editor
	ColorNone  ColorIdiomaticEnum = iota
	ColorRed                      // 1
	ColorGreen                    // 2
	ColorBlue                     // 3
)

// String coverts a ColorIdiomaticEnum's value to its equivalent "symbol"
// NOTE: Changing symbol names requires manually updating this method.
func (c ColorIdiomaticEnum) String() string {
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

// Parse sets c if s matches a "symbol"
// NOTE: Changing symbol names requires manually updating this method.
func (c *ColorIdiomaticEnum) Parse(s string) error {
	*c = ColorNone // Default unless overridden
	switch s {
	case "Red":
		*c = ColorRed
	case "Green":
		*c = ColorGreen
	case "Blue":
		*c = ColorBlue
	default:
		return fmt.Errorf("couldn't parse %q into a \"Color\"", s)
	}
	return nil
}

// ColorMenu returns the set of ColorIdiomaticEnum's symbols and values
// NOTE: Changing symbol names requires manually updating this method.
func ColorMenu() map[string]ColorIdiomaticEnum {
	return map[string]ColorIdiomaticEnum{
		ColorNone.String():  ColorNone,
		ColorRed.String():   ColorRed,
		ColorGreen.String(): ColorGreen,
		ColorBlue.String():  ColorBlue,
	}
}

func ExampleIdiomaticEnum() {
	var c ColorIdiomaticEnum = ColorRed
	printf("Color: %s\n", c) // Calls String()

	if err := c.Parse("Blue"); err == nil { // Blue is a valid color
		printf("Color: %s\n", c)
	} else {
		printf("Parse error: %s\n", err)
	}

	if err := c.Parse("Purple"); err == nil { // Purple is not a valid color
		printf("Color: %s\n", c)
	} else {
		printf("Parse error: %s\n", err)
	}

	c = ColorIdiomaticEnum(123) // A value with no matching symbol
	switch c {
	case ColorRed:
		printf("Painting the image red\n")
	case ColorBlue:
		printf("Painting the image blue\n")
	case ColorGreen:
		printf("Painting the image green\n")
	case 123:
		printf("No symbol: %d\n", c)
	}

	printf("\nColor menu:\n")
	for symbol, value := range ColorMenu() {
		printf("   %-6s %d\n", symbol, value)
	}

	// Unordered Output:
	// Color: Red
	// Color: Blue
	// Parse error: couldn't parse "Purple" into a "Color"
	// No symbol: 123
	//
	// Color menu:
	//    None   0
	//    Red    1
	//    Green  2
	//    Blue   3
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// A ColorStruct value is a signed 16-bit integer
type ColorStruct int16

// This pattern scopes the enum "symbols" to this global EColorStruct variable
// NOTE: EColorStruct or any of its fields could be overwritten breaking your code
var EColorStruct = struct {
	None, Red, Green, Blue ColorStruct
}{ColorStruct(0), ColorStruct(1), ColorStruct(2), ColorStruct(3)}

// String coverts a ColorStruct's value to its equivalent "symbol"
func (c ColorStruct) String() string {
	switch c {
	case EColorStruct.None:
		return "None"
	case EColorStruct.Red:
		return "Red"
	case EColorStruct.Green:
		return "Green"
	case EColorStruct.Blue:
		return "Blue"
	default:
		return "Unknown"
	}
}

// Parse sets c if s matches a "symbol"
func (c *ColorStruct) Parse(s string) error {
	*c = EColorStruct.None // Default unless overridden
	switch s {
	case "Red":
		*c = EColorStruct.Red
	case "Green":
		*c = EColorStruct.Green
	case "Blue":
		*c = EColorStruct.Blue
	default:
		return fmt.Errorf("couldn't parse %q into a \"Color\"", s)
	}
	return nil
}

func ExampleStructEnum() {
	var c ColorStruct = EColorStruct.Red
	printf("Color: %s\n", c) // Calls String()

	if err := c.Parse("Blue"); err == nil { // Blue is a valid color
		printf("Color: %s\n", c)
	} else {
		printf("Parse error: %s\n", err)
	}

	if err := c.Parse("Purple"); err == nil { // Purple is not a valid color
		printf("Color: %s\n", c)
	} else {
		printf("Parse error: %s\n", err)
	}
	EColorStruct.Red = ColorStruct(123) // This should not be allowed at all!

	c = ColorStruct(123) // A value with no matching symbol
	switch c {
	case EColorStruct.Red:
		printf("Painting the image red\n")
	case EColorStruct.Blue:
		printf("Painting the image blue\n")
	case EColorStruct.Green:
		printf("Painting the image green\n")
	case 123:
		printf("No symbol: %d\n", c)
	}

	// Output:
	// Color: Red
	// Color: Blue
	// Parse error: couldn't parse "Purple" into a "Color"
	// Painting the image red
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var EColor = Color(0).None() // Helper variable used by consuming code (improves cross-package consumption)
type Color int16             // I want Color enum variables to be signed 16-bit values

// Define Color's "symbols" and their values:
// NOTE: The compiler inlines calls to these methods (verify with go test -gcflags -m)
func (Color) None() Color  { return Color(0) }
func (Color) Red() Color   { return Color(1) }
func (Color) Green() Color { return Color(2) }
func (Color) Blue() Color  { return Color(3) }

// String coverts a Color enum value to its equivalent "symbol" or
// a string with an integer value if value has no matching symbol
func (c Color) String() string {
	// Calls each of Color’s methods that take no arguments & returns a Color
	// If returned value matches c’s value, return method’s name
	// Else , return c’s integer value as string
	return enum.StringInt(c, reflect.TypeOf(c))
}

// Parse sets c if s matches a symbol or is a number which can be parsed.
func (c *Color) Parse(s string) error {
	// Finds a Color method named s (optionally case-insensitive).
	// If found, calls it and sets c to its value & returns
	// Else if strict is true, returns error
	// Else (strict is off), parses s as integer
	//    If OK, set c to integer & returns; else returns error
	enumVal, err := enum.ParseInt(reflect.TypeOf(c), s, true, false)
	if enumVal != nil {
		*c = enumVal.(Color) // If no error, type assert to Color and set c
	}
	return err
}

// If you're willing to allocate some memory and decrease app startup time,
// You can allocate & initialize these maps to make String/Parse methods faster.
// Useful if String/Parse are called frequently in time-sensitive areas of your code.
var colorTables = struct {
	stringTable map[Color]string
	parseTable  map[string]Color
}{}

func init() {
	colorTables.stringTable = map[Color]string{}
	colorTables.parseTable = map[string]Color{}
	enum.GetSymbols(reflect.TypeOf(EColor), func(enumSymbolName string, enumSymbolValue interface{}) bool {
		colorTables.stringTable[enumSymbolValue.(Color)] = enumSymbolName
		colorTables.parseTable[enumSymbolName] = enumSymbolValue.(Color)
		return false
	})
}

// String coverts a Color enum value to its equivalent "symbol" or
// a string with an integer value if value has no matching symbol
func (c Color) StringFast() string {
	if s, ok := colorTables.stringTable[c]; ok {
		return s
	}
	return strconv.Itoa(int(c))
}

// Parse sets c if s matches a symbol or is a number which can be parsed.
func (c *Color) ParseFast(s string) error {
	if enumVal, ok := colorTables.parseTable[s]; ok {
		*c = enumVal
		return nil
	}
	if n, err := strconv.Atoi(s); err == nil {
		*c = Color(n)
		return nil
	} else {
		return err
	}
}

func ExampleReflectionEnum() {
	var c Color = EColor.Red()
	printf("Color: %s\n", c) // Calls String()

	if err := c.Parse("blue"); err == nil { // Blue is a valid color (case-insensitive)
		printf("Color: %s\n", c)
	} else {
		printf("Parse error: %s\n", err)
	}

	if err := c.Parse("Purple"); err == nil { // Purple is not a valid color
		printf("Color: %s\n", c)
	} else {
		printf("Parse error: %s\n", err)
	}

	c = Color(123) // A value with no matching symbol
	switch c {
	case EColor.Red():
		printf("Painting the image red\n")
	case EColor.Blue():
		printf("Painting the image blue\n")
	case EColor.Green():
		printf("Painting the image green\n")
	case 123:
		printf("No symbol: %d\n", c)
		c.Parse("0x2") // Parsing an integer string works (even hex)!
		printf("Painting the image: %s\n", c)

	}

	// BONUS: Display all the enum's symbols and each symbol's value
	enum.GetSymbols(reflect.TypeOf(EColor),
		func(enumSymbolName string, enumSymbolValue interface{}) (stop bool) {
			printf("%-6s %d\n", enumSymbolName, enumSymbolValue)
			return false
		})

	// Unordered output:
	// Color: Red
	// Color: Blue
	// Parse error: couldn't parse "Purple" into a "Color"
	// No symbol: 123
	// Painting the image: Green
	// None   0
	// Red    1
	// Green  2
	// Blue   3
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var EProtocol = Protocol("").None() // Helper variable used by consuming code (improves cross-package consumption)
type Protocol string                // I want Protocol enum variables to be string values

// Define Protocol's "symbols" and their values:
func (Protocol) None() Protocol { return Protocol("(none)") }
func (Protocol) UDP() Protocol  { return Protocol("User Datagram Protocol") }
func (Protocol) TCP() Protocol  { return Protocol("Transmission Control Protocol") }

// String coverts a Protocol enum value to its equivalent "symbol"
func (p Protocol) String() string {
	return enum.String(p, reflect.TypeOf(p))
}

// Parse sets p if s matches a symbol
func (p *Protocol) Parse(s string) error {
	v, err := enum.Parse(reflect.TypeOf(p), s, true)
	if err == nil {
		*p = v.(Protocol)
	}
	return err
}

func ExampleStringEnum() {
	var p Protocol = EProtocol.TCP()
	printf("Protocol: %s\n", p) // Calls String()

	if err := p.Parse("udp"); err == nil { // Is a valid protocol (case-insensitive)
		printf("Protocol: %s\n", p)
	} else {
		printf("Parse error: %s\n", err)
	}

	if err := p.Parse("rdp"); err == nil { // rdp is not a valid protocol
		printf("Protocol: %s\n", p)
	} else {
		printf("Parse error: %s\n", err)
	}

	switch p {
	case EProtocol.UDP():
		printf("Using " + string(p) + "\n")	// Shows string value instead of symbol
	case EProtocol.TCP():
		printf("Using " + string(p) + "\n")
	}

	// BONUS: Display all the enum's symbols and each symbol's value
	enum.GetSymbols(reflect.TypeOf(EProtocol),
		func(enumSymbolName string, enumSymbolValue interface{}) (stop bool) {
			printf("%-6s %s\n", enumSymbolName, string(enumSymbolValue.(Protocol)))
			return false
		})

	// Unordered output:
	// Protocol: TCP
	// Protocol: UDP
	// Parse error: couldn't parse "rdp" into a "Protocol"
	// Using User Datagram Protocol
	// None   (none)
	// UDP    User Datagram Protocol
	// TCP    Transmission Control Protocol
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var EAccess = Access(0).None() // Helper variable used by consuming code (improves cross-package consumption)
type Access uint32             // I want Access enum variables to flags (MUST be an unsigned integer)

// Define Access' "symbols" and their values (Note that each symbol is represented by a bit):
func (Access) None() Access           { return Access(0x00) }
func (Access) Read() Access           { return Access(0x01) }
func (Access) Write() Access          { return Access(0x02) }
func (Access) Execute() Access        { return Access(0x04) }
func (a Access) IsSet(a2 Access) bool { return (uint32(a) & uint32(a2)) != 0 } // Optional helper method if you'd like

// String coverts an Access enum value to its equivalent "symbols" (comma separated)
func (a Access) String() string {
	// Call Access' methods that take non arguments & return Access
	// if value == 0, return symbol/method that returns 0
	// Else, skip any method/symbol that returns 0;
	//    append to string any method where return value & f == method's return value
	// return string
	return enum.StringUintFlags(uint64(a), reflect.TypeOf(a), 16)
}

// Parse sets a if s matches 1+ symbols separated by commas (,)
func (a *Access) Parse(s string) error {
	// For each symbol (separated by ',') ...
	// Call enum type's method matching symbol; if found, OR value
	// Else, parse string as uint64
	//    If parsed, OR value; else return error
	v, err := enum.ParseUintFlags(reflect.TypeOf(a), s, true)
	if err == nil {
		*a = Access(v) // If no error, convert integer to Access and set a's value
	}
	return err
}

func ExampleUintFlags() {
	var a Access = EAccess.Write() | EAccess.Read()
	printf("Access: %s\n", a) // Calls String()

	if err := a.Parse("read, execute"); err == nil { // Is a valid protocol (case-insensitive)
		printf("Access: %s\n", a)
	} else {
		printf("Parse error: %s\n", err)
	}

	a = Access(0x106) // Write, Execute & 0x100
	printf("Access: %s\n", a)

	// Demonstrate that this can be round-tripped
	var b Access
	if err := b.Parse(a.String()); err == nil {
		printf("Access value: 0x%x\n", uint32(b))
	} else {
		printf("Error: %s\n", err)
	}

	if (a & EAccess.Write()) != 0 {
		printf("Use write access\n")
	}

	if a.IsSet(EAccess.Write()) { // Alternate syntax using enum method
		printf("Use write access\n")
	}

	// BONUS: Display all the enum's symbols and each symbol's value
	enum.GetSymbols(reflect.TypeOf(EAccess),
		func(enumSymbolName string, enumSymbolValue interface{}) (stop bool) {
			printf("%-8s 0x%x\n", enumSymbolName, uint32(enumSymbolValue.(Access)))
			return false
		})

	// Unordered output:
	// Access: Read, Write
	// Access: Execute, Read
	// Access: Execute, Write, 0x100
	// Access value: 0x106
	// Use write access
	// Use write access
	// None     0x0
	// Read     0x1
	// Write    0x2
	// Execute  0x4
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func init() {
	log.SetFlags(0)
}

func printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
	log.Printf(format, v...)
}
