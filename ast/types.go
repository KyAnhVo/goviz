package ast

type TypeParam struct {
	Name       string
	Constraint TypeNode
}

type TypeNode interface {
	implType()
}

type TypeName struct {
	Name    string
	TypeArg []TypeNode
}

func (t *TypeName) implType() {}

type DefinedType struct {
	Name       string
	Underlying TypeNode
	TypeArg    []TypeNode
}

func (t *DefinedType) implType() {}

type Array struct {
	ElementType TypeNode
	Size        Expression
}

func (arr *Array) implType() {}

type Slice struct {
	ElementType TypeNode
}

func (s *Slice) implType() {}

type Struct struct {
	UnembeddedFields []struct {
		Name string
		Type TypeNode
		Tag  string
	}
	EmbeddedFields []TypeNode
}

func (s *Struct) implType() {}

type Pointer struct {
	To TypeNode
}

func (p *Pointer) implType() {}

type Function struct {
	Params []struct {
		Name       string
		Type       TypeNode
		IsVariadic bool
	}
	// nil indicates return nothing
	Result []TypeNode
}

func (f *Function) implType() {}

type Interface struct {
	RequiredMethods  []Function
	UnionFrom        []*Interface
	TypeRestrictions []struct {
		Type TypeNode
		// ~T is not exact (underlying), T is exact
		Exact bool
	}
}

func (t *Interface) implType() {}

type Map struct {
	KeyType   TypeNode
	ValueType TypeNode
}

func (t *Map) implType() {}

type Channel struct {
	ElemType  TypeNode
	Direction ChannelDirection
}

func (t *Channel) implType() {}

type ChannelDirection int

const (
	ChanSend ChannelDirection = iota
	ChanRecv
	ChanBoth
)
