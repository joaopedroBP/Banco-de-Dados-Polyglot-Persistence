package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Track struct {
	TrackID       int    `bson:"track_id"`
	TrackDuration int    `bson:"track_duration"`
	TrackName     string `bson:"track_name"`
	TrackGenre    string `bson:"track_genre"`
	TrackAlbum    string `bson:"track_album_name"`
}

type Playlist struct {
	PlaylistID primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name"`
	Creator    string             `bson:"creator"`
	Tracks     []Track            `bson:"tracks"`
	CreatedAt  time.Time          `bson:"created_at"`
	Public     bool               `bson:"public"`
}

func CreateTrack(id int, duration int, name string, genre string, album string) Track {
	return Track{
		TrackID:       id,
		TrackDuration: duration,
		TrackName:     name,
		TrackGenre:    genre,
		TrackAlbum:    album,
	}
}

func AddPlaylist(db *mongo.Database, playlistName string, creator string, public bool) error {
	collection := db.Collection("Playlists")

	filter := bson.M{"name": playlistName}
	var existing Playlist
	err := collection.FindOne(context.TODO(), filter).Decode(&existing)

	if err == nil {
		return fmt.Errorf("já existe uma playlist com o nome '%s'", playlistName)
	} else if err != mongo.ErrNoDocuments {
		return err
	}

	playlist := Playlist{
		PlaylistID: primitive.NewObjectID(),
		Name:       playlistName,
		Creator:    creator,
		Tracks:     []Track{},
		CreatedAt:  time.Now(),
		Public:     public,
	}

	_, err2 := collection.InsertOne(context.TODO(), playlist)
	if err2 != nil {
		return err2
	}

	log.Printf("Playlist %s criada com sucesso!\n", playlistName)
	return nil
}

func add_track_to_playlist(db *mongo.Database, playlistName string, id int, duration int, name string, genre string, album string) error {
	collection := db.Collection("Playlists")

	newTrack := CreateTrack(id, duration, name, genre, album)
	filter := bson.M{"name": playlistName}
	update := bson.M{"$push": bson.M{"tracks": newTrack}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("erro ao adicionar track: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("nenhuma playlist encontrada com o nome %s", playlistName)
	}

	log.Printf("Faixa '%s' adicionada à playlist '%s' com sucesso!\n", name, playlistName)
	return nil
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("falha ao conctar")
	}
	fmt.Println("conectado")

	database := client.Database("bancosproj")
	colectionName := "Playlists"
	collections, err := database.ListCollectionNames(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	exists := false
	for _, c := range collections {
		if c == colectionName {
			exists = true
			break
		}
	}

	if !exists {
		err = database.CreateCollection(context.TODO(), colectionName)
		if err != nil {
			log.Fatalf("Falha ao criar coleção")
		}
	}

}
