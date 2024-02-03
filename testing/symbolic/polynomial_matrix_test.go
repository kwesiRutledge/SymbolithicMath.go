package symbolic_test

import (
	getKMatrix "github.com/MatProGo-dev/SymbolicMath.go/get/KMatrix"
	getKVector "github.com/MatProGo-dev/SymbolicMath.go/get/KVector"
	"github.com/MatProGo-dev/SymbolicMath.go/smErrors"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"strings"
	"testing"
)

/*
polynomial_matrix_test.go
Description:

	Tests for the functions mentioned in the polynomial_matrix.go file.
*/

/*
TestPolynomialMatrix_Check1
Description:

	Tests that the Check() method properly catches an improperly initialized
	matrix of Polynomials (i.e., no polynomials are given).
*/
func TestPolynomialMatrix_Check1(t *testing.T) {
	// Constants
	pm := symbolic.PolynomialMatrix{}
	expectedError := smErrors.EmptyMatrixError{pm}

	// Test
	err := pm.Check()
	if err.Error() != expectedError.Error() {
		t.Errorf(
			"expected Check() to return \"%v\"; received %v",
			expectedError,
			err,
		)
	}
}

/*
TestPolynomialMatrix_Check2
Description:

	Tests that the Check() method properly catches an improperly initialized
	matrix of Polynomials (i.e., a polynomial is given with an improper number of
	columns in row 2 than in every other row.)
*/
func TestPolynomialMatrix_Check2(t *testing.T) {
	// Constants
	v1 := symbolic.NewVariable()
	m1 := symbolic.Monomial{
		Coefficient:     1.0,
		VariableFactors: []symbolic.Variable{v1},
		Exponents:       []int{1, 2},
	}
	p1 := symbolic.Polynomial{
		[]symbolic.Monomial{m1},
	}
	var pm symbolic.PolynomialMatrix = [][]symbolic.Polynomial{
		{p1, p1},
		{p1},
		{p1, p1},
	}

	expectedError := smErrors.MatrixColumnMismatchError{
		ExpectedNColumns: 2,
		ActualNColumns:   1,
		Row:              1,
	}

	// Test
	err := pm.Check()
	if err.Error() != expectedError.Error() {
		t.Errorf(
			"expected Check() to return \"%v\"; received %v",
			expectedError,
			err,
		)
	}
}

/*
TestPolynomialMatrix_Check3
Description:

	Tests that the Check() method properly catches an improperly initialized
	matrix of Polynomials (i.e., a polynomial is given with an improper number of
	degrees in a monomial in the third row.)
*/
func TestPolynomialMatrix_Check3(t *testing.T) {
	// Constants
	v1 := symbolic.NewVariable()
	m2 := symbolic.Monomial{
		Coefficient:     1.0,
		VariableFactors: []symbolic.Variable{v1},
		Exponents:       []int{1, 2},
	}

	// Create polynomials
	p1 := v1.ToPolynomial()
	p2 := symbolic.Polynomial{[]symbolic.Monomial{m2}}

	// Construct matrix
	var pm symbolic.PolynomialMatrix = [][]symbolic.Polynomial{
		{p1, p1},
		{p1, p1},
		{p1, p2},
	}

	// Test
	err := pm.Check()
	if !strings.Contains(
		err.Error(),
		"2,1", // coordinate of bad polynomial
	) {
		t.Errorf(
			"expected Check() to return error containing \"%v\"; received %v",
			"2,1",
			err,
		)
	}
}

/*
TestPolynomialMatrix_Check4
Description:

	Tests that the Check() method properly returns nil when a matrix of
	Polynomials is properly initialized. (in this case a 3 x 2 matrix of
	Polynomials is given).
*/
func TestPolynomialMatrix_Check4(t *testing.T) {
	// Constants
	v1 := symbolic.NewVariable()
	p1 := v1.ToPolynomial()

	// Construct matrix
	var pm symbolic.PolynomialMatrix = [][]symbolic.Polynomial{
		{p1, p1},
		{p1, p1},
		{p1, p1},
	}

	// Test
	if pm.Check() != nil {
		t.Errorf(
			"expected Check() to return nil; received %v",
			pm.Check(),
		)
	}
}

/*
TestPolynomialMatrix_Variables1
Description:

	Tests that the Variables() method properly returns the variables in a
	PolynomialMatrix. In this case, the polynomials all contain
	a single variable and so the length of Variables should be 1.
*/
func TestPolynomialMatrix_Variables1(t *testing.T) {
	// Constants
	v1 := symbolic.NewVariable()
	p1 := v1.ToPolynomial()

	// Construct matrix
	var pm symbolic.PolynomialMatrix = [][]symbolic.Polynomial{
		{p1, p1},
		{p1, p1},
		{p1, p1},
	}

	// Test
	vars := pm.Variables()
	if len(vars) != 1 {
		t.Errorf(
			"expected len(vars) to be 1; received %v",
			len(vars),
		)
	}
}

