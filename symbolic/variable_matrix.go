package symbolic

import (
	"fmt"
	"github.com/MatProGo-dev/SymbolicMath.go/smErrors"
	"gonum.org/v1/gonum/mat"
)

/*
variable_matrix.go
Description:

	This file contains functions that are used to create and
	manipulate matrices of variables (i.e., VariableMatrix) objects.
*/

// Types
// =====

type VariableMatrix [][]Variable

// =======
// Methods
// =======

/*
Check
Description:

	This method is used to make sure that the variable matrix is well-defined.
	It checks:
	- That the matrix is not empty.
	- That all of the rows have the same number of columns.
	- That all of the variables are well-defined.
*/
func (vm VariableMatrix) Check() error {
	// Input Processing

	// Check to see if the matrix is empty
	if len(vm) == 0 {
		return smErrors.EmptyMatrixError{
			Expression: vm,
		}
	}

	// Check the number of columns in each row
	var numCols int = len(vm[0])
	for ii, vmRow := range vm {
		// Check the length of row i
		if len(vmRow) != numCols {
			return smErrors.MatrixColumnMismatchError{
				ExpectedNColumns: numCols,
				ActualNColumns:   len(vmRow),
				Row:              ii,
			}
		}
	}

	// Check the variables
	for ii, vmRow := range vm {
		for jj, v := range vmRow {
			err := v.Check()
			if err != nil {
				return fmt.Errorf(
					"error in entry (%v, %v): %v",
					ii, jj,
					err,
				)
			}
		}
	}

	// All checks passed
	return nil
}

/*
Variables
Description:

	This function returns the list of variables in the matrix.
*/
func (vm VariableMatrix) Variables() []Variable {
	// Input Processing
	err := vm.Check()
	if err != nil {
		panic(err)
	}

	// Algorithm
	var variables []Variable
	for _, vmRow := range vm {
		for _, v := range vmRow {
			variables = append(variables, v)
		}
	}

	return UniqueVars(variables)
}

/*
Dims
Description:

	This function returns the dimensions of the variable matrix.
*/
func (vm VariableMatrix) Dims() []int {
	// Input Processing
	err := vm.Check()
	if err != nil {
		panic(err)
	}

	// Algorithm
	return []int{len(vm), len(vm[0])}
}

/*
Plus
Description:

	This function adds two variable matrices together.
*/
func (vm VariableMatrix) Plus(e interface{}) Expression {
	// Input Processing
	err := vm.Check()
	if err != nil {
		panic(err)
	}

	if IsExpression(e) {
		// Convert e to an expression
		eAsE, _ := e.(Expression)
		err = eAsE.Check()
		if err != nil {
			panic(
				fmt.Errorf("error in second argument to VariableMatrix.Plus: %v", err),
			)
		}

		err := CheckDimensionsInAddition(vm, eAsE)
		if err != nil {
			panic(err)
		}
	}

	// Algorithm
	switch right := e.(type) {
	case float64:
		// Use K case
		return vm.Plus(K(right))
	case K:
		// Create a new matrix of polynomials.
		var pmOut PolynomialMatrix
		for _, vmRow := range vm {
			var pmRow []Polynomial
			for _, v := range vmRow {
				pmRow = append(pmRow, v.Plus(right).(Polynomial))
			}
			pmOut = append(pmOut, pmRow)
		}
	}

	// panic if the type is not recognized
	panic(
		smErrors.UnsupportedInputError{
			FunctionName: "VariableMatrix.Plus",
			Input:        e,
		},
	)
}

/*
Multiply
Description:

	This function multiplies a variable matrix by another term.
*/
func (vm VariableMatrix) Multiply(e interface{}) Expression {
	// Input Processing
	err := vm.Check()
	if err != nil {
		panic(err)
	}

	if IsExpression(e) {
		// Convert e to an expression
		eAsE, _ := e.(Expression)
		err = eAsE.Check()
		if err != nil {
			panic(
				fmt.Errorf("error in second argument to VariableMatrix.Multiply: %v", err),
			)
		}

		err := CheckDimensionsInMultiplication(vm, eAsE)
		if err != nil {
			panic(err)
		}
	}

	// Algorithm
	switch right := e.(type) {
	case float64:
		// Use K case
		return vm.Multiply(K(right))
	case K:
		// Create a new matrix of polynomials.
		var mmOut MonomialMatrix
		for _, vmRow := range vm {
			var mmRow []Monomial
			for _, v := range vmRow {
				mmRow = append(mmRow, v.Multiply(right).(Monomial))
			}
			mmOut = append(mmOut, mmRow)
		}
	}

	// panic if the type is not recognized
	panic(
		smErrors.UnsupportedInputError{
			FunctionName: "VariableMatrix.Multiply",
			Input:        e,
		},
	)
}

