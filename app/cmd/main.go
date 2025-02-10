package main

import (
	"context"
	"log"

	"chatbot/cmd/api"
	"chatbot/config"
	"chatbot/db"

	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// 1. Setup up database connection.
	// MySQL
	// sqlDB, err := db.MySQLDriver(mysql.Config{
	// 	User:                 config.Envs.DBUser,
	// 	Passwd:               config.Envs.DBPassword,
	// 	Addr:                 config.Envs.DBAddress,
	// 	DBName:               config.Envs.DBName,
	// 	Net:                  "tcp",
	// 	AllowNativePasswords: true,
	// 	ParseTime:            true,
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Mongo
	mongoClient, err := db.MongoDriver(config.Envs.MongoUri)
	if err != nil {
		log.Fatal(err)
	}
	// 2. Initialize database
	// initMySQLStorage(sqlDB)
	initMongoDBStorage(mongoClient)

	// 3. Set up the server
	server := api.NewAPIServer(":8080", mongoClient)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

// func initMySQLStorage(db *sql.DB) {
// 	err := db.Ping()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Println("MYSQL DB: Successfully Connected")
// }

func initMongoDBStorage(client *mongo.Client) {
	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatal("Could not connect to MongoDB:", err)
	}
	log.Println("MONGO DB: Successfully Connected")
}
