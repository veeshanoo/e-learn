package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	uuid := primitive.NewObjectID()
	fmt.Println(uuid)
	fmt.Println(uuid.String())
}
