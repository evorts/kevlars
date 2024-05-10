package rules

type F func()
type FR1[T any] func() T
type FR2[T1 any, T2 any] func() (T1, T2)
