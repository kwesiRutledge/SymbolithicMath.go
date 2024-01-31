package symbolic_test

import (
	"fmt"
	getKVector "github.com/MatProGo-dev/SymbolicMath.go/get/KVector"
	"github.com/MatProGo-dev/SymbolicMath.go/smErrors"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"strings"
	"testing"
)

/*
monomial_vector_test.go
Description:

	Tests for the functions mentioned in the monomial_vector.go file.
*/

/*
TestMonomialVector_Check1
Description:

	Tests that the Check() method properly catches an improperly initialized
	vector of Monomials (i.e., no monomials are given).
*/
func TestMonomialVector_Check1(t *testing.T) {
	// Constants
	mv := symbolic.MonomialVector{}
	expectedError := smErrors.EmptyVectorError{mv}

	// Test
	err := mv.Check()
	if err.Error() != expectedError.Error() {
		t.Errorf(
			"expected Check() to return \"%v\"; received %v",
			expectedError,
			err,
		)
	}
}

/*
TestMonomialVector_Check2
Description:

	Tests that the Check() method properly catches an improperly initialized
	vector of Monomials (i.e., a monomial is given with an improper number of
	degrees).
*/
func TestMonomialVector_Check2(t *testing.T) {
	// Constants
	v1 := symbolic.NewVariable()
	m1 := symbolic.Monomial{
		Coefficient:     1.0,
		VariableFactors: []symbolic.Variable{v1},
		Exponents:       []int{1, 2},
	}
	mv := symbolic.MonomialVector{m1}
	expectedError := fmt.Errorf(
		"the number of degrees (%v) does not match the number of variables (%v)",
		len(m1.Exponents),
		len(m1.VariableFactors),
	)

	// Test
	err := mv.Check()
	if !strings.Contains(
		err.Error(),
		expectedError.Error(),
	) {
		t.Errorf(
			"expected Check() to return \"%v\"; received %v",
			expectedError,
			err,
		)
	}
}

/*
TestMonomialVector_Check3
Description:

	Verifies that the Check() method returns nil when a constant is
	given as a monomial.
*/
func TestMonomialVector_Check3(t *testing.T) {
	// Constants
	m1 := symbolic.Monomial{
		Coefficient:     3.14,
		VariableFactors: []symbolic.Variable{},
		Exponents:       []int{},
	}
	mv := symbolic.MonomialVector{m1}

	// Test
	if mv.Check() != nil {
		t.Errorf(
			"expected Check() to return nil; received %v",
			mv.Check(),
		)
	}
}

/*
TestMonomialVector_Variables1
Description:

	Verifies that the Variables() method panics when an improperly initialized
	vector of Monomials is given.
*/
func TestMonomialVector_Variables1(t *testing.T) {
	// Constants
	mv := symbolic.MonomialVector{}

	// Test
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf(
				"Expected mv.Variables() to panic; received %v",
				mv.Variables(),
			)
		}
	}()

	mv.Variables()
}

/*
TestMonomialVector_Variables2
Description:

	Verifies that the Variables() method returns the correct value when a
	vector of monomials of length 2 is given, with no repeated variables.
*/
func TestMonomialVector_Variables2(t *testing.T) {
	// Constants
	v1 := symbolic.NewVariable()
	v2 := symbolic.NewVariable()
	m1 := symbolic.Monomial{
		Coefficient:     1.0,
		VariableFactors: []symbolic.Variable{v1},
		Exponents:       []int{1},
	}
	m2 := symbolic.Monomial{
		Coefficient:     1.0,
		VariableFactors: []symbolic.Variable{v2},
		Exponents:       []int{1},
	}
	mv := symbolic.MonomialVector{m1, m2}

	// Test
	variables := mv.Variables()
	if len(variables) != 2 {
		t.Errorf(
			"expected len(variables) to be 2; received %v",
			len(variables),
		)
	}
}

