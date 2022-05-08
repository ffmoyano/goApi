package controller

import (
	"encoding/json"
	"github.com/ffmoyano/goApi/logger"
	"github.com/ffmoyano/goApi/model"
	"github.com/ffmoyano/goApi/service"
	"github.com/ffmoyano/goApi/util"
	"math/rand"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"time"
)

type response map[string]string

var err error
var tokenResponse model.TokenResponse

// SignUp generates a new user with the provided data
func SignUp(w http.ResponseWriter, r *http.Request) {
	var dto model.SignUpRequest
	if err = json.NewDecoder(r.Body).Decode(&dto); err != nil {
		logger.ErrorLogger.Printf("Error decoding body request: %s", err)
		util.Response(w, http.StatusBadRequest, "Error decoding json request")
	}
	errors := validateSignUp(dto, w)
	// if there is no errors we hash the password and
	if len(errors) == 0 {
		var hashedPassword string
		if hashedPassword, err = util.HashPassword(dto.Password); err != nil {
			logger.ErrorLogger.Printf("Error hashing password: %s", err)
			util.Response(w, http.StatusInternalServerError, response{"Error": "The user could not be created"})
		} else {
			dto.Password = hashedPassword
		}
		if err = service.SignUp(dto); err != nil {
			logger.ErrorLogger.Printf("Error in SignUp service: %s", err)

			util.Response(w, http.StatusConflict, response{"Error": "The user could not be created"})
		} else {
			util.Response(w, http.StatusCreated, response{"Success": "User successfully created"})
		}

	} else {
		logger.ErrorLogger.Printf("Error in the information provided by the client: %s", err)
		util.Response(w, http.StatusBadRequest, errors)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var dto model.LoginRequest
	var user model.User

	if err = json.NewDecoder(r.Body).Decode(&dto); err != nil {
		logger.ErrorLogger.Printf("Error decoding body request: %s", err)
		util.Response(w, http.StatusInternalServerError,
			response{"Error": "Something happened, please try again later"})
	}
	if user, err = service.FindByUsername(dto.Username); err != nil {
		logger.ErrorLogger.Printf("Error in FindByUsername service: %s", err)
		util.Response(w, http.StatusInternalServerError,
			response{"Error": "Something happened, please try again later"})
	}
	if user.Username != "" && util.CheckPasswordHash(dto.Password, user.Password) {
		dto.Password = user.Password

		if tokenResponse, err = createTokenResponse(); err != nil {
			logger.ErrorLogger.Printf("Error in createTokenResponse: %s", err)
			util.Response(w, http.StatusInternalServerError,
				response{"Error": "The jwt couldn't be formed"})
		}

		err = service.InsertRefreshToken(user, tokenResponse)
		if err != nil {
			logger.ErrorLogger.Printf("Error in InsertRefreshToken service: %s", err)
			util.Response(w, http.StatusInternalServerError,
				response{"Error": "Error inserting refresh token in db"})
		}

		util.Response(w, http.StatusOK, tokenResponse)
	} else {
		logger.InfoLogger.Printf("Username or password incorrect: %s", err)
		util.Response(w, http.StatusNotFound,
			response{"Error": "There was no user found with that username and password"})
	}
}

func validateSignUp(dto model.SignUpRequest, w http.ResponseWriter) map[string]string {
	errors := make(map[string]string)

	if len(dto.Name) == 0 {
		errors["NameError"] = "The name can't be empty."
	}

	if len(dto.Username) == 0 {
		errors["UserNameError"] = "The username can't be empty."
	} else {
		var user model.User
		user, err = service.FindByUsername(dto.Username)
		if err != nil {
			logger.InfoLogger.Printf("Username already exists: %s", err)
			util.Response(w, http.StatusInternalServerError,
				response{"Error": "Error checking username duplicate"})
		} else {
			if user.ID != 0 {
				errors["UserNameError"] = "The username already exists."
			}
		}
	}

	if _, err = mail.ParseAddress(dto.Email); err != nil {
		errors["EmailError"] = err.Error()
	}

	if len(dto.Password) < 8 || len(dto.Password) > 40 {
		errors["PasswordError"] = "The password must have between 8 and 40 characters."
	}
	return errors
}

func createTokenResponse() (model.TokenResponse, error) {
	var token string

	expirationTime := time.Now().Add(30 * time.Minute)
	if token, err = util.AssembleJWT(util.Payload{Sub: "ffmoyano", Exp: expirationTime.Unix()},
		os.Getenv("secret")); err != nil {
		return tokenResponse, err
	}
	tokenResponse.Jwt = "Bearer " + token
	tokenResponse.RefreshToken = generateRefreshToken()
	return tokenResponse, nil

}

func generateRefreshToken() string {
	rand.Seed(time.Now().Unix())
	var refresh strings.Builder
	charSet := "abcdedfghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPRSTUVXYZ"
	length := 26
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		refresh.WriteString(string(randomChar))
	}
	return refresh.String()
}
