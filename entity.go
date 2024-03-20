package acmelib

import (
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
	EntityID   EntityID
	Name       string
	Desc       string
	CreateTime time.Time
	UpdateTime time.Time
}

func newEntity(name, desc string) *entity {
	id := newEntityID()
	createTime := time.Now()

	return &entity{
		EntityID:   id,
		Name:       name,
		Desc:       desc,
		CreateTime: createTime,
		UpdateTime: createTime,
	}
}

func (e *entity) getEntityID() EntityID {
	return e.EntityID
}

func (e *entity) getName() string {
	return e.Name
}

func (e *entity) getCreateTime() time.Time {
	return e.CreateTime
}

func (e *entity) getUpdateTime() time.Time {
	return e.UpdateTime
}

func (e *entity) setUpdateTimeNow() {
	e.UpdateTime = time.Now()
}

func (e *entity) UpdateDesc(desc string) error {
	if e.Desc != desc {
		e.Desc = desc
		e.setUpdateTimeNow()
	}

	return nil
}

func (e *entity) UpdateName(name string) error {
	if e.Name != name {
		e.Name = name
		e.setUpdateTimeNow()
	}

	return nil
}

func (e *entity) toString() string {
	var builder strings.Builder

	builder.WriteString("entity_id: " + e.EntityID.String() + "\n")
	builder.WriteString("name: " + e.Name + "\n")
	builder.WriteString("description: " + e.Desc + "\n")
	builder.WriteString("create_time: " + e.CreateTime.String() + "\n")
	builder.WriteString("update_time: " + e.UpdateTime.String() + "\n")

	return builder.String()
}
