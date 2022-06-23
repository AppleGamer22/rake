package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/AppleGamer22/rake/server/cleaner"
	"github.com/AppleGamer22/rake/server/db"
	"github.com/AppleGamer22/rake/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func VSCOPage(writer http.ResponseWriter, request *http.Request) {
	user, err := Verify(request)
	if err != nil {
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
		log.Println(err)
		return
	}

	if err := request.ParseForm(); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	history := db.History{
		Type: db.VSCO,
	}
	owner := cleaner.Line(request.Form.Get("owner"))
	post := cleaner.Line(request.Form.Get("post"))
	if post != "" {
		filter := bson.M{
			"post": post,
			"type": db.VSCO,
		}
		if err := db.Histories.FindOne(context.Background(), filter).Decode(&history); err != nil {
			urlString, username, err := shared.VSCO(owner, post)
			if err != nil {
				historyDisplay(user, history, []error{err}, writer)
				return
			}

			URL, err := url.Parse(urlString)
			if err != nil {
				log.Println(err)
				historyDisplay(user, history, []error{err}, writer)
				return
			}
			fileName := fmt.Sprintf("%s_%s", post, path.Base(URL.Path))

			if err := StorageHandler.Save(user, db.VSCO, username, fileName, urlString); err != nil {
				log.Println(err)
				historyDisplay(user, history, []error{err}, writer)
				return
			}

			history = db.History{
				ID:    primitive.NewObjectID().Hex(),
				U_ID:  user.ID.Hex(),
				URLs:  []string{fmt.Sprintf("storage/%s/%s/%s", db.VSCO, username, fileName)},
				Type:  db.VSCO,
				Owner: username,
				Post:  post,
				Date:  time.Now(),
			}

			if _, err := db.Histories.InsertOne(context.Background(), history); err != nil {
				historyDisplay(user, history, []error{err}, writer)
				return
			}
		}
	}

	historyDisplay(user, history, nil, writer)
}
