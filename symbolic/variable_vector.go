package symbolic

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
)

/*
variable_vector.go
Description:
	The VariableVector type will represent a
*/

/*
VariableVector
Description:

	Represnts a variable in a optimization problem. The variable is
*/
type VariableVector struct {
	Elements []Variable
}

// =========
// Functions
// =========

/*
Length
Description:

	Returns the length of the vector of optimization variables.
*/
func (vv VariableVector) Length() int {
	return len(vv.Elements)
}

/*
Len
Description:

	This function is created to mirror the GoNum Vector API. Does the same thing as Length.
*/
func (vv VariableVector) Len() int {
	return vv.Length()
}

/*
At
Description:

	Mirrors the gonum api for vectors. This extracts the element of the variable vector at the index x.
*/
func (vv VariableVector) AtVec(idx int) ScalarExpression {
	// Constants

	// Algorithm
	return vv.Elements[idx]
}

/*
Variables
Description:

	Returns the slice of all variables in the vector.
*/
func (vv VariableVector) Variables() []Variable {
	return vv.Elements
}

/*
Constant
Description:

	Returns an all zeros vector as output from the method.
*/
func (vv VariableVector) Constant() mat.VecDense {
	zerosOut := ZerosVector(vv.Len())
	return zerosOut
}

/*
LinearCoeff
Description:

	Returns the matrix which is multiplied by Variables to get the current "expression".
	For a single vector, this is an identity matrix.
*/
func (vv VariableVector) LinearCoeff() mat.Dense {
	return Identity(vv.Len())
}

/*
Plus
Description:

	This member function computes the addition of the receiver vector var with the
	incoming vector expression ve.
*/
func (vv VariableVector) Plus(rightIn interface{}) Expression {
	// Constants
	// vvLen := vv.Len()

	// Processing Errors
	err := vv.Check()
	if err != nil {
		panic(err)
	}

	// Algorithm
	switch right := rightIn.(type) {

	default:
		errString := fmt.Sprintf(
			"Unrecognized expression type %T for addition of VariableVector vv.Plus(%v)!",
			right, right,
		)
		panic(fmt.Errorf(errString))
	}
}

/*
Mult
Description:

	This member function computest the multiplication of the receiver vector var with some
	incoming vector expression (may result in quadratic?).
*/
func (vv VariableVector) Mult(c float64) (VectorExpression, error) {
	return vv, fmt.Errorf("The Mult() method for VariableVector is not implemented yet!")
}

/*
Multiply
Description:

	Multiplication of a VariableVector with another expression.
*/
func (vv VariableVector) Multiply(rightIn interface{}) Expression {
	//Input Processing
	err := vv.Check()
	if err != nil {
		panic(err)
	}

	if IsExpression(rightIn) {
		rightAsE, _ := ToExpression(rightIn)
		err = CheckDimensionsInMultiplication(vv, rightAsE)
		if err != nil {
			panic(err)
		}
	}

	switch right := rightIn.(type) {
	default:
		panic(fmt.Errorf(
			"The input to VariableVector's Multiply() method (%v) has unexpected type: %T",
			right, rightIn,
		))
	}
}

/*
LessEq
Description:

	This method creates a less than or equal to vector constraint using the receiver as the left hand side and the
	input rhs as the right hand side if it is valid.
*/
func (vv VariableVector) LessEq(rightIn interface{}) Constraint {
	return vv.Comparison(rightIn, SenseLessThanEqual)
}

/*
GreaterEq
Description:

	This method creates a greater than or equal to vector constraint using the receiver as the left hand side and the
	input rhs as the right hand side if it is valid.
*/
func (vv VariableVector) GreaterEq(rightIn interface{}) Constraint {
	return vv.Comparison(rightIn, SenseGreaterThanEqual)
}

/*
Eq
Description:

	This method creates an equal to vector constraint using the receiver as the left hand side and the
	input rhs as the right hand side if it is valid.
*/
func (vv VariableVector) Eq(rightIn interface{}) Constraint {
	// Constants

	// Algorithm
	return vv.Comparison(rightIn, SenseEqual)

}

/*
Comparison
Description:

	This method creates a constraint of type sense between
	the receiver (as left hand side) and rhs (as right hand side) if both are valid.
*/
func (vv VariableVector) Comparison(rhs interface{}, sense ConstrSense) Constraint {
	// Constants

	// Algorithm
	switch rhsConverted := rhs.(type) {
	case KVector:
		// Check length of input and output.
		if vv.Len() != rhsConverted.Len() {
			panic(
				fmt.Errorf(
					"The two inputs to comparison '%v' must have the same dimension, but #1 has dimension %v and #2 has dimension %v!",
					sense,
					vv.Len(),
					rhsConverted.Len(),
				),
			)
		}
		return VectorConstraint{vv, rhsConverted, sense}
	case mat.VecDense:
		rhsAsKVector := KVector(rhsConverted)

		return vv.Comparison(rhsAsKVector, sense)

	case VariableVector:
		// Check length of input and output.
		if vv.Len() != rhsConverted.Len() {
			panic(
				fmt.Errorf(
					"The two inputs to comparison '%v' must have the same dimension, but #1 has dimension %v and #2 has dimension %v!",
					sense,
					vv.Len(),
					rhsConverted.Len(),
				),
			)
		}
		// Do Computation
		return VectorConstraint{vv, rhsConverted, sense}

	default:
		panic(fmt.Errorf("The Eq() method for VariableVector is not implemented yet for type %T!", rhs))
	}
}

func (vv VariableVector) Copy() VariableVector {
	// Constants

	// Algorithm
	newVarSlice := []Variable{}
	for varIndex := 0; varIndex < vv.Len(); varIndex++ {
		// Append to newVar Slice
		newVarSlice = append(newVarSlice, vv.Elements[varIndex])
	}

	return VariableVector{newVarSlice}

}

/*
Transpose
Description:

	This method creates the transpose of the current vector and returns it.
*/
func (vv VariableVector) Transpose() Expression {
	//vvCopy := vv.Copy()
	return vv //VarVectorTranspose(vvCopy) // TODO: Fix This.
}

/*
Dims
Description:

	Dimensions of the variable vector.
*/
func (vv VariableVector) Dims() []int {
	return []int{vv.Len(), 1}
}

/*
Check
Description:

	Checks whether or not the VariableVector has a sensible initialization.
*/
func (vv VariableVector) Check() error {
	// Check that each variable is properly defined
	for ii, element := range vv.Elements {
		err := element.Check()
		if err != nil {
			return fmt.Errorf(
				"element %v has an issue: %v",
				ii, err,
			)
		}
	}

	// If nothing was thrown, then return nil!
	return nil
}
