package entities

type TempoResponseDto struct {
	Kelvin     float64 `json:"temp_K"`
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	City       string  `json:"city"`
}