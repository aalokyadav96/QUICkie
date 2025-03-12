package handlers

import (
	"encoding/json"
	"log"
	"naevis/structs"
	"net/http"
	"strings"
)

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

// GetEventsByTypeHandler handles requests to /events/{ENTITY_TYPE}?query=QUERY
func GetEventsByTypeHandler(w http.ResponseWriter, r *http.Request) {
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
