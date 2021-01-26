use crate::uconfig;

use std::io;
use std::str;
use std::process::{Command, Stdio};
use log::*;

//Execute kill and stop command for containers.
fn ctr_cmd(runtime: &str, id: &str, cmd :&str) -> Result <(), io::Error> {
    let output = match Command::new(runtime)
        .arg(cmd)
        .arg(id)
        .stdout(Stdio::piped())
        .output()
    {
        Ok(output) => output,
        Err(err) => {
            error!("Error while {} container {} : {} ", cmd, id, err);
            return Err(err);
        }
    };
    debug!("Status {} container {} :: {} ",cmd, id, output.status);
    
    Ok(())
}

//Clean containers
fn ctr_clean(runtime: &str) -> Result<i32, io::Error > {
    
    //Read container list and status
    let output = Command::new(runtime).arg("list").output()?;
    if !output.status.success() {
        error!("Command executed with failing error code");
        return Ok(-1);
    }
    trace!("Container list is :: {:?}", output.stdout);
    
    //stdout => str slice
    let listdata = match str::from_utf8(&output.stdout) {
        Ok(v) => v,
        Err(err) => {
            error!("Inavlid UTF-8 sequrnce read: {}", err);
            return Ok(-1);
        }
    };
    trace!("Container List Data is :: {:?}", listdata);

    //Parsing str slice    
    let ctrdata: Vec<&str> = listdata.split("\n").collect();
    trace!("Container Data is :: {:?}", ctrdata);
    
    //Parsing each record
    for fields in ctrdata {
        let field: Vec<&str> = fields.split_whitespace().collect();
        trace!("Container data fields are :: {:?}", field);
        if field.len() > 0 {
            match field[2] {
                "stopped" => {
                    // delete container if it's stopped
                    ctr_cmd( runtime, field[0], "delete")?;
                },
                 "running" => {
                    // kill contaier if it's running
                    if let Ok(_) = ctr_cmd( runtime, field[0], "kill") {
                        ctr_cmd( runtime, field[0], "delete")?;
                    }
                 },
                 _ => continue,
            }
        }
    }
    
    Ok(0)
}

//Start service
pub fn start_cmd(ctrname: &str) {
    debug!("Dummy {} container start.", ctrname);
}

//Stop service
pub fn stop_cmd(ctrname: &str) {
    debug!("Dummy {} container stop.", ctrname);
}

//Restart service
pub fn restart_cmd(ctrname: &str) {
    debug!("Dummy {} container restart.", ctrname);
}

//Clean containers
pub fn clean_cmd() {
    debug!("Cleaning conatiners.");
    
     // Read OCI Config
     let init = uconfig::Config::get_oci_runtime();
     let init = match init {
         Some(init) => init,
         None => {
             panic!("No OCI runtime config provided in config.");
         }
     };

    //Read path if provided otherwise default
     let oci_path = match &init.path {
        Some(path) => path,
        None => "/usr/bin/crun",
    };
    debug!("OCI runtime :: {}", oci_path);
    
    //Clean containers
    match ctr_clean(&oci_path) {
        Ok(_) => info!("Cleaning completetd"),
        Err(err) => error!("Cleaning completetd with error:: {} ", err), 
    };
}

// Restart microCE
pub fn restart_u_ce() {
    debug!("Dummy uContainer Engine restart.");
}
