package fn

func Curry[A any, B any, C any](f func(A, B) C) func(A) func(B) C {
	return func(a A) func(B) C {
		return func(b B) C {
			return f(a, b)
		}
	}
}

func Pipe2[T1, T2, T3 any](f1 func(T1) T2, f2 func(T2) T3) func(T1) T3 {
	return func(t1 T1) T3 {
		return f2(f1(t1))
	}
}

func Pipe3[T1, T2, T3, T4 any](f1 func(T1) T2, f2 func(T2) T3, f3 func(T3) T4) func(T1) T4 {
	return func(t1 T1) T4 {
		return f3(f2(f1(t1)))
	}
}

func AsPointer[T any](x T) *T {
	return &x
}
