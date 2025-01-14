package symbolic

import (
	"fmt"

	"github.com/MatProGo-dev/SymbolicMath.go/smErrors"
)

/*
Expression
Description:

	This interface should be implemented by and ScalarExpression and VectorExpression
*/
type Expression interface {
	// NumVars returns the number of variables in the expression
	Variables() []Variable

	// Dims returns a slice describing the true dimensions of a given expression (scalar, vector, or matrix)
	Dims() []int

	// Plus adds the current expression to another and returns the resulting
	// expression
	Plus(rightIn interface{}) Expression

	// Minus subtracts an expression from the current one and returns the resulting
	// expression
	Minus(rightIn interface{}) Expression

	// Multiply multiplies the current expression to another and returns the
	// resulting expression
	Multiply(rightIn interface{}) Expression

	// Transpose transposes the given expression
	Transpose() Expression

	// LessEq returns a less than or equal to (<=) constraint between the
	// current expression and another
	LessEq(rightIn interface{}) Constraint

	// GreaterEq returns a greater than or equal to (>=) constraint between the
	// current expression and another
	GreaterEq(rightIn interface{}) Constraint

	// Eq returns an equality (==) constraint between the current expression
	// and another
	Eq(rightIn interface{}) Constraint

	// Comparison
	Comparison(rightIn interface{}, sense ConstrSense) Constraint

	// DerivativeWrt returns the derivative of the expression with respect to the input variable vIn.
	DerivativeWrt(vIn Variable) Expression

	// Check
	Check() error

	// String returns a string representation of the expression
	String() string

	// Substitute returns the expression with the variable vIn replaced with the expression eIn
	Substitute(vIn Variable, eIn ScalarExpression) Expression

	// SubstituteAccordingToMap returns the expression with the variables in the map replaced with the corresponding expressions
	SubstituteAccordingTo(subMap map[Variable]Expression) Expression

	// Power
	// Raises the scalar expression to the power of the input integer
	Power(exponent int) Expression

	// At returns the value at the given row and column index
	At(ii, jj int) ScalarExpression
}

/*
NumVariables
Description:

	The number of distinct variables.
*/
func NumVariables(e Expression) int {
	return len(e.Variables())
}

/*
VariableIDs
Description:

	Returns a list of ids associated with each variable.
*/
func VariableIDs(e Expression) []uint64 {
	vSlice := e.Variables()

	var idSlice []uint64
	for _, v := range vSlice {
		idSlice = append(idSlice, v.ID)
	}

	return idSlice
}

/*
IsExpression
Description:

	Tests whether or not the input variable is one of the expression types.
*/
func IsExpression(e interface{}) bool {
	return IsScalarExpression(e) || IsVectorExpression(e) || IsMatrixExpression(e)
}

func ToExpression(e interface{}) (Expression, error) {
	switch {
	case IsScalarExpression(e):
		return ToScalarExpression(e)
	case IsVectorExpression(e):
		return ToVectorExpression(e)
	case IsMatrixExpression(e):
		return ToMatrixExpression(e)
	}

	// If the input is not a valid expression, panic
	panic(
		smErrors.UnsupportedInputError{
			FunctionName: "ToExpression",
			Input:        e,
		},
	)
}

/*
Minus
Description:

	subtracts the current expression from another and returns the resulting expression
*/
func Minus(left, right Expression) Expression {
	return left.Plus(
		right.Multiply(-1.0),
	)
}

/*
IsLinear
Description:

	Determines whether an input object is a
	valid linear expression.
	In math, this means that the polynomial like expression
	has a degree less than or equal to 1.
*/
func IsLinear(e Expression) bool {
	// Input Processing
	if !IsPolynomialLike(e) {
		return false // Not a polynomial like expression, so it can't be linear
	}

	eAsPL, _ := ToPolynomialLike(e)

	return eAsPL.Degree() <= 1
}

/*
IsQuadratic
Description:

	Determines whether or not an input object is a
	valid Quadratic Expression.
	In math, this means that the polynomial like expression
	has a degree less than or equal to 2.
*/
func IsQuadratic(e Expression) bool {
	// Input Processing
	if !IsPolynomialLike(e) {
		return false // Not a polynomial like expression, so it can't be linear
	}

	eAsPL, _ := ToPolynomialLike(e)

	return eAsPL.Degree() <= 2
}

