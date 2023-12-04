package symbolic

import "fmt"

/*
polynomial.go
Description:
	This file defines the function associated with the Polynomial object.
*/

/*
Type Definition
*/
type Polynomial struct {
	Monomials []Monomial
}

// Member Methods

/*
Check
Description:

	Verifies that all elements of the polynomial are defined correctly.
*/
func (p Polynomial) Check() error {
	for _, monomial := range p.Monomials {
		err := monomial.Check()
		if err != nil {
			return err
		}
	}

	// All checks passed
	return nil
}

/*
Variables
Description:

	The unique variables used to define the polynomial.
*/
func (p Polynomial) Variables() []Variable {
	var variables []Variable // The variables in the polynomial
	for _, monomial := range p.Monomials {
		variables = append(variables, monomial.Variables()...)
	}
	return UniqueVars(variables)
}

/*
Dims
Description:

	The scalar polynomial should have dimensions [1,1].
*/
func (p Polynomial) Dims() []int {
	return []int{1, 1}
}

/*
Plus
Description:

	Defines an addition between the polynomial and another expression.
*/
func (p Polynomial) Plus(e interface{}) Expression {
	// Input Processing
	err := p.Check()
	if err != nil {
		panic(err)
	}

	if IsExpression(e) {
		eAsE, _ := ToExpression(e)
		err := CheckDimensionsInAddition(p, eAsE)
		if err != nil {
			panic(err)
		}
	}

	// Constants
	switch right := e.(type) {
	case float64:
		return p.Plus(K(right))
	case K:
		pCopy := p

		// Algorithm
		constantIndex := pCopy.ConstantMonomialIndex()
		if constantIndex != -1 {
			// Monomial does not contain a constant,
			// so add a new monomial.
			rightAsMonom := right.ToMonomial()
			pCopy.Monomials = append(pCopy.Monomials, rightAsMonom)
		} else {
			// Monomial does contain a constant, so
			// modify the monomial which represents that constant.
			newMonomial := pCopy.Monomials[constantIndex]
			newMonomial.Coefficient += float64(right)
			pCopy.Monomials[constantIndex] = newMonomial
		}
		return pCopy

	case Variable:
		pCopy := p

		// Check to see if the variable is already in the polynomial
		variableIndex := pCopy.VariableMonomialIndex(right)
		if variableIndex != -1 {
			// Monomial does not contain the variable,
			// so add a new monomial.
			rightAsMonom := right.ToMonomial()
			pCopy.Monomials = append(pCopy.Monomials, rightAsMonom)
		} else {
			// Monomial does contain the variable, so
			// modify the monomial which represents that variable.
			newMonomial := pCopy.Monomials[variableIndex]
			newMonomial.Coefficient += 1.0
			pCopy.Monomials[variableIndex] = newMonomial
		}
	}

	// Unrecognized response is a panic
	panic(
		fmt.Errorf("Unexpected type of right in the Plus() method: %T (%v)", e, e),
	)
}

/*
ConstantMonomialIndex
Description:

	Returns the index of the monomial in the polynomial which is a constant.
	If none are found, then this returns -1.
*/
func (p Polynomial) ConstantMonomialIndex() int {
	// Input Processing
	err := p.Check()
	if err != nil {
		panic(err)
	}

	// Algorithm
	for ii, monomial := range p.Monomials {
		if monomial.IsConstant() {
			return ii
		}
	}

	// No constant monomial found
	return -1
}

/*
VariableMonomialIndex
Description:

	Returns the index of the monomial which represents the variable given as vIn.
*/
func (p Polynomial) VariableMonomialIndex(vIn Variable) int {
	// Input Processing
	err := p.Check()
	if err != nil {
		panic(err)
	}

	err = vIn.Check()
	if err != nil {
		panic(err)
	}

	// Algorithm
	for ii, monomial := range p.Monomials {
		if monomial.IsVariable(vIn) {
			return ii
		}
	}

	// No variable monomial found
	return -1
}

/*
Multiply
Description:

	Implements the multiplication operator between a polynomial and another expression.
*/
func (p Polynomial) Multiply(e interface{}) Expression {
	// Input Processing
	err := p.Check()
	if err != nil {
		panic(err)
	}

	if IsExpression(e) {
		eAsE, _ := ToExpression(e)
		err := CheckDimensionsInMultiplication(p, eAsE)
		if err != nil {
			panic(err)
		}
	}

	// Algorithm
	switch right := e.(type) {
	case float64:
		return p.Multiply(K(right))
	case K:
		pCopy := p
		for ii, _ := range pCopy.Monomials {
			product_ii := pCopy.Monomials[ii].Multiply(right)
			pCopy.Monomials[ii] = product_ii.(Monomial) // Convert to Monomial
		}
		return pCopy
	case Variable:
		pCopy := p
		for ii, _ := range pCopy.Monomials {
			product_ii := pCopy.Monomials[ii].Multiply(right)
			pCopy.Monomials[ii] = product_ii.(Monomial) // Convert to Monomial
		}
		return pCopy
	case Monomial:
		pCopy := p
		for ii, _ := range pCopy.Monomials {
			product_ii := pCopy.Monomials[ii].Multiply(right)
			pCopy.Monomials[ii] = product_ii.(Monomial)
		}
		return pCopy
	}

	// Unrecognized response is a panic
	panic(
		fmt.Errorf("Unexpected type of right in the Multiply() method: %T (%v)", e, e),
	)
}

/*
Transpose
Description:

	The transpose operator when applied to a scalar is just the same scalar object.
*/
func (p Polynomial) Transpose() Expression {
	return p
}

/*
Comparison
Description:

	Creates a constraint between the polynomial and another expression
	of the sense provided in Sense.
*/
func (p Polynomial) Comparison(rightIn interface{}, sense ConstrSense) Constraint {
	// Input Processing
	err := p.Check()
	if err != nil {
		panic(err)
	}

	right, err := ToScalarExpression(rightIn)
	if err != nil {
		panic(err)
	}

	// Algorithm
	return ScalarConstraint{p, right, sense}
}

/*
LessEq
Description:

	Creates a less than equal constraint between the polynomial and another expression.
*/
func (p Polynomial) LessEq(rightIn interface{}) Constraint {
	return p.Comparison(rightIn, SenseLessThanEqual)
}

/*
GreaterEq
Description:

	Creates a greater than equal constraint between the polynomial and another expression.
*/
func (p Polynomial) GreaterEq(rightIn interface{}) Constraint {
	return p.Comparison(rightIn, SenseGreaterThanEqual)
}

/*
Eq
Description:

	Creates an equality constraint between the polynomial and another expression.
*/
func (p Polynomial) Eq(rightIn interface{}) Constraint {
	return p.Comparison(rightIn, SenseEqual)
}

/*
Constant
Description:

	Retrieves the constant component of the polynomial if there is one.
*/
func (p Polynomial) Constant() float64 {
	// Input Processing
	err := p.Check()
	if err != nil {
		panic(err)
	}

	// Algorithm
	constantIndex := p.ConstantMonomialIndex()
	if constantIndex == -1 {
		return 0.0
	} else {
		return p.Monomials[constantIndex].Coefficient
	}
}
