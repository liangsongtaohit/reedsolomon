package reedsolomon

import (
	"errors"
)

type matrix [][]byte // byte[row][col]

func newMatrix(rows, cols int) matrix {
	m := matrix(make([][]byte, rows))
	for i := range m {
		m[i] = make([]byte, cols)
	}
	return m
}

// return identity matrix(upper) cauchy matrix(lower)
func genEncodeMatrix(rows, cols int) matrix {
	m := newMatrix(rows, cols)
	// identity matrix
	for j := 0; j < cols; j++ {
		m[j][j] = byte(1)
	}
	// cauchy matrix
	for i := cols; i < rows; i++ {
		for j := 0; j < cols; j++ {
			d := i ^ j
			a := inverseTable[d]
			m[i][j] = byte(a)
		}
	}
	return m
}

func (m matrix) invert() (matrix, error) {
	size := len(m)
	iM := identityMatrix(size)
	mIM, _ := m.augIM(iM)

	err := mIM.gaussJordan()
	if err != nil {
		return nil, err
	}
	return mIM.subMatrix(size), nil
}

// IN -> (IN|I)
func (m matrix) augIM(iM matrix) (matrix, error) {
	result := newMatrix(len(m), len(m[0])+len(iM[0]))
	for r, row := range m {
		for c := range row {
			result[r][c] = m[r][c]
		}
		cols := len(m[0])
		for c := range iM[0] {
			result[r][cols+c] = iM[r][c]
		}
	}
	return result, nil
}

var ErrSingular = errors.New("reedsolomon: matrix is singular")

// (IN|I) -> (I|OUT)
func (m matrix) gaussJordan() error {
	rows := len(m)
	columns := len(m[0])
	// Clear out the part below the main diagonal and scale the main
	// diagonal to be 1.
	for r := 0; r < rows; r++ {
		// If the element on the diagonal is 0, find a row below
		// that has a non-zero and swap them.
		if m[r][r] == 0 {
			for rowBelow := r + 1; rowBelow < rows; rowBelow++ {
				if m[rowBelow][r] != 0 {
					m.swapRows(r, rowBelow)
					break
				}
			}
		}
		// After swap, if we find all elements in this column is 0, it means the matrix's det is 0
		if m[r][r] == 0 {
			return ErrSingular
		}
		// Scale to 1.
		if m[r][r] != 1 {
			d := m[r][r]
			scale := inverseTable[d]
			// every element(this column) * m[r][r]'s inverse
			for c := 0; c < columns; c++ {
				m[r][c] = gfMul(m[r][c], scale)
			}
		}
		//Make everything below the 1 be a 0 by subtracting a multiple of it
		for rowBelow := r + 1; rowBelow < rows; rowBelow++ {
			if m[rowBelow][r] != 0 {
				// scale * m[r][r] = scale, scale + scale = 0
				// makes m[r][r+1] = 0 , then calc left elements
				scale := m[rowBelow][r]
				for c := 0; c < columns; c++ {
					m[rowBelow][c] ^= gfMul(scale, m[r][c])
				}
			}
		}
	}
	// Now clear the part above the main diagonal.
	// same logic with clean upper
	for d := 0; d < rows; d++ {
		for rowAbove := 0; rowAbove < d; rowAbove++ {
			if m[rowAbove][d] != 0 {
				scale := m[rowAbove][d]
				for c := 0; c < columns; c++ {
					m[rowAbove][c] ^= gfMul(scale, m[d][c])
				}
			}
		}
	}
	return nil
}

func identityMatrix(n int) matrix {
	m := newMatrix(n, n)
	for i := 0; i < n; i++ {
		m[i][i] = byte(1)
	}
	return m
}

// (I|OUT) -> OUT
func (m matrix) subMatrix(size int) matrix {
	result := newMatrix(size, size)
	for r := 0; r < size; r++ {
		for c := size; c < size*2; c++ {
			result[r][c-size] = m[r][c]
		}
	}
	return result
}

var ErrInvalidRowSize = errors.New("reedsolomon: invalid row size")

// SwapRows Exchanges two rows in the matrix.
func (m matrix) swapRows(r1, r2 int) {
	m[r2], m[r1] = m[r1], m[r2]
}

func gfMul(a, b byte) byte {
	return mulTable[a][b]
}
