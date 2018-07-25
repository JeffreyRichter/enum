/*
Package enum simplifies the creation of enumerated types (which Go does not natively support).

There are many benefits to defining & using enumerated types in your code

 - Enforces compile-time type-safety resulting in more robust code
   - An enum is a data type as opposed to just using integers, strings, etc.
   - If symbols are scoped to an enum type, symbol discovery is improved
 - Using enum symbols in code makes the code self-documenting
   - For example, "var color Color = Color.Red" is better than using an integer (like 1)
 - Restricts values to small set of legal values
   - For example, "var color Color = 217" is bad if only Red, Green, & Blue are supported colors
 - Can offer String/Parse conversions
   - Useful for command-line arguments, JSON/XML values, output/logging, etc.
 - Can return complete set of legal values
   - Useful for showing “menu” of legal set of values to a user or client package

Defining an Enumerated Type

Defining an enumerated type is as simple as coming up with a name for your type (like Color), choosing an underlying
data type for it (like int16) and then defining your desired symbols (like None, Red, Green, and Blue) and each symbol's
value (like 0, 1, 2, 3). For each symbol, you define a method that takes no parameters and returns the
enumerated type. Here is an example:

 var EColor = Color(0).None() // Helper variable used by consuming code (improves cross-package consumption)

 type Color int16             // I want Color enum variables to be signed 16-bit values

 // Define Color's "symbols" and their values:
 // NOTE: The compiler inlines calls to these methods so the doign this is very efficient
 func (Color) None() Color  { return Color(0) }
 func (Color) Red() Color   { return Color(1) }
 func (Color) Green() Color { return Color(2) }
 func (Color) Blue() Color  { return Color(3) }

Using an enumerated type is easy. Use an instance of a Color type to call one of its symbol methods:

 c := Color(0).Red() // Sets the variable c to Red (1)

To simplify this code even more (and to make using an enum defined in one package more easily usable in a code in
another package, I recommend defining a public global variable in your enum-defining package. The EColor variable
shown above is an example of this. It allows you to write code like this:

 c := EColor.Red()  // Sets the variable c to Red (1)

Implementing String and Parse Methods

So far, nothing shown above requires the use of anything in this enum package. What this package provides you is an
easy way to implement String and Parse methods on your enum types. The String method takes an enum variable (like c)
and returns its symbol as a string (Red). Conversely, Parse accepts a string (like Red) and sets an enum type's variable
to its value (1). The code below demonstrates how to implement String and Parse methods for the Color enum type:

 // String coverts a Color enum value to its equivalent "symbol" or
 // a string with an integer value if value has no matching symbol
 func (c Color) String() string {
    return enum.StringInt(c, reflect.TypeOf(c))
 }

 // Parse sets c if s matches a symbol or is a number which can be parsed.
 func (c *Color) Parse(s string) error {
    enumVal, err := enum.ParseInt(reflect.TypeOf(c), s, true, false)
    if enumVal != nil {
       *c = enumVal.(Color) // If no error, type assert to Color and set c
    }
    return err
 }

The great thing about these methods is that you can add, remove, or rename any of your enum type's symbol methods and
these methods require no change at all; they just work! In addition, Parse optionally supports case-insensitive string
matching and optionally supports strict parsing. With strict parsing off (false), parsing a string of "123" will set the
color variable to a value of 123 instead of returning an error. When String is called, if it cannot find a matching
symbol method, it returns a string with the number "123". Unstrict parsing allows round-tripping of data (a number
string from XML, JSON, or whatever) and being able to parse it. And then later, String converts it back to a number
string without any loss of information.

Getting all of an Enumerated Types's Symbols and Values

This enum package offers a GetSymbols function that invokes your callback method once for each of your
enumerated type's symbols. Your callback is called once per symbol and is passed the symbol's string and its value
(as an interface{}); your callback can process these however it likes. Your callback returns false to continue
enumerating symbols or it can return true to prematurely stop the enumeration if it has found whatever it is looking
for.

Below is an example of code calling this package's GetSymbols method. The callback method simply displays each
symbol's string along with its numeric value.

 enum.GetSymbols(reflect.TypeOf(EColor),
    func(enumSymbolName string, enumSymbolValue interface{}) (stop bool) {
       fmt.Printf("%-6s %d\n", enumSymbolName, enumSymbolValue)
       return false
    })

Working with Bit Flag Enumerated Types

You can also define enumerated types that consist of bit flags (symbols) that you can bitwise-OR together. Note that
the enumerated type underlying type MUST be an unsigned integer (like uint32). Here is an example of an enumerated
type that defines a set of potential access conditions:

 var EAccess = Access(0).None() // Helper variable used by consuming code (improves cross-package consumption)
 type Access uint32             // I want Access enum variables to flags (MUST be an unsigned integer)

 // Define Access' "symbols" and their values (Note that each symbol is represented by a bit):
 func (Access) None() Access           { return Access(0x00) }
 func (Access) Read() Access           { return Access(0x01) }
 func (Access) Write() Access          { return Access(0x02) }
 func (Access) Execute() Access        { return Access(0x04) }

 // String coverts an Access enum value to its equivalent "symbols" (comma separated)
 func (a Access) String() string {
    return enum.StringUintFlags(uint64(a), reflect.TypeOf(a), 16)
 }

 // Parse sets a if s matches 1+ symbols separated by commas (,)
 func (a *Access) Parse(s string) error {
    v, err := enum.ParseUintFlags(reflect.TypeOf(a), s, true)
    if err == nil {
       *a = Access(v) // If no error, convert integer to Access and set a's value
    }
    return err
 }

Here is code showing how to use String and Parse with this enumerated type:

 var a Access = EAccess.Write() | EAccess.Read()
 printf("%s\n", a) // Calls String() which returns "Read, Write"

 var b Access
 if err := b.Parse("write, execute"); err == nil {	// Note optional case-insensitive matching
    printf("%s", b)	// Returns "Write, Execute"
 } else {
    printf("Error: %v\n", err)
 }
*/
package enum