/*
TestMonomialVector_Variables3
Description:

	Verifies that the Variables() method returns the correct value when a
	vector of monomials of length 2 is given, with repeated variables in
	each monomial.
*/
func TestMonomialVector_Variables3(t *testing.T) {
	// Constants
	v1 := symbolic.NewVariable()
	m1 := symbolic.Monomial{
		Coefficient:     1.0,
		VariableFactors: []symbolic.Variable{v1},
		Exponents:       []int{1},
	}
	m2 := symbolic.Monomial{
		Coefficient:     1.0,
		VariableFactors: []symbolic.Variable{v1},
		Exponents:       []int{1},
	}
	mv := symbolic.MonomialVector{m1, m2}

	// Test
	variables := mv.Variables()
	if len(variables) != 1 {
		t.Errorf(
			"expected len(variables) to be 1; received %v",
			len(variables),
		)
	}
}

/*
TestMonomialVector_Len1
Description:

	Verifies that the Len() method returns the correct value when a
	vector of monomials of length 2 is given.
*/
func TestMonomialVector_Len1(t *testing.T) {
	// Constants
	v1 := symbolic.NewVariable()
	m1 := symbolic.Monomial{
		Coefficient:     1.0,
		VariableFactors: []symbolic.Variable{v1},
		Exponents:       []int{1},
	}
	m2 := symbolic.Monomial{
		Coefficient:     1.0,
		VariableFactors: []symbolic.Variable{v1},
		Exponents:       []int{1},
	}
	mv := symbolic.MonomialVector{m1, m2}

	// Test
	if mv.Len() != 2 {
		t.Errorf(
			"expected mv.Len() to be 2; received %v",
			mv.Len(),
		)
	}
}

/*
TestMonomialVector_Dims1
Description:

	Verifies that the Dims() method returns the correct value when a
	vector of monomials of length 2 is given.
*/
func TestMonomialVector_Dims1(t *testing.T) {
	// Constants
	v1 := symbolic.NewVariable()
	m1 := symbolic.Monomial{
		Coefficient:     1.0,
		VariableFactors: []symbolic.Variable{v1},
		Exponents:       []int{1},
	}
	m2 := symbolic.Monomial{
		Coefficient:     1.0,
		VariableFactors: []symbolic.Variable{v1},
		Exponents:       []int{1},
	}
	mv := symbolic.MonomialVector{m1, m2}

	// Test
	if mv.Dims()[0] != 2 || mv.Dims()[1] != 1 {
		t.Errorf(
			"expected mv.Dims() to be [2,1]; received %v",
			mv.Dims(),
		)
	}
}

/*
TestMonomialVector_Constant1
Description:

	This test verifies that the Constant() method throws a panic
	when an improperly initialized vector of monomials is given.
*/
func TestMonomialVector_Constant1(t *testing.T) {
	// Constants
	mv := symbolic.MonomialVector{}

	// Test
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf(
				"Expected mv.Constant() to panic; received %v",
				mv.Constant(),
			)
		}
	}()

	mv.Constant()
}

/*
TestMonomialVector_Constant2
Description:

	This test verifies that the Constant() method returns the correct
	value (all zeros) when a vector of 4 monomials (each with nonzero
	number of variablefactors is given).
*/
func TestMonomialVector_Constant2(t *testing.T) {
	// Constants
	v1 := symbolic.NewVariable()
	v2 := symbolic.NewVariable()
	v3 := symbolic.NewVariable()
	v4 := symbolic.NewVariable()
	m1 := symbolic.Monomial{
		Coefficient:     1.0,
		VariableFactors: []symbolic.Variable{v3},
		Exponents:       []int{4},
	}
	m2 := symbolic.Monomial{
		Coefficient:     1.0,
		VariableFactors: []symbolic.Variable{v1},
		Exponents:       []int{1},
	}
	m3 := symbolic.Monomial{
		Coefficient:     1.0,
		VariableFactors: []symbolic.Variable{v1, v2},
		Exponents:       []int{1, 1},
	}
	m4 := symbolic.Monomial{
		Coefficient:     1.0,
		VariableFactors: []symbolic.Variable{v1, v2, v3, v4},
		Exponents:       []int{1, 1, 1, 1},
	}
	mv := symbolic.MonomialVector{m1, m2, m3, m4}

	// Test
	constant5 := mv.Constant()
	for ii := 0; ii < len(mv); ii++ {
		if constant5.AtVec(ii) != 0 {
			t.Errorf(
				"Expected mv.Constant() to be [0,0,0,0]; received %v at index %v",
				constant5.AtVec(ii),
				ii,
			)
		}

	}
}