/*
TestPolynomialMatrix_Variables2
Description:

	Test that the variables() method properly panics
	when applied to an improperly initialized matrix of polynomials.
*/
func TestPolynomialMatrix_Variables2(t *testing.T) {
	// Constants
	var pm symbolic.PolynomialMatrix

	// Test
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf(
				"expected pm.Variables() to panic; received %v",
				pm.Variables(),
			)
		}
	}()

	pm.Variables()
}

/*
TestPolynomialMatrix_Dims1
Description:

	Tests that the Dims() method properly returns the dimensions of a
	PolynomialMatrix. In this case, the matrix is 3 x 2.
*/
func TestPolynomialMatrix_Dims1(t *testing.T) {
	// Constants
	v1 := symbolic.NewVariable()
	p1 := v1.ToPolynomial()

	// Construct matrix
	var pm symbolic.PolynomialMatrix = [][]symbolic.Polynomial{
		{p1, p1},
		{p1, p1},
		{p1, p1},
	}

	// Test
	dim1, dim2 := pm.Dims()[0], pm.Dims()[1]
	if dim1 != 3 {
		t.Errorf(
			"expected pm.Dims() to be [3,2]; received %v rows",
			dim1,
		)
	}

	if dim2 != 2 {
		t.Errorf(
			"expected pm.Dims() to be [3,2]; received %v columns",
			dim2,
		)
	}
}

/*
TestPolynomialMatrix_Dims2
Description:

	Tests that the Dims() method properly panics when applied to an
	improperly initialized matrix of polynomials.
*/
func TestPolynomialMatrix_Dims2(t *testing.T) {
	// Constants
	var pm symbolic.PolynomialMatrix

	// Test
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf(
				"expected pm.Dims() to panic; received %v",
				pm.Dims(),
			)
		}
	}()

	pm.Dims()
}

/*
TestPolynomialMatrix_Plus1
Description:

	Tests that the Plus() panics when an improperly initialized
	matrix of polynomials is used to call it..
*/
func TestPolynomialMatrix_Plus1(t *testing.T) {
	// Constants
	var pm symbolic.PolynomialMatrix

	// Test
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf(
				"expected pm.Plus(pm) to panic; received %v",
				pm.Plus(pm),
			)
		}
	}()

	pm.Plus(pm)
}

/*
TestPolynomialMatrix_Plus2
Description:

	Tests that the Plus() method properly panics an error if two
	matrices of different dimensions are given to the method.
*/
func TestPolynomialMatrix_Plus2(t *testing.T) {
	// Constants
	v1 := symbolic.NewVariable()
	p1 := v1.ToPolynomial()
	var pm1 symbolic.PolynomialMatrix = [][]symbolic.Polynomial{
		{p1, p1},
		{p1, p1},
		{p1, p1},
	}
	var pm2 symbolic.PolynomialMatrix = [][]symbolic.Polynomial{
		{p1, p1},
		{p1, p1},
	}

	// Create panic checking logic
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf(
				"expected pm1.Plus(pm2) to panic; received %v",
				pm1.Plus(pm2),
			)
		}

		rAsError, tf := r.(error)
		if !tf {
			t.Errorf(
				"expected r to be an error; received %v of type %T",
				r, r,
			)
		}

		if rAsError.Error() != (smErrors.DimensionError{
			Operation: "Plus",
			Arg1:      pm1,
			Arg2:      pm2,
		}).Error() {
			t.Errorf(
				"expected r to be a DimensionError; received %v",
				r,
			)
		}
	}()

	// Test
	pm1.Plus(pm2)

}

/*
TestPolynomialMatrix_Plus3
Description:

	Tests that the Plus() method properly adds a polynomial matrix
	with a float64. The result should be a polynomial matrix with
	each polynomial containing two monomials.
*/
func TestPolynomialMatrix_Plus3(t *testing.T) {
	// Constants
	v1 := symbolic.NewVariable()
	p1 := v1.ToPolynomial()
	var pm1 symbolic.PolynomialMatrix = [][]symbolic.Polynomial{
		{p1, p1},
		{p1, p1},
		{p1, p1},
	}

	// Test
	pm2 := pm1.Plus(1.0)

	pm2AsPM, tf := pm2.(symbolic.PolynomialMatrix)
	if !tf {
		t.Errorf(
			"expected pm2 to be a PolynomialMatrix; received %v",
			pm2,
		)
	}

	for _, pmRow := range pm2AsPM {
		for _, p := range pmRow {
			if len(p.Monomials) != 2 {
				t.Errorf(
					"expected len(p.Monomials) to be 2; received %v",
					len(p.Monomials),
				)
			}
		}
	}
}

