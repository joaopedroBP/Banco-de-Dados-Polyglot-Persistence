package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gocql/gocql"
)

func addToHistory(session *gocql.Session, userId string, trackId int, trackName string, trackGenre string) error {
	return session.Query(`
		INSERT INTO listening_history (user_id, ts, trackid, trackname, trackgenre)
		VALUES (?, ?, ?, ?, ?)
	`, userId, time.Now(), trackId, trackName, trackGenre).Exec()
}

func addListeningTime(session *gocql.Session, userId string, trackName string, listeningTime int) error {
	return session.Query(`
		UPDATE listening_time
		SET duration = duration + ?
		WHERE user_id = ? AND trackname = ?
	`, listeningTime, userId, trackName).Exec()
}

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatal("Uso: go run main.go <func1> <func2> [args...]")
	}

	func_call := args[0] + args[1]

	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Port = 9042
	cluster.Keyspace = "s2_scylla"
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Não foi possível criar sessão no keyspace s2_scylla: %v", err)
	}
	defer session.Close()

	if err := session.Query(`
		CREATE TABLE IF NOT EXISTS listening_history (
			user_id text,
			ts timestamp,
			trackid int,
			trackname text,
			trackgenre text,
			PRIMARY KEY (user_id, ts)
		) WITH CLUSTERING ORDER BY (ts DESC)
	`).Exec(); err != nil {
		log.Fatalf("Não foi possível criar a tabela listening_history: %v", err)
	}

	if err := session.Query(`
		CREATE TABLE IF NOT EXISTS listening_time (
			user_id text,
			trackname text,
			duration int,
			PRIMARY KEY (user_id, trackname)
		)
	`).Exec(); err != nil {
		log.Fatalf("Não foi possível criar a tabela listening_time: %v", err)
	}

	switch func_call {
	case "addhistory":
		if len(args) < 5 {
			log.Fatal("Uso: add history <userId> <trackId> <trackName> <trackGenre>")
		}
		trackId, err := strconv.Atoi(args[2])
		if err != nil {
			log.Fatal("trackId precisa ser um número")
		}
		if err := addToHistory(session, args[1], trackId, args[3], args[4]); err != nil {
			log.Fatalf("Erro ao adicionar histórico: %v", err)
		}
		fmt.Println("Música adicionada ao histórico!")

	case "addtime":
		if len(args) < 4 {
			log.Fatal("Uso: add time <userId> <trackName> <listeningTime>")
		}
		listeningTime, err := strconv.Atoi(args[3])
		if err != nil {
			log.Fatal("listeningTime precisa ser um número")
		}
		if err := addListeningTime(session, args[1], args[2], listeningTime); err != nil {
			log.Fatalf("Erro ao atualizar listening_time: %v", err)
		}
		fmt.Println("Tempo de reprodução atualizado com sucesso!")

	default:
		log.Fatal("Função desconhecida. Use 'add history' ou 'add time'")
	}
}
