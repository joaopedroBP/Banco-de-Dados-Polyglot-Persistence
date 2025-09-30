use crate::result::Result;
mod db;
mod result;
#[tokio::main]
async fn main() -> Result<()> {
    println!("connecting to db");
    let uri = std::env::var("SCYLLA_URI").unwrap_or_else(|_| "127.0.0.1:9042".to_string());
    let _session = db::create_session(&uri).await?;

    println!("Conectado no scylla em {}", uri);
    Ok(())
}
