package functions

import (
	"fmt"
	"reflect"
)

// RunFuncWithArgs runs a sequence of functions with their corresponding arguments.
//
// It takes a variadic parameter of interface{} type, where each item in the
// sequence represents either a function or an argument for the preceding
// function.
//
// The function iterates over the sequence and performs the following steps for
// each item:
//  1. Checks if the item is a function. If not, it returns an error indicating
//     that a function item was encountered without a preceding function.
//  2. Stores the function value in the currentFunc variable.
//  3. Determines the number of input arguments the function expects.
//  4. Creates a slice to store the input arguments.
//  5. Iterates over the input arguments for the function and performs the
//     following steps for each argument:
//     a. Checks if the next item in the sequence is assignable to the current
//     function's input argument type. If not, it returns an error indicating
//     an invalid argument type.
//     b. Appends the argument value to the args slice.
//     c. Increments the index variable to move to the next item in the sequence.
//  6. Checks if the number of input arguments collected matches the expected
//     number. If not, it returns an error indicating that not enough arguments
//     were provided for the function.
//  7. Calls the function with the collected input arguments and stores the
//     results in a slice.
//  8. Checks if the first result (assuming there is only one result) is not nil.
//     If it's not nil, it returns the first result as an error.
//  9. If no error was encountered, continues to the next item in the sequence.
//  10. If all items in the sequence have been processed without encountering
//     any error, returns nil.
func RunFuncWithArgs(funcsAndArgs ...interface{}) error {
	var currentFunc reflect.Value // Stores the current function value

	for i := 0; i < len(funcsAndArgs); i++ {
		// Get the current item from the sequence
		item := funcsAndArgs[i]
		itemValue := reflect.ValueOf(item)

		// Check if the item is a function
		if itemValue.Kind() != reflect.Func {
			// If not, return an error
			return fmt.Errorf("invalid sequence: function item encountered without preceding function: %v", itemValue.Type())
		}

		// Store the function value in the currentFunc variable
		currentFunc = itemValue

		// Determine the number of input arguments the function expects
		numIn := currentFunc.Type().NumIn()

		// Create a slice to store the input arguments
		args := make([]reflect.Value, 0, numIn)

		// Iterate over the input arguments for the function
		for j := 0; j < numIn && i+1 < len(funcsAndArgs); j++ {
			// Get the next item in the sequence
			nextItem := funcsAndArgs[i+1]
			nextItemValue := reflect.ValueOf(nextItem)

			// Check if the next item is assignable to the current function's input argument type
			if !nextItemValue.Type().AssignableTo(currentFunc.Type().In(j)) {
				// If not, return an error
				return fmt.Errorf("invalid argument type for function: %v", nextItemValue.Type())
			}

			// Append the argument value to the args slice
			args = append(args, nextItemValue)

			// Increment the index variable to move to the next item in the sequence
			i++
		}

		// Check if the number of input arguments collected matches the expected number
		if len(args) != numIn {
			// If not, return an error
			return fmt.Errorf("not enough arguments for function: expected %d, got %d", numIn, len(args))
		}

		// Call the function with the collected input arguments and store the results in a slice
		results := currentFunc.Call(args)

		// Check if the first result (assuming there is only one result) is not nil
		if len(results) > 0 && !results[0].IsNil() {
			// If it's not nil, return the first result as an error
			return results[0].Interface().(error)
		}
	}

	// If no error was encountered, return nil
	return nil
}
