/*
 *Compile com: javac -cp "postgresql-42.7.3.jar" Main.java
 *Rode com java -cp ".:postgresql-42.7.3.jar" Main
*/

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.Statement;
import java.sql.PreparedStatement;  
import java.sql.SQLException;     


public class Main{
    public static String createUser(Connection con, String userId, String email, String senha) {
        String sql = "INSERT INTO users (user_id, email, senha) VALUES (?, ?, ?);";
        try (PreparedStatement pstmt = con.prepareStatement(sql)) {
            pstmt.setString(1, userId);
            pstmt.setString(2, email);
            pstmt.setString(3, senha);
            pstmt.executeUpdate();
            return userId;
        } catch (SQLException e) {
            System.out.println("Erro ao criar usuário '" + userId + "':"); 
            e.printStackTrace();
            return "";
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

                    String createArtists = """
                        CREATE TABLE IF NOT EXISTS artists (
                            artist_id VARCHAR(100) PRIMARY KEY,  -- nome do artista
                            user_id VARCHAR(100) REFERENCES users(user_id) ON DELETE CASCADE,
                            descricao TEXT
                        );
                        """;

                    String createGenres = """
                        CREATE TABLE IF NOT EXISTS genres (
                            genre_id SERIAL PRIMARY KEY,
                            nome VARCHAR(100) NOT NULL UNIQUE,
                            descricao TEXT 
                        );
                        """;

                    String createAlbums = """
                        CREATE TABLE IF NOT EXISTS albums (
                            album_id SERIAL PRIMARY KEY,
                            artist_id VARCHAR(100) REFERENCES artists(artist_id) ON DELETE CASCADE,
                            genre_id INT REFERENCES genres(genre_id),
                            titulo VARCHAR(100) NOT NULL,
                            data_lancamento DATE,
                            capa_url TEXT
                        );
                        """;

                    String createTracks = """
                        CREATE TABLE IF NOT EXISTS tracks (
                            track_id SERIAL PRIMARY KEY,
                            album_id INT REFERENCES albums(album_id) ON DELETE CASCADE,
                            titulo VARCHAR(100) NOT NULL,
                            duracao_segundos INT, 
                            numero_faixa INT 
                        );
                        """;

                    stmt.execute(createUsers);
                    stmt.execute(createArtists);
                    stmt.execute(createGenres);
                    stmt.execute(createAlbums);
                    stmt.execute(createTracks);

          }       
      }
    }catch (Exception e){
      System.out.println("Erro ao conectar-se ao PostgreSQL");
      e.printStackTrace();
    }
  }
}

