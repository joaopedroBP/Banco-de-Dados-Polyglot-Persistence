/*
 *Compile com: javac -cp "postgresql-42.7.3.jar" Main.java
 *Rode com java -cp ".:postgresql-42.7.3.jar" Main
*/


import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.Statement;
import java.sql.PreparedStatement;  
import java.sql.SQLException;     
import java.sql.ResultSet;
import java.sql.Timestamp;

public class Main{

  // FUNÇÕES DO USUÁRIO
  public static String createUser(Connection con, String userId, String email, String senha) {
    String sql = "INSERT INTO users (user_name, email, senha) VALUES (?, ?, ?);";
    try (PreparedStatement pstmt = con.prepareStatement(sql)) {
      pstmt.setString(1, userId);
      pstmt.setString(2, email);
      pstmt.setString(3, senha);
      pstmt.executeUpdate();
      return userId;
    } catch (SQLException e) {
      System.out.println("Erro ao criar usuário '" + userId + "'!"); 
      e.printStackTrace();
      return "";
    }
  }

  public static void listUsers(Connection con) {
    String sql = "SELECT user_name, email, data_criacao FROM users;";
    try (PreparedStatement pstmt = con.prepareStatement(sql);
         ResultSet rs = pstmt.executeQuery()) {

      System.out.println("Lista de usuários:");
      while (rs.next()) {
        String userName = rs.getString("user_name"); 
        String email = rs.getString("email");
        Timestamp dataCriacao = rs.getTimestamp("data_criacao");

        System.out.println("User: " + userName);
        System.out.println("Email: " + email);
        System.out.println("Data de criação: " + dataCriacao);
        System.out.println("---------------------------");
      }

    } catch (SQLException e) {
      System.out.println("Erro ao listar usuarios!");
      e.printStackTrace();
    }
  }

  public static boolean checkUser(Connection con, String userName) {
    String sql = "SELECT user_name FROM users WHERE user_name = ?;";
    try (PreparedStatement pstmt = con.prepareStatement(sql)) {
      pstmt.setString(1, userName); 
      try (ResultSet rs = pstmt.executeQuery()) {
        return rs.next(); 
      }
    } catch (SQLException e) {
      System.out.println("Erro ao verificar usuário: " + userName + "!");
      e.printStackTrace();
      return false;
    }
  }

  public static void removeUser(Connection con, String userName){
    if(!checkUser(con,userName)){
      System.out.println("Usuário: " + userName + " não existe!");
      return;
    }

    String sql = "DELETE FROM users WHERE user_name = ?";
    try (PreparedStatement pstmt = con.prepareStatement(sql)){
      pstmt.setString(1,userName);
      int mudadas = pstmt.executeUpdate();
      if(mudadas > 0){ 
        System.out.println("Usuário: " + userName + " Removido !");
      }
    }catch (SQLException e){
      System.out.println("Erro ao remover usuário: " + userName + "!");
      e.printStackTrace();
    }
  }

  // FUNÇÕES DO ARTISTA
  public static String createArtist(Connection con, String artistName, String desc) {
    String sql = "INSERT INTO artists (artist_name, descricao) VALUES (?, ?);";
    try (PreparedStatement pstmt = con.prepareStatement(sql)) {
      pstmt.setString(1, artistName);
      pstmt.setString(2, desc);
      pstmt.executeUpdate();
      return artistName;
    } catch (SQLException e) {
      System.out.println("Erro ao criar artista '" + artistName + "'!"); 
      e.printStackTrace();
      return "";
    }
  }

  public static void listArtists(Connection con) {
    String sql = "SELECT artist_id, artist_name, descricao FROM artists;";
    try (PreparedStatement pstmt = con.prepareStatement(sql);
         ResultSet rs = pstmt.executeQuery()) {

      System.out.println("Artist list:");
      while (rs.next()) {
        int artistId = rs.getInt("artist_id");
        String artName = rs.getString("artist_name");
        String desc = rs.getString("descricao");
        System.out.println("ID: " + artistId);
        System.out.println("Artist: " + artName);
        System.out.println("Description: " + desc);
        System.out.println("---------------------------");
      }

    } catch (SQLException e) {
      System.out.println("Erro ao listar artistas!");
      e.printStackTrace();
    }
  }

  public static boolean checkArtist(Connection con, String artistId) {
    String sql = "SELECT artist_id FROM artists WHERE artist_id = ?;";
    try (PreparedStatement pstmt = con.prepareStatement(sql)) {
      pstmt.setInt(1, Integer.parseInt(artistId));
      try (ResultSet rs = pstmt.executeQuery()) {
        return rs.next();
      }

    } catch (SQLException e) {
      System.out.println("Erro ao verificar artista!");
      e.printStackTrace();
      return false;
    }
  }

