package interservice

import (
	"encoding/json"
	"fmt"
)

type UserRequest struct {
	Type      string                 `json:"type"`
	Token     string                 `json:"token"`
	RequestID float64                `json:"request_id,omitempty"`
	Page      float64                `json:"page,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	MessageID float64                `json:"message_id,omitempty"`
}

type ServerResponse struct {
	Users [2]float64 `json:"channels"`
}

func Responde(msg []byte) ([]byte, [2]float64, error) {
	var request UserRequest
	var response []byte
	var userIDS [2]float64
	var err error

	if err = json.Unmarshal(msg, &request); err != nil {
		err = fmt.Errorf("invalid bytes to unmarshall")
		return response, userIDS, err
	}

	if request.Token == "" {
		err = fmt.Errorf("token field is required")
		return response, userIDS, err
	}

	switch request.Type {
	case "load_chats":
		response, err = LoadChats(request.Token, request.Page)

	case "load_messages":
		response, err = LoadMessages(request.Token, request.RequestID, request.Page)

	case "new_message":
		response, err = NewMessage(request.Token, request.RequestID, request.Data)
		if err == nil {
			userIDS, err = getUserIDs(response)
		}

	case "delete_message":
		response, err = DeleteMessage(request.Token, request.MessageID)
		if err == nil {
			userIDS, err = getUserIDs(response)
		}

	default:
		err = fmt.Errorf("unknown request type")
	}

	return response, userIDS, err
}

func getUserIDs(rsp []byte) ([2]float64, error) {
	var serverResponse ServerResponse
	var userIDs [2]float64
	var err error

	if err = json.Unmarshal(rsp, &serverResponse); err != nil {
		err = fmt.Errorf(("invalid bytes(server response) to unmarshal"))
		return userIDs, err
	}

	userIDs = serverResponse.Users
	return userIDs, err
}
