package lines

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
)

type converter func(string) ([]reflect.Value, error)

func Each(in io.Reader, cb interface{}) error {
	// Check callback looks correct
	funcType := reflect.TypeOf(cb)
	if funcType.Kind() != reflect.Func || funcType.NumIn() != 1 || funcType.NumOut() != 0 {
		return fmt.Errorf("callback function has incorrect signature")
	}
	funcArg := funcType.In(0)
	funcValue := reflect.ValueOf(cb)

	// Find conversion function
	var conv converter
	switch funcArg.Kind() {
	case reflect.Bool:
		conv = func(s string) ([]reflect.Value, error) {
			parsed, err := strconv.ParseBool(s)
			// Must convert to funcArg incase it is a user defined type (e.g. type MyInt int)
			value := reflect.ValueOf(parsed).Convert(funcArg)
			return []reflect.Value{value}, err
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		conv = func(s string) ([]reflect.Value, error) {
			parsed, err := strconv.ParseInt(s, 10, funcArg.Bits())
			value := reflect.ValueOf(parsed).Convert(funcArg)
			return []reflect.Value{value}, err
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		conv = func(s string) ([]reflect.Value, error) {
			parsed, err := strconv.ParseUint(s, 10, funcArg.Bits())
			value := reflect.ValueOf(parsed).Convert(funcArg)
			return []reflect.Value{value}, err
		}
	case reflect.Float32, reflect.Float64:
		conv = func(s string) ([]reflect.Value, error) {
			parsed, err := strconv.ParseFloat(s, funcArg.Bits())
			value := reflect.ValueOf(parsed).Convert(funcArg)
			return []reflect.Value{value}, err
		}
	case reflect.String:
		conv = func(s string) ([]reflect.Value, error) {
			value := reflect.ValueOf(s).Convert(funcArg)
			return []reflect.Value{value}, nil
		}
	default:
		return fmt.Errorf("unable to convert argument of type %s", funcArg)
	}

	// Convert lines and call function
	line := 0
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line++
		parsed, err := conv(scanner.Text())
		if err != nil {
			return fmt.Errorf("failed to convert line %d: %w", line, err)
		}
		funcValue.Call(parsed)
	}

	return nil
}

func EachStdin(cb interface{}) error {
	return Each(os.Stdin, cb)
}
