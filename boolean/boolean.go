package boolean

func And(left bool, right bool) bool {
	return left && right
}

func Or(left bool, right bool) bool {
	return left || right
}

func Negate(x bool) bool {
	return !x
}
