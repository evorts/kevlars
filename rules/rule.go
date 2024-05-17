package rules

import (
	"github.com/evorts/kevlars/rules/eval"
)

func WhenTrue(expression bool, run F) {
	WhenTrueE(expression, run, nil)
}

func WhenTrueE(expression bool, run F, orElse F) {
	if expression {
		run()
		return
	}
	if orElse != nil {
		orElse()
	}
}

func WhenTrueR1[T any](expression bool, run FR1[T], orElse FR1[T]) T {
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

func WhenError(err error, run F) {
	WhenTrue(err != nil, run)
}

func WhenErrorE(err error, run F, orElse F) {
	WhenTrueE(err != nil, run, orElse)
}

func WhenNil[T any](v *T, run F) {
	WhenNilE(v, run, nil)
}

func WhenNilE[T any](v *T, run F, orElse F) {
	if v == nil {
		run()
		return
	}
	if orElse != nil {
		orElse()
	}
}

func WhenNilR1[T any](v *T, run FR1[T]) T {
	return WhenNilRE1(v, run, func() T {
		return *v
	})
}

func WhenNilRE1[T any](v *T, run FR1[T], orElse FR1[T]) T {
	if v == nil {
		return run()
	}
	if orElse != nil {
		return orElse()
	}
	panic("or else is not defined")
}

func WhenEmpty[T any](v *T, run F) {
	WhenEmptyE(v, run, nil)
}

func WhenEmptyE[T any](v *T, run F, orElse F) {
	if eval.IsEmpty(v) {
		run()
		return
	}
	if orElse != nil {
		orElse()
	}
}

func WhenEmptyR1[T any](v *T, run FR1[T], orElse FR1[T]) T {
	if eval.IsEmpty(v) {
		return run()
	}
	if orElse != nil {
		return orElse()
	}
	panic("or else is not defined")
}

func WhenNotEmpty[T any](v *T, run F) {
	WhenNotEmptyE(v, run, nil)
}

func WhenNotEmptyE[T any](v *T, run F, orElse F) {
	if !eval.IsEmpty(v) {
		run()
		return
	}
	if orElse != nil {
		orElse()
	}
}

func WhenNotEmptyR1[T any](v *T, run FR1[T], orElse FR1[T]) T {
	if !eval.IsEmpty(v) {
		return run()
	}
	if orElse != nil {
		return orElse()
	}
	panic("or else is not defined")
}

func WhenSliceNotEmpty[T any](items []T, run F) {
	WhenSliceNotEmptyE(items, run, nil)
}

func WhenSliceNotEmptyE[T any](items []T, run F, orElse F) {
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

func Iif[T1 any](expression bool, whenTrue T1, orElse T1) T1 {
	if expression {
		return whenTrue
	}
	return orElse
}

func IfNil[T any](v *T, defaultValue T) T {
	return IfNilE(v, defaultValue, *v)
}

func IfNilE[T any](v *T, whenTrue T, orElse T) T {
	if v == nil {
		return whenTrue
	}
	return orElse
}

func IfEmpty[T comparable](v T, defaultValue T) T {
	return IfEmptyE(v, defaultValue, v)
}

func IfEmptyE[T comparable](v T, whenTrue T, orElse T) T {
	if eval.IsEmpty(v) {
		return whenTrue
	}
	return orElse
}

func IfNotEmpty[T comparable](v T, defaultValue T) T {
	return IfNotEmptyE(v, defaultValue, v)
}

func IfNotEmptyE[T comparable](v T, whenTrue T, orElse T) T {
	if !eval.IsEmpty(v) {
		return whenTrue
	}
	return orElse
}
