package smErrors

import "fmt"

/*
dimension_error.go
Description:
	Defining the dimension error object and all of its associated functions.
*/

/*
DimensionError
Description:
*/
type DimensionError struct {
	Arg1      MatrixLike
	Arg2      MatrixLike
	Operation string // Either multiply or Plus
}

func (de DimensionError) Error() string {
	dimStrings := de.ArgDimsAsStrings()
	return fmt.Sprintf(
		"dimension error: Cannot perform %v between expression of dimension %v and expression of dimension %v",
		de.Operation,
		dimStrings[0],
		dimStrings[1],
	)
}

func (de DimensionError) ArgDimsAsStrings() []string {
	// Create string for arg 1
	arg1DimsAsString := "("
	for ii, dimValue := range de.Arg1.Dims() {
		arg1DimsAsString += fmt.Sprintf("%v", dimValue)
		if ii != len(de.Arg1.Dims())-1 { // If this isn't the last element of size
			arg1DimsAsString += ","
		}
	}
	arg1DimsAsString += ")"

	// Create string for arg 2
	arg2DimsAsString := "("
	for ii, dimValue := range de.Arg2.Dims() {
		arg2DimsAsString += fmt.Sprintf("%v", dimValue)
		if ii != len(de.Arg2.Dims())-1 { // If this isn't the last element of size
			arg2DimsAsString += ","
		}
	}
	arg2DimsAsString += ")"

	return []string{arg1DimsAsString, arg2DimsAsString}

}

func CheckDimensionsInAddition(left, right MatrixLike) error {
	// Check that the size of columns in left and right agree
	dimsAreMatched := (left.Dims()[0] == right.Dims()[0]) && (left.Dims()[1] == right.Dims()[1])
	dimsAreMatched = dimsAreMatched || IsScalarExpression(left)
	dimsAreMatched = dimsAreMatched || IsScalarExpression(right)

	if !dimsAreMatched {
		return DimensionError{
			Operation: "Plus",
			Arg1:      left,
			Arg2:      right,
		}
	}
	// If dimensions match, then return nothing.
	return nil
}

/*
CheckDimensionsInSubtraction
Description:

	This function checks that the dimensions of the left and right expressions
*/
func CheckDimensionsInSubtraction(left, right MatrixLike) error {
	// Check that the size of columns in left and right agree
	dimsAreMatched := (left.Dims()[0] == right.Dims()[0]) && (left.Dims()[1] == right.Dims()[1])
	dimsAreMatched = dimsAreMatched || IsScalarExpression(left)
	dimsAreMatched = dimsAreMatched || IsScalarExpression(right)

	if !dimsAreMatched {
		return DimensionError{
			Operation: "Minus",
			Arg1:      left,
			Arg2:      right,
		}
	}
	// If dimensions match, then return nothing.
	return nil
}

/*
CheckDimensionsInMultiplication
Description:

	This function checks that the dimensions of the left and right expressions
	are compatible for multiplication.
	We allow:
	- Multiplication if Dimensions Match OR
	- Multiplication if one of the expressions is a scalar
*/
func CheckDimensionsInMultiplication(left, right MatrixLike) error {
	// Check that dimensions match
	dimsMatch := left.Dims()[1] == right.Dims()[0]

	// Check that one of the expressions is a scalar
	leftIsScalar := IsScalarExpression(left)
	rightIsScalar := IsScalarExpression(right)

	multiplicationIsAllowed := dimsMatch || leftIsScalar || rightIsScalar

	// Check that the # of columns in left
	// matches the # of rows in right
	if !multiplicationIsAllowed {
		return DimensionError{
			Operation: "Multiply",
			Arg1:      left,
			Arg2:      right,
		}
	}
	// If dimensions match, then return nothing.
	return nil
}

/*
CheckDimensionsInHStack
Description:

	This function checks that the dimensions of the left and right expressions
	are compatible for horizontal stacking.
	We allow:
	- Stacking if the number of rows match
*/
func CheckDimensionsInHStack(sliceToStack ...MatrixLike) error {
	// Check that the size of columns in left and right agree
	var nRowsInSlice []int
	for _, slice := range sliceToStack {
		nRowsInSlice = append(nRowsInSlice, slice.Dims()[0])
	}

	// Check that the number of rows in each slice is the same
	for ii := 1; ii < len(nRowsInSlice); ii++ {
		// If the number of rows in the slice is not the same as the previous slice,
		// then return an error
		dimsAreMatched := nRowsInSlice[ii] == nRowsInSlice[ii-1]
		if !dimsAreMatched {
			return DimensionError{
				Operation: "HStack",
				Arg1:      sliceToStack[ii-1],
				Arg2:      sliceToStack[ii],
			}
		}
	}

	// If dimensions match, then return nothing.
	return nil
}

/*
CheckDimensionsInVStack
Description:

	This function checks that the dimensions of the left and right expressions
	are compatible for vertical stacking.
	We allow:
	- Stacking if the number of columns match
*/
func CheckDimensionsInVStack(sliceToStack ...MatrixLike) error {
	// Check that the size of columns in left and right agree
	var nColsInSlice []int
	for _, slice := range sliceToStack {
		nColsInSlice = append(nColsInSlice, slice.Dims()[1])
	}

	// Check that the number of rows in each slice is the same
	for ii := 1; ii < len(nColsInSlice); ii++ {
		// If the number of rows in the slice is not the same as the previous slice,
		// then return an error
		dimsAreMatched := nColsInSlice[ii] == nColsInSlice[ii-1]
		if !dimsAreMatched {
			return DimensionError{
				Operation: "VStack",
				Arg1:      sliceToStack[ii-1],
				Arg2:      sliceToStack[ii],
			}
		}
	}

	// If dimensions match, then return nothing.
	return nil
}
