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

fn ctr_clean() -> Result<i32, io::Error > {
    let runtime = "runc";
    
    //Read container list and status
    let output = Command::new("runc").arg("list").output()?;
    if !output.status.success() {
        error!("Command executed with failing error code");
        return Ok(-1);
    }
    println!("runc list is :: {:?}", output.stdout);
    
    //stdout => str slice
    let listdata = match str::from_utf8(&output.stdout) {
        Ok(v) => v,
        Err(err) => {
            println!("Inavlid UTF-8 sequrnce read: {}", err);
            return Ok(-1);
        }
    };
    trace!("Container List Data is :: {:?}", listdata);

    //Parsing str slice    
    let ctrdata: Vec<&str> = listdata.split("\n").collect();
    trace!("Conatiner Data is :: {:?}", ctrdata);
    //Parsing each record
    for fields in ctrdata {
        let field: Vec<&str> = fields.split_whitespace().collect();
        trace!("Conatiner data fields are :: {:?}", field);
        if field.len() > 0 {
            match field[2] {
                "stopped" => {
                    ctr_cmd( runtime, field[0], "delete")?;
                },
                 "running" => {
                    if let Ok(_) = ctr_cmd( runtime, field[0], "kill") {
                        ctr_cmd( runtime, field[0], "delete")?;
                    }
                 },
                 _ => continue,
            }
        }
    }
    return Ok(0);
}

pub fn start_cmd(ctrname: &str) {
    debug!("Dummy {} container start.", ctrname);
}

pub fn stop_cmd(ctrname: &str) {
    debug!("Dummy {} container stop.", ctrname);
}

pub fn restart_cmd(ctrname: &str) {
    debug!("Dummy {} container restart.", ctrname);
}

pub fn clean_cmd() {
    debug!("Cleaning conatiners.");
    match ctr_clean() {
        Ok(_) => info!("Cleaning completetd"),
        Err(err) => error!("Cleaning completetd with error:: {} ", err), 
    };
}

pub fn restart_u_ce() {
    debug!("Dummy uContainer Engine restart.");
}
