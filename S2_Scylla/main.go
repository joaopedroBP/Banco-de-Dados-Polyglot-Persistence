package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/gocql/gocql"
)

type TrackTime struct {
	TrackID  int
	Duration int
}

func addToHistory(session *gocql.Session, userId string, trackId int) error {
	return session.Query(`
		INSERT INTO listening_history (user_id, ts, trackid)
		VALUES (?, ?, ?)
	`, userId, time.Now(), trackId).Exec()
}

func addListeningTime(session *gocql.Session, userId string, trackID int, listeningTime int) error {
	return session.Query(`
		UPDATE listening_time
		SET duration = duration + ?
		WHERE user_id = ? AND trackid = ?
	`, listeningTime, userId, trackID).Exec()
}

func topNTracks(session *gocql.Session, userId string, n int) error {
	iter := session.Query(`
		SELECT trackid, duration FROM listening_time
		WHERE user_id = ?
	`, userId).Iter()

	var trackID int
	var duration int
	tracks := []TrackTime{}

	for iter.Scan(&trackID, &duration) {
		tracks = append(tracks, TrackTime{TrackID: trackID, Duration: duration})
	}

	if err := iter.Close(); err != nil {
		return fmt.Errorf("Erro ao ler listening_time: %v", err)
	}

	sort.Slice(tracks, func(i, j int) bool {
		return tracks[i].Duration > tracks[j].Duration
	})

	if n > len(tracks) {
		n = len(tracks)
	}

	fmt.Printf("Top %d músicas mais ouvidas do usuário %s:\n", n, userId)
	for i := 0; i < n; i++ {
		fmt.Printf("%d. %d - %d segundos\n", i+1, tracks[i].TrackID, tracks[i].Duration)
	}

	return nil
}

func recentNTracks(session *gocql.Session, userId string, n int) error {
	var trackID int
	var ts time.Time

	iter := session.Query(`
		SELECT trackid,ts FROM listening_history
		WHERE user_id = ?
		LIMIT ?
	`, userId, n).Iter()

	fmt.Printf("Últimas %d músicas do usuário %s:\n", n, userId)
	count := 1
	for iter.Scan(&trackID, &ts) {
		fmt.Printf("%d. TrackID: %d | Data: %s\n", count, ts.Format("2006-01-02 15:04:05"))
		count++
	}

	if err := iter.Close(); err != nil {
		return fmt.Errorf("Erro ao ler listening_history: %v", err)
	}

	return nil
}

func main() {
	args := os.Args[1:]
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
			PRIMARY KEY (user_id, ts)
		) WITH CLUSTERING ORDER BY (ts DESC)
	`).Exec(); err != nil {
		log.Fatalf("Não foi possível criar a tabela listening_history: %v", err)
	}

	if err := session.Query(`
		CREATE TABLE IF NOT EXISTS listening_time (
			user_id text,
			trackid int,
			duration int,
			PRIMARY KEY (user_id, trackid)
		)
	`).Exec(); err != nil {
		log.Fatalf("Não foi possível criar a tabela listening_time: %v", err)
	}

	switch func_call {
	case "addhistory":
		trackId, err := strconv.Atoi(args[2])
		if err != nil {
			log.Fatal("trackId precisa ser um número")
		}
		if err := addToHistory(session, args[1], trackId); err != nil {
			log.Fatalf("Erro ao adicionar histórico: %v", err)
		}
		fmt.Println("Música adicionada ao histórico!")

	case "addtime":
		listeningTime, err := strconv.Atoi(args[4])
		if err != nil {
			log.Fatal("listeningTime precisa ser um número")
		}

		trackId, err := strconv.Atoi(args[3])
		if err != nil {
			log.Fatal("trackId precisa ser um número")
		}

		if err := addListeningTime(session, args[2], trackId, listeningTime); err != nil {
			log.Fatalf("Erro ao atualizar listening_time: %v", err)
		}
		fmt.Println("Tempo de reprodução atualizado com sucesso!")
	case "toptracks":
		n, err := strconv.Atoi(args[2])
		if err != nil {
			log.Fatal("O parâmetro n precisa ser um número")
		}
		if err := topNTracks(session, args[1], n); err != nil {
			log.Fatalf("Erro ao buscar musicas mais ouvidas: %v", err)
		}
	case "recenttracks":
		n, err := strconv.Atoi(args[2])
		if err != nil {
			log.Fatal("O parâmetro n precisa ser um número")
		}
		if err := recentNTracks(session, args[1], n); err != nil {
			log.Fatalf("Erro ao buscar musicas recentes: %v", err)
		}
	}
}
