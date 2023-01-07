use actix_web::{App, HttpServer, Responder, get, web};
use sea_orm::{DatabaseConnection, Database};

#[get("/test")]
async fn test(data: web::Data<AppState>) -> impl Responder {
    "Hello world!"
}

// Database config
const DATABASE_URL: &str = "postgres://postgres:deinemutter123@localhost/learn";

#[derive(Debug, Clone)]
struct AppState {
    conn: DatabaseConnection,
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {

    let conn = Database::connect(DATABASE_URL).await.unwrap();

    println!("Connected to database: {}", DATABASE_URL);

    let state = AppState { conn };

    HttpServer::new(move || {
        App::new()
            .app_data(web::Data::new(state.clone()))
            .service(test)
    })
    .bind(("127.0.0.1", 8080))?
    .run()
    .await
}