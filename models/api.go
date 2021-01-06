package models

// JSONAPIData represent object in response
type JSONAPIData struct {
	Type          string      `json:"type"`
	Id            string      `json:"id"`
	Attributes    interface{} `json:"attributes"`
	Relationships interface{} `json:"relationships,omitempty"`
}

// JSONAPIError represent error in response
type JSONAPIError struct {
	Status string      `json:"status"`
	Source interface{} `json:"source,omitempty"`
	Title  string      `json:"title"`
	Detail string      `json:"detail"`
}

// JSONAPI represent response
type JSONAPI struct {
	Data   interface{}    `json:"data,omitempty"`
	Errors []JSONAPIError `json:"errors,omitempty"`
	Meta   interface{}    `json:"meta,omitempty"`
}

// NewJSONAPIerror permit to forge new error response
func NewJSONAPIerror(status string, title string, detail string, source interface{}) *JSONAPI {
	return &JSONAPI{
		Errors: []JSONAPIError{
			{
				Status: status,
				Title:  title,
				Detail: detail,
				Source: source,
			},
		},
	}
}

// NewJSONAPIData permit to forge new data response
func NewJSONAPIData(data interface{}) *JSONAPI {
	return &JSONAPI{
		Data: &JSONAPIData{
			Attributes: data,
		},
	}
}

// ResponseError reresent error
type ResponseError struct {
	Message string `json:"error"`
	Code    int    `json:"error_code"`
}
