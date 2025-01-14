package symbolic

import (
	"fmt"

	"github.com/MatProGo-dev/SymbolicMath.go/smErrors"
	"gonum.org/v1/gonum/mat"
)

/*
Integer constants representing commonly used numbers. Makes for better
readability
*/
const (
	Zero     = K(0)
	One      = K(1)
	Infinity = K(1e100)
)

// K is a constant expression type for an MIP. K for short ¯\_(ツ)_/¯
type K float64

/*
Check
Description:

	Checks to make sure that the constant is initialized properly.
	Constants are always initialized properly, so this should always return
	no error.
*/
func (c K) Check() error {
	return nil
}

/*
Variables
Description:

	Shares all variables included in the expression that is K.
	It is a constant, so there are none.
*/
func (c K) Variables() []Variable {
	return []Variable{}
}

// Constant returns the constant additive value in the expression. For
// constants, this is just the constants value
func (c K) Constant() float64 {
	return float64(c)
}

/*
LinearCoeff
Description

	Returns the coefficient of the linear term in the expression. For a constant,
	this is always a matrix of zeros.
*/
func (c K) LinearCoeff(wrt ...[]Variable) mat.VecDense {
	// Constants

	// Input Processing
	var wrtVars []Variable = []Variable{}
	if len(wrt) > 0 {
		// If the user provided a slice of variables, use that instead
		wrtVars = wrt[0]
	}

	if len(wrtVars) == 0 {
		// If the user didn't provide any variables, then panic!
		// We cannot construct zero length vectors in gonum
		panic(
			smErrors.EmptyLinearCoeffsError{Expression: c},
		)
	}

	// Algorithm
	return ZerosVector(len(wrtVars))
}

/*
Plus
Description:

	adds the current expression to another and returns the resulting expression
*/
func (c K) Plus(rightIn interface{}) Expression {
	// Input Processing
	if IsExpression(rightIn) {
		rightAsE, _ := ToExpression(rightIn)
		err := rightAsE.Check()
		if err != nil {
			panic(fmt.Errorf("error in second argument to %v.Plus: %v", c, err))
		}

		// Dimension checks should be fine because c is a scalar.
		//err = smErrors.CheckDimensionsInAddition(c, rightAsE)
		//if err != nil {
		//	panic(err)
		//}
	}

	// Switching based on input type
	switch right := rightIn.(type) {
	case int:
		return c.Plus(K(right))
	case float64:
		return c.Plus(K(right))
	case K:
		return K(c.Constant() + right.Constant())
	case Variable:
		return right.Plus(c)
	case Monomial:
		return right.Plus(c)
	case Polynomial:
		return right.Plus(c)
	case mat.VecDense:
		return c.Plus(VecDenseToKVector(right))
	case *mat.VecDense:
		return c.Plus(VecDenseToKVector(*right))
	case KVector, VariableVector, MonomialVector, PolynomialVector:
		// Convert to VectorExpression
		ve, _ := ToVectorExpression(right)
		return ve.Plus(c)
	case mat.Dense:
		return c.Plus(DenseToKMatrix(right))
	case *mat.Dense:
		return c.Plus(DenseToKMatrix(*right))
	case KMatrix, VariableMatrix, MonomialMatrix, PolynomialMatrix:
		// Convert to MatrixExpression
		me, _ := ToMatrixExpression(right)
		return me.Plus(c)
	}

	// Default response is a panic
	panic(
		smErrors.UnsupportedInputError{
			FunctionName: "K.Plus",
			Input:        rightIn,
		},
	)
}

/*
Minus
Description:

	This function subtracts the current expression from another and returns the resulting expression.
*/
func (c K) Minus(rightIn interface{}) Expression {
	// Input Processing
	if IsExpression(rightIn) {
		rightAsE, _ := ToExpression(rightIn)
		err := rightAsE.Check()
		if err != nil {
			panic(fmt.Errorf("error in second argument to %v.Minus: %v", c, err))
		}

		//err = smErrors.CheckDimensionsInSubtraction(c, rightAsE)
		//if err != nil {
		//	panic(err)
		//}

		// Use Minus function
		return Minus(c, rightAsE)
	}

	// Switching based on input type
	switch right := rightIn.(type) {
	case float64:
		return c.Minus(K(right))
	}

	// Default response is a panic
	panic(
		smErrors.UnsupportedInputError{
			FunctionName: "K.Minus",
			Input:        rightIn,
		},
	)
}

// LessEq returns a less than or equal to (<=) constraint between the
// current expression and another
func (c K) LessEq(rightIn interface{}) Constraint {
	return c.Comparison(rightIn, SenseLessThanEqual)
}

// GreaterEq returns a greater than or equal to (>=) constraint between the
// current expression and another
func (c K) GreaterEq(rightIn interface{}) Constraint {
	return c.Comparison(rightIn, SenseGreaterThanEqual)
}

// Eq returns an equality (==) constraint between the current expression
// and another
func (c K) Eq(rightIn interface{}) Constraint {
	return c.Comparison(rightIn, SenseEqual)
}

