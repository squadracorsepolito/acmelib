package acmelib

import "fmt"

type MessageSortingMethod string

const (
	MessagesByName       MessageSortingMethod = "messages_by_name"
	MessagesByCreateTime MessageSortingMethod = "messages_by_create_time"
	MessagesByUpdateTime MessageSortingMethod = "messages_by_update_time"
)

var messagesSorter = newEntitySorter(
	newEntitySorterMethod(MessagesByName, func(messages []*Message) []*Message { return sortByName(messages) }),
	newEntitySorterMethod(MessagesByCreateTime, func(messages []*Message) []*Message { return sortByCreateTime(messages) }),
	newEntitySorterMethod(MessagesByUpdateTime, func(messages []*Message) []*Message { return sortByUpdateTime(messages) }),
)

func SelectMessagesSortingMethod(methos MessageSortingMethod) {
	messagesSorter.selectSortingMethod(methos)
}

type Message struct {
	*entity
	ParentMessage *Node

	signals *entityCollection[*Signal]

	Size int
}

func NewMessage(name, desc string, size int) *Message {
	return &Message{
		entity: newEntity(name, desc),

		signals: newEntityCollection[*Signal](),

		Size: size,
	}
}

func (m *Message) errorf(err error) error {
	return m.ParentMessage.errorf(fmt.Errorf("message %s: %v", m.Name, err))
}

func (m *Message) UpdateName(name string) error {
	if err := m.ParentMessage.messages.updateName(m.ID, m.Name, name); err != nil {
		return m.errorf(err)
	}

	return m.entity.UpdateName(name)
}

func (m *Message) AddSignal(signal *Signal) error {
	if err := m.signals.addEntity(signal); err != nil {
		return m.errorf(err)
	}

	signal.ParentMessage = m
	m.setUpdateTimeNow()

	return nil
}

func (m *Message) ListSignals() []*Signal {
	return signalsSorter.sortEntities(m.signals.listEntities())
}

func (m *Message) RemoveSignal(signalID EntityID) error {
	if err := m.signals.removeEntity(signalID); err != nil {
		return m.errorf(err)
	}

	m.setUpdateTimeNow()

	return nil
}