/*
Transpose
Description:

	This function returns the transpose of the variable matrix.
*/
func (vm VariableMatrix) Transpose() Expression {
	// Input Processing
	err := vm.Check()
	if err != nil {
		panic(err)
	}

	// Create blank variable matrix
	nRows, nCols := vm.Dims()[0], vm.Dims()[1]
	var vmOut VariableMatrix
	for ii := 0; ii < nCols; ii++ {
		vmOut = append(vmOut, make([]Variable, nRows))
	}

	// Fill in the elements of the new matrix
	for ii := 0; ii < nRows; ii++ {
		for jj := 0; jj < nCols; jj++ {
			vmOut[jj][ii] = vm[ii][jj]
		}
	}

	return vmOut
}

/*
Comparison
Description:

	This function compares the variable matrix to another expression.
*/
func (vm VariableMatrix) Comparison(rightIn interface{}, sense ConstrSense) Constraint {
	// Input Processing
	err := vm.Check()
	if err != nil {
		panic(err)
	}

	if IsExpression(rightIn) {
		// Convert e to an expression
		rightAsE, _ := rightIn.(Expression)
		err = rightAsE.Check()
		if err != nil {
			panic(
				fmt.Errorf("error in second argument to VariableMatrix.Comparison: %v", err),
			)
		}

		err := CheckDimensionsInComparison(vm, rightAsE, sense)
		if err != nil {
			panic(err)
		}
	}

	// Algorithm
	switch right := rightIn.(type) {
	case float64:
		// Use K case
		return vm.Comparison(K(right), sense)
	case K:
		// Create a new matrix of polynomials.
		onesMat := OnesMatrix(vm.Dims()[0], vm.Dims()[1])
		var KAsDense mat.Dense
		KAsDense.Scale(float64(right), &onesMat)

		return MatrixConstraint{
			LeftHandSide:  vm,
			RightHandSide: DenseToKMatrix(KAsDense),
			Sense:         sense,
		}
	}

	// If the type is not recognized, panic
	panic(
		smErrors.UnsupportedInputError{
			FunctionName: "VariableMatrix.Comparison",
			Input:        rightIn,
		},
	)
}

/*
LessEq
Description:

	Returns a less than or equal to (<=) constraint between
	the VariableMatrix and another expression.
*/
func (vm VariableMatrix) LessEq(rightIn interface{}) Constraint {
	return vm.Comparison(rightIn, SenseLessThanEqual)
}

/*
GreaterEq
Description:

	Returns a greater than or equal to (>=) constraint between
	the VariableMatrix and another expression.
*/
func (vm VariableMatrix) GreaterEq(rightIn interface{}) Constraint {
	return vm.Comparison(rightIn, SenseGreaterThanEqual)
}

/*
Eq
Description:

	Returns an equality (==) constraint between
	the VariableMatrix and another expression.
*/
func (vm VariableMatrix) Eq(rightIn interface{}) Constraint {
	return vm.Comparison(rightIn, SenseEqual)
}

/*
DerivativeWrt
Description:

	This function returns the derivative of the variable matrix
	with respect to a given variable.
*/
func (vm VariableMatrix) DerivativeWrt(v Variable) Expression {
	// Input Processing
	err := vm.Check()
	if err != nil {
		panic(err)
	}

	// Algorithm
	kmOut := DenseToKMatrix(
		ZerosMatrix(vm.Dims()[0], vm.Dims()[1]),
	)

	// Iterate through matrix and if there is one
	// element that is equal to the variable, set the
	// corresponding element in the output matrix to 1.
	for ii, vmRow := range vm {
		for jj, vmElt := range vmRow {
			if vmElt.ID == v.ID {
				kmOut[ii][jj] = K(1)
			}
		}
	}

	return kmOut
}

/*
At
Description:

	This function returns the element of the variable matrix at the given indices.
*/
func (vm VariableMatrix) At(ii, jj int) ScalarExpression {
	// Input Processing
	err := vm.Check()
	if err != nil {
		panic(err)
	}

	// Algorithm
	return vm[ii][jj]
}

/*
Constant
Description:

	This function returns the constant value in the variable matrix.
	this should always be zero.
*/
func (vm VariableMatrix) Constant() mat.Dense {
	return ZerosMatrix(vm.Dims()[0], vm.Dims()[1])
}

/*
String
Description:

	This function returns a string representation of the variable matrix.
*/
func (vm VariableMatrix) String() string {
	// Input Processing
	err := vm.Check()
	if err != nil {
		panic(err)
	}

	// Algorithm
	var out string = "["
	for ii, vmRow := range vm {
		out += "["
		for jj, v := range vmRow {
			out += v.String()
			if jj < len(vmRow)-1 {
				out += ", "
			}
		}
		out += "]"
		if ii < len(vm)-1 {
			out += ", "
		}
	}
	out += "]"
	return out
}