package internal

type Base struct {
	Relationships []*RelationshipRegister

	Err error
}

func (t *Base) NewRelationship() *RelationshipRegister {
	ret := &RelationshipRegister{
		owner:      t,
		fromFields: []string{},
		toFields:   []string{},
		fromTable:  "",
		toTable:    "",
	}
	t.Relationships = append(t.Relationships, ret)
	return ret

}
