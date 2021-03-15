package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thecasualcoder/dobby/pkg/model"
	"net/http"
)

// Call allows to make a http call to another service and send the response
// @Summary Call a http endpoint
// @Description Make a http call to another service and send the response
// @Description Supports all REST operations
// @Tags Feature
// @Accept json
// @Produce json
// @Success 200 {object} interface{}
// @Router /call [post]
// @Param body body model.CallRequest true "'{url: http://httpbin.org/post, method: POST, body: {key: value}}' will make a post request to http://httpbin.org/post"
// https://github.com/swaggo/swag/blob/3d90fc0a5c6ef9566df81fe34425b0b35b0f651e/operation.go#L184
func (h *Handler) Call(c Context) {
	decoder := json.NewDecoder(c.GetRequestBody())
	var callRequest model.CallRequest
	err := decoder.Decode(&callRequest)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("error when decoding request: %s", err.Error())})
		return
	}
	response, err := h.makeCall(callRequest)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("error when making request to %s: %s", callRequest.URL, err.Error())})
		return
	}
	sendResponse(c, response, callRequest.URL)
}

func (h *Handler) makeCall(callRequest model.CallRequest) (*http.Response, error) {
	marshal, err := json.Marshal(callRequest.Body)
	if err != nil {
		return nil, fmt.Errorf("error when marshalling request body: %s", err)
	}
	request, err := http.NewRequest(callRequest.Method, callRequest.URL, bytes.NewBuffer(marshal))
	if err != nil {
		return nil, fmt.Errorf("error when creating new request to %s: %s", callRequest.URL, err)
	}
	return h.client.Do(request)
}
