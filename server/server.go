package server

import (
	"context"
	"e-learn/dbg"
	"e-learn/mongodb"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

type Server struct {
	MongoClient *mongodb.MongoClient
	Router      *mux.Router
	Config      Config
}

type Config struct {
	Address string `json:"address"`
	Port    string `json:"port"`
}

var DefaultConfig = Config{
	Address: "0.0.0.0",
	Port:    "9999",
}

func (s *Server) getConfig() {
	defer dbg.MonitorFunc("server config init")()

	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		dbg.ConsoleLog(err)
		s.Config = DefaultConfig
		return
	}

	jsonFile, err := os.Open(fmt.Sprintf("%s%cserver%cconfig.json", path, os.PathSeparator, os.PathSeparator))
	if err != nil {
		dbg.ConsoleLog(err)
		s.Config = DefaultConfig
		return
	}

	defer func() {
		if err := jsonFile.Close(); err != nil {
			dbg.ConsoleLog("bad config file close")
		}
	}()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		dbg.ConsoleLog(err)
		s.Config = DefaultConfig
		return
	}

	if err := json.Unmarshal(bytes, &s.Config); err != nil {
		// If we have trouble reading config file we use default config
		dbg.ConsoleLog(err)
		s.Config = DefaultConfig
		return
	}
	dbg.ConsoleLog(s.Config)
}

func (s *Server) Init() error {
	s.getConfig()
	s.Router = mux.NewRouter()
	s.MongoClient = &mongodb.MongoClient{}
	return s.MongoClient.InitConn()
}

func (s *Server) BuildHandlers() {
	//s.Router.HandleFunc("/test", s.authWrapper(s.Test)).Methods("POST")
	s.Router.HandleFunc("/register", s.Register).Methods("POST")
	s.Router.HandleFunc("/login", s.Login).Methods("POST")
	s.Router.HandleFunc("/workspaces/get", s.authWrapper(s.GetWorkspaces)).Methods("POST")
	s.Router.HandleFunc("/workspaces/insert", s.authWrapper(s.InsertWorkspace)).Methods("POST")
	s.Router.HandleFunc("/categories/get", s.authWrapper(s.GetCategories)).Methods("POST")
	s.Router.HandleFunc("/categories/insert", s.authWrapper(s.InsertCategory)).Methods("POST")
	s.Router.HandleFunc("/courses/get", s.authWrapper(s.GetCourses)).Methods("POST")
	s.Router.HandleFunc("/courses/insert", s.authWrapper(s.InsertCourse)).Methods("POST")
	s.Router.HandleFunc("/courses/join", s.authWrapper(s.JoinCourse)).Methods("POST")
	s.Router.HandleFunc("/courses/get/user", s.authWrapper(s.GetCoursesForUser)).Methods("POST")
	s.Router.HandleFunc("/students/get", s.authWrapper(s.GetStudent)).Methods("POST")
	s.Router.HandleFunc("/students/insert", s.authWrapper(s.InsertStudent)).Methods("POST")
	s.Router.HandleFunc("/teachers/get", s.authWrapper(s.GetTeacher)).Methods("POST")
	s.Router.HandleFunc("/teachers/insert", s.authWrapper(s.InsertTeacher)).Methods("POST")
	s.Router.HandleFunc("/srt/gen", s.authWrapper(s.GenerateSrt)).Methods("POST")
	s.Router.HandleFunc("/srt/get", s.authWrapper(s.GetSrt)).Methods("POST")
	s.Router.HandleFunc("/lessons/insert", s.authWrapper(s.InsertLesson)).Methods("POST")
	s.Router.HandleFunc("/reviews/add", s.authWrapper(s.AddReview)).Methods("POST")
}

func (s *Server) Run() {
	mc := &mongodb.MongoClient{}
	if err := mc.InitConn(); err != nil {
		log.Fatal(err)
	} else {

	}

	if err := s.Init(); err != nil {
		log.Fatal(err)
	} else {
		dbg.ConsoleLog("successful rest server init")
	}

	s.BuildHandlers()
	//
	//var ids []string
	//var names []string
	//var course_ids []string
	//var ratings []int
	//var feedbacks []string

	//ids = append(ids, "5faf24e46d91d1bdda7b7346")
	//ids = append(ids, "5fafaf5bcdec8a9ac01877dd")
	//ids = append(ids, "5fafaf6acdec8a9ac01877df")
	//ids = append(ids, "5fafaf71cdec8a9ac01877e2")
	//
	//names = append(names, "Antonio Bandericescu")
	//names = append(names, "Veaceslav Stanislav")
	//names = append(names, "Marek Sokol")
	//names = append(names, "Ionut Vasile")
	//
	//course_ids = append(course_ids, "5fafaa5f59d46a1c3ada1904")
	//course_ids = append(course_ids, "5fafaa2359d46a1c3ada1903")
	//course_ids = append(course_ids, "5fafa9d459d46a1c3ada1902")
	//course_ids = append(course_ids, "5faf4e986fae12f8e0821378")
	//course_ids = append(course_ids, "5faf4df16fae12f8e0821377")
	//
	//ratings = append(ratings, 5)
	//ratings = append(ratings, 4)
	//ratings = append(ratings, 4)
	//ratings = append(ratings, 3)
	//ratings = append(ratings, 5)
	//
	//feedbacks = append(feedbacks, "Un curs nemaipomenit.")
	//feedbacks = append(feedbacks, "Acest curs mi-a schimbat viziunea asupra lumii.")
	//feedbacks = append(feedbacks, "Recomand.")
	//feedbacks = append(feedbacks, "A fost ok. Se putea si mai bine.")
	//feedbacks = append(feedbacks, "Doamna Chimichiuri este o profesoara extraordinara.")
	//
	//for i := 0; i < 5; i++ {
	//	for j := 0; j < 5; j++ {
	//		for k := 0; k < 4; k++ {
	//			review := &mongodb.Review{
	//				Name:   names[k],
	//				Rating: ratings[j],
	//				Feedback: feedbacks[j],
	//				StudentId: ids[k],
	//			}
	//
	//			if err := s.MongoClient.AddReview(review, course_ids[i]); err != nil {
	//				dbg.ConsoleLog("EROARE")
	//			}
	//		}
	//	}
	//}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    s.Config.Address + ":" + s.Config.Port,
		Handler: s.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		time.Sleep(2 * time.Second)
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
