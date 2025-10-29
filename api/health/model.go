package health

type Health struct {
	Failing 		bool `json:"failing"`
	MinResponseTime int  `json:"minResponseTime"`
}