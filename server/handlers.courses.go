package server

import (
	"e-learn/dbg"
	"e-learn/mongodb"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (s *Server) GetCourses(res http.ResponseWriter, req *http.Request) {
	defer dbg.MonitorFunc("api get courses")()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, res)
		return
	}

	type Query struct {
		Data struct {
			WorkspaceId string `json:"workspace_id"`
			CategoryId  string `json:"cat_id"`
		} `json:"data"`
	}

	var query Query
	if err = json.Unmarshal(body, &query); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	courses, err := s.MongoClient.GetCourses(query.Data.WorkspaceId, query.Data.CategoryId)
	dbg.ConsoleLog(courses)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	type Response struct {
		Data []*mongodb.Course `json:"data"`
	}
	response := Response{
		Data: courses,
	}

	if err = json.NewEncoder(res).Encode(response); err != nil {
		respondWithError(err, http.StatusInternalServerError, res)
		return
	}

	return
}

func (s *Server) InsertCourse(res http.ResponseWriter, req *http.Request) {
	defer dbg.MonitorFunc("api insert course")()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, res)
		return
	}

	type Query struct {
		Data *mongodb.Course `json:"data"`
	}

	var query Query
	if err = json.Unmarshal(body, &query); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	if err := s.MongoClient.InsertCourse(query.Data); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	return
}
