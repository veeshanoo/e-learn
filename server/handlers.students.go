package server

import (
	"e-learn/dbg"
	"e-learn/mongodb"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"net/http"
)

func (s *Server) GetStudent(res http.ResponseWriter, req *http.Request) {
	defer dbg.MonitorFunc("api get student")()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, res)
		return
	}

	type Query struct {
		Data string `json:"data"`
	}

	var query Query
	if err = json.Unmarshal(body, &query); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	student, err := s.MongoClient.GetStudent(bson.M{"email": query.Data})
	if err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	if err = json.NewEncoder(res).Encode(student); err != nil {
		respondWithError(err, http.StatusInternalServerError, res)
		return
	}

	return
}

func (s *Server) InsertStudent(res http.ResponseWriter, req *http.Request) {
	defer dbg.MonitorFunc("api insert student")()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, res)
		return
	}

	type Query struct {
		Data *mongodb.Student `json:"data"`
	}

	var query Query
	if err = json.Unmarshal(body, &query); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	if err := s.MongoClient.InsertStudent(query.Data); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	return
}
