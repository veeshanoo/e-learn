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

	workspaces, err := s.MongoClient.GetWorkspaces()
	if err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	type Response struct {
		Data []*mongodb.Workspace `json:"data"`
	}
	response := Response{
		Data: workspaces,
	}

	if err = json.NewEncoder(res).Encode(response); err != nil {
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
	if err = json.Unmarshal(body, &query); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	if err := s.MongoClient.InsertWorkspace(query.Data); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	return
}
