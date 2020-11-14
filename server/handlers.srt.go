package server

import (
	"e-learn/dbg"
	"e-learn/mongodb"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (s *Server) GenerateSrt(res http.ResponseWriter, req *http.Request) {
	defer dbg.MonitorFunc("api generate srt")()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, res)
		return
	}

	type Query struct {
		Data struct {
			URL      string `json:"url"`
			CourseId string `json:"course_id"`
			LessonId string `json:"lesson_id"`
		}
	}

	var query Query
	if err = json.Unmarshal(body, &query); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	id, err := s.MongoClient.GetSrtLink(query.Data.URL)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	video := &mongodb.Video{
		URL:       query.Data.URL,
		Subtitles: id,
	}

	if err := s.MongoClient.UpdateLesson(query.Data.CourseId, query.Data.LessonId, video); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	return
}

func (s *Server) GetSrt(res http.ResponseWriter, req *http.Request) {
	defer dbg.MonitorFunc("api get srt")()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, res)
		return
	}

	type Query struct {
		Data struct {
			CourseId string `json:"course_id"`
			LessonId string `json:"lesson_id"`
		}
	}

	var query Query
	if err = json.Unmarshal(body, &query); err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	video, err := s.MongoClient.GetVideoWithSrt(query.Data.CourseId, query.Data.LessonId)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, res)
		return
	}

	if err = json.NewEncoder(res).Encode(video); err != nil {
		respondWithError(err, http.StatusInternalServerError, res)
		return
	}

	return
}
