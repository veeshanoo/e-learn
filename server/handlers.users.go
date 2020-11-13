package server

import (
	"e-learn/dbg"
	"e-learn/mongodb"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func respondWithError(err error, statusCode int, response http.ResponseWriter) {
	dbg.ConsoleLog(err)

	type ErrorStatus struct {
		Error string `json:"status"`
	}
	errStatus := ErrorStatus{
		Error: fmt.Sprintf("%s", err),
	}

	response.WriteHeader(statusCode)
	if tmpErr := json.NewEncoder(response).Encode(errStatus); tmpErr != nil {
		dbg.ConsoleLog(tmpErr)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) authWrapper(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		f(w, r)
		return
	}
}

func (s *Server) Register(res http.ResponseWriter, req *http.Request) {
	defer dbg.MonitorFunc("api register")()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, res)
		return
	}

	var query mongodb.UserAuth
	if err = json.Unmarshal(body, &query); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	session, err := s.MongoClient.InsertNewUser(&query)
	if err != nil {
		respondWithError(err, http.StatusForbidden, res)
		return
	}

	if err = json.NewEncoder(res).Encode(session); err != nil {
		respondWithError(err, http.StatusInternalServerError, res)
		return
	}

	return
}

func (s *Server) Login(res http.ResponseWriter, req *http.Request) {
	defer dbg.MonitorFunc("api login")()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, res)
		return
	}

	var query mongodb.UserAuth
	if err = json.Unmarshal(body, &query); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	if _, err := s.MongoClient.GetUser(query.Email, query.Password, true); err != nil {
		respondWithError(err, http.StatusUnauthorized, res)
		return
	}

	session, err := s.MongoClient.InsertSession(query.Email)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, res)
		return
	}

	if err = json.NewEncoder(res).Encode(session); err != nil {
		respondWithError(err, http.StatusInternalServerError, res)
		return
	}

	return
}