/*
HStack
Description:

	Stacks the input expressions horizontally.
*/
func HStack(eIn ...Expression) Expression {
	// Input Checking

	// TODO: Panic if there are 0 expressions in the input
	if len(eIn) == 0 {
		panic(
			fmt.Errorf("HStack: There must be at least one expression in the input; received 0"),
		)
	}

	// Check that all the expressions have the same number of rows
	var mlSlice []smErrors.MatrixLike // First convert expression slice to matrix like slice
	for _, e := range eIn {
		mlSlice = append(mlSlice, e)
	}

	err := smErrors.CheckDimensionsInHStack(mlSlice...)
	if err != nil {
		panic(err)
	}

	// Setup
	var nCols []int
	for _, e := range eIn {
		nCols = append(nCols, e.Dims()[1])
	}

	// Create the resulting Matrix's shape
	var result [][]ScalarExpression
	for rowIndex := 0; rowIndex < eIn[0].Dims()[0]; rowIndex++ {
		var tempRow []ScalarExpression
		for stackIndexII := 0; stackIndexII < len(eIn); stackIndexII++ {
			nCols_ii := nCols[stackIndexII]
			// Add all of the columns from the current expression to the row
			for colIndex := 0; colIndex < nCols_ii; colIndex++ {
				tempRow = append(tempRow, eIn[stackIndexII].At(rowIndex, colIndex))
			}
		}
		// Add the row to the result
		result = append(result, tempRow)
	}

	// Return the simplified form of the expression
	return ConcretizeExpression(result)
}

/*
VStack
Description:

	Stacks the input expressions vertically.
*/
func VStack(eIn ...Expression) Expression {
	// Input Checking

	// Panic if there are 0 expressions in the input
	if len(eIn) == 0 {
		panic(
			fmt.Errorf("VStack: There must be at least one expression in the input; received 0"),
		)
	}

	// Check that all the expressions have the same number of columns
	var mlSlice []smErrors.MatrixLike // First convert expression slice to matrix like slice
	for _, e := range eIn {
		mlSlice = append(mlSlice, e)
	}

	err := smErrors.CheckDimensionsInVStack(mlSlice...)
	if err != nil {
		panic(err)
	}

	// Setup

	// Create the resulting Matrix's shape
	var result [][]ScalarExpression
	for stackIndexII := 0; stackIndexII < len(eIn); stackIndexII++ {
		// Each row will be made from the rows of the current expression
		eII := eIn[stackIndexII]
		for rowIndex := 0; rowIndex < eII.Dims()[0]; rowIndex++ {
			var tempRow []ScalarExpression
			for colIndex := 0; colIndex < eII.Dims()[1]; colIndex++ {
				tempRow = append(tempRow, eII.At(rowIndex, colIndex))
			}
			// Add the row to the result
			result = append(result, tempRow)
		}
	}

	// Return the simplified form of the expression
	return ConcretizeExpression(result)
}

/*
ConcretizeExpression
Description:

	Converts the input expression to a valid type that implements "Expression".
*/
func ConcretizeExpression(e interface{}) Expression {
	// Input Processing

	// Convert
	var (
		concrete Expression
	)
	switch concreteVal := e.(type) {
	case []ScalarExpression:
		concreteVectorE := ConcretizeVectorExpression(concreteVal)
		// If vector expression is a scalar (i.e., has 1 row), return the scalar expression
		if concreteVectorE.Dims()[0] == 1 {
			concrete = concreteVectorE.At(0, 0)
		} else {
			concrete = concreteVectorE
		}

	case [][]ScalarExpression:
		concreteMatrixE := ConcretizeMatrixExpression(concreteVal)
		// If matrix expression is a scalar (i.e., has 1 row and 1 column), return the scalar expression
		switch {
		case concreteMatrixE.Dims()[0] == 1 && concreteMatrixE.Dims()[1] == 1: // If the matrix is a scalar
			concrete = concreteMatrixE.At(0, 0)
		case concreteMatrixE.Dims()[1] == 1: // If the matrix is a column vector
			interm := make([]ScalarExpression, concreteMatrixE.Dims()[0])
			for ii := 0; ii < concreteMatrixE.Dims()[0]; ii++ {
				interm[ii] = concreteMatrixE.At(ii, 0)
			}
			concrete = ConcretizeVectorExpression(interm)
		default:
			concrete = concreteMatrixE
		}
	default:
		panic(
			smErrors.UnsupportedInputError{
				FunctionName: "ConcretizeExpression",
				Input:        e,
			},
		)
	}

	return concrete
}
