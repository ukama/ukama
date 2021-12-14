use crate::uconfig;
use std::fs;
use std::process::{Command, Stdio};
use procfs::process::Process;
use log::*;
use std::{thread, time};

//Read microCE config
fn read_micro_ce_config () -> Option<&'static uconfig::AppConfig> {
    let initcfg = uconfig::Config::get_init_config();
    trace!("microCE init:: {:?}", initcfg);

    //Capture micro_ce config
    let mut cecfg: Option<&uconfig::AppConfig> = None;
    for initapp in initcfg {
        if initapp.name == "microCE" {
            cecfg = Some(initapp);
            break;
        }
    }
    cecfg
}

// Read process ID for microCE
fn read_pid(pidfile: &str) -> Option<i32> {
    let mut pid : Option<i32> = None;
    //Read pid file for the container
    let contents = match fs::read_to_string(pidfile) {
        Ok(str) => str,
        Err(err) => {
            error!("Error reading pid file {}: {}", pidfile, err);
            return pid;
        }
    };

    //convert to integer
    let pidval: i32 = match contents.parse() {
        Ok(val) => val,
        Err(err) => {
            error!("Error reading pid file {} : {}", pidfile, err);
            return pid;
        }
    };

    pid = Some(pidval);
    pid
}

// Starting microCE
pub fn init() {
    
    let pidfile ="/var/log/microCE.pid";
    let ucfg = match read_micro_ce_config() {
        Some(cfg) => cfg,
        None => {
            error!("No init config for microCE found.");
            return;
        },
    };
    debug!("microCE config:: {:?}", ucfg);
    
    let micro_ce_path = match &ucfg.path {
        None => "/usr/bin/microCE.d".to_string(),
        Some(path) => path.to_string(),
    };

    //starting micro_ce
    let _child = match Command::new(&micro_ce_path)
        .arg("--pid-file")
        .arg(&pidfile)
        .stdout(Stdio::piped())
        .spawn()
    {
        Ok(output) => output,
        Err(err) => {
            error!("Error while starting {} : {} ", micro_ce_path, err);
            return;
        }
    };
    
    //Waiting for the pid fo the micro_ce process.
    for wait in 1..3 {
        let sec = time::Duration::from_millis(1000);
        //let now = time::Instant::now();
        thread::sleep(5*sec);
        let pid = match read_pid(&pidfile) {
        Some(p) => p,
        None => 0,
        };
        
        if pid == 0 {
            debug!("No pid found for microCE process after {} seconds.", wait*5 );
        } else {
            //Check process status
            let proc = Process::new(pid).expect("Failed to read micro_ce status");
            if proc.is_alive() {
                info!("microCE engine is running with process is {}.", pid);
            }
            break;
        }
    };

    debug!("microCE initialization complete");
}