/*
TestPolynomialMatrix_Plus4
Description:

	Tests that the Plus() method properly adds a polynomial matrix
	with THE SAME polynomial matrix. The result should be a polynomial
	matrix with each polynomial containing one monomial.
*/
func TestPolynomialMatrix_Plus4(t *testing.T) {
	// Constants
	v1 := symbolic.NewVariable()
	p1 := v1.ToPolynomial()
	var pm1 symbolic.PolynomialMatrix = [][]symbolic.Polynomial{
		{p1, p1},
		{p1, p1},
		{p1, p1},
	}

	// Test
	pm2 := pm1.Plus(pm1)

	pm2AsPM, tf := pm2.(symbolic.PolynomialMatrix)
	if !tf {
		t.Errorf(
			"expected pm2 to be a PolynomialMatrix; received %v",
			pm2,
		)
	}

	for _, pmRow := range pm2AsPM {
		for _, p := range pmRow {
			if len(p.Monomials) != 1 {
				t.Errorf(
					"expected len(p.Monomials) to be 1; received %v",
					len(p.Monomials),
				)
			}
		}
	}
}

/*
TestPolynomialMatrix_Plus5
Description:

	Tests that the Plus() method properly adds a polynomial matrix
	with a polynomial. The result should be a polynomial matrix with
	each polynomial containing three monomials.
*/
func TestPolynomialMatrix_Plus5(t *testing.T) {
	// Constants
	v1 := symbolic.NewVariable()
	p1 := v1.ToPolynomial()
	var pm1 symbolic.PolynomialMatrix = [][]symbolic.Polynomial{
		{p1, p1},
		{p1, p1},
		{p1, p1},
	}

	p2 := v1.Plus(symbolic.NewVariable()).Plus(symbolic.NewVariable()).(symbolic.Polynomial)

	// Test
	pm2 := pm1.Plus(p2)

	pm2AsPM, tf := pm2.(symbolic.PolynomialMatrix)
	if !tf {
		t.Errorf(
			"expected pm2 to be a PolynomialMatrix; received %v",
			pm2,
		)
	}

	for _, pmRow := range pm2AsPM {
		for _, p := range pmRow {
			if len(p.Monomials) != 3 {
				t.Errorf(
					"expected len(p.Monomials) to be 3; received %v",
					len(p.Monomials),
				)
			}
		}
	}
}

/*
TestPolynomialMatrix_Plus6
Description:

	Tests that the Plus() method properly adds a polynomial matrix
	to a monomial. The result should be a polynomial matrix with
	each polynomial containing two monomials.
*/
func TestPolynomialMatrix_Plus6(t *testing.T) {
	// Constants
	v1 := symbolic.NewVariable()
	p1 := v1.ToPolynomial()
	var pm1 symbolic.PolynomialMatrix = [][]symbolic.Polynomial{
		{p1, p1},
		{p1, p1},
		{p1, p1},
	}

	m1 := symbolic.NewVariable().ToMonomial()

	// Test
	pm2 := pm1.Plus(m1)

	pm2AsPM, tf := pm2.(symbolic.PolynomialMatrix)
	if !tf {
		t.Errorf(
			"expected pm2 to be a PolynomialMatrix; received %v",
			pm2,
		)
	}

	for _, pmRow := range pm2AsPM {
		for _, p := range pmRow {
			if len(p.Monomials) != 2 {
				t.Errorf(
					"expected len(p.Monomials) to be 2; received %v",
					len(p.Monomials),
				)
			}
		}
	}
}

/*
TestPolynomialMatrix_Plus7
Description:

	Tests that the Plus() method properly adds a polynomial matrix
	to a matrix of constants (K). The result should be a polynomial
	matrix with each polynomial containing two monomials.
*/
func TestPolynomialMatrix_Plus7(t *testing.T) {
	// Constants
	v1 := symbolic.NewVariable()
	p1 := v1.ToPolynomial()
	var pm1 symbolic.PolynomialMatrix = [][]symbolic.Polynomial{
		{p1, p1},
		{p1, p1},
		{p1, p1},
	}

	km1 := getKMatrix.From(symbolic.OnesMatrix(3, 2))

	// Test
	pm2 := pm1.Plus(km1)

	pm2AsPM, tf := pm2.(symbolic.PolynomialMatrix)
	if !tf {
		t.Errorf(
			"expected pm2 to be a PolynomialMatrix; received %v",
			pm2,
		)
	}

	for _, pmRow := range pm2AsPM {
		for _, p := range pmRow {
			if len(p.Monomials) != 2 {
				t.Errorf(
					"expected len(p.Monomials) to be 2; received %v",
					len(p.Monomials),
				)
			}
		}
	}
}

