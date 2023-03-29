package commands

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sceary-expert/memory-key-value-database-in-go/configs"
	"github.com/sceary-expert/memory-key-value-database-in-go/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type QueueRequestBody struct {
	QueueName     string `json:"queue-name"`
	QueueElements string `json:"queue-elements"`
}

var queueCollection *mongo.Collection = configs.GetCollection(configs.DB, "queue")

func Qpush() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()

		var queueRequestBody QueueRequestBody
		if err := json.NewDecoder(r.Body).Decode(&queueRequestBody); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&queueRequestBody); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		targetQueueName := queueRequestBody.QueueName
		filter := bson.D{{"queuename", targetQueueName}}

		var result QueueRequestBody
		err := queueCollection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {

			// ErrNoDocuments means that the filter did not match any documents in the collection
			if err == mongo.ErrNoDocuments {
				_, err := memoryCollection.InsertOne(ctx, queueRequestBody)
				if err != nil {
					rw.WriteHeader(http.StatusInternalServerError)
					response := responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
					json.NewEncoder(rw).Encode(response)
					return
				}
				rw.WriteHeader(http.StatusInternalServerError)
				response := responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(rw).Encode(response)
				return
			}

			rw.WriteHeader(http.StatusBadRequest)
			response := responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		oldQueueElements := queueRequestBody.QueueElements
		newQueueElements := oldQueueElements + result.QueueElements

		filter = bson.D{{"queuename", queueRequestBody.QueueName}}

		update := bson.D{{"$set", bson.D{{"queueelements", newQueueElements}}}}
		result2, err2 := queueCollection.UpdateOne(context.TODO(), filter, update)
		if err2 != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": result2}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		rw.WriteHeader(http.StatusCreated)
		response := responses.Response{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result2}}
		json.NewEncoder(rw).Encode(response)
	}
}
