use actix_web::{App, HttpResponse, HttpServer, post, Responder, web};
use sea_orm::{DatabaseConnection, Database};
use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize)]
struct LoginRequest {

    #[serde(default = "default_username")]
    username: String,

    #[serde(default)]
    password: String,
}

pub fn default_username() -> String {
    "default".to_string()
}

#[post("/login")]
async fn test(req: web::Json<LoginRequest>, data: web::Data<AppState>) -> impl Responder {
    HttpResponse::Ok().body(req.username.to_owned())
}

// Database configuration
const DATABASE_URL: &str = "postgres://postgres:deinemutter123@localhost/rust_test";

#[derive(Debug, Clone)]
struct AppState {
    conn: DatabaseConnection
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {

    println!("Starting connection to database");
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