/*
TestMonomialVector_Constant3
Description:

	This test verifies that the Constant() method returns the correct
	value (first two element nonzero) when a vector of 4 monomials
	is given and the first two elements are constants.
*/
func TestMonomialVector_Constant3(t *testing.T) {
	// Constants
	v1 := symbolic.NewVariable()
	v2 := symbolic.NewVariable()
	v3 := symbolic.NewVariable()
	v4 := symbolic.NewVariable()
	m1 := symbolic.Monomial{
		Coefficient:     3.14,
		VariableFactors: []symbolic.Variable{},
		Exponents:       []int{},
	}
	m2 := symbolic.Monomial{
		Coefficient:     3.14,
		VariableFactors: []symbolic.Variable{},
		Exponents:       []int{},
	}
	m3 := symbolic.Monomial{
		Coefficient:     1.0,
		VariableFactors: []symbolic.Variable{v1, v2},
		Exponents:       []int{1, 1},
	}
	m4 := symbolic.Monomial{
		Coefficient:     1.0,
		VariableFactors: []symbolic.Variable{v1, v2, v3, v4},
		Exponents:       []int{1, 1, 1, 1},
	}
	mv := symbolic.MonomialVector{m1, m2, m3, m4}

	// Test
	constant5 := mv.Constant()
	for ii := 0; ii < len(mv); ii++ {
		if ii < 2 {
			if constant5.AtVec(ii) != 3.14 {
				t.Errorf(
					"Expected mv.Constant() to be [1,1,0,0]; received %v at index %v",
					constant5.AtVec(ii),
					ii,
				)
			}
		} else {
			if constant5.AtVec(ii) != 0 {
				t.Errorf(
					"Expected mv.Constant() to be [1,1,0,0]; received %v at index %v",
					constant5.AtVec(ii),
					ii,
				)
			}
		}
	}
}

/*
TestMonomialVector_Plus1
Description:

	Verifies that the Plus() method throws a panic when an improperly
	initialized vector of monomials is given.
*/
func TestMonomialVector_Plus1(t *testing.T) {
	// Constants
	mv := symbolic.MonomialVector{}

	// Test
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf(
				"Expected mv.Plus(3.14) to panic; received %v",
				mv.Plus(3.14),
			)
		}
	}()

	mv.Plus(3.14)
}

/*
TestMonomialVector_Plus2
Description:

	Verifies that the Plus() method throws a panic when a well-formed
	vector of monomials is added to an improperly initialized expression
	(in this case a monomial matrix).
*/
func TestMonomialVector_Plus2(t *testing.T) {
	// Constants
	mv := symbolic.MonomialVector{
		symbolic.NewVariable().ToMonomial(),
		symbolic.NewVariable().ToMonomial(),
	}
	pm := symbolic.MonomialMatrix{}

	// Test
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf(
				"Expected mv.Plus(pm) to panic; received %v",
				mv.Plus(pm),
			)
		}
	}()

	mv.Plus(pm)
}

/*
TestMonomialVector_Plus3
Description:

	Verifies that the Plus() method throws a panic when a well-formed
	vector of monomials is added to an well formed
	matrix of polynomials that do not have identical dimensions.
*/
func TestMonomialVector_Plus3(t *testing.T) {
	// Constants
	mv := symbolic.MonomialVector{
		symbolic.NewVariable().ToMonomial(),
		symbolic.NewVariable().ToMonomial(),
	}
	pm := symbolic.PolynomialMatrix{
		{
			symbolic.NewVariable().ToPolynomial(),
			symbolic.NewVariable().ToPolynomial(),
			symbolic.NewVariable().ToPolynomial(),
		},
	}

	// Test
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf(
				"Expected mv.Plus(pm) to panic; received %v",
				mv.Plus(pm),
			)
		}

		rAsE, tf := r.(error)
		if !tf {
			t.Errorf(
				"Expected mv.Plus(pm) to panic with an error; received %v",
				r,
			)
		}

		if !strings.Contains(
			rAsE.Error(),
			smErrors.DimensionError{
				Operation: "Plus",
				Arg1:      mv,
				Arg2:      pm,
			}.Error(),
		) {
			t.Errorf(
				"Expected mv.Plus(pm) to panic with an error containing \"dimensions\"; received %v",
				rAsE.Error(),
			)
		}
	}()

	mv.Plus(pm)
}

