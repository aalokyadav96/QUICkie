package structs

// Index represents the incoming JSON event structure.
type Index struct {
	EntityType string `json:"entity_type"`
	Action     string `json:"action"`
	EntityId   string `json:"entity_id"`
	ItemId     string `json:"item_id"`
	ItemType   string `json:"item_type"`
}

// MongoData is a dummy structure for the additional data
// fetched from MongoDB.
type MongoData struct {
	AdditionalInfo string
}

// Result represents a single search result.
type Result struct {
	Placeid     string `json:"placeid" bson:"placeid"`
	Eventid     string `json:"eventid" bson:"eventid"`
	Businessid  string `json:"businessid" bson:"businessid"`
	Peopleid    string `json:"peopleid" bson:"peopleid"`
	Type        string `json:"type" bson:"type"`
	Location    string `json:"location" bson:"location"`
	Category    string `json:"category" bson:"category"`
	Date        string `json:"date" bson:"date"`
	Price       string `json:"price" bson:"price"`
	Description string `json:"description" bson:"description"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Rating      string `json:"rating,omitempty"`
	Contact     string `json:"contact,omitempty"`
	Image       string `json:"image,omitempty"`
	Link        string `json:"link,omitempty"`
}