  public static void removeArtist(Connection con, String artistID){
    if(!checkArtist(con,artistID)){
      System.out.println("Artista não existe!");
      return;
    }

    String sql = "DELETE FROM artists WHERE artist_id = ?";
    try (PreparedStatement pstmt = con.prepareStatement(sql)){
      pstmt.setInt(1,Integer.parseInt(artistID));
      int mudadas = pstmt.executeUpdate();
      if(mudadas > 0){ 
        System.out.println("Artista removido!");
      }
    }catch (SQLException e){
      System.out.println("Erro ao remover artista!");
      e.printStackTrace();
    }
  }

  // FUNÇÕES DO GÊNERO
  public static String createGenre(Connection con, String genreName) {
    String sql = "INSERT INTO genres (genre_name) VALUES (?);";
    try (PreparedStatement pstmt = con.prepareStatement(sql)) {
      pstmt.setString(1, genreName);
      pstmt.executeUpdate();
      return genreName;
    } catch (SQLException e) {
      System.out.println("Erro ao criar gênero '" + genreName + "'!"); 
      e.printStackTrace();
      return "";
    }
  }

  public static void listGenres(Connection con) {
    String sql = "SELECT genre_id, genre_name FROM genres;";
    try (PreparedStatement pstmt = con.prepareStatement(sql);
         ResultSet rs = pstmt.executeQuery()) {

      System.out.println("Genre list:");
      while (rs.next()) {
        int genreId = rs.getInt("genre_id");
        String genreName = rs.getString("genre_name");

        System.out.println("ID: " + genreId);
        System.out.println("Genre: " + genreName);
        System.out.println("---------------------------");
      }

    } catch (SQLException e) {
      System.out.println("Erro ao listar gêneros!");
      e.printStackTrace();
    }
  }

  public static boolean checkGenre(Connection con, String genreId) {
    String sql = "SELECT genre_id FROM genres WHERE genre_id = ?;";
    try (PreparedStatement pstmt = con.prepareStatement(sql)) {
      pstmt.setInt(1, Integer.parseInt(genreId));
      try (ResultSet rs = pstmt.executeQuery()) {
        return rs.next(); 
      }

    } catch (SQLException e) {
      System.out.println("Erro ao verificar gênero!");
      e.printStackTrace();
      return false;
    }
  }

  public static void removeGenre(Connection con, String genreID){
    if(!checkGenre(con,genreID)){
      System.out.println("Gênero não existe!");
      return;
    }

    String sql = "DELETE FROM genres WHERE genre_id = ?";
    try (PreparedStatement pstmt = con.prepareStatement(sql)){
      pstmt.setInt(1,Integer.parseInt(genreID));
      int mudadas = pstmt.executeUpdate();
      if(mudadas > 0){ 
        System.out.println("Gênero removido!");
      }
    }catch (SQLException e){
      System.out.println("Erro ao remover gênero!");
      e.printStackTrace();
    }
  }

  // FUNÇÕES DO ALBUM E DA TRACK
  public static String createAlbum(Connection con, String albumName, String artistId, String releaseDate, String genreId) {
    if (!checkArtist(con, artistId)) {
      System.out.println("Erro: artista com ID " + artistId + " não existe!");
      return "";
    }

    if (!checkGenre(con, genreId)) {
      System.out.println("Erro: gênero com ID " + genreId + " não existe!");
      return "";
    }

    String sql = "INSERT INTO albums (artist_id, genre_id, titulo, data_lancamento) VALUES (?, ?, ?, ?);";
    try (PreparedStatement pstmt = con.prepareStatement(sql)) {
      pstmt.setInt(1, Integer.parseInt(artistId));
      pstmt.setInt(2, Integer.parseInt(genreId));
      pstmt.setString(3, albumName);
      pstmt.setDate(4, java.sql.Date.valueOf(releaseDate));
      pstmt.executeUpdate();
      System.out.println("Álbum '" + albumName + "' criado com sucesso!");
      return albumName;
    } catch (SQLException e) {
      System.out.println("Erro ao criar álbum '" + albumName + "!");
      e.printStackTrace();
      return "";
    }
  }

