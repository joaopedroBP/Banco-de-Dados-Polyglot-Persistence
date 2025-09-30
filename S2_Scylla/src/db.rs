use scylla::client::session::Session;
use scylla::client::session_builder::SessionBuilder;

use crate::Result;

pub async fn create_session(uri: &str) -> Result<Session> {
    SessionBuilder::new()
        .known_node(uri)
        .build()
        .await
        .map_err(From::from)
}
