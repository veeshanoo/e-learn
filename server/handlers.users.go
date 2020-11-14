package server

import (
	"bytes"
	"e-learn/dbg"
	"e-learn/mongodb"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (s *Server) Test(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, res)
		return
	}

	type Ceva struct {
		Unu int `json:"unu"`
		Doi int `json:"doi"`
	}

	type Query struct {
		Token string `json:"token"`
		Asd   Ceva   `json:"data"`
	}

	var query Query
	if err = json.Unmarshal(body, &query); err != nil {
		fmt.Println(err)
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	dbg.ConsoleLog(query)
}

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

func (s *Server) authWrapper(f func(res http.ResponseWriter, req *http.Request)) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Header().Set("Content-Type", "application/json")

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			respondWithError(err, http.StatusInternalServerError, res)
			return
		}

		type Query struct {
			Token string `json:"token"`
		}

		var query Query
		if err = json.Unmarshal(body, &query); err != nil {
			respondWithError(err, http.StatusBadRequest, res)
			return
		}

		if _, err := s.MongoClient.GetSession(query.Token); err != nil {
			respondWithError(err, http.StatusUnauthorized, res)
			return
		}

		req.Body.Close()
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		f(res, req)
		return
	}
}

func (s *Server) Register(res http.ResponseWriter, req *http.Request) {
	defer dbg.MonitorFunc("api register")()

	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")

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

	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")

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

func (s *Server) JoinCourse(res http.ResponseWriter, req *http.Request) {
	defer dbg.MonitorFunc("api join course")()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, res)
		return
	}

	type Query struct {
		Token    string `json:"token"`
		CourseId string `json:"data"`
	}

	var query Query
	if err = json.Unmarshal(body, &query); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	if err := s.MongoClient.JoinCourse(query.CourseId, query.Token); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	return
}
