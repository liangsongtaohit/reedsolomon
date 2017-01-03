/**
 * Matrix Algebra over an 8-bit Galois Field
 *
 * Copyright 2015, Klaus Post
 * Copyright 2015, Backblaze, Inc.
 */

package reedsolomon

import (
	//"errors"
	//"fmt"
	//"strconv"
	//"strings"
)

// byte[row][col]
type matrix [][]byte

// newMatrix returns a matrix of zeros.
func newMatrix(rows, cols int) (matrix, error) {
	if rows <= 0 {
		return nil, errInvalidRowSize
	}
	if cols <= 0 {
		return nil, errInvalidColSize
	}

	m := matrix(make([][]byte, rows))
	for i := range m {
		m[i] = make([]byte, cols)
	}
	return m, nil
}

// TODO 测试完后改小写
func genEncodeMatrix(rows, cols int) (matrix, error) {
	m, err := newMatrix(rows, cols)
	if err != nil {
		return nil, err
	}

	// identity matrix
	for j := 0; j < cols; j++ {
		m[j][j] = byte(1)
	}
	// cauchy matrix
	for i := cols; i < rows; i++ {
		for j := 0; j < cols; j++ {
			dividend := i ^ j
			a := InverseTable[dividend]
			m[i][j] = byte(a)
		}
	}
	return m, nil
}
