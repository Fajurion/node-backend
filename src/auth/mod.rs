mod auth;

use actix_web::web;

pub fn config(cfg: &mut web::ServiceConfig) {
    cfg.service(auth::login)
    .service(auth::register);
}