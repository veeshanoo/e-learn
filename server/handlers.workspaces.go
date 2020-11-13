package server

import (
	"e-learn/dbg"
	"e-learn/mongodb"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (s *Server) GetWorkspaces(res http.ResponseWriter, req *http.Request) {
	defer dbg.MonitorFunc("api get workspaces")()

	categories, err := s.MongoClient.GetWorkspaces()
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

func (s *Server) InsertWorkspace(res http.ResponseWriter, req *http.Request) {
	defer dbg.MonitorFunc("api insert workspace")()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, res)
		return
	}

	type Query struct {
		Data *mongodb.Workspace `json:"data"`
	}

	var query Query
	if err = json.Unmarshal(body, &query.Data); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	if err := s.MongoClient.InsertWorkspace(query.Data); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	return
}
