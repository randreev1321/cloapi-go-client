package request_tools

import (
	"encoding/json"
	"fmt"
)

type DefaultError struct {
	Code         int      `json:"code"`
	Title        string   `json:"title"`
	ErrorMessage ErrorMsg `json:"error"`
}

type ErrorMsg struct {
	Message string `json:"message"`
}

func (de DefaultError) Error() string {
	bt, e := json.Marshal(de)
	if e != nil {
		return fmt.Sprintf("can't decode a error body, %v", e)
	}
	return string(bt)
}
