package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
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

func AddPlaylist(db *mongo.Database, playlistName string, creator string, public string) error {
	collection := db.Collection("Playlists")
	filter := bson.M{"name": playlistName}
	var existing Playlist

	singleResult := collection.FindOne(context.TODO(), filter)

	err := singleResult.Decode(&existing)

	if err == nil {
		return fmt.Errorf("já existe uma playlist com o nome '%s'", playlistName)
	} else if err != mongo.ErrNoDocuments {
		return fmt.Errorf("erro ao buscar playlist: %w", err)
	}

	publicBool := false
	if public == "true" {
		publicBool = true
	}

	playlist := Playlist{
		PlaylistID: primitive.NewObjectID(),
		Name:       playlistName,
		Creator:    creator,
		Tracks:     []Track{},
		CreatedAt:  time.Now(),
		Public:     publicBool,
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

func listPlaylist(db *mongo.Database, playlistName string) error {
	collection := db.Collection("Playlists")
	var playlist Playlist
	filter := bson.M{"name": playlistName}
	err := collection.FindOne(context.TODO(), filter).Decode(&playlist)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("nenhuma playlist com o nome %s encontrada", playlistName)
		}
		return err
	}

	fmt.Printf("Playlist: %s\n", playlist.Name)
	fmt.Printf("Criador: %s\n", playlist.Creator)
	fmt.Printf("Publica: %t\n", playlist.Public)
	fmt.Printf("Criada em: %s\n", playlist.CreatedAt.Format("02/01/2006 15:04:05"))
	fmt.Println("Tracks:")

	if len(playlist.Tracks) == 0 {
		fmt.Printf(" Nenhuma faixa na playlist")
	} else {
		for i, track := range playlist.Tracks {
			fmt.Printf("  %d) %s\n", i+1, track.TrackName)
			fmt.Printf("     ID: %d\n", track.TrackID)
			fmt.Printf("     Duração: %d segundos\n", track.TrackDuration)
			fmt.Printf("     Gênero: %s\n", track.TrackGenre)
			fmt.Printf("     Álbum: %s\n", track.TrackAlbum)
		}
	}
	return nil
}

func main() {
	args := os.Args[1:]
	s1 := args[0]
	s2 := args[1]
	func_call := s1 + s2
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("falha ao conctar")
	}

	database := client.Database("bancosproj")
	colectionName := "Playlists"

	collections, err := database.ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		log.Fatalf("Erro listando coleções: %v", err)
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

	switch func_call {
	case "addplaylist":
		err := AddPlaylist(database, args[2], args[3], args[4])
		if err != nil {
			log.Fatal(err)
		}
	case "addtrack":
		num_args3, err := strconv.Atoi(args[3])
		if err != nil {
			log.Fatal("erro ao converter string")
		}
		num_args4, err2 := strconv.Atoi(args[4])
		if err2 != nil {
			log.Fatal("erro ao converter string")
		}
		add_track_to_playlist(database, args[2], num_args3, num_args4, args[5], args[6], args[7])
	case "listplaylist":
		listPlaylist(database, args[2])
	}

}
