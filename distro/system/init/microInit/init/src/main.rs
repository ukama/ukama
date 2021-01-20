extern crate exec;
extern crate libc;
extern crate nix;

use log::*;
use simplelog::*;
use nix::mount;
use nix::unistd;
use walkdir::Error;
use walkdir::WalkDir;
use std::ffi::CString;
use std::{ fs, process, thread, time, fs::File, process::Command };


// Copy file from initramfs to new mount.
fn copyfs(newroot: &str) -> Result<(), Error> {
    let cwd = "/";
    trace!("Copying the files to path {} from {}", newroot, cwd);

    //Walk through each directory under /
    for entry in WalkDir::new(cwd)
        .min_depth(1)
        .max_depth(1)
        .same_file_system(true)
    {
        // Avoid recursive copy of /mnt
        let file = entry.unwrap();
        if file.file_type().is_dir() {
            if file.path().to_str() == Some(newroot) {
                trace!("Skipping {} ", file.path().display());
                continue;
            }
        }

        //Copy rest of the files with permissions.
        let fpath = file.path().as_os_str().to_string_lossy().into_owned();
        let destpath = format!("{}{}", newroot, fpath);
        trace!(
            "Directory {} getting copied from {} to New root {} at DestPath {}",
            cwd, fpath, newroot, destpath
        );

        //Copy Command
        let _copy_process = Command::new("cp")
            .arg("-rp")
            .arg(&fpath)
            .arg(&destpath)
            .spawn()
            .expect("Init:: Failure to copy data to root.");
    }

    //Sync filesystem
    unistd::sync();

    Ok(())
}

// Preparing new rootfs
fn preparefs() {
    let mnt_point = "/mnt";
    let none: Option<&str> = None;

    //Mount Flags.
    let _proc_mount_flags = mount::MsFlags::MS_NOSUID
        | mount::MsFlags::MS_NODEV
        | mount::MsFlags::MS_NOEXEC
        | mount::MsFlags::MS_RELATIME;
    
    //New tmpfs root 
    match mount::mount(
        Some("rootfs"),
        mnt_point,
        Some("tmpfs"),
        mount::MsFlags::empty(),
        none,
    ) {
        Ok(_) => info!("Mounted {} on {}.", "rootfs", mnt_point),
        Err(err) => panic!("Mounting {} on {} failed :: {}.", "rootfs", mnt_point, err),
    }

    //Copy files
    match copyfs(&mnt_point) {
        Ok(_) => info!("Copy Completed to new rootfs {}.", mnt_point),
        Err(err) => panic!("Failed to copy files to new rootfs{}: {}.", mnt_point, err),
    }

    //change directory
    match unistd::chdir(mnt_point) {
        Ok(_) => info!("Change directory to {} done.", mnt_point),
        Err(err) => panic!("Change directory to {} failed: {}.", mnt_point, err),
    }

    //change root
    let changeroot = ".";
    match unistd::chroot(changeroot) {
        Ok(_) => debug!("Change root {} done ", changeroot),
        Err(err) => panic!("Change root to {} failed : {}.", changeroot, err),
    }

    // Only for debug
    let lpaths = fs::read_dir("/").unwrap();
    for lpath in lpaths {
        trace!("Name: {}", lpath.unwrap().path().display())
    }
}

// Init function
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
    
    //Check if booting from TMPFS or from RAMFS
    let ramfs_magic = 0x858458f6;
    let tmpfs_magic = 0x01021994;
    let fstype;
    unsafe {
        let root = CString::new("/").unwrap();
        let mut statfs: libc::statfs = std::mem::zeroed();
        if libc::statfs(root.as_ptr(), &mut statfs) >= 0 {
            println!("{:#x}", statfs.f_type);
        }
        fstype = statfs.f_type;
    }
    if fstype == ramfs_magic || fstype == tmpfs_magic {
        debug!("Init:: FS found is {:#x} Copy filesystem to /mnt/", fstype);
        preparefs();
    }
    
    //busybox init /sbin/init 
    let exec_prg = "/sbin/init";
    debug!("Init:: Starting {} after 10 secs", exec_prg);

    //Give time to sync
    let ten_millis = time::Duration::from_millis(10000);
    //let now = time::Instant::now();

    thread::sleep(ten_millis);

    // Exec the specified program.  If all goes well, this will never
    // return.  If it does return, it will always retun an error.
    let err = exec::Command::new(exec_prg).exec();
    error!("Error: {} Can't exec {}", err, exec_prg);
    process::exit(1);
}
