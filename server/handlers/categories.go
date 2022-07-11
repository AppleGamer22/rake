package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/AppleGamer22/rake/server/cleaner"
	"github.com/AppleGamer22/rake/server/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func Categories(writer http.ResponseWriter, request *http.Request) {
	user, err := Verify(request)
	if err != nil {
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
		log.Println(err, user.ID.Hex())
		return
	}

	if err := request.ParseForm(); err != nil {
		http.Error(writer, "failed to read request form", http.StatusBadRequest)
		return
	}

	for _, category := range user.Categories {
		editedCategory := cleaner.Line(request.Form.Get(category))
		filter := bson.M{
			"$or": bson.A{
				bson.M{
					"_id": user.ID,
				},
				bson.M{
					"U_ID": user.ID.Hex(),
				},
			},
		}
		operations := []mongo.WriteModel{}
		switch editedCategory {
		case "":
			continue
		case http.MethodDelete:
			updateOperation := mongo.NewUpdateManyModel()
			filter["categories"] = category
			updateOperation.SetFilter(filter)
			updateOperation.SetUpdate(bson.M{
				"$pull": bson.M{
					"categories": category,
				},
			})
			operations = append(operations, updateOperation)
		default:
			updateOperation := mongo.NewUpdateOneModel()
			filter["categories"] = category
			updateOperation.SetFilter(filter)
			updateOperation.SetUpdate(bson.M{
				"$set": bson.M{
					"categories.$": editedCategory,
				},
			})
			operations = append(operations, updateOperation)

			sortOperation := mongo.NewUpdateManyModel()
			filter["categories"] = editedCategory
			sortOperation.SetFilter(filter)
			sortOperation.SetUpdate(bson.M{
				"$push": bson.M{
					"$each": bson.A{},
					"$sort": 1,
				},
			})
			operations = append(operations, sortOperation)
		}

		bulkOptions := options.BulkWriteOptions{}
		bulkOptions.SetOrdered(true)

		writeConcern := writeconcern.New(writeconcern.WMajority())
		readConcern := readconcern.Snapshot()
		transactionOptions := options.Transaction().SetWriteConcern(writeConcern).SetReadConcern(readConcern)
		session, err := db.Client.StartSession()
		if err != nil {
			log.Println(err, category, editedCategory)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		defer session.EndSession(context.Background())

		_, err = session.WithTransaction(context.Background(), func(ctx mongo.SessionContext) (interface{}, error) {
			if _, err := db.Users.BulkWrite(ctx, operations, &bulkOptions); err != nil {
				return nil, err
			}

			_, err := db.Histories.BulkWrite(ctx, operations, &bulkOptions)

			return nil, err
		}, transactionOptions)

		if err != nil {
			log.Println(err, category, editedCategory)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(writer, request, request.Referer(), http.StatusTemporaryRedirect)
}
