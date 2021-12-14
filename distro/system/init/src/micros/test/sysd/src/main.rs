
use log::*;
use simplelog::*;
use std::{env, fs::File, process};
use std::io::prelude::*;
use std::{thread, time};
use procfs::process::Process;

fn write_pid(path: &str, pid: &i32) -> std::io::Result<()> {
        let mut file = File::create(path)?;
        file.write(pid.to_string().as_bytes())?;
        Ok(())
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
        
        //Argument parsing
        let args: Vec<String> = env::args().collect();
        trace!("{} Arguments to sysd service {:?}", args.len(), args);
        
        if args.len() < 3 {
            error!("Usage : sysd --pid-file /var/log/microCE.pid");
            process::exit(0);
        }

        let proc = Process::myself().expect("Failure to get process");
        let pid = proc.pid();
        info!("Pid for {} is {}", args[0], pid);
        
        //Write a pid file
        match write_pid(&args[2], &pid) {
            Ok(_) => debug!("PID {} Write successfull to {} for process {}.", pid, args[2], args[0]),
            Err(err) => error!("Write failed for pid {} ti file {} for process {} : {}.", pid, args[2], args[0], err),
        }

        loop {
            let ten_millis = time::Duration::from_millis(1000);
            //let now = time::Instant::now();
            thread::sleep(ten_millis);
        }
}
