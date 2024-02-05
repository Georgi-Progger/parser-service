package annoucement

type Annoucement struct {
	Id           uint64
	Model        string `json:"model"`
	Price        string `json:"price"`
	Year         string `json:"year"`
	Generation   string `json:"generation"`
	Mileage      string `json:"mileage"`
	History      string `json:"history"`
	PTS          string `json:"pts"`
	Owners       string `json:"owners"`
	Condition    string `json:"condition"`
	Modification string `json:"modification"`
	EngineVolume string `json:"engine_volume"`
	EngineType   string `json:"engine_type"`
	Transmission string `json:"transmission"`
	Drive        string `json:"drive"`
	Equipment    string `json:"equipment"`
	BodyType     string `json:"body_type"`
	Color        string `json:"color"`
	Steering     string `json:"steering"`
	VIN          string `json:"vin"`
	Exchange     string `json:"exchange"`
}
