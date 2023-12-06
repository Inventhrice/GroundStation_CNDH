package main

type Coordinates struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

type Rotations struct {
	P float32 `json:"p"`
	Y float32 `json:"y"`
	R float32 `json:"r"`
}

type Status struct {
	PayloadPower bool    `json:"payloadPower"`
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

type ShipData struct {
	Coordinates Coordinates `json:"coordinate"`
	Rotations   Rotations   `json:"rotation"`
}