/*
TestPolynomialMatrix_Multiply1
Description:

	Tests that the Multiply() method properly panics when an improperly
	initialized matrix of polynomials is used to call it.
*/
func TestPolynomialMatrix_Multiply1(t *testing.T) {
	// Constants
	var pm symbolic.PolynomialMatrix

	// Test
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf(
				"expected pm.Multiply(pm) to panic; received %v",
				pm.Multiply(pm),
			)
		}
	}()

	pm.Multiply(pm)
}

/*
TestPolynomialMatrix_Multiply2
Description:

	Tests that the Multiply() method properly panics an error if it
	is called with a properly defined matrix of Polynomials and an improperly
	initialized matrix of polynomials.
*/
func TestPolynomialMatrix_Multiply2(t *testing.T) {
	// Constants
	v1 := symbolic.NewVariable()
	p1 := v1.ToPolynomial()
	var pm1 symbolic.PolynomialMatrix = [][]symbolic.Polynomial{
		{p1, p1},
		{p1, p1},
	}

	var pm2 symbolic.PolynomialMatrix

	// Create panic checking logic
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf(
				"expected pm1.Multiply(pm2) to panic; received %v",
				pm1.Multiply(pm2),
			)
		}

		rAsError, tf := r.(error)
		if !tf {
			t.Errorf(
				"expected r to be an error; received %v of type %T",
				r, r,
			)
		}

		if !strings.Contains(
			rAsError.Error(),
			pm2.Check().Error(),
		) {
			t.Errorf(
				"expected r to be a DimensionError; received %v",
				r,
			)
		}
	}()

	// Test
	pm1.Multiply(pm2)
}

/*
TestPolynomialMatrix_Multiply3
Description:

	Tests that the Multiply() method properly panics when a polynomial matrix
	with a KVector (vector of constants) that does not match in size.
	A DimensionError should be thrown.
*/
func TestPolynomialMatrix_Multiply3(t *testing.T) {
	// Constants
	v1 := symbolic.NewVariable()
	p1 := v1.ToPolynomial()
	var pm1 symbolic.PolynomialMatrix = [][]symbolic.Polynomial{
		{p1, p1},
		{p1, p1},
	}

	kv1 := getKVector.From([]float64{1.0, 2.1, 3.14})

	// Create panic checking logic
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf(
				"expected pm1.Multiply(km1) to panic; received %v",
				pm1.Multiply(kv1),
			)
		}

		rAsError, tf := r.(error)
		if !tf {
			t.Errorf(
				"expected r to be an error; received %v of type %T",
				r, r,
			)
		}

		if !strings.Contains(
			rAsError.Error(),
			smErrors.DimensionError{
				Operation: "Multiply",
				Arg1:      pm1,
				Arg2:      kv1,
			}.Error(),
		) {
			t.Errorf(
				"expected r to be a DimensionError; received %v",
				r,
			)
		}
	}()

	// Test
	pm1.Multiply(kv1)
}

/*
TestPolynomialMatrix_Multiply4
Description:

	Tests that the Multiply() method properly multiplies a polynomial matrix
	with a single constant. The result should be a polynomial matrix with
	the same number of monomials as the original. But the coefficient
	of each monomial should be multiplied by the constant.
*/
func TestPolynomialMatrix_Multiply4(t *testing.T) {
	// Constants
	v1 := symbolic.NewVariable()
	p1 := v1.ToPolynomial()
	var pm1 symbolic.PolynomialMatrix = [][]symbolic.Polynomial{
		{p1},
		{p1},
		{p1},
	}

	// Test
	pm2 := pm1.Multiply(3.14)

	pm2AsPM, tf := pm2.(symbolic.PolynomialMatrix)
	if !tf {
		t.Errorf(
			"expected pm2 to be a PolynomialMatrix; received %v",
			pm2,
		)
	}

	for _, pmRow := range pm2AsPM {
		for _, p := range pmRow {
			if len(p.Monomials) != 1 {
				t.Errorf(
					"expected len(p.Monomials) to be 1; received %v",
					len(p.Monomials),
				)
			}

			if p.Monomials[0].Coefficient != 3.14 {
				t.Errorf(
					"expected p.Monomials[0].Coefficient to be 3.14; received %v",
					p.Monomials[0].Coefficient,
				)
			}
		}
	}
}
