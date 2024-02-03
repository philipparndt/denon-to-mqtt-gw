package denon

type Message struct {
	Power     string  `json:"power"`
	Volume    float64 `json:"volume"`
	VolumeMax float64 `json:"volumeMax"`
}
