package api

import (
	"log"
	"net/http"

	"chatbot/cmd/service/chatbot"
	"chatbot/cmd/service/user"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type APIServer struct {
	addr        string
	mongoClient *mongo.Client
}

func NewAPIServer(addr string, mongoClient *mongo.Client) *APIServer {
	return &APIServer{
		addr:        addr,
		mongoClient: mongoClient,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()

	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// User Router
	userStore := user.NewStore(s.mongoClient)
	userHandler := *user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	// Chatbot Router
	chatbotStore := chatbot.NewStore(s.mongoClient)
	chatbotHanlder := *chatbot.NewHandler(chatbotStore, userStore)
	chatbotHanlder.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)
}
