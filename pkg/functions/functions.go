package functions

import (
	"context"
	"fmt"
	"reflect"
)

// runFuncWithArgs runs a sequence of functions with their corresponding arguments using recursion.
func runFuncWithArgs(ctx context.Context, funcsAndArgs []interface{}, index int) error {
	if index >= len(funcsAndArgs) {
		return nil
	}
	// Check if context is done
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Proceed only if context is not done

		item := funcsAndArgs[index]
		itemValue := reflect.ValueOf(item)

		if itemValue.Kind() != reflect.Func {
			return fmt.Errorf("invalid sequence: function item encountered without preceding function: %v", itemValue.Type())
		}

		numIn := itemValue.Type().NumIn()
		isVariadic := itemValue.Type().IsVariadic()

		args := []reflect.Value{}
		for j := 0; j < numIn && index+1 < len(funcsAndArgs); j++ {
			if isVariadic && j >= numIn-1 {
				for ; index+1 < len(funcsAndArgs); index++ {
					nextItem := funcsAndArgs[index+1]
					nextItemValue := reflect.ValueOf(nextItem)
					if !nextItemValue.Type().AssignableTo(itemValue.Type().In(numIn - 1).Elem()) {
						break
					}
					args = append(args, nextItemValue)
				}
				break
			}

			nextItem := funcsAndArgs[index+1]
			nextItemValue := reflect.ValueOf(nextItem)
			if !nextItemValue.Type().AssignableTo(itemValue.Type().In(j)) {
				return fmt.Errorf("invalid argument type for function: %v", nextItemValue.Type())
			}

			args = append(args, nextItemValue)
			index++
		}

		if len(args) < numIn {
			return fmt.Errorf("not enough arguments for function: expected %d, got %d", numIn, len(args))
		}

		results := itemValue.Call(args)
		if len(results) > 0 && !results[0].IsNil() {
			return results[0].Interface().(error)
		}

		return runFuncWithArgs(ctx, funcsAndArgs, index+1)
	}
}

// RunFuncWithArgs runs a sequence of functions with their corresponding arguments.
//
// It takes a variadic parameter of interface{} type, where each item in the
// sequence represents either a function or an argument for the preceding
// function.
//
// The function calls the internal recursive function runFuncWithArgs to iterate
// over the sequence and perform the necessary steps for each item.
func RunFuncWithArgs(funcsAndArgs ...interface{}) error {
	return runFuncWithArgs(context.Background(), funcsAndArgs, 0)
}

// RunFuncWithContext runs a sequence of functions with their corresponding arguments.
//
// It takes a variadic parameter of interface{} type, where each item in the
// sequence represents either a function or an argument for the preceding
// function.
//
// The function calls the internal recursive function runFuncWithArgs to iterate
// over the sequence and perform the necessary steps for each item.
func RunFuncWithContext(ctx context.Context, funcsAndArgs ...interface{}) error {
	return runFuncWithArgs(ctx, funcsAndArgs, 0)
}
