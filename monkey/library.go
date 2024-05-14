package monkey

import "fmt"

var Universe = map[string]*BuiltinFunction{
	"len": NewBuiltinFunction("len", func(args ...Value) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("wrong number of arguments. got=%d, want=1", len(args))
		}
		var arg0 = args[0]
		switch arg0 := arg0.(type) {
		case Indexable:
			return Int(arg0.Len()), nil
		default:
			return nil, fmt.Errorf("argument to `len` not supported, got %s", arg0.Type())
		}
	}),
	"print": NewBuiltinFunction("print", func(args ...Value) (Value, error) {
		str := make([]any, len(args))
		for i, arg := range args {
			str[i] = arg.String()
		}
		fmt.Println(str...)
		return Null, nil
	}),
}
