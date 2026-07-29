package rotate

import "math"

type matrix [][]float64

func (m matrix) row() int {
	return len(m)
}

func (m matrix) col() int {
	if len(m) == 0 {
		return 0
	}
	return len(m[0])
}

func (m matrix) mul(n matrix) matrix {
	if m.col() != n.row() {
		panic("Invalid matrix multiplication.")
	}
	l := make(matrix, m.row())
	for i := 0; i < m.row(); i++ {
		l[i] = make([]float64, n.col())
		for j := 0; j < n.col(); j++ {
			var s float64
			for k := 0; k < m.col(); k++ {
				s += m[i][k] * n[k][j]
			}
			l[i][j] = s
		}
	}
	return l
}

func identity(n int) matrix {
	m := make(matrix, n)
	for i := 0; i < n; i++ {
		m[i] = make([]float64, n)
		m[i][i] = 1
	}
	return m
}

func rotation3d(idx [2]int, t float64) matrix {
	m := identity(3)
	i, j := idx[0], idx[1]
	m[i][i] = math.Cos(t)
	m[i][j] = math.Sin(t)
	m[j][i] = -math.Sin(t)
	m[j][j] = math.Cos(t)
	return m
}

func (m matrix) t() matrix {
	n := make(matrix, m.col())
	for i := 0; i < 3; i++ {
		n[i] = make([]float64, m.row())
		for j := 0; j < 3; j++ {
			n[i][j] = m[j][i]
		}
	}
	return n
}
