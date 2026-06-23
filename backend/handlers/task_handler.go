package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskHandler struct {
	Collection *mongo.Collection
}

func NewTaskHandler(db *mongo.Database) *TaskHandler {
	return &TaskHandler{
		Collection: db.Collection("tasks"),
	}
}

// GetTasks obtiene todas las tareas
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var tasks []models.Task = make([]models.Task, 0)
	cursor, err := h.Collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var task models.Task
		if err := cursor.Decode(&task); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// CreateTask crea una nueva tarea
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Cuerpo de petición inválido", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, "El título es obligatorio", http.StatusBadRequest)
		return
	}

	task.ID = primitive.NewObjectID()
	task.Completed = false
	task.CreatedAt = time.Now()

	_, err := h.Collection.InsertOne(ctx, task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// UpdateTask actualiza el estado de completado o contenido de una tarea (?id=...)
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Se requiere el parámetro 'id'", http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Cuerpo de petición inválido", http.StatusBadRequest)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"completed": task.Completed,
		},
	}
	// Si viene un título o descripción en la petición de actualización, también los actualizamos
	if task.Title != "" {
		update["$set"].(bson.M)["title"] = task.Title
	}
	if task.Description != "" {
		update["$set"].(bson.M)["description"] = task.Description
	}

	_, err = h.Collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Tarea actualizada correctamente"})
}

// DeleteTask elimina una tarea (?id=...)
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Se requiere el parámetro 'id'", http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	_, err = h.Collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Tarea eliminada correctamente"})
}
