# Banco de Dados Polyglot Persistence

Augusto Pereira Teixeira - 24.123.008-5 <br>
João Pedro Bazoli Palma - 24.123.041-6

## Dependências
- Docker
- Docker Compose
- Java (JDK mais recente)
- Driver JDBC do PostgreSQL (baixar [aqui](https://jdbc.postgresql.org/download))
- GoLang (versão mais recente)

## Uso
Inicie os bancos com:

```bash
docker compose up -d
```

Teste a conexão do PostgreSQL no java executando os seguintes comandos na pasta 'Postgre_Connection'

```bash
javac -cp "postgresql-42.7.3.jar" Main.java
java -cp ".:postgresql-42.7.3.jar" Main
```

## Configuração inicial do keyspace do Scylla 
Após subir os containers com o docker espere de 30-60 segundos e execute os seguintes comandos em bash
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
se tudo ocorreu bem s2_scylla aparece na lista e você já pode rodar os codigo relacionados ao S2_Scylla


