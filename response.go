// Server response
package main

import (
	"encoding/json"
)

type Response struct {
	// true if the request finished without errors.
	// message then contains the response. If false
	// message contains the error.
	Success bool `json:"success"`

	Message string `json:"message"`
}

func NewResponse(msg string, rc bool) *Response {
	r := Response{Success: rc, Message: msg}
	return &r
}

func (r *Response) toJSON() ([]byte, error) {
	return json.Marshal(r)
}
