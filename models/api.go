package models

type JSONAPIData struct {
	Type          string      `json:"type"`
	Id            string      `json:"id"`
	Attributes    interface{} `json:"attributes"`
	Relationships interface{} `json:"relationships,omitempty"`
}

type JSONAPIError struct {
	Status string      `json:"status"`
	Source interface{} `json:"source,omitempty"`
	Title  string      `json:"title"`
	Detail string      `json:"detail"`
}

type JSONAPI struct {
	Data   interface{}    `json:"data,omitempty"`
	Errors []JSONAPIError `json:"errors,omitempty"`
	Meta   interface{}    `json:"meta,omitempty"`
}

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

func NewJSONAPIData(data interface{}) *JSONAPI {
	return &JSONAPI{
		Data: &JSONAPIData{
			Attributes: data,
		},
	}
}

type ResponseError struct {
	Message string `json:"error"`
	Code    int    `json:"error_code"`
}
