# Banco-de-Dados-Polyglot-Persistence.

Augusto Pereira Teixeira - 24.123.008-5
João Pedro Bazoli Palma - 24.123.041-6

## Instalação e uso
Tenha Docker Instalado
Tenha Java instalado
(De preferência ambos em versões recentes)

inicie os bancos com:

```bash
docker compose up -d
```

Teste a conexão do PostgreSQL no java executando os seguintes comandos na pasta 'Postgre_Connection'

```bash
javac -cp "postgresql-42.7.3.jar" Main.java
java -cp ".:postgresql-42.7.3.jar" Main
```
