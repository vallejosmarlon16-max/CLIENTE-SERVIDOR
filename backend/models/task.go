package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Task representa el modelo de una tarea en la base de datos MongoDB.
type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Completed   bool               `bson:"completed" json:"completed"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}
