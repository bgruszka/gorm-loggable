package loggable

import "time"

type ChangeLog struct {
	ID         uint      `gorm:"primaryKey"`
	CreatedAt  time.Time `gorm:"DEFAULT:current_timestamp"`
	Action     string    `gorm:"type:VARCHAR(10)"`
	ObjectID   string    `gorm:"index;type:VARCHAR(30)"`
	ObjectType string    `gorm:"index;type:VARCHAR(50)"`
	RawObject  string    `gorm:"type:JSON"`
}
