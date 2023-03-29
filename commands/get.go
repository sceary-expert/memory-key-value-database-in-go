package commands

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sceary-expert/memory-key-value-database-in-go/responses"
	"go.mongodb.org/mongo-driver/bson"
)

// var memoryCollection *mongo.Collection = configs.GetCollection(configs.DB, "memory")
// var validate = validator.New()

type GetRequestBody struct {

	// <key>
	// The key under which the given value will be stored. (required)
	Key string `json:"key" validate:"required"`
}

func Get() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()

		var getRequestBody GetRequestBody
		if err := json.NewDecoder(r.Body).Decode(&getRequestBody); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&getRequestBody); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		// key := setRequestBody.Key
		// value := setRequestBody.Value
		// expiryTime : setRequestBody.ExpiryTime
		// condition := setRequestBody.Condition
		targetKey := getRequestBody.Key
		filter := bson.D{{"key", targetKey}}

		var result SetRequestBody
		err := memoryCollection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusCreated)
		response := responses.Response{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result.Value}}
		json.NewEncoder(rw).Encode(response)

	}
}
