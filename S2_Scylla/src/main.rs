use crate::result::Result;
mod result;
#[tokio::main]

async fn main() -> Result<()> {
    println!("connecting to db");
}
