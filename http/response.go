package http

// We can define a base response struct that will be used by all other response structs.
// This will help us to maintain a consistent response format across the application.
// This will also help us to easily change the response format in the future.
// Example:
//  type TestResponse struct {
// 	*BaseResponse `json:"response,omitempty"`
// 	Data interface{} `json:"data"`
//  }
type BaseResponse struct {
	Status  int     `json:"status"`
	Message *string `json:"message,omitempty"`
	Error   string  `json:"error"`
}
