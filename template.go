package acmelib

// type template[Ent any, Ref templateRef] struct {
// 	entity Ent
// 	refs   *set[EntityID, Ref]
// }

// func newTemplate[Ent any, Ref templateRef](e Ent) *template[Ent, Ref] {
// 	return &template[Ent, Ref]{
// 		entity: e,

// 		refs: newSet[EntityID, Ref](),
// 	}
// }

// func (t *template[Ent, Ref]) addRef(ref Ref) {
// 	t.refs.add(ref.EntityID(), ref)
// }

// func (t *template[Ent, Ref]) removeRef(refID EntityID) {
// 	t.refs.remove(refID)
// }

// type SignalUnitTemplate = template[*SignalUnit, *StandardSignal]

type templateRef interface {
	EntityID() EntityID
}

// type TemplateKind int

// const (
// 	TemplateKindCANIDBuilder TemplateKind = iota
// 	TemplateKindSignalType
// 	TemplateKindSignalUnit
// 	TemplateKindSignalEnum
// )

// func (tk TemplateKind) String() string {
// 	switch tk {
// 	case TemplateKindCANIDBuilder:
// 		return "can-id-builder"
// 	case TemplateKindSignalType:
// 		return "signal-type"
// 	case TemplateKindSignalUnit:
// 		return "signal-unit"
// 	case TemplateKindSignalEnum:
// 		return "signal-enum"
// 	default:
// 		return "unknown"
// 	}
// }

// type Template interface {
// 	EntityID() EntityID
// 	Name() string
// 	Desc() string
// 	CreateTime() time.Time

// 	setParentNetwork(parentNetwork *Network)
// 	ParentNetwork() *Network

// 	TemplateKind() TemplateKind

// 	ToCANIDBuilder() (*CANIDBuilder, error)
// 	ToSignalType() (*SignalType, error)
// 	ToSignalUnit() (*SignalUnit, error)
// 	ToSignalEnum() (*SignalEnum, error)
// }

type withTemplateRefs[R templateRef] struct {
	parentNetwork *Network

	// tplKind TemplateKind

	refs *set[EntityID, R]
}

func newWithTemplateRefs[R templateRef]() *withTemplateRefs[R] {
	return &withTemplateRefs[R]{
		parentNetwork: nil,

		// tplKind: tplKind,

		refs: newSet[EntityID, R](),
	}
}

func (t *withTemplateRefs[R]) isTemplate() bool {
	return t.parentNetwork != nil
}

func (t *withTemplateRefs[R]) addRef(ref R) {
	t.refs.add(ref.EntityID(), ref)
}

func (t *withTemplateRefs[R]) removeRef(refID EntityID) {
	t.refs.remove(refID)
}

// func (t *withTemplateRefs[R]) setParentNetwork(parentNetwork *Network) {
// 	t.parentNetwork = parentNetwork
// }

// func (t *withTemplateRefs[R]) ParentNetwork() *Network {
// 	return t.parentNetwork
// }

// func (t *withTemplateRefs[R]) TemplateKind() TemplateKind {
// 	return t.tplKind
// }

// func (t *withTemplateRefs[R]) References() []R {
// 	return t.refs.getValues()
// }

// func (t *withTemplateRefs[R]) ToCANIDBuilder() (*CANIDBuilder, error) {
// 	return nil, &ConversionError{
// 		From: t.tplKind.String(),
// 		To:   TemplateKindCANIDBuilder.String(),
// 	}
// }

// func (t *withTemplateRefs[R]) ToSignalType() (*SignalType, error) {
// 	return nil, &ConversionError{
// 		From: t.tplKind.String(),
// 		To:   TemplateKindSignalType.String(),
// 	}
// }

// func (t *withTemplateRefs[R]) ToSignalUnit() (*SignalUnit, error) {
// 	return nil, &ConversionError{
// 		From: t.tplKind.String(),
// 		To:   TemplateKindSignalUnit.String(),
// 	}
// }

// func (t *withTemplateRefs[R]) ToSignalEnum() (*SignalEnum, error) {
// 	return nil, &ConversionError{
// 		From: t.tplKind.String(),
// 		To:   TemplateKindSignalEnum.String(),
// 	}
// }
