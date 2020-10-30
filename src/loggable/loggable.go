package loggable

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

type LoggableModelInterface interface {
	Meta() interface{}
}

type LoggableModel struct{}

func (LoggableModel) Meta() interface{} { return nil }

type Loggable struct {
	*gorm.DB
}

func (l *Loggable) Name() string {
	return "gorm:loggable"
}

func (l *Loggable) Initialize(db *gorm.DB) error { //can be called repeatedly
	l.DB = db

	l.setupCallbacks()

	return nil
}

func (l *Loggable) setupCallbacks() {
	l.DB.Callback().Create().After("gorm:create").Register("loggable:create", l.addCreated)
	l.DB.Callback().Update().After("gorm:update").Register("loggable:update", l.addUpdated)
}

func (l *Loggable) addCreated(db *gorm.DB) {
	if l.isLoggable(db) && db.Error == nil && db.Statement.Schema != nil {
		err := l.add(db, "create")

		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (l *Loggable) addUpdated(db *gorm.DB) {
	if l.isLoggable(db) && db.Error == nil && db.Statement.Schema != nil {
		err := l.add(db, "update")

		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (l *Loggable) add(scope *gorm.DB, action string) error {
	valueOfModel := reflect.ValueOf(scope.Statement.Model).Elem()

	objectID := l.interfaceToString(valueOfModel.FieldByName(scope.Statement.Schema.PrioritizedPrimaryField.Name).Interface())
	objectType := scope.Statement.Schema.Name

	return l.addRecord(action, scope.Statement.Model, objectID, objectType)
}
func (l *Loggable) addRecord(action string, model interface{}, objectID string, objectType string) error {
	fields := l.processFields(model, action)

	var objectToJSON = model

	if len(fields) > 0 {
		objectToJSON = fields
	}

	rawObject, err := json.Marshal(objectToJSON)

	if err != nil {
		return err
	}

	cl := ChangeLog{
		Action:     action,
		ObjectID:   objectID,
		ObjectType: objectType,
		RawObject:  string(rawObject),
	}

	return l.DB.Session(&gorm.Session{}).Create(&cl).Error
}

func (l Loggable) processFields(model interface{}, action string) map[string]interface{} {
	valueOf := reflect.ValueOf(model).Elem()
	typeOf := valueOf.Type()

	tagName := "loggable"

	var fields = make(map[string]interface{})

	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		tag := field.Tag.Get(tagName)

		tags := strings.Split(tag, ",")

		var foundTag bool

		for _, tag := range tags {
			if tag == action {
				foundTag = true
			}
		}

		if foundTag {
			fields[field.Name] = valueOf.FieldByName(field.Name).Interface()
		}
	}

	return fields
}

func (l *Loggable) isLoggable(scope *gorm.DB) bool {
	_, ok := scope.Statement.Model.(LoggableModelInterface)

	return ok
}

func (l *Loggable) interfaceToString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	default:
		return fmt.Sprint(v)
	}
}

func New() *Loggable {
	return &Loggable{}
}
