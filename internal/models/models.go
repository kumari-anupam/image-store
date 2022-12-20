package models

// ResponseError model for err response.
type ResponseError struct {
	HTTPStatusCode int      `json:"httpStatusCode"`
	ErrorCode      string   `json:"errorCode"`
	Recommendation []string `json:"recommendation"`
	MessageDetails string   `json:"messageDetails"`
}
