package models

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type eventDB struct {
	//we always make copies of this session for parallelization purposes.
	session    *mgo.Session
	Database   string
	Collection string
}

//EventDB contract
type EventDB interface {
	AddEvent(*Event) (string, error)
	FindEventByID(string) (*Event, error)
	FindEventByName(string) (*Event, error)
	FindAllEvents() ([]Event, error)
}

//Contract to control the underlying sessions of mongodb
type MongoLayer interface {
	//we will replicate eventDB's session everywhere for parallelism
	GetNewSession() *mgo.Session
}

//EventService : Interface abstraction that can be passed to
// handlers to be able to call underlying methods and allowing us to swap logic
// as long as the input output contract isn't violated.
type EventService interface {
	EventDB    //This specifies the contracts specific to the events.
	MongoLayer //This specifies the contracts specific to mongodb
}

type eventService struct {
	EventDB
	MongoLayer
}

//NewEventService : creates a new EventService when required.
func NewEventService(session *mgo.Session, database string, collection string) EventService {
	edb := &eventDB{session: session, Database: database, Collection: collection}
	return &eventService{EventDB: edb, MongoLayer: edb}
}

/*
Functions that make the eventDB struct fit the EventDB interface
*/

func (edb *eventDB) AddEvent(event *Event) (string, error) {
	//get a new mongo session.
	session := edb.GetNewSession()
	defer session.Close()

	if !event.ID.Valid() {
		event.ID = bson.NewObjectId()
	}

	if !event.Location.ID.Valid() {
		event.Location.ID = bson.NewObjectId()
	}

	//ensures that you cannot add two events with the same name
	eventCollection := session.DB(edb.Database).C(edb.Collection)
	err := eventCollection.EnsureIndex(mgo.Index{
		Key:    []string{"name"},
		Unique: true,
	})

	if err != nil {
		return "", err
	}

	err = eventCollection.Insert(event)
	if err != nil {
		return "", err
	}

	return event.ID.String(), nil
}

func (edb *eventDB) FindEventByID(id string) (*Event, error) {
	session := edb.GetNewSession()
	defer session.Close()

	event := &Event{}

	err := session.DB(edb.Database).C(edb.Collection).Find(bson.M{"_id": bson.ObjectId(bson.ObjectIdHex(id))}).One(&event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (edb *eventDB) FindEventByName(name string) (*Event, error) {
	session := edb.GetNewSession()
	defer session.Close()

	event := &Event{}

	err := session.DB(edb.Database).C(edb.Collection).Find(bson.M{"name": name}).One(&event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (edb *eventDB) FindAllEvents() ([]Event, error) {
	session := edb.GetNewSession()
	defer session.Close()

	events := []Event{}

	err := session.DB(edb.Database).C(edb.Collection).Find(nil).All(&events)
	if err != nil {
		return events, nil
	}
	return events, err
}

/*
These function here helps eventDB fulfil the MongoLayer interface as well
*/
func (edb *eventDB) GetNewSession() *mgo.Session {
	return edb.session.Copy()
}

//NewMongoConnection : spins up a new mongo connection
func NewMongoConnection(connection string) (*mgo.Session, error) {
	session, err := mgo.Dial(connection)
	if err != nil {
		return nil, err
	}
	return session, err
}
