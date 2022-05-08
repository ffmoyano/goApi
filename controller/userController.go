package controller

import (
	"github.com/ffmoyano/goApi/logger"
	"github.com/ffmoyano/goApi/model"
	"github.com/ffmoyano/goApi/service"
	"github.com/ffmoyano/goApi/util"
	"net/http"
)

// GetUsers returns the global user list to the client
func GetUsers(w http.ResponseWriter, _ *http.Request) {
	var users []model.User
	var dtoResponse []model.UserResponse

	// if service is successful it collects the users, if not it returns an error response to the client
	if users, err = service.FindAllUsers(); err != nil {
		logger.ErrorLogger.Printf("There was an error in FindAllUsers service: %s", err)
		util.Response(w, http.StatusInternalServerError,
			response{"Error": "Something happened, please try again later"})
	}
	// convert the users to appropriate dto model for the client and append to response
	for _, user := range users {
		dtoResponse = append(dtoResponse, model.UserToUserResponse(user))
	}
	// send response with ok status and the dto model
	util.Response(w, http.StatusOK, dtoResponse)
}
