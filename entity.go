package acmelib

import (
	"fmt"
	"strings"
	"time"

	"github.com/jaevor/go-nanoid"
)

type EntityID string

func newEntityID() EntityID {
	gen, err := nanoid.Standard(21)
	if err != nil {
		panic(err)
	}
	return EntityID(gen())
}

func (id EntityID) String() string {
	return string(id)
}

type entity struct {
	entityID   EntityID
	name       string
	desc       string
	createTime time.Time
}

func newEntity(name, desc string) *entity {
	id := newEntityID()
	createTime := time.Now()

	return &entity{
		entityID:   id,
		name:       name,
		desc:       desc,
		createTime: createTime,
	}
}

func (e *entity) EntityID() EntityID {
	return e.entityID
}

func (e *entity) Name() string {
	return e.name
}

func (e *entity) Desc() string {
	return e.desc
}

func (e *entity) CreateTime() time.Time {
	return e.createTime
}

func (e *entity) UpdateDesc(desc string) error {
	if e.desc != desc {
		e.desc = desc
	}

	return nil
}

func (e *entity) UpdateName(name string) error {
	if e.name != name {
		e.name = name
	}

	return nil
}

func (e *entity) toString() string {
	var builder strings.Builder

	builder.WriteString("entity_id: " + e.entityID.String() + "\n")
	builder.WriteString("name: " + e.name + "\n")

	if len(e.desc) > 0 {
		builder.WriteString(fmt.Sprintf("desc: %s\n", e.desc))
	}

	builder.WriteString(fmt.Sprintf("create_time: %s\n", e.createTime.Format(time.RFC3339)))

	return builder.String()
}
