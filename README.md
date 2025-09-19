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