/*
TestMonomialVector_Plus4
Description:

	Verifies that the Plus() method returns the correct value when a
	well-formed vector of monomials is added to a well-formed
	vector of monomials.
*/
func TestMonomialVector_Plus4(t *testing.T) {
	// Constants
	v1 := symbolic.NewVariable()
	v2 := symbolic.NewVariable()
	m1 := symbolic.Monomial{
		Coefficient:     1.0,
		VariableFactors: []symbolic.Variable{v1},
		Exponents:       []int{1},
	}
	m2 := symbolic.Monomial{
		Coefficient:     1.0,
		VariableFactors: []symbolic.Variable{v2},
		Exponents:       []int{1},
	}
	mv1 := symbolic.MonomialVector{m1, m2}
	mv2 := symbolic.MonomialVector{m1, m2}

	// Test
	sum := mv1.Plus(mv2)

	sumAsPV, tf := sum.(symbolic.PolynomialVector)
	if !tf {
		t.Errorf(
			"expected sum to be a PolynomialVector; received %T",
			sum,
		)
	}

	for _, polynomial := range sumAsPV {
		if len(polynomial.Monomials) != 1 {
			t.Errorf(
				"expected len(polynomial.Monomials) to be 1; received %v",
				len(polynomial.Monomials),
			)
		}
	}
}

/*
TestMonomialVector_Plus5
Description:

	This test verifies that the method properly panics if a valid
	vector of monomials is multiplied by an invalid expression
	(in this case a matrix of monomials).
*/
func TestMonomialVector_Plus5(t *testing.T) {
	// Constants
	v1 := symbolic.NewVariable()
	v2 := symbolic.NewVariable()
	m1 := symbolic.Monomial{
		Coefficient:     1.0,
		VariableFactors: []symbolic.Variable{v1},
		Exponents:       []int{1},
	}
	m2 := symbolic.Monomial{
		Coefficient:     1.0,
		VariableFactors: []symbolic.Variable{v2},
		Exponents:       []int{1},
	}
	mv1 := symbolic.MonomialVector{m1, m2}

	var mm symbolic.MonomialMatrix

	// Setup defer function for catching panic
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf(
				"Expected mv1.Plus(mm) to panic; received %v",
				mv1.Plus(mm),
			)
		}

		rAsE, tf := r.(error)
		if !tf {
			t.Errorf(
				"Expected mv1.Plus(mm) to panic with an error; received %v",
				r,
			)
		}

		if !strings.Contains(
			rAsE.Error(),
			mm.Check().Error(),
		) {
			t.Errorf(
				"Expected mv1.Plus(mm) to panic with an error containing \"%v\"; received %v",
				mm.Check().Error(),
				rAsE.Error(),
			)
		}
	}()

	// Test
	mv1.Plus(mm)
}

/*
TestMonomialVector_Plus6
Description:

	This test verifies that the method properly creates a monomial vector,
	if the current monomial vector contains all constants.
*/
func TestMonomialVector_Plus6(t *testing.T) {
	// Constants
	k2 := symbolic.K(3.14)

	// Create a monomial vector of constants
	kv1 := getKVector.From([]float64{1, 2, 3, 4, 5})
	mv1 := kv1.ToMonomialVector()

	// Compute Sum
	sum := mv1.Plus(k2)

	// Verify that the sum is a monomial vector
	if _, tf := sum.(symbolic.KVector); !tf {
		t.Errorf(
			"expected sum to be a MonomialVector; received %T",
			sum,
		)
	}

}

