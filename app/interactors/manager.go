package interactors

import (
	"dominus/app/domain/entities"
	"dominus/app/domain/event"
	"dominus/app/domain/rules"
	"dominus/app/domain/topic"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type manager struct {
	r          rules.RuleInt
	t          topic.TopicInt
	repo       RepositoryInt
	tdb        *entities.Topic
	collection string
	lg         event.LogsInt
}

func (m *manager) CreateTopic(ctx RestContextInt) error {

	if err := ctx.BodyParser(m.tdb); err != nil {
		m.lg.WriteLog("InsertObject", err.Error())
		return err
	}

	if err := m.r.ValidateStruct(m.tdb); err != nil {
		m.lg.WriteLog("InsertObject", err.Error())
		return err
	}

	m.t.CreateTopic(m.tdb.Topic, m.tdb.Subscribers)
	m.tdb.CreatedAt = time.Now()

	if _, err := m.repo.InsertObject(m.tdb, m.collection); err != nil {
		m.lg.WriteLog("InsertObject", err.Error())
		return err
	}
	
	return nil 
}

func (m *manager) GetTopic(ctx RestContextInt) (map[string]any, error) {
	var (
		topics  []entities.Topic
		message = make(map[string]any)
	)

	if err := m.repo.FindObjects(m.collection, &topics, bson.D{}); err != nil {
		return nil, err
	}

	message["topics"] = topics
	return message, nil
}

func (m *manager) UpdateTopic(ctx RestContextInt) error {

	if err := ctx.BodyParser(m.tdb); err != nil {
		return err
	}

	if err := m.r.ValidateStruct(m.tdb); err != nil {
		return err
	}

	m.t.CreateTopic(m.tdb.Topic, m.tdb.Subscribers)
	filter := primitive.D{primitive.E{Key: "topic", Value: m.tdb.Topic}}
	update := primitive.D{primitive.E{Key: "$set", Value: m.tdb}}
	m.tdb.UpdatedAt = time.Now()
	
	if _, err := m.repo.UpdateObject(filter, update, m.collection); err != nil {
		return err
	}

	return nil 
}

func (m *manager) DeleteTopic(ctx RestContextInt) error {

	if err := ctx.BodyParser(m.tdb); err != nil {
		return err
	}

	if err := m.r.ValidateStruct(m.tdb); err != nil {
		return err
	}

	if err := m.t.DeleteTopic(m.tdb.Topic); err != nil {
		return err
	}

	filter := primitive.D{primitive.E{Key: "topic", Value: m.tdb.Topic}}

	if _, err := m.repo.DeleteObject(filter, m.collection); err != nil {
		return err
	}

	return nil 
}

type ManagerInt interface {
	//Create a new topic in dominus
	//
	//Parameters
	//
	//-> ctx: fasthttp context
	CreateTopic(ctx RestContextInt) error
	//Return all topic in dominus
	//
	//Parameters
	//
	//-> ctx: fasthttp context
	GetTopic(ctx RestContextInt) (map[string]any, error)
	//Delete a topic using a key selected
	//
	//Parameters
	//
	//-> ctx: fasthttp context
	DeleteTopic(ctx RestContextInt) error
	//Update a topic in dominus
	//
	//Parameters
	//
	//-> ctx: fasthttp context
	UpdateTopic(ctx RestContextInt) error
}
