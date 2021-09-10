package tibero

import (
	"time"
)

type IModel interface {
	setID(id int64)
	BeforeCreate()
	BeforeUpdate()
	BeforeDelete()
	GetDeletedAt() time.Time
}

type BaseModel struct {
	ID        int64       `rei:"id,primary key" json:"id"`
	CreatedAt time.Time   `rei:"created_at" json:"createdAt"`
	UpdatedAt time.Time   `rei:"updated_at" json:"updatedAt"`
	DeletedAt interface{} `rei:"deleted_at,omitempty" json:"deletedAt,omitempty"`
}

func (b *BaseModel) setID(id int64) {
	b.ID = id
}

func (b *BaseModel) BeforeCreate() {
	var timeNow = time.Now()
	b.CreatedAt = timeNow
	b.UpdatedAt = timeNow
}

func (b *BaseModel) GetDeletedAt() time.Time {
	switch d := (b.DeletedAt).(type) {
	case time.Time:
		return d
	default:
		return time.Time{}
	}
}

func (b *BaseModel) BeforeUpdate() {
	b.UpdatedAt = time.Now()
}

func (b *BaseModel) BeforeDelete() {
	b.DeletedAt = time.Now()
}
