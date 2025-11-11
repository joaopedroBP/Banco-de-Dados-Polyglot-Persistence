/*
 *Compile com: javac -cp "postgresql-42.7.3.jar" Main.java
 *Rode com java -cp ".:postgresql-42.7.3.jar" Main
 * SOMENTE LINUX
 */

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.Statement;
import java.sql.PreparedStatement;  
import java.sql.SQLException;     


public class Main{
    public static void createUser(Connection con, String userId, String email, String senha) {
        String sql = "INSERT INTO users (user_id, email, senha) VALUES (?, ?, ?);";
        try (PreparedStatement pstmt = con.prepareStatement(sql)) {
            pstmt.setString(1, userId);
            pstmt.setString(2, email);
            pstmt.setString(3, senha);
            pstmt.executeUpdate();
            System.out.println("Usuário '" + userId + "' criado com sucesso!");
        } catch (SQLException e) {
            System.out.println("Erro ao criar usuário '" + userId + "':");
            e.printStackTrace();
        }
    }
  
  public static void main(String[] args){
    String url = "jdbc:postgresql://localhost:5432/postgres";
    String usuario = "postgres";
    String senha = "postgres";

    String function_called = "" + args[0] + args[1];
    System.out.println("função chamada = " + function_called);

    try (Connection con = DriverManager.getConnection(url,usuario,senha)){
      if(con != null){
        System.out.println("Conectado ao PostgreSQL");
            try (Statement stmt = con.createStatement()) {
                    String createUsers = """
                        CREATE TABLE IF NOT EXISTS users (
                            user_id VARCHAR(100) PRIMARY KEY,  -- nome do usuário
                            email VARCHAR(100) NOT NULL UNIQUE,
                            senha VARCHAR(100) NOT NULL,
                            data_criacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                        );
                        """;

                    String createArtist = """
                        CREATE TABLE IF NOT EXISTS artists (
                            artist_id VARCHAR(100) PRIMARY KEY,  -- nome do artista
                            user_id VARCHAR(100) REFERENCES users(user_id) ON DELETE CASCADE,
                            descricao TEXTj  
                        );
                        """;

                    String createGenre = """
                        CREATE TABLE IF NOT EXISTS genres (
                            genre_id SERIAL PRIMARY KEY,
                            nome VARCHAR(100) NOT NULL UNIQUE,
                            descricao TEXT 
                        );
                        """;

                    String createAlbum = """
                        CREATE TABLE IF NOT EXISTS albums (
                            album_id SERIAL PRIMARY KEY,
                            artist_id VARCHAR(100) REFERENCES artist(artist_id) ON DELETE CASCADE,
                            genre_id INT REFERENCES genre(genre_id),
                            titulo VARCHAR(100) NOT NULL,
                            data_lancamento DATE,
                            capa_url TEXT
                        );
                        """;

                    String createTrack = """
                        CREATE TABLE IF NOT EXISTS tracks (
                            track_id SERIAL PRIMARY KEY,
                            album_id INT REFERENCES album(album_id) ON DELETE CASCADE,
                            titulo VARCHAR(100) NOT NULL,
                            duracao_segundos INT, 
                            numero_faixa INT 
                        );
                        """;

                    stmt.execute(createUsers);
                    stmt.execute(createArtist);
                    stmt.execute(createGenre);
                    stmt.execute(createAlbum);
                    stmt.execute(createTrack);

          }       
      }
    }catch (Exception e){
      System.out.println("Erro ao conectar-se ao PostgreSQL");
      e.printStackTrace();
    }
  }
}

