package acmelib

type templateRef interface {
	EntityID() EntityID
}

type template[Ent any, Ref templateRef] struct {
	entity Ent
	refs   *set[EntityID, Ref]
}

func newTemplate[Ent any, Ref templateRef](e Ent) *template[Ent, Ref] {
	return &template[Ent, Ref]{
		entity: e,

		refs: newSet[EntityID, Ref](),
	}
}

func (t *template[Ent, Ref]) addRef(ref Ref) {
	t.refs.add(ref.EntityID(), ref)
}

func (t *template[Ent, Ref]) removeRef(refID EntityID) {
	t.refs.remove(refID)
}
