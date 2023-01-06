pub mod auth;

use actix_web::{web, App, HttpServer};

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        App::new()
            .service(web::scope("/auth").configure(auth::config))
    })
    .bind(("127.0.0.1", 8080))?
    .run()
    .await
}