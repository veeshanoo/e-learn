package mongodb

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
	"time"
)

func (mc *MongoClient) InsertLesson(lesson *Lesson, courseId string) error {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	collection := mc.Client.Database(MyDb.DbName).Collection(MyDb.Courses)

	course, err := mc.GetCourse(courseId)
	if err != nil {
		return err
	}

	lesson.Id = primitive.NewObjectID().Hex()
	course.Lessons = append(course.Lessons, lesson)

	filter := bson.M{"_id": courseId}
	update := bson.M{"$set": bson.M{"lessons": course.Lessons}}

	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (mc *MongoClient) UpdateLesson(courseId, lessonId string, video *Video) error {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	collection := mc.Client.Database(MyDb.DbName).Collection(MyDb.Courses)

	course, err := mc.GetCourse(courseId)
	if err != nil {
		return err
	}

	for i := 0; i < len(course.Lessons); i++ {
		if course.Lessons[i].Id == lessonId {
			course.Lessons[i].Video = video
			break
		}
	}

	filter := bson.M{"_id": courseId}
	update := bson.M{"$set": bson.M{"lessons": course.Lessons}}

	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

type URLParam struct {
	Key   string
	Value string
}

func makeRequest(url string, reqType string, reqBody []byte, urlParams []URLParam) ([]byte, error) {
	allowedMethods := map[string]bool{
		"GET":    true,
		"POST":   true,
		"DELETE": true,
		"PUT":    true,
	}

	if !allowedMethods[reqType] {
		return nil, errors.New("unknown method")
	}

	timeout := 1200 * time.Second
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest(reqType, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-type", "application/json")

	// Here we build query parameters
	params := req.URL.Query()
	for _, x := range urlParams {
		params.Add(x.Key, x.Value)
	}

	req.URL.RawQuery = params.Encode()

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return resBody, nil
}

func (mc *MongoClient) GetVideoWithSrt(courseId, lessonId string) (*Video, error) {
	course, err := mc.GetCourse(courseId)
	if err != nil {
		return nil, err
	}

	video := &Video{}
	for i := 0; i < len(course.Lessons); i++ {
		if course.Lessons[i].Id == lessonId {
			video = course.Lessons[i].Video
			break
		}
	}

	srt, err := mc.GetSrt(video.Subtitles)
	if err != nil {
		return nil, err
	}

	video.Subtitles = srt
	return video, nil
}

func (mc *MongoClient) GetSrtLink(link string) (string, error) {
	type Query struct {
		Url string `json:"url"`
	}

	address := fmt.Sprintf("http://192.168.1.130:3000/generate")
	query := Query{
		Url: link,
	}

	reqBody, err := json.Marshal(query)
	if err != nil {
		return "", err
	}

	resBody, err := makeRequest(address, "POST", reqBody, nil)
	if err != nil {
		return "", err
	}

	type Response struct {
		Data string `json:"data"`
	}

	response := Response{}
	if err := json.Unmarshal(resBody, &response); err != nil {
		return "", err
	}

	return response.Data, nil
}

func (mc *MongoClient) GetSrt(id string) (string, error) {
	type Query struct {
		Id string `json:"id"`
	}

	address := fmt.Sprintf("http://192.168.1.130:3000/srt")
	query := Query{
		Id: id,
	}

	reqBody, err := json.Marshal(query)
	if err != nil {
		return "", err
	}

	resBody, err := makeRequest(address, "POST", reqBody, nil)
	if err != nil {
		return "", err
	}

	type Response struct {
		Data string `json:"data"`
	}

	response := Response{}
	if err := json.Unmarshal(resBody, &response); err != nil {
		return "", err
	}

	return response.Data, nil
}