/*
TestMonomialVector_Plus7
Description:

	This test verifies that the method properly creates a polynomial vector
	when the monomial vector is added to a constant AND the monomial vector
	is not already a constant vector.
*/
func TestMonomialVector_Plus7(t *testing.T) {
	// Constants
	N := 10
	vv1 := symbolic.NewVariableVector(N)
	mv1 := vv1.ToMonomialVector()
	k2 := symbolic.K(3.14)

	// Compute Sum
	sum := mv1.Plus(k2)

	// Verify that the sum is a polynomial vector
	if _, tf := sum.(symbolic.PolynomialVector); !tf {
		t.Errorf(
			"expected sum to be a PolynomialVector; received %T",
			sum,
		)
	}

	// Check that each element of the polynomial vector
	// contains 2 monomials
	for _, polynomial := range sum.(symbolic.PolynomialVector) {
		if len(polynomial.Monomials) != 2 {
			t.Errorf(
				"expected len(polynomial.Monomials) to be 2; received %v",
				len(polynomial.Monomials),
			)
		}
	}
}

/*
TestMonomialVector_Multiply1
Description:

	Verifies that the Multiply() method throws a panic when an improperly
	initialized vector of monomials is given.
*/
func TestMonomialVector_Multiply1(t *testing.T) {
	// Constants
	mv := symbolic.MonomialVector{}

	// Test
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf(
				"Expected mv.Multiply(3.14) to panic; received %v",
				mv.Multiply(3.14),
			)
		}
	}()

	mv.Multiply(3.14)
}

/*
TestMonomialVector_Multiply2
Description:

	Verifies that the Multiply() method throws a panic when a well-formed
	vector of monomials is multiplied by an improperly initialized expression
	(in this case a monomial matrix).
*/
func TestMonomialVector_Multiply2(t *testing.T) {
	// Constants
	mv := symbolic.MonomialVector{
		symbolic.NewVariable().ToMonomial(),
		symbolic.NewVariable().ToMonomial(),
	}
	mm2 := symbolic.MonomialMatrix{}

	// Test
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf(
				"Expected mv.Multiply(mm2) to panic; received %v",
				mv.Multiply(mm2),
			)
		}
	}()

	mv.Multiply(mm2)
}

/*
TestMonomialVector_Multiply3
Description:

	Verifies that the Multiply() method throws a panic when a well-formed
	vector of monomials is multiplied by a well formed monomial matrix
	that does not have compatible dimensions.
*/
func TestMonomialVector_Multiply3(t *testing.T) {
	// Constants
	mv := symbolic.MonomialVector{
		symbolic.NewVariable().ToMonomial(),
		symbolic.NewVariable().ToMonomial(),
	}
	mm2 := symbolic.MonomialMatrix{
		{
			symbolic.NewVariable().ToMonomial(),
			symbolic.NewVariable().ToMonomial(),
			symbolic.NewVariable().ToMonomial(),
		},
		{
			symbolic.NewVariable().ToMonomial(),
			symbolic.NewVariable().ToMonomial(),
			symbolic.NewVariable().ToMonomial(),
		},
	}

	// Test
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf(
				"Expected mv.Multiply(mm2) to panic; received %v",
				mv.Multiply(mm2),
			)
		}
	}()

	mv.Multiply(mm2)
}

/*
TestMonomialVector_Multiply4
Description:

	Verifies that the Multiply() method returns the correct value when a
	well-formed vector of monomials is multiplied by a float64.
	We will verify that the result is a monomial vector where each
	monomials coefficient is multiplied by the float64.
*/
func TestMonomialVector_Multiply4(t *testing.T) {
	// Constants
	mv := symbolic.MonomialVector{
		symbolic.NewVariable().ToMonomial(),
		symbolic.NewVariable().ToMonomial(),
	}
	f2 := 3.14

	// Test
	product := mv.Multiply(f2)

	// Verify that the product is a monomial vector
	if _, tf := product.(symbolic.MonomialVector); !tf {
		t.Errorf(
			"expected product to be a MonomialVector; received %T",
			product,
		)
	}

	// Verify that the coefficients are correct
	for _, monomial := range product.(symbolic.MonomialVector) {
		if monomial.Coefficient != 3.14 {
			t.Errorf(
				"expected monomial.Coefficient to be 3.14; received %v",
				monomial.Coefficient,
			)
		}
	}
}

