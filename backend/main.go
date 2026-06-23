package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"backend/handlers"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// enableCORS es un middleware sencillo para habilitar CORS
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Manejar peticiones preflight OPTIONS
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Conectando a MongoDB en: %s", mongoURI)

	// Conectar a MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Error al conectar a MongoDB: %v", err)
	}

	// Verificar conexión
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("No se pudo hacer ping a MongoDB: %v", err)
	}

	log.Println("Conexión a MongoDB establecida exitosamente.")
	db := client.Database("spa_db")
	taskHandler := handlers.NewTaskHandler(db)

	// Configurar rutas
	mux := http.NewServeMux()
	mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.GetTasks(w, r)
		case http.MethodPost:
			taskHandler.CreateTask(w, r)
		case http.MethodPut:
			taskHandler.UpdateTask(w, r)
		case http.MethodDelete:
			taskHandler.DeleteTask(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// Envolver el multiplexor de rutas con el middleware de CORS
	handler := enableCORS(mux)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Servidor corriendo en el puerto %s", port)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Error al iniciar el servidor HTTP: %v", err)
	}
}