  public static boolean checkAlbum(Connection con, String albumId) {
    String sql = "SELECT album_id FROM albums WHERE album_id = ?;";
    try (PreparedStatement pstmt = con.prepareStatement(sql)) {
      pstmt.setInt(1, Integer.parseInt(albumId));
      try (ResultSet rs = pstmt.executeQuery()) {
        return rs.next();
      }
    } catch (SQLException e) {
      System.out.println("Erro ao verificar álbum");
      e.printStackTrace();
      return false;
    }
  }

  public static void removeAlbum(Connection con, String albumId){
    if(!checkAlbum(con, albumId)){
      System.out.println("Álbum não encontrado!");
      return;
    }

    String sql = "DELETE FROM albums WHERE album_id = ?";
    try (PreparedStatement pstmt = con.prepareStatement(sql)){
      pstmt.setInt(1, Integer.parseInt(albumId));
      int mudadas = pstmt.executeUpdate();
      if(mudadas > 0){
        System.out.println("Álbum removido!");
      }
    } catch (SQLException e){
      System.out.println("Erro ao remover álbum!");
      e.printStackTrace();
    }
  }

  public static boolean checkTrackOrder(Connection con, String albumId, String albumNum) {
    String sql = "SELECT 1 FROM tracks WHERE album_id = ? AND numero_faixa = ?;";
    try (PreparedStatement pstmt = con.prepareStatement(sql)) {
      pstmt.setInt(1, Integer.parseInt(albumId));
      pstmt.setInt(2, Integer.parseInt(albumNum));
      try (ResultSet rs = pstmt.executeQuery()) {
        return !rs.next();
      }
    } catch (SQLException e) {
      System.out.println("Erro ao verificar ordem da faixa!");
      e.printStackTrace();
      return false;
    }
  }

  public static String createTrack(Connection con, String trackName, String albumId, String duration, String albumNum) {
    if (!checkAlbum(con, albumId) ) {
      System.out.println("Album não existe!");
      return "";
    }
    if (!checkTrackOrder(con, albumId, albumNum)){
      System.out.println("Track nessa posição no album já existe!");
      return "";
    }
    String sql = "INSERT INTO tracks (album_id, titulo, duracao_segundos, numero_faixa) VALUES (?, ?, ?, ?);";
    try (PreparedStatement pstmt = con.prepareStatement(sql)) {
      pstmt.setInt(1, Integer.parseInt(albumId));
      pstmt.setString(2, trackName);
      pstmt.setInt(3, Integer.parseInt(duration));
      pstmt.setInt(4, Integer.parseInt(albumNum));
      pstmt.executeUpdate();
      return trackName;
    } catch (SQLException e) {
      System.out.println("Erro ao criar track!");
      e.printStackTrace();
      return "";
    }
  }

  public static boolean checkTrack(Connection con, String trackId) {
    String sql = "SELECT track_id FROM tracks WHERE track_id = ?;";
    try (PreparedStatement pstmt = con.prepareStatement(sql)) {
      pstmt.setInt(1, Integer.parseInt(trackId));
      try (ResultSet rs = pstmt.executeQuery()) {
        return rs.next();
      }
    } catch (SQLException e) {
      System.out.println("Erro ao verificar track!");
      e.printStackTrace();
      return false;
    }
  }

  public static void removeTrack(Connection con, String trackId){
    if(!checkTrack(con, trackId)){
      System.out.println("Track não encontrada!");
      return;
    }

    String sql = "DELETE FROM tracks WHERE track_id = ?";
    try (PreparedStatement pstmt = con.prepareStatement(sql)){
      pstmt.setInt(1, Integer.parseInt(trackId));
      int mudadas = pstmt.executeUpdate();
      if(mudadas > 0){
        System.out.println("Track removida!");
      }
    } catch (SQLException e){
      System.out.println("Erro ao remover track!");
      e.printStackTrace();
    }
  }

  public static void listAlbumsAndTracks(Connection con) {
    String sql = """
      SELECT 
        albums.titulo AS album_title,
        genres.genre_name AS genre_name,
        tracks.titulo AS track_title,
        tracks.duracao_segundos AS track_duration,
        tracks.numero_faixa AS track_number
      FROM albums
      LEFT JOIN genres ON albums.genre_id = genres.genre_id
      LEFT JOIN tracks ON albums.album_id = tracks.album_id
      ORDER BY albums.album_id, tracks.numero_faixa;
      """;

    try (PreparedStatement pstmt = con.prepareStatement(sql);
         ResultSet rs = pstmt.executeQuery()) {

      String lastAlbum = "";

      while (rs.next()) {
        String albumTitle = rs.getString("album_title");
        String genreName = rs.getString("genre_name");
        String trackTitle = rs.getString("track_title");
        int trackDuration = rs.getInt("track_duration");
        int trackNumber = rs.getInt("track_number");

        if (!albumTitle.equals(lastAlbum)) {
          System.out.println();
          System.out.println("Album: " + albumTitle);
          System.out.println("Gênero: " + genreName);
          System.out.println("Músicas:");
          lastAlbum = albumTitle;
        }

        if (trackTitle != null) {
          System.out.println("  " + trackNumber + ". " + trackTitle + " (" + trackDuration + "s)");
        }
      }

    } catch (SQLException e) {
      System.out.println("Erro ao listar álbuns e faixas");
      e.printStackTrace();
    }
  }

