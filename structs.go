package main

type Coordinates struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

type Rotations struct {
	P int `json:"p"`
	Y int `json:"y"`
	R int `json:"r"`
}

type Status struct {
	PayloadPower string  `json:"payloadPower"`
	DataWaiting  bool    `json:"dataWaiting"`
	ChargeStatus bool    `json:"chargeStatus"`
	Voltage      float32 `json:"voltage"`
}

type TelemetryData struct {
	Coordinates Coordinates `json:"coordinate"`
	Rotations   Rotations   `json:"rotation"`
	Fuel        int         `json:"fuel"`
	Temp        float32     `json:"temp"`
	Status      Status      `json:"status"`
}

// Request represents the JSON data structure for incoming requests.
type RedirectRequest struct {
	Verb string `json:"verb"`
	URI  string `json:"uri"`
	Data string `json:"data"`
}
