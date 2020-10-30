# gorm-loggable
Loggable plugin for gorm v2

#How to use it?
```
go get github.com/bgruszka/gorm-loggable
```

then after db init:
```go
db, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})
```

register loggable plugin:

```go
DB.Use(loggable.New())
```

In order to track changes of your model you have to mark a model with `loggable.LoggableModel`.

```go
type Order struct {
    ...
    loggable.LoggableModel
}
```

Next you have to specify which fields should be tracked in inserts and updates.

There are to options to use:
* to track creates:
  ```
  loggable:"create"
  ```

* to track updates:
  ```
  loggable:"update"
  ```

It's of course possible to track both:
```
loggable:"create,update"
```

More detailed example:

```go
type Order struct {
    ID         uint `gorm:"primary_key"`
    UserID     uint `gorm:"not null" loggable:"create"`
    User       User
    State      OrderState      `gorm:"not null" loggable:"create,update"`
    TotalPrice decimal.Decimal `gorm:"not null" loggable:"create"`
    Items      []OrderItem
    CreatedAt  time.Time
    UpdatedAt  time.Time
    DeletedAt  *time.Time `sql:"index"`
    loggable.LoggableModel
}
```

All the logs will be stored in `change_logs` table using this model:
```go
type ChangeLog struct {
    ID         uint      `gorm:"primaryKey"`
    CreatedAt  time.Time `gorm:"DEFAULT:current_timestamp"`
    Action     string    `gorm:"type:VARCHAR(10)"`
    ObjectID   string    `gorm:"index;type:VARCHAR(30)"`
    ObjectType string    `gorm:"index;type:VARCHAR(50)"`
    RawObject  string    `gorm:"type:JSON"`
}
```