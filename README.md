# Banco de Dados Polyglot Persistence

- Augusto Pereira Teixeira - 24.123.008-5

- João Pedro Bazoli Palma - 24.123.041-6

## Uso

Os bancos de dados são iniciados via Docker Compose. Na pasta raiz do projeto, execute o seguinte comando:

```bash
docker compose up -d
```

### Configuração inicial do keyspace do Scylla 

Após iniciar o container do ScyllaDB, é necessário criar um keyspace para o projeto. Siga os passos abaixo:

- **assegure-se de que o container do ScyllaDB está em execução**

```bash
docker exec -it scylla cqlsh 172.18.0.3
CREATE KEYSPACE IF NOT EXISTS s2_scylla
WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
exit
```
Teste se deu certo rodando os seguintes comandos
```bash
docker exec -it scylla cqlsh 172.18.0.3
DESCRIBE KEYSPACES;
```

Se o keyspace `s2_scylla` aparecer na lista, a configuração foi feita com sucesso.

### Compilação do módulo PostgreSQL

```bash
javac -cp "postgresql-42.7.3.jar" Main.java
java -cp ".:postgresql-42.7.3.jar" Main
```

### Compilação dos módulos MongoDB e ScyllaDB

Dentro dos diretórios de ambos os módulos, execute `go build` para compilar os programas.

### Uso do S1 (módulo principal que interage com o usuário)

```
# Dentro do diretório S1:
$ ./s1.sh help
Usage: ./s1.sh <command> <subcommand> [args...]

Commands:
  add      Add a new entry to the database
  rm       Remove an entry from the database
  list     List entries from the database
  get      Retrieve a specific entry

Subcommands:
  user
  artist
  album
  track
  playlist
  genre
  history
  like
  mostp
  time

Each subcommand may require additional arguments. You can get more details by running:
  ./s1.sh command subcommand

Examples:
  ./s1.sh add user <username> <email> <password>
  ./s1.sh add artist <name> <description>
  ./s1.sh list track
  ./s1.sh get album <album_id>
  ./s1.sh rm playlist <user_id> <playlist_id>

Notes:
  - 'mostp' is only available under: list mostp <user_id> <top_num>
  - 'time' can be added and removed, and listed via 'list time'
```

## Dependências

- Docker
- Docker Compose
- Java (JDK mais recente)
- Driver JDBC do PostgreSQL (baixar [aqui](https://jdbc.postgresql.org/download))
- GoLang (versão mais recente)