/*
Comparison
Description:

	This method compares the receiver with expression rhs in the sense provided by sense.
*/
func (c K) Comparison(rhsIn interface{}, sense ConstrSense) Constraint {
	// InputProcessing
	if IsExpression(rhsIn) {
		rhs, _ := ToExpression(rhsIn)
		err := rhs.Check()
		if err != nil {
			panic(err)
		}

		// Normally, we would check the dimensions here, but since K is a scalar, we don't need to.
	}

	// Constants

	// Algorithm
	switch right := rhsIn.(type) {
	case float64:
		// Use the version of Comparison for K
		return c.Comparison(K(right), sense)
	case K, Variable, Monomial, Polynomial:
		// Cast right to scalar expression
		se, _ := ToScalarExpression(right)
		return ScalarConstraint{c, se, sense}
	case mat.VecDense:
		// Convert to KVector
		return c.Comparison(VecDenseToKVector(right), sense)
	case *mat.VecDense:
		// Convert to KVector
		return c.Comparison(VecDenseToKVector(*right), sense)
	case KVector:
		// Transform right into a KVector as well
		var kAsKVector KVector
		for ii := 0; ii < len(right); ii++ {
			kAsKVector = append(kAsKVector, c)
		}

		// Create vector comparison
		return VectorConstraint{
			LeftHandSide:  kAsKVector,
			RightHandSide: right,
			Sense:         sense,
		}
	case VariableVector:
		// Transform right into a KVector as well
		var kAsKVector KVector
		for ii := 0; ii < len(right); ii++ {
			kAsKVector = append(kAsKVector, c)
		}

		// Create vector comparison
		return VectorConstraint{
			LeftHandSide:  kAsKVector,
			RightHandSide: right,
			Sense:         sense,
		}
	}

	// Panic if the input is not recognized
	panic(
		smErrors.UnsupportedInputError{
			FunctionName: "K.Comparison",
			Input:        rhsIn,
		},
	)

}

/*
Multiply
Description:

	This method multiplies the input constant by another expression.
*/
func (c K) Multiply(term1 interface{}) Expression {
	// Constants

	// Input Processing
	if IsExpression(term1) {
		// Cast to expression
		term1AsE, _ := ToExpression(term1)
		err := term1AsE.Check()
		if err != nil {
			panic(err)
		}

		//// Check dimensions (not necessary for scalar K)
		//err := smErrors.CheckDimensionsInMultiplication(c, term1AsE)
		//if err != nil {
		//	panic(err)
		//}
	}

	// Algorithm
	switch right := term1.(type) {
	case float64:
		return c.Multiply(K(right))
	case int:
		return c.Multiply(K(right))
	case K:
		return c * right
	case Variable:
		return right.Multiply(c)
	case Monomial:
		return right.Multiply(c)
	case Polynomial:
		return right.Multiply(c)
	}

	// Unrecornized response is a panic
	panic(
		smErrors.UnsupportedInputError{
			FunctionName: "K.Multiply",
			Input:        term1,
		},
	)
}

func (c K) Dims() []int {
	return []int{1, 1} // Signifies scalar
}

func (c K) Transpose() Expression {
	return c
}

/*
ToMonomial
Description:

	Converts the constant into a monomial.
*/
func (c K) ToMonomial() Monomial {
	return Monomial{
		Coefficient:     float64(c),
		VariableFactors: []Variable{},
		Exponents:       []int{},
	}
}

/*
ToPolynomial
Description:

	Converts the constant into a polynomial.
*/
func (c K) ToPolynomial() Polynomial {
	return Polynomial{
		Monomials: []Monomial{c.ToMonomial()},
	}
}

/*
DerivativeWrt
Description:

	Computes the derivative of a constant, which should be 0 for any constant.
*/
func (c K) DerivativeWrt(vIn Variable) Expression {
	return Zero
}

/*
Degree
Description:

	The degree of a constant is always 0.
*/
func (c K) Degree() int {
	return 0
}

/*
String
Description:

	Returns a string representation of the constant.
*/
func (c K) String() string {
	return fmt.Sprintf("%v", float64(c))
}

/*
Substitute
Description:

	Substitutes the variable vIn with the expression eIn.
*/
func (c K) Substitute(vIn Variable, eIn ScalarExpression) Expression {
	return c
}

/*
SubstituteAccordingTo
Description:

	Substitutes the variables in the map with the corresponding expressions.
*/
func (c K) SubstituteAccordingTo(subMap map[Variable]Expression) Expression {
	return c
}

/*
Power
Description:

	Computes the power of the constant.
*/
func (c K) Power(exponent int) Expression {
	return ScalarPowerTemplate(c, exponent)
}

/*
At
Description:

	Returns the value at the given row and column index.

Note:

	For a constant, this is always the constant itself.
	The value of ii and jj should always be 0.
*/
func (c K) At(ii, jj int) ScalarExpression {
	// Input Processing

	// Check to see whether or not the index is valid.
	err := smErrors.CheckIndexOnMatrix(ii, jj, c)
	if err != nil {
		panic(err)
	}

	return c
}
