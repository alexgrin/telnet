package database

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
    "mud/utils"
    //"fmt"
)

type dbError struct {
	message string
}

func (self dbError) Error() string {
	return self.message
}

func newDbError(message string) dbError {
	var err dbError
	err.message = message
	return err
}

func getCollection(session *mgo.Session, collection string) *mgo.Collection {
	return session.DB("mud").C(collection)
}

func FindUser(session *mgo.Session, name string) (bool, error) {
	c := getCollection(session, "users")
	q := c.Find(bson.M{"name": name})

	count, err := q.Count()

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func NewUser(session *mgo.Session, name string) error {

	found, err := FindUser(session, name)

	if err != nil {
		return err
	}

	if found {
		return newDbError("That user already exists")
	}

	c := getCollection(session, "users")
	c.Insert(bson.M{"name": name})

	return nil
}

func GetCharacterLocation(session *mgo.Session, name string) (string, error) {
	c := getCollection(session, "characters")
	q := c.Find(bson.M{"name": name})

	result := map[string]string{}
	err := q.One(&result)

	if err != nil {
		return "", err
	}

	return result["location"], nil
}

func SetUserLocation(session *mgo.Session, name string, locationId string) error {
	c := getCollection(session, "users")
	c.Update(bson.M{"name": name}, bson.M{"location": locationId})
	return nil
}

func GetUserCharacters(session *mgo.Session, name string) ([]string, error) {
	c := getCollection(session, "users")
	q := c.Find(bson.M{"name": name})

	result := map[string][]string{}
	err := q.One(&result)

	return result["characters"], err
}

func NewCharacter(session *mgo.Session, user string, character string) error {
	c := getCollection(session, "users")
	c.Update(bson.M{"name": user}, bson.M{"$push": bson.M{"characters": character}})

    c = getCollection(session, "characters")
    c.Insert(bson.M{"name": character})

	return nil
}

func DeleteCharacter(session *mgo.Session, user string, character string) error {
	c := getCollection(session, "users")
	c.Update(bson.M{"user": user}, bson.M{"$pull": bson.M{"characters": utils.Simplify(character)}})

    c = getCollection(session, "characters")
    c.Remove(bson.M{"name": character, "user":user})

	return nil
}

// vim: nocindent
