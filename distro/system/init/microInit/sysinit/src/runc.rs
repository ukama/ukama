use crate::uconfig;

use prctl;
use log::*;
use wait_timeout::ChildExt;
use std::{ fs, io, path::Path,process::Command,
           process::Stdio, time::Duration };


//Wiat for child process.
fn wait_for_child(mut child: std::process::Child) -> Result<Option<i32>, io::Error> {
    //Waited time out as sometime your conatiner may get stuck
    let one_sec = Duration::from_secs(10);
    let status_code = match child.wait_timeout(one_sec).unwrap() {
        Some(status) => status.code(),
        None => {
            // child hasn't exited yet
            error!("Having issues with container.");
            child.kill().unwrap();
            child.wait().unwrap().code()
        }
    };

    Ok(status_code)
}

pub fn runc_init(onboot: &Vec<uconfig::AppConfig>) -> Result<i32, io::Error> {
    let runc_bin = "/usr/bin/runc";
    let logdir = "/run/log/onboot";
    let varlog = "/var/log/onboot";

    //Set Subreaper
    let _ = prctl::set_child_subreaper(true).unwrap();

    //Make log dir
    match fs::create_dir_all(logdir) {
        Ok(_) => debug!("PID directory {} set for onboot containers.", logdir),
        Err(err) => {
            error!(
                "PID directory {} couldn't be created for onboot containers.",
                err
            );
            return Err(err);
        }
    }

    match fs::create_dir_all(varlog) {
        Ok(_) => debug!("Log directory {} set for onboot containers.", varlog),
        Err(err) => {
            error!(
                "Log directory {} couldn't be created for onboot containers.",
                err
            );
            return Err(err);
        }
    }

    let mut status: i32 = 0;

    let onboot_iter = onboot.iter();
    for ctr in onboot_iter {
        debug!("Preparing setup for onboot container {:?} ", ctr);
        let lpath = match ctr.path {
            Some(ref path) => path,
            None => {
                println!("Missing {} container path in uConfig.ml.", ctr.name);
                continue;
            }
        };

        // Check if container directory exist or not.
        if !Path::new(lpath).exists() {
            error!("Container {} rootfs missing.", ctr.name);
            continue;
        }

        //pid file
        let pidfile = format!("{}{}{}", logdir, "/", ctr.name);
        
        //Not taking care of stdio and err right now
        //TODO: Error and logging
        let child = match Command::new(runc_bin)
            .arg("create")
            .arg("--bundle")
            .arg(lpath)
            .arg("--pid-file")
            .arg(&pidfile)
            .arg(&ctr.name)
            .stdout(Stdio::piped())
            .spawn()
        {
            Ok(output) => output,
            Err(err) => {
                error!("Error while creating container {} : {} ", ctr.name, err);
                status = 1;
                continue;
            }
        };

        trace!(
            "{} create --bundle {} --pidifle {} {}",
            runc_bin, lpath, pidfile, ctr.name
        );

        match wait_for_child(child) {
            Ok(val) => match val {
                Some(ret) => info!("Container {} created successfully ret {}.", ctr.name, ret),
                None => error!("Having issues while conatiner {} creation.", ctr.name),
            },
            Err(err) => error!(
                "Error while check container {} creation  process : {}",
                ctr.name, err
            ),
        }

        //Read pid file for the container
        let contents = match fs::read_to_string(&pidfile) {
            Ok(str) => str,
            Err(err) => {
                error!("Error reading pid file for container {}: {}", ctr.name, err);
                status = 1;
                continue;
            }
        };

        //convert to integer
        let _pid: i32 = match contents.parse() {
            Ok(val) => val,
            Err(err) => {
                error!("Error reading pid file for {} : {}", ctr.name, err);
                status = 1;
                continue;
            }
        };

        //Start container
        let schild = match Command::new(runc_bin)
            .arg("start")
            .arg(&ctr.name)
            .stdout(Stdio::piped())
            .spawn()
        {
            Ok(output) => output,
            Err(err) => {
                error!("Error while starting container {} : {} ", ctr.name, err);
                status = 1;
                continue;
            }
        };

        trace!("{} start {}", runc_bin, ctr.name);

        match wait_for_child(schild) {
            Ok(val) => match val {
                Some(ret) => debug!("Container {} created successfully ret {}.", ctr.name, ret),
                None => error!("Having issues while conatiner {} creation.", ctr.name),
            },
            Err(err) => error!(
                "Error while check container {} creation  process : {}",
                ctr.name, err
            ),
        }

        //Cleaning pid file
        match fs::remove_file(&pidfile) {
            Ok(_) => debug!(
                "PID file {} cleaned for onboot container {}.",
                pidfile, ctr.name
            ),
            Err(err) => {
                error!(
                    "PID file {} couldn't be cleaned for onboot container {} : {}",
                    pidfile, ctr.name, err
                );
            }
        }
    }

    Ok(status)
}

pub fn onboot() -> i32 {
    let onboot = uconfig::Config::get_onboot_config();
    trace!("Onboot Container list {:?}", onboot);
    //println!("{:?}", uconfig::Config::get_onboot_config());
    let _ = runc_init(onboot).unwrap();
    0
}

pub fn onshutdown() -> i32 {
    trace!("{:?}", uconfig::Config::get_onshutdown_config());
    0
}
