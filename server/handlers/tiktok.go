package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AppleGamer22/rake/server/cleaner"
	"github.com/AppleGamer22/rake/server/db"
	"github.com/AppleGamer22/rake/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TikTokPage(writer http.ResponseWriter, request *http.Request) {
	user, err := Verify(request)
	if err != nil {
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
		log.Println(err)
		return
	}

	if err := request.ParseForm(); err != nil {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	history := db.History{
		Type: db.TikTok,
	}
	owner := cleaner.Line(request.Form.Get("owner"))
	post := cleaner.Line(request.Form.Get("post"))
	if post != "" {
		filter := bson.M{
			"post": post,
			"type": db.TikTok,
		}
		if err := db.Histories.FindOne(context.Background(), filter).Decode(&history); err != nil {
			URL, username, err := shared.TikTok(owner, post)

			if err != nil {
				log.Println(err)
				historyDisplay(user, history, []error{err}, writer)
				return
			}

			fileName := fmt.Sprintf("%s.mp4", post)
			if err := StorageHandler.Save(db.TikTok, username, fileName, URL); err != nil {
				log.Println(err)
				log.Println(err)
				historyDisplay(user, history, []error{err}, writer)
				return
			}

			history = db.History{
				ID:    primitive.NewObjectID().Hex(),
				U_ID:  user.ID.Hex(),
				URLs:  []string{fmt.Sprintf("storage/%s/%s/%s", db.TikTok, username, fileName)},
				Type:  db.TikTok,
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
