package main

import (
	"log"
	"time"
	"net/http"
	"encoding/json"
)

type TokenResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ResponseData struct {
	Data string `json:"data"`
}

func authenticationHandler(writer http.ResponseWriter, request *http.Request) {
	var credentials Credentials

	err := json.NewDecoder(request.Body).Decode(&credentials)
	if err != nil {
		log.Println("No credentials found in request")
		sendJsonResponse(ErrorResponse{"Please provide credentials to login"}, http.StatusUnauthorized, writer)
		return
	}

	tokenString, expires, err := createToken(credentials)
	if err != nil {
		log.Println("Error while creating/signing the token")
		sendJsonResponse(ErrorResponse{"Please try again later"}, http.StatusInternalServerError, writer)
		return
	}

	token := userToken{}
	token.usercode = credentials.Username
	token.usertoken = tokenString
	token.expires = expires
	err = createUserToken(token)
	if err != nil {
		log.Println("Error while persisting the token: ", err)
		sendJsonResponse(ErrorResponse{"Please try again later"}, http.StatusInternalServerError, writer)
		return
	}

	response := TokenResponse{tokenString}
	sendJsonResponse(response, http.StatusOK, writer)
}

func validateTokenHandler(writer http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token := r.Header.Get(UserTokenHeader)
	if token == "" {
		log.Println("No header %s found within request headers", UserTokenHeader)
		sendJsonResponse(ErrorResponse{"Must login to access API"}, http.StatusUnauthorized, writer)
		return
	}

	userToken, err := queryUserToken(token)
	if err != nil {
		log.Println("Error while querying the token: ", err)
		sendJsonResponse(ErrorResponse{"You must login to access API"}, http.StatusUnauthorized, writer)
		return
	} else if userToken == nil {
		log.Println("Token not found: [%s]", token)
		sendJsonResponse(ErrorResponse{"You must login to access API"}, http.StatusUnauthorized, writer)
		return
	} else if  time.Now().After(userToken.expires) {
		log.Println("Token [%s] has expired", token)
		sendJsonResponse(ErrorResponse{"Your token has expired, please login again"}, http.StatusUnauthorized, writer)
		return
	}

	next(writer, r)
}

func apiHandler(writer http.ResponseWriter, r *http.Request) {
	// TODO: implement the API
	response := ResponseData{"Gained access to protected resource"}
	sendJsonResponse(response, http.StatusOK, writer)
}

func sendJsonResponse(response interface{}, httpStatus int, writer http.ResponseWriter) {
	body, err := json.Marshal(response)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(httpStatus)
	writer.Write(body)
}