  public static void main(String[] args){
    String url = "jdbc:postgresql://localhost:5432/postgres";
    String usuario = "postgres";
    String senha = "postgres";

    String function_called = "" + args[0] + args[1];

    try (Connection con = DriverManager.getConnection(url,usuario,senha)){
      if(con != null){
        System.out.println("Conectado ao PostgreSQL");
        try (Statement stmt = con.createStatement()) {
          String createUsersTable = """
              CREATE TABLE IF NOT EXISTS users (
                  user_name VARCHAR(100) PRIMARY KEY,
                  email VARCHAR(100) NOT NULL UNIQUE,
                  senha VARCHAR(100) NOT NULL,
                  data_criacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP
              );
              """;

          String createArtistsTable = """
              CREATE TABLE IF NOT EXISTS artists (
                  artist_id SERIAL PRIMARY KEY,  
                  artist_name VARCHAR(100),
                  descricao TEXT
              );
              """;

          String createGenresTable = """
              CREATE TABLE IF NOT EXISTS genres (
                  genre_id SERIAL PRIMARY KEY,
                  genre_name VARCHAR(100) NOT NULL UNIQUE
              );
              """;

          String createAlbumsTable = """
              CREATE TABLE IF NOT EXISTS albums (
                  album_id SERIAL PRIMARY KEY,
                  artist_id INT REFERENCES artists(artist_id) ON DELETE CASCADE,
                  genre_id INT REFERENCES genres(genre_id),
                  titulo VARCHAR(100) NOT NULL,
                  data_lancamento DATE
              );
              """;

          String createTracksTable = """
              CREATE TABLE IF NOT EXISTS tracks (
                  track_id SERIAL PRIMARY KEY,
                  album_id INT REFERENCES albums(album_id) ON DELETE CASCADE,
                  titulo VARCHAR(100) NOT NULL,
                  duracao_segundos INT, 
                  numero_faixa INT
              );
              """;

          stmt.execute(createUsersTable);
          stmt.execute(createArtistsTable);
          stmt.execute(createGenresTable);
          stmt.execute(createAlbumsTable);
          stmt.execute(createTracksTable);
        }

        switch (function_called){
          case "useradd":
            createUser(con,args[2],args[3],args[4]);
            break;
          case "artistadd":
            createArtist(con,args[2],args[3]);
            break;
          case "genreadd":
            createGenre(con,args[2]);
            break;
          case "albumadd":
            createAlbum(con,args[2],args[3],args[4],args[5]);
            break;
          case "trackadd":
            createTrack(con,args[2],args[3],args[4],args[5]);
            break;
          case "userlist":
            listUsers(con);
            break;
          case "artistlist":
            listArtists(con);
            break;
          case "genrelist":
            listGenres(con);
            break;
          case "albumlist":
            listAlbumsAndTracks(con);
            break;
          case  "userrmv":
            removeUser(con,args[2]);
            break;
          case  "artistrmv":
            removeArtist(con,args[2]);
            break;
          case  "genrermv":
            removeGenre(con,args[2]);
            break;
          case  "trackrmv":
            removeTrack(con,args[2]);
            break;
          case  "albumrmv":
            removeAlbum(con,args[2]);
            break;
          case "usercheck":
            boolean check = checkUser(con,args[2]);
            break;
          case "artistcheck":
            boolean check = checkArtist(con,args[2]);
            break;
          case "genrecheck":
            boolean check = checkGenre(con,args[2]);
            break;
          case "trackcheck":
            boolean check = checkTrack(con,args[2]);
            break;
          case "albumcheck":
            boolean check = checkAlbum(con,args[2]);
            break;
        }
      }
    }catch (Exception e){
      System.out.println("Erro ao conectar-se ao PostgreSQL");
      e.printStackTrace();
    }
  }
}
