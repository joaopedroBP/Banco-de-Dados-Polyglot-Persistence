/*
 *Compile com: javac -cp "postgresql-42.7.3.jar" Main.java
 *Rode com java -cp ".:postgresql-42.7.3.jar" Main
 * SOMENTE LINUX
 */

import java.sql.Connection;
import java.sql.DriverManager;

public class Main{
  public static void main(String[] args){
    String url = "jdbc:postgresql://localhost:5432/postgres";
    String usuario = "postgres";
    String senha = "postgres";

    try (Connection con = DriverManager.getConnection(url,usuario,senha)){
      if(con != null){
        System.out.println("Conectado ao PostgreSQL");
      }
    }catch (Exception e){
      System.out.println("Erro ao conectar-se ao PostgreSQL");
      e.printStackTrace();
    }
  }
}
