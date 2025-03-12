package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"naevis/initdb"
	"naevis/mongops"
	"naevis/structs"
	"net/http"
	"strings"

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
	mux.HandleFunc("/events/", srv.GetEventsByTypeHandler) // Matches /events/{ENTITY_TYPE}

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

// GetEventsByTypeHandler handles requests to /events/{ENTITY_TYPE}?query=QUERY
func (s *Server) GetEventsByTypeHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ENTITY_TYPE from the URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/events/"), "/")
	if len(pathParts) == 0 || pathParts[0] == "" {
		http.Error(w, "Missing ENTITY_TYPE in URL", http.StatusBadRequest)
		return
	}
	entityType := pathParts[0]

	log.Println(entityType)

	// Get query parameter
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests allowed", http.StatusMethodNotAllowed)
		return
	}

	// Convert the events slice to JSON.
	response, err := json.Marshal(GetResultsOfType(entityType, query))
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	// Send JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

// Function to get results based on entity type
func GetResultsOfType(entityType string, query string) []structs.Result {
	var resarr []structs.Result

	switch entityType {
	case "events":
		resarr = append(resarr,
			structs.Result{
				Type:        "event",
				ID:          "event123",
				Name:        "Tech Conference 2025",
				Location:    "Conference Hall A",
				Category:    "Technology",
				Date:        "2025-06-15",
				Price:       "100",
				Description: "A conference on Go and Zig programming languages.",
				Image:       "https://example.com/event.jpg",
				Link:        "https://eventsite.com/register",
			},
			structs.Result{
				Type:        "event",
				ID:          "event456",
				Name:        "AI Summit",
				Location:    "Silicon Valley",
				Category:    "Artificial Intelligence",
				Date:        "2025-07-10",
				Price:       "200",
				Description: "The biggest AI event of the year!",
				Image:       "https://example.com/ai_summit.jpg",
				Link:        "https://aisummit.com",
			},
		)

	case "places":
		resarr = append(resarr,
			structs.Result{
				Type:        "place",
				ID:          "place789",
				Name:        "Central Park",
				Location:    "New York City",
				Category:    "Public Park",
				Rating:      "4.7",
				Description: "A beautiful park in the city center.",
				Image:       "https://example.com/central_park.jpg",
				Link:        "https://maps.google.com?q=Central+Park",
			},
			structs.Result{
				Type:        "place",
				ID:          "place101",
				Name:        "Grand Canyon",
				Location:    "Arizona, USA",
				Category:    "Natural Wonder",
				Rating:      "4.9",
				Description: "One of the most breathtaking canyons in the world.",
				Image:       "https://example.com/grand_canyon.jpg",
				Link:        "https://maps.google.com?q=Grand+Canyon",
			},
		)

	case "people":
		resarr = append(resarr,
			structs.Result{
				Type:        "people",
				ID:          "people123",
				Name:        "Alice Johnson",
				Location:    "San Francisco",
				Category:    "Software Engineer",
				Description: "An experienced developer specializing in Go and AI.",
				Image:       "https://example.com/alice.jpg",
				Link:        "https://linkedin.com/in/alicejohnson",
			},
			structs.Result{
				Type:        "people",
				ID:          "people456",
				Name:        "John Doe",
				Location:    "New York",
				Category:    "Machine Learning Expert",
				Description: "ML researcher focusing on deep learning advancements.",
				Image:       "https://example.com/johndoe.jpg",
				Link:        "https://linkedin.com/in/johndoe",
			},
		)

	case "businesses":
		resarr = append(resarr,
			structs.Result{
				Type:        "business",
				ID:          "business789",
				Name:        "TechNova",
				Location:    "Silicon Valley",
				Category:    "Tech Startup",
				Rating:      "4.8",
				Contact:     "+1 555-1234",
				Description: "A startup focused on AI and cloud computing.",
				Image:       "https://example.com/technova.jpg",
				Link:        "https://technova.com",
			},
			structs.Result{
				Type:        "business",
				ID:          "business101",
				Name:        "GreenFoods",
				Location:    "Los Angeles",
				Category:    "Organic Food Company",
				Rating:      "4.5",
				Contact:     "+1 555-5678",
				Description: "Leading organic food supplier with sustainable farming practices.",
				Image:       "https://example.com/greenfoods.jpg",
				Link:        "https://greenfoods.com",
			},
		)

	default:
		resarr = append(resarr, structs.Result{
			Type:        "unknown",
			Description: "Invalid entity type.",
		})
	}

	return resarr
}
