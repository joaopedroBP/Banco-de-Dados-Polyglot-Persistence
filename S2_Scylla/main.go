package main

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
)

func main() {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Port = 9042
	cluster.Consistency = gocql.Quorum
	cluster.DisableInitialHostLookup = true

	// Cria sessão sem keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Não foi possível conectar ao Scylla: %v", err)
	}
	defer session.Close()

	fmt.Println("Conectado ao Scylla!")

	err = session.Query(`
    CREATE KEYSPACE IF NOT EXISTS S2_Scylla
    WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}
  `).Exec()
	if err != nil {
		log.Fatalf("Não foi possível criar a keyspace: %v", err)
	}

	fmt.Println("Keyspace S2_Scylla está pronta!")
}
