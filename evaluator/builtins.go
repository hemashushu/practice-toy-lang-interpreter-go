// original from https://interpreterbook.com/

package evaluator

import (
	"fmt"
	"interpreter/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("number of arguments for `len` expected 1, actual %d", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}

			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}

			default:
				return newError("argument type of `len` expected STRING or ARRAY, actual %s", args[0].Type())
			}
		},
	},

	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("number of arguments for `len` expected 1, actual %d",
					len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument type of `first` expected ARRAY, actual %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return NULL
		},
	},

	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("number of arguments for `last` expected 1, actual %d",
					len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument type of `last` expected ARRAY, actual %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}
			return NULL
		},
	},

	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("number of arguments for `rest` expected 1, actual %d",
					len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument type of `rest` expected ARRAY, actual %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			if length > 0 {
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}
			return NULL
		},
	},

	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("number of arguments for `push` expected 2, actual %d",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument type of `push` expected ARRAY, actual %s",
					args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, arr.Elements)

			newElements[length] = args[1]
			return &object.Array{Elements: newElements}
		},
	},

	"pop": {
		// TODO
		Fn: func(args ...object.Object) object.Object {
			return NULL
		},
	},

	"puts": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
}
