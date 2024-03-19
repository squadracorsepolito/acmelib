package acmelib

import (
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

type entity struct {
	ID         EntityID
	Name       string
	Desc       string
	CreateTime time.Time
	UpdateTime time.Time
}

func newEntity(name, desc string) *entity {
	id := newEntityID()
	createTime := time.Now()

	return &entity{
		ID:         id,
		Name:       name,
		Desc:       desc,
		CreateTime: createTime,
		UpdateTime: createTime,
	}
}

func (be *entity) getID() EntityID {
	return be.ID
}

func (be *entity) getName() string {
	return be.Name
}

func (be *entity) getCreateTime() time.Time {
	return be.CreateTime
}

func (be *entity) getUpdateTime() time.Time {
	return be.UpdateTime
}

func (be *entity) setUpdateTimeNow() {
	be.UpdateTime = time.Now()
}

func (be *entity) UpdateDesc(desc string) error {
	if be.Desc != desc {
		be.Desc = desc
		be.setUpdateTimeNow()
	}

	return nil
}

func (be *entity) UpdateName(name string) error {
	if be.Name != name {
		be.Name = name
		be.setUpdateTimeNow()
	}

	return nil
}
