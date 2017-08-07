extern crate prometheus;
extern crate rocket;
#[macro_use]
extern crate serde_derive;
extern crate toml;
#[macro_use]
extern crate clap;

use std::fs::File;
use std::io::{Read, Cursor};

use clap::{App, Arg};
use prometheus::{Opts, Gauge, Encoder};
use rocket::State;
use rocket::handler::Outcome;
use rocket::http::Method;
use rocket::{Request, Route};

#[derive(Deserialize)]
struct Server {
    name: String,
    host: String,
    ports: Option<Vec<u16>>,
    tls_ports: Option<Vec<u16>>,
}

#[derive(Deserialize)]
struct Config {
    servers: Vec<Server>,
}

fn main() {
    let app = App::new(crate_name!())
        .author(crate_authors!())
        .version(crate_version!())
        .arg(Arg::with_name("config")
            .help("config file path")
            .takes_value(true)
            .short("c")
            .required(true)
            .long("config"))
        .get_matches();

    let config = app.value_of("config").unwrap();

    let file = File::open(config).expect("could not open config file");

    let contents =
        file.bytes().collect::<Result<Vec<u8>, std::io::Error>>().expect("could not read file");

    let cfg = match toml::from_str::<Config>(&String::from_utf8(contents).unwrap()) {
        Ok(c) => c,
        Err(e) => panic!("could not parse toml: {}", e),
    };

    let metrics_route = Route::new(Method::Get, "/metrics", metrics);

    rocket::ignite()
        .manage(cfg)
        .mount("/", vec![metrics_route])
        .launch();
}

fn metrics<'r>(request: &'r Request, _: rocket::Data) -> rocket::handler::Outcome<'r> {
    let cfg = request.guard::<State<Config>>().unwrap().inner();

    let r = prometheus::Registry::new();

    for server in &cfg.servers {
        let up = Gauge::with_opts(Opts::new(format!("{}_up", server.name),
                                            "server is up".to_string()))
            .unwrap();
        // TODO: add actual metrics here
        up.inc();
        r.register(Box::new(up)).unwrap();
    }

    let encoder = prometheus::TextEncoder::new();
    let mut resp = rocket::response::Response::build();
    let mut metric_body = vec![];
    encoder.encode(&r.gather(), &mut metric_body).unwrap();
    resp.streamed_body(Cursor::new(metric_body));

    Outcome::from(request, resp.finalize())
}
