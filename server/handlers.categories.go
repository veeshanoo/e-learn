package server

import (
	"e-learn/dbg"
	"e-learn/mongodb"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
)

func (s *Server) GetCategories(res http.ResponseWriter, req *http.Request) {
	defer dbg.MonitorFunc("api get categories")()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, res)
		return
	}

	type Query struct {
		Data primitive.ObjectID `json:"data"`
	}

	var query Query
	if err = json.Unmarshal(body, &query.Data); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	categories, err := s.MongoClient.GetCategories(query.Data)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	if err = json.NewEncoder(res).Encode(categories); err != nil {
		respondWithError(err, http.StatusInternalServerError, res)
		return
	}

	return
}

func (s *Server) InsertCategory(res http.ResponseWriter, req *http.Request) {
	defer dbg.MonitorFunc("api insert category")()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, res)
		return
	}

	type Query struct {
		Data *mongodb.Category `json:"data"`
	}

	var query Query
	if err = json.Unmarshal(body, &query.Data); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	if err := s.MongoClient.InsertCategory(query.Data); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	return
}
