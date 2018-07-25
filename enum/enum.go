package enum

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// SymbolInfo defines a callback function that is invoked once per an enum type's symbol.
// The callback is passed the enum's symbol and its value.
// Return false to continue enumerating enum symbols/values or false to prematurely stop enumeration.
type SymbolInfo func(enumSymbolName string, enumSymbolValue interface{}) (stop bool)

// isValidEnumSymbolMethod is an internal function that returns true if an enum type's
// method represents a symbol.
func isValidEnumSymbolMethod(enumType reflect.Type, m reflect.Method) bool {
	// A symbol method must take 1 arg (the receiver) and return 1 value whose type matches the enum's type
	return m.Type.NumIn() == 1 && m.Type.NumOut() == 1 && m.Type.Out(0) == enumType
}

// GetSymbols invokes the SymbolInfo callback method once for each symbol defined on the enum type.
func GetSymbols(enumType reflect.Type, esi SymbolInfo) {
	// Pass 1 argument that is a zero-value of t
	args := [1]reflect.Value{reflect.Zero(enumType)}

	// Call enum methods looking for one that returns the same value we have
	for m := 0; m < enumType.NumMethod(); m++ {
		method := enumType.Method(m)
		if !isValidEnumSymbolMethod(enumType, method) {
			continue
		}
		// Call the enum method, convert the result to the enumType interface
		value := method.Func.Call(args[:])[0].Convert(enumType).Interface()
		// Pass the symbol name & value to the callback; stop enumeration if the callback returns true
		if esi(method.Name, value) {
			return
		}
	}
}

// String returns the symbol for a enum type's value. If the value has no symbol, "" is returned.
func String(enumValue interface{}, enumType reflect.Type) string {
	symbolResult := ""
	// Get symbols; if symbol's value matches enumValue, return symbol's name & stop enumeration
	GetSymbols(enumType, func(symbol string, value interface{}) bool {
		if value == enumValue {
			symbolResult = symbol
			return true
		}
		return false
	})
	return symbolResult // Returns "" if no matching symbol found
}

// StringInt returns the symbol for a enum type's value. If the value has no symbol,
// a string containing the integer value (in decimal) is returned.
func StringInt(intValue interface{}, enumType reflect.Type) string {
	// Calls enumType’s methods that return an enumType
	// If returned value matches intValue, return method’s name; else return intValue as string
	if symbolName := String(intValue, enumType); symbolName != "" {
		return symbolName // Returns matching symbol (if found)
	}
	return fmt.Sprintf("%d", intValue) // No match, return the number as a string
}

// StringUintFlags considers intValue as a bit of bit flags OR'd together and returns the
// comma-separated symbols whose bits are present. If the value has bits set which do not
// correspond to any symbol, then the remaining integer value (in intBase) is concatenated
// to the string.
func StringUintFlags(intValue uint64, enumType reflect.Type, intBase int) string {
	// Call flag's methods that return a flag
	// if flag == 0, return symbol/method that returns 0
	// else skip any method/symbol that returns 0; concatenate to string any method whose return value & f == method's return value
	// return string
	bitsFound := uint64(0)
	symbolNames := strings.Builder{}
	GetSymbols(enumType, func(symbolName string, symbolValue interface{}) bool {
		symVal := reflect.ValueOf(symbolValue).Uint()
		if intValue == 0 && symVal == 0 {
			symbolNames.WriteString(symbolName) // We found a match, return the method's name (the enum's symbol)
			return true                         // Stop
		}
		if symVal != 0 && (intValue&symVal == symVal) {
			bitsFound |= symVal
			if symbolNames.Len() > 0 {
				symbolNames.WriteString(", ")
			}
			symbolNames.WriteString(symbolName)
		}
		return false // Continue symbol enumeration
	})
	if bitsFound != intValue {
		// Some bits in the original value were not accounted for, append the remaining decimal value
		if symbolNames.Len() > 0 {
			symbolNames.WriteString(", ")
		}
		symbolNames.WriteString("0x")	// Prefix base-16 integer with "0x"
		symbolNames.WriteString(strconv.FormatUint(intValue^bitsFound, intBase))
	}
	return symbolNames.String() // Returns matching symbol (if found)
}

