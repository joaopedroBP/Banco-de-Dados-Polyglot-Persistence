package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Track struct {
	TrackID int `bson:"track_id"`
}

type Playlist struct {
	PlaylistID primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name"`
	Creator    string             `bson:"creator"`
	Tracks     []Track            `bson:"tracks"`
	CreatedAt  time.Time          `bson:"created_at"`
	Public     bool               `bson:"public"`
}

func AddPlaylist(db *mongo.Database, creator string, playlistName string, public string) error {
	collection := db.Collection("Playlists")
	filter := bson.M{"name": playlistName, "creator": creator}
	var existing Playlist

	singleResult := collection.FindOne(context.TODO(), filter)

	err := singleResult.Decode(&existing)

	if err == nil {
		return fmt.Errorf("já existe uma playlist com o nome '%s' para o criador '%s'", playlistName, creator)
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

func add_track_to_playlist(db *mongo.Database, creator string, playlistName string, id int) error {
	collection := db.Collection("Playlists")

	filter := bson.M{"name": playlistName, "creator": creator}
	update := bson.M{"$push": bson.M{"tracks": Track{TrackID: id}}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("erro ao adicionar track: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("nenhuma playlist encontrada com o nome '%s' para o criador '%s'", playlistName, creator)
	}

	log.Printf("Faixa '%d' adicionada à playlist '%s' com sucesso!\n", id, playlistName)
	return nil
}

func listPlaylist(db *mongo.Database, creator string, playlistName string) error {
	collection := db.Collection("Playlists")
	var playlist Playlist
	filter := bson.M{"name": playlistName, "creator": creator}
	err := collection.FindOne(context.TODO(), filter).Decode(&playlist)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("nenhuma playlist com o nome '%s' e criador '%s' encontrada", playlistName, creator)
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
			fmt.Printf("  %d) %s\n", i+1, track.TrackID)
		}
	}
	return nil
}

func remove_track(db *mongo.Database, creator string, playlistName string, trackID int) error {
	collection := db.Collection("Playlists")

	filter := bson.M{"name": playlistName, "creator": creator}
	update := bson.M{"$pull": bson.M{"tracks": bson.M{"track_id": trackID}}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("erro ao remover track: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("nenhuma playlist encontrada com o nome '%s'", playlistName)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("nenhuma track com o id %d encontrada na playlist '%s'", trackID, playlistName)
	}

	log.Printf("Track com ID %d removida da playlist '%s' com sucesso!\n", trackID, playlistName)
	return nil
}

func remove_playlist(db *mongo.Database, creator string, playlistName string) error {
	colection := db.Collection("Playlists")

	result, err := colection.DeleteOne(context.TODO(), bson.M{"name": playlistName, "creator": creator})
	if err != nil {
		return fmt.Errorf("erro ao remover playlist: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("nenhuma playlist encontrada com o nome '%s'", playlistName)
	}

	log.Printf("Playlist '%s' removida com sucesso!\n", playlistName)

	return nil
}

func listUserPlaylists(db *mongo.Database, userID string) error {
	collection := db.Collection("Playlists")

	filter := bson.M{"creator": userID}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("erro ao buscar playlists do usuário '%s': %v", userID, err)
	}
	defer cursor.Close(context.TODO())

	playlistsExist := false

	for cursor.Next(context.TODO()) {
		var playlist Playlist
		if err := cursor.Decode(&playlist); err != nil {
			return fmt.Errorf("erro ao decodificar playlist: %v", err)
		}

		playlistsExist = true
		fmt.Printf("Playlist: %s\n", playlist.Name)
		fmt.Printf("  Publica: %t\n", playlist.Public)
		fmt.Printf("  Criada em: %s\n", playlist.CreatedAt.Format("02/01/2006 15:04:05"))
		fmt.Printf("  Número de faixas: %d\n", len(playlist.Tracks))
		fmt.Println("---------------------------")
	}

	if !playlistsExist {
		fmt.Printf("Nenhuma playlist encontrada para o usuário '%s'\n", userID)
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
	if slices.Contains(collections, colectionName) {
		exists = true
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
			log.Fatalf("Falha ao adicionar playlist: %v", err)
		}

	case "addtrack":
		// add track usuário playlist trackID
		num_args4, err := strconv.Atoi(args[4])
		if err != nil {
			log.Fatalf("Erro ao converter TrackID '%s': %v", args[4], err)
		}

		err = add_track_to_playlist(database, args[2], args[3], num_args4)
		if err != nil {
			log.Fatalf("Falha ao adicionar track: %v", err)
		}
	case "listplaylist":
		// list playlist usuário playlist
		err := listUserPlaylists(database, args[2])

		if err != nil {
			log.Fatalf("Falha ao listar playlist: %v", err)
		}
	case "listplaylists":
		err := listUserPlaylists(database, args[2])
		if err != nil {
			log.Fatalf("Falha ao listar playlists do usuário: %v", err)
		}

	case "rmtrack":
		// rm track usuário playlist trackID
		num_args3, err := strconv.Atoi(args[3])
		if err != nil {
			log.Fatalf("Erro ao converter TrackID '%s': %v", args[3], err)
		}

		err = remove_track(database, args[2], args[3], num_args3)
		if err != nil {
			log.Fatalf("Falha ao remover track: %v", err)
		}

	case "rmplaylist":
		err := remove_playlist(database, args[2], args[3])
		if err != nil {
			log.Fatalf("Falha ao remover playlist: %v", err)
		}
	}

}
