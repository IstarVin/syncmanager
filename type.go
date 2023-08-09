package syncmanager

type Queue struct {
	Func func(...any)
	Args []any
}
