package instrument

type InstrumentUpdateRequest struct {
	Symbol         string  `json:"symbol" validate:"omitempty,min=1,max=20"`
	Name           string  `json:"name" validate:"omitempty,min=2,max=100"`
	InstrumentType string  `json:"instrument_type" validate:"omitempty,oneof=Equity Bond ETF Crypto Forex"`
	Exchange       string  `json:"exchange" validate:"omitempty,alphanum,min=2,max=10"`
	LastPrice      float64 `json:"last_price" validate:"omitempty,gt=0"`
}