// Parse converts an enum type's symbol to its corresponding value.
func ParseInt(enumTypePtr reflect.Type, s string, caseInsensitive bool, strict bool) (enumVal interface{}, err error) {
	enumVal, err = Parse(enumTypePtr, s, caseInsensitive)
	if err == nil || strict {
		return // If no error or strict parsing, return Parse's results
	}

	// strict is off: Try to parse s as a string of digits into a 64-bit integer & return its value
	value := reflect.New(enumTypePtr.Elem()).Elem() // Create an enumType & get its underlying value
	if kind := value.Kind(); kind == reflect.Int || kind == reflect.Int8 || kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64 {
		intVal, parseErr := strconv.ParseInt(s, 0, int(enumTypePtr.Elem().Size())*8)
		if parseErr == nil {
			value.SetInt(intVal)        // Set the underlying value to the parsed integer
			enumVal = value.Interface() // Return the underlying value
			err = nil                   // If ParseUint returned no error, return intVal and err = nil
		}
	} else {
		intVal, parseErr := strconv.ParseUint(s, 0, int(enumTypePtr.Elem().Size())*8)
		if parseErr == nil {
			value.SetUint(intVal)       // Set the underlying value to the parsed integer
			enumVal = value.Interface() // Return the underlying value
			err = nil                   // If ParseUint returned no error, return intVal and err = nil
		}
	}
	return
}

// Parse converts an enum type's symbol to its corresponding value.
func Parse(enumTypePtr reflect.Type, s string, caseInsensitive bool) (interface{}, error) {
	// Finds enumType's method named s (optionally case-insensitive).
	// If found, calls it and returns its value; else returns error
	// sets c to its value & returns
	// If strict, return error
	// Parses s as integer; if OK, set c to int & returns; else returns error

	enumType := enumTypePtr.Elem() // Convert from *T to T
	// Look for a method name that matches the string we're trying to parse
	if method, found := findMethod(enumType, s, caseInsensitive); found {
		// Pass 1 argument that is a zero-value of t.
		args := [1]reflect.Value{reflect.Zero(enumType)}

		// Call the enum type's method passing in the arg receiver; the returned t is converted to an EnumInt32
		// The caller must convert this to their exact type
		return method.Func.Call(args[:])[0].Convert(enumType).Interface(), nil
	}
	return nil, fmt.Errorf("couldn't parse %q into a %q", s, enumType.Name())
}

// findMethod is an internal function that looks up an enum type's method (symbol) by name.
func findMethod(enumType reflect.Type, methodName string, caseInsensitive bool) (reflect.Method, bool) {
	if !caseInsensitive {
		return enumType.MethodByName(methodName) // Look up the method by exact name and case
	}
	methodName = strings.ToLower(methodName)    // lowercase the passed method name
	for m := 0; m < enumType.NumMethod(); m++ { // Iterate through all the methods matching their lowercase equivalents
		method := enumType.Method(m)
		if strings.ToLower(method.Name) == methodName {
			return method, true
		}
	}
	return reflect.Method{}, false
}

// ParseUintFlags parses a comma-separated string of symbols OR-ing each symbol's value. The
// final value is returned.
func ParseUintFlags(enumTypePtr reflect.Type, s string, caseInsensitive bool) (uint64, error) {
	val := uint64(0)
	for _, f := range strings.Split(s, ",") {
		f = strings.TrimSpace(f)
		v, err := Parse(enumTypePtr, f, caseInsensitive)
		if err == nil {
			val |= reflect.ValueOf(v).Uint() // Symbol found, OR its value
		} else {
			// strict is off: Try to parse f as a string of digits into a uint64 & return its value
			i, err := strconv.ParseUint(f, 0, int(enumTypePtr.Elem().Size())*8)
			if err == nil {
				val |= i // Successful parse, OR its value
			} else {
				return 0, fmt.Errorf("couldn't parse %q into a %q", f, enumTypePtr.Elem().Name())
			}
		}
	}
	return val, nil
}
