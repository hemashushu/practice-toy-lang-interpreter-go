// original from https://interpreterbook.com/

package object

import "fmt"

type ObjectType string

// ObjectType 可能的值
const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE" // 包裹其他 Object 的 Object，用于 return 语句
	ERROR_OBJ        = "ERROR"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return ObjectType(INTEGER_OBJ)
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return ObjectType(BOOLEAN_OBJ)
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type Null struct {
	//
}

func (n *Null) Type() ObjectType {
	return ObjectType(NULL_OBJ)
}

func (n *Null) Inspect() string {
	return "null"
}

// 包裹其他 Object 的 Object，用于 return 语句
type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (r *ReturnValue) Inspect() string  { return r.Value.Inspect() }
