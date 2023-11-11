package main

type Coordinates struct {
	X string `json:"x"`
	Y string `json:"y"`
	Z string `json:"z"`
}

type Rotations struct {
	P string `json:"p"`
	Y string `json:"y"`
	R string `json:"r"`
}

type Status struct {
	PayloadPower string `json:"payloadPower"`
	DataWaiting  string `json:"dataWaiting"`
	ChargeStatus string `json:"chargeStatus"`
	Voltage      string `json:"voltage"`
}

type TelemetryData struct {
	Coordinates Coordinates `json:"coordinates"`
	Rotations   Rotations   `json:"rotations"`
	Fuel        string      `json:"fuel"`
	Temp        string      `json:"temp"`
	Status      Status      `json:"status"`
}

// Request represents the JSON data structure for incoming requests.
type RedirectRequest struct {
	Verb string `json:"verb"`
	URI  string `json:"uri"`
	Data string `json:"data"`
}
