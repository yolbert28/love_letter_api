package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/yolbert28/deliver_love_letter/internal/store"
	"github.com/yolbert28/deliver_love_letter/internal/transport"
)

var db *sql.DB

func initDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: No se encontró archivo .env")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("La variable DATABASE_URL no está configurada")
	}

	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error al abrir la conexión:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error al conectar con Supabase:", err)
	}

	fmt.Println("¡Conectado exitosamente a Supabase!")
}

func main() {
	initDB()
	defer db.Close()

	repo := store.NewLetterRepository(db)
	handler := transport.NewLetterHandler(repo)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /letters", handler.Create)
	mux.HandleFunc("GET /letters", handler.GetAll)
	mux.HandleFunc("GET /letters/date", handler.GetByDate)
	mux.HandleFunc("POST /letters/tap", handler.IncrementTapCount)
	mux.HandleFunc("PUT /letters/{id}", handler.Update)
	mux.HandleFunc("DELETE /letters/{id}", handler.Delete)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"https://love-letter-frontend-roan.vercel.app",
			"http://localhost:3000",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		AllowCredentials: true,
	})

	handlerWithCORS := c.Handler(mux)

	log.Println("Servidor corriendo en :8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithCORS))
}
