package entities

type GetTempoResponseDto struct {
	Kelvin float64 `json:"temp_K"`
	Celsius float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
}
