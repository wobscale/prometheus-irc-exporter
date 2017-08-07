extern crate prometheus;
extern crate rocket;
extern crate url;
extern crate irc;
extern crate rand;

use std::io::Cursor;
use std::default::Default;

use prometheus::{Opts, Gauge, Encoder};
use rocket::handler::Outcome;
use rocket::http::{Status, Method};
use rocket::{Request, Route};

fn main() {
    let metrics_route = Route::new(Method::Get, "/metrics", metrics);

    rocket::ignite().mount("/", vec![metrics_route]).launch();
}

fn metrics<'r>(request: &'r Request, _: rocket::Data) -> rocket::handler::Outcome<'r> {
    let query = request.uri().query().unwrap();
    let target = rocket::request::FormItems::from(query)
        .find(|&el| el.0 == "target")
        .map(|s| s.1.url_decode());

    let target = match target {
        Some(t) => t,
        None => {
            return Outcome::failure(Status::new(400, "no target"));
        }
    };

    let target = match target {
        Ok(t) => t,
        Err(_) => {
            return Outcome::failure(Status::new(400, "invalid target"));
        }
    };

    let irc_url = match url::Url::parse(&target) {
        Ok(i) => i,
        Err(_) => {
            return Outcome::failure(Status::new(400, "target must be a valid irc uri"));
        }
    };

    if irc_url.scheme() != "irc" && irc_url.scheme() != "ircs" {
        return Outcome::failure(Status::new(400, "target must be a valid irc uri"));
    }

    let r = prometheus::Registry::new();

    for metrics in scrape_irc_server(irc_url) {
        r.register(metrics).unwrap();
    }

    let encoder = prometheus::TextEncoder::new();
    let mut resp = rocket::response::Response::build();
    let mut metric_body = vec![];
    encoder.encode(&r.gather(), &mut metric_body).unwrap();
    resp.streamed_body(Cursor::new(metric_body));

    Outcome::from(request, resp.finalize())
}

fn scrape_irc_server(target: url::Url) -> Vec<Box<prometheus::Collector>> {
    let nick = format!("promirc_{}", rand::random::<i32>());

    let up = Gauge::with_opts(Opts::new("up", "server is up")).unwrap();
    let tls_expiration = Gauge::with_opts(Opts::new("tls_expiration", "server tls expiration"))
        .unwrap();
    let global_users = Gauge::with_opts(Opts::new("global_users", "number of global users"))
        .unwrap();
    let local_users = Gauge::with_opts(Opts::new("local_users", "number of local users")).unwrap();
    let uptime = Gauge::with_opts(Opts::new("uptime", "server uptime")).unwrap();

    let ssl = target.scheme() == "ircs";
    let port = target.port().unwrap_or(if ssl { 6697 } else { 6667 });

    let cfg = irc::client::prelude::Config{
        nickname: Some(nick),
        server: target.host().map(|h| h.to_string()),
        port: Some(port),
        use_ssl: Some(target.scheme() == "ircs"),
        .. Default::default()
    };

    let server = match irc::client::prelude::IrcServer::from_config(cfg) {
        Ok(s) => {
            up.inc();
            s
        },
        Err(s) => {
            return vec![Box::new(up)];
        },
    };

    vec![Box::new(up)]
}
