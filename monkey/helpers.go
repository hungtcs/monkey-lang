package monkey

import "github.com/hungtcs/monkey-lang/syntax"

func isSameType(x, y Value) bool {
	return x.Type() == y.Type()
}

func b2i(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

func threeway(op syntax.Token, cmp int) bool {
	switch op {
	case syntax.EQ:
		return cmp == 0
	case syntax.NE:
		return cmp != 0
	case syntax.LE:
		return cmp <= 0
	case syntax.LT:
		return cmp < 0
	case syntax.GE:
		return cmp >= 0
	case syntax.GT:
		return cmp > 0
	}
	panic(op)
}