/*
TestMonomialVector_Multiply5
Description:

	Verifies that the Multiply() method returns the correct value when a
	well-formed vector of monomials is multiplied by an invalid non-expression
	(string).
*/
func TestMonomialVector_Multiply5(t *testing.T) {
	// Constants
	mv := symbolic.MonomialVector{
		symbolic.NewVariable().ToMonomial(),
		symbolic.NewVariable().ToMonomial(),
	}
	s2 := "This is a test string."

	// Test
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf(
				"Expected mv.Multiply(s2) to panic; received %v",
				mv.Multiply(s2),
			)
		}
	}()

	mv.Multiply(s2)
}

/*
TestMonomialVector_Transpose1
Description:

	Verifies that the Transpose() method throws a panic when an improperly
	initialized vector of monomials is given.
*/
func TestMonomialVector_Transpose1(t *testing.T) {
	// Constants
	mv := symbolic.MonomialVector{}

	// Test
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf(
				"Expected mv.Transpose() to panic; received %v",
				mv.Transpose(),
			)
		}
	}()

	mv.Transpose()
}

/*
TestMonomialVector_Transpose2
Description:

	Verifies that the Transpose() method returns the correct value when a
	well-formed vector of monomials is given.
	Checks that the dimensions of the transposed vector are correct.
*/
func TestMonomialVector_Transpose2(t *testing.T) {
	// Constants
	mv := symbolic.MonomialVector{
		symbolic.NewVariable().ToMonomial(),
		symbolic.NewVariable().ToMonomial(),
		symbolic.NewVariable().ToMonomial(),
	}

	// Test
	transposed := mv.Transpose()
	mvT, tf := transposed.(symbolic.MonomialMatrix)
	if !tf {
		t.Errorf(
			"expected mvT to be a MonomialMatrix; received %T",
			transposed,
		)
	}

	if mvT.Dims()[0] != mv.Dims()[1] || mvT.Dims()[1] != mv.Dims()[0] {
		t.Errorf(
			"expected mvT.Dims() to be [%v,%v]; received %v",
			mv.Dims()[1],
			mv.Dims()[0],
			mvT.Dims(),
		)
	}
}

/*
TestMonomial_LessEq1
Description:

	Verifies that the LessEq() method returns the correct value when
	compared to a KVector.
*/
func TestMonomial_LessEq1(t *testing.T) {
	// Constants
	N := 4
	vv1 := symbolic.NewVariableVector(N)
	mv1 := vv1.ToMonomialVector()

	kv2 := getKVector.From([]float64{1, 2, 3, 4})

	// Create a vector constraint
	constraint := mv1.LessEq(kv2)

	// verify that the cosntraint is a vector constraint
	vc2, tf := constraint.(symbolic.VectorConstraint)
	if !tf {
		t.Errorf("expected constraint to be a VectorConstraint; received %T", constraint)
	}

	// Verify that the constraint's dimensions are correct
	if vc2.Dims()[0] != N {
		t.Errorf(
			"expected constraint.Dims()[0] to be %v; received %v",
			N,
			vc2.Dims()[0],
		)
	}
}

/*
TestMonomial_GreaterEq1
Description:

	Verifies that the GreaterEq() method returns the correct value when
	compared to a K.
*/
func TestMonomial_GreaterEq1(t *testing.T) {
	// Constants
	N := 5
	vv1 := symbolic.NewVariableVector(N)
	mv1 := vv1.ToMonomialVector()

	k2 := symbolic.K(3.14)

	// Create a vector constraint
	constraint := mv1.GreaterEq(k2)

	// verify that the cosntraint is a vector constraint
	vc2, tf := constraint.(symbolic.VectorConstraint)
	if !tf {
		t.Errorf("expected constraint to be a VectorConstraint; received %T", constraint)
	}

	// Verify that the constraint's dimensions are correct
	if vc2.Dims()[0] != N {
		t.Errorf(
			"expected constraint.Dims()[0] to be %v; received %v",
			N,
			vc2.Dims()[0],
		)
	}
}
