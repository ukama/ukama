mod cmd;
mod runc;
mod micro_ce;
mod uconfig;

#[macro_use]
extern crate lazy_static;
use log::*;
use simplelog::*;
use std::{env, fs::File, process};

fn usage(prg: &str) {
    println!("USAGE: {} [options] COMMAND", prg);
    println!("Commands:");
    println!("  uInit Prepare the system at start of day");
    println!("  stop        Stop a service");
    println!("  start       Start a service");
    println!("  restart     Restart a service");
    println!("  help        Print this message");
    println!(
        "Run '{} COMMAND --help' for more information on the command\n",
        prg
    );
    println!("Options:");
}

fn main() {
    //Logger
    CombinedLogger::init(vec![
        #[cfg(feature = "termcolor")]
        TermLogger::new(LevelFilter::Trace, Config::default(), TerminalMode::Mixed),
        #[cfg(not(feature = "termcolor"))]
        SimpleLogger::new(LevelFilter::Trace, Config::default()),
        WriteLogger::new(
            LevelFilter::Debug,
            Config::default(),
            File::create("/var/log/init.log").unwrap(),
        ),
    ])
    .unwrap();

    trace!("Starting rc.init service");
    let args: Vec<String> = env::args().collect();
    trace!("{} Arguments to rc.init service {:?}", args.len(), args);

    //parse uconfig
    let uconfig = uconfig::Config::get();
    trace!("{:#?}", uconfig);

    if args.len() <= 1 {
        if args[0].contains("onboot") {
            process::exit(runc::onboot());
        } else if args[0].contains("onshutdown") {
            process::exit(runc::onshutdown());
        } else if args[0].contains("microCE") {
            micro_ce::init();
            process::exit(0);
        } else {
            usage(&args[0]);
            process::exit(1);
        }
    } else {
        match args[1].as_str() {
            "stop" => cmd::stop_cmd(&args[2]),
            "start" => cmd::start_cmd(&args[2]),
            "restart" => cmd::restart_cmd(&args[2]),
            "uInit" => cmd::restart_u_ce(),
            "clean" => cmd::clean_cmd(),
            _ => usage(&args[0]),
        }
    }
}
