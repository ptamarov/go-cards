package orders

type Comparable interface {
	gt(t Comparable) bool
}

func IsGreaterThanLex[T Comparable](x, y []T) bool {
	if len(y) == 0 {
		return false
	}
	if len(x) == 0 {
		return true
	}
	if x[0].gt(y[0]) {
		return false
	}
	if y[0].gt(x[0]) {
		return true
	} else {
		return IsGreaterThanLex(x[1:], y[1:])
	}
}
