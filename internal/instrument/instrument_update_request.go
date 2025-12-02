package instrument

type InstrumentUpdateRequest struct {
	Symbol          string  `json:"symbol" validate:"required,min=2,max=50"`
	Name            string  `json:"name" validate:"required,min=2,max=50"`
	Instrument_Type string  `json:"type" validate:"required,email"`
	Exchange        string  `json:"exchange" validate:"omitempty,e164"`
	Last_Price      float64 `json:"last_price" validate:"omitempty,gt=0"`
}
