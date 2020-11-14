package server

import (
	"e-learn/dbg"
	"e-learn/mongodb"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (s *Server) InsertLesson(res http.ResponseWriter, req *http.Request) {
	defer dbg.MonitorFunc("api insert lesson")()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, res)
		return
	}

	type Query struct {
		Data struct {
			Lesson   *mongodb.Lesson `json:"lesson"`
			CourseId string          `json:"course_id"`
		} `json:"data"`
	}

	var query Query
	if err = json.Unmarshal(body, &query); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	if err := s.MongoClient.InsertLesson(query.Data.Lesson, query.Data.CourseId); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	return
}
