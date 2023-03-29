package commands

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/sceary-expert/memory-key-value-database-in-go/configs"
	"github.com/sceary-expert/memory-key-value-database-in-go/responses"
	"go.mongodb.org/mongo-driver/mongo"
)

var memoryCollection *mongo.Collection = configs.GetCollection(configs.DB, "memory")
var validate = validator.New()

type SetRequestBody struct {

	// <key>
	// The key under which the given value will be stored. (required)
	Key string `json:"key" validate:"required"`
	// <value>
	// The value to be stored. (required)

	Value string `json:"value" validate:"required"`
	// <expiry time>
	// Specifies the expiry time of the key in seconds.(optional) (must be an integer value)
	ExpiryTime int `json:"expiry-time"`
	// <condition>
	// Specifies the decision to take if the key already exists(optional)
	// Accepts either NX or XX.
	// NX -- Only set the key if it does not already exist.
	// XX -- Only set the key if it already exists.
	Condition string `json:"condition"`
}

func Set() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()

		var setRequestBody SetRequestBody
		if err := json.NewDecoder(r.Body).Decode(&setRequestBody); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&setRequestBody); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		// key := setRequestBody.Key
		// value := setRequestBody.Value
		// expiryTime : setRequestBody.ExpiryTime
		// condition := setRequestBody.Condition
		result, err := memoryCollection.InsertOne(ctx, setRequestBody)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}
		rw.WriteHeader(http.StatusCreated)
		response := responses.Response{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}}
		json.NewEncoder(rw).Encode(response)

	}
}
