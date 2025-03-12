package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"naevis/handlers"
	"naevis/initdb"
	"naevis/mongops"
	"naevis/structs"
	"net/http"

	"github.com/quic-go/quic-go/http3"
	_ "modernc.org/sqlite"
)

// Server holds our dependencies such as the SQLite DB.
type Server struct {
	db *sql.DB
}

func main() {
	// Initialize SQLite DB.
	db, err := initdb.InitDB("events.db")
	if err != nil {
		log.Fatalf("Failed to initialize DB: %v", err)
	}
	defer db.Close()

	// Create our server instance.
	srv := &Server{db: db}

	// Set up HTTP mux with our event handler.
	mux := http.NewServeMux()
	mux.HandleFunc("/event", srv.EventHandler)
	mux.HandleFunc("/events/", handlers.GetEventsByTypeHandler) // Matches /events/{ENTITY_TYPE}

	// Start the QUIC server using TLS.
	quicServer := &http3.Server{
		Addr:    ":4433",
		Handler: mux,
	}

	log.Println("QUIC server listening on port 4433...")
	log.Fatal(quicServer.ListenAndServeTLS("cert.pem", "key.pem"))
}

// eventHandler receives and processes incoming event POST requests.
func (s *Server) EventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read request body.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse JSON into an Index instance.
	var event structs.Index
	if err := json.Unmarshal(body, &event); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Received event: %+v", event)

	// Fetch additional data from MongoDB (dummy implementation).
	mongoData, err := mongops.FetchDataFromMongoDB(event)
	if err != nil {
		// Log the error; you can decide whether to fail the request or continue.
		log.Printf("Error fetching MongoDB data: %v", err)
		// In this example, we continue without the additional info.
	}

	// Store the event and additional MongoDB data in SQLite.
	if err := s.storeEvent(event, mongoData); err != nil {
		http.Error(w, "Failed to store event", http.StatusInternalServerError)
		log.Printf("Error storing event: %v", err)
		return
	}

	// Send a success response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"message": "Event received and stored successfully"}`)
}

// storeEvent inserts the event data along with MongoDB data into the SQLite database.
func (s *Server) storeEvent(event structs.Index, mongoData structs.MongoData) error {
	insertSQL := `
	INSERT INTO events (entity_type, action, entity_id, item_id, item_type, additional_info)
	VALUES (?, ?, ?, ?, ?, ?);`
	_, err := s.db.Exec(insertSQL,
		event.EntityType,
		event.Action,
		event.EntityId,
		event.ItemId,
		event.ItemType,
		mongoData.AdditionalInfo,
	)
	return err
}
