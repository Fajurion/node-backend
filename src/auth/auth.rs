use actix_web::{post, Responder, HttpResponse, web, HttpRequest};
use serde::Deserialize;

#[derive(Deserialize)]
struct LoginRequest {
    #[serde(default)]
    username: String,

    #[serde(default)]
    password: String,
}

#[post("/login")]
async fn login(req: web::Json<LoginRequest>) -> impl Responder {

    HttpResponse::Ok().body(format!("username: {}, password: {}", req.username, req.password))
}

#[derive(Deserialize)]
struct RegisterRequest {
    #[serde(default)]
    username: String,

    #[serde(default)]
    email: String,

    #[serde(default)]
    password: String,

    #[serde(default)]
    confirm_password: String,
}

#[post("/register")]
async fn register(req: web::Json<RegisterRequest>) -> impl Responder {
    HttpResponse::Ok().body("Register page")
}