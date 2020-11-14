package server

import (
	"e-learn/dbg"
	"e-learn/mongodb"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"net/http"
)

func (s *Server) GetTeacher(res http.ResponseWriter, req *http.Request) {
	defer dbg.MonitorFunc("api get teacher")()

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

	teacher, err := s.MongoClient.GetTeacher(bson.M{"email": query.Data})
	if err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	if err = json.NewEncoder(res).Encode(teacher); err != nil {
		respondWithError(err, http.StatusInternalServerError, res)
		return
	}

	return
}

func (s *Server) InsertTeacher(res http.ResponseWriter, req *http.Request) {
	defer dbg.MonitorFunc("api insert teacher")()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, res)
		return
	}

	type Query struct {
		Data *mongodb.Teacher `json:"data"`
	}

	var query Query
	if err = json.Unmarshal(body, &query); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	if err := s.MongoClient.InsertTeacher(query.Data); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	return
}
