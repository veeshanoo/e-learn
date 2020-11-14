package server

import (
	"e-learn/dbg"
	"e-learn/mongodb"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (s *Server) AddReview(res http.ResponseWriter, req *http.Request) {
	defer dbg.MonitorFunc("api insert course")()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, res)
		return
	}

	type Query struct {
		Data struct {
			Review   *mongodb.Review `json:"review"`
			CourseId string          `json:"course_id"`
		} `json:"data"`
	}

	var query Query
	if err = json.Unmarshal(body, &query); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	if err := s.MongoClient.AddReview(query.Data.Review, query.Data.CourseId); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	return
}
