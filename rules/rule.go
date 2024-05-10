package rules

func WhenSliceNotEmpty[T any](items []T, run F, orElse F) {
	if items == nil || len(items) == 0 {
		if orElse != nil {
			orElse()
		}
		return
	}
	run()
}

func WhenSliceNotEmptyR1[IT any, T1 any](items []IT, run FR1[T1], orElse FR1[T1]) T1 {
	if items == nil || len(items) == 0 {
		if orElse != nil {
			return orElse()
		}
		panic("or else is not defined")
	}
	return run()
}

func WhenTrue(expression bool, run F, orElse F) {
	if expression {
		run()
		return
	}
	if orElse != nil {
		orElse()
	}
}

func WhenTrueR1[T1 any](expression bool, run FR1[T1], orElse FR1[T1]) T1 {
	if expression {
		return run()
	}
	if orElse != nil {
		return orElse()
	}
	panic("or else is not defined")
}

func WhenTrueR2[T1 any, T2 any](expression bool, run FR2[T1, T2], orElse FR2[T1, T2]) (T1, T2) {
	if expression {
		return run()
	}
	if orElse != nil {
		return orElse()
	}
	panic("or else is not defined")
}

func WhenNil[T any](v T, run F, orElse F) {
	if v == nil {
		run()
		return
	}
	if orElse != nil {
		orElse()
	}
}

func WhenNilR1[T any, T1 any](v T, run FR1[T1], orElse FR1[T1]) T1 {
	if v == nil {
		return run()
	}
	if orElse != nil {
		return orElse()
	}
	panic("or else is not defined")
}
