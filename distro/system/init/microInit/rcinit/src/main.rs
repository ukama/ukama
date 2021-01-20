
use prctl;
use log::*;
use simplelog::*;
use csv::Error;
use csv::ReaderBuilder;
use glob::glob;
use serde::Deserialize;
use system_shutdown::{reboot, shutdown};
use rlimit::{setrlimit, Resource, Rlim};
use nix::{ libc, mount, sys, unistd, sys::stat };
use std::{ env, fs, io, fs::File, io::Write,
           path::Path, process::Command, io::prelude::*,
           os::unix::fs::PermissionsExt};


//Cgroup info
#[derive(Debug, Deserialize)]
struct Record {
    subsys_name: String,
    hierarchy: Option<String>,
    num_cgroups: Option<String>,
    enabled: Option<String>,
}

//Mount info
#[derive(Debug, Deserialize)]
struct Mount {
    device: String,
    mount_point: String,
    file_system: String,
    flags: String,
    cmpfst1: String,
    cmpfst2: String,
}

//Read files in the directory.
fn readdir(dir: &str) -> Result<Vec<std::path::PathBuf>, Error> {
    let mut files = fs::read_dir(dir)?
        .map(|res| res.map(|e| e.path()))
        .collect::<Result<Vec<_>, io::Error>>()?;

    files.sort();

    Ok(files)
}

//Write to a file
fn write(path: &str, data: &str) -> io::Result<()> {
    let mut file = File::create(path)?;
    let metadata = file.metadata()?;
    
    let mut permissions = metadata.permissions();
    permissions.set_mode(0o600);
    
    file.write_all(data.as_bytes())?;
    file.sync_all()?;
    
    Ok(())
}

//Read from a file
fn read(path: &str) -> Result<String, io::Error> {
    let mut file = File::open(path)?;
    let mut data = String::new();
    
    file.read_to_string(&mut data)?;
    
    Ok(data)
}

//Parsing Cgroups from /proc
fn cgroup_parse() -> Result<Vec<Box<String>>, Error> {
    let mut vec: Vec<Box<String>> = vec![];
    let mut reader = ReaderBuilder::new()
        .delimiter(b'\t')
        .has_headers(false)
        .from_path("/proc/cgroups")?;
        
        //Deserialize the enteries read. 
        for result in reader.deserialize::<Record>() {
        let record = result?;
        println!("{:?}", record);
        let data = match record.enabled {
            None => continue,
            Some(s) => s,
        };

        //Store the enabled CGroups.
        if data == "1" {
            vec.push(Box::new(record.subsys_name));
            trace!("Record pushed.");
        }
    }
    
    return Ok(vec);
}

//Loading driver using modprobe
fn modalias(path: &str) {
    let alias = match read(&path) {
        Ok(str) => str,
        Err(_)=> return,
    };

    let output = Command::new("/sbin/modprobe")
    .arg("-abq")
    .arg(&alias)
    .output()
    .expect("failed to execute process");
    
    trace!("/sbin/modprobe -abq {} status: {}", alias, output.status);

}

//Check if file exists
fn exists(path: &str) -> bool {
    match sys::stat::stat(path) {
        Ok(_) => true,
        Err(_) => false,
    }
}

// Make directory
fn mkdir(dir: &str, mode: stat::Mode) {
    match unistd::mkdir(dir, mode) {
        Ok(_) => trace!("created {:?}", dir),
        Err(err) => error!("Error creating directory: {}", err),
    }
}

// Make device
fn mkchar(dir: &str, perm: stat::Mode, major: u64, minor: u64) {
    let mode = stat::SFlag::S_IFCHR;
    let dev = sys::stat::makedev(major, minor);
    
    match sys::stat::mknod(dir, mode, perm, dev) {
        Ok(_) => trace!("created device  {:?} file ", dir),
        Err(err) => error!("Error creating dev {:?} file: {}", dir, err),
    }
}

// Make symlink
fn symlink(path: &str, newpath: &str) {
    match unistd::symlinkat(path, None, newpath) {
        Ok(_) => trace!("created link  for {} file at {} ", path, newpath),
        Err(err) => error!(
            "Error creating link {} file failed at {} :: {}",
            path, newpath, err
        ),
    }
}

// set as a subreaper
fn set_subreaper(subreaper: bool) {
    match prctl::set_child_subreaper(subreaper) {
        Ok(_) => debug!("Set as a subreaper."),
        Err(err) => error!("Error setting as a subreaper :: {}", err),
    }
}

// Mount with error return
fn mount(source: &str, target: &str, fstype: &str, flags: mount::MsFlags, data: &str) -> bool {
    let mnt = match mount::mount(Some(source), target, Some(fstype), flags, Some(data)) {
        Ok(_) => true,
        Err(err) => {
            error!(
                "Mount {} failed on {} {:?} flags: {:?} : {}",
                source, fstype, target, flags, err
            );
            false
        }
    };

    if !mnt {
        error!("Error mounting {} to {}", source, target);
    }

    mnt
}

// Mount in some cases, do not even log an error
fn mount_silent(source: &str, target: &str, fstype: &str, flags: mount::MsFlags, data: &str) {
    match mount::mount(Some(source), target, Some(fstype), flags, Some(data)) {
        Ok(_) => trace!("Mounting {} to {}", source, target),
        Err(err) => error!("Error mounting {} to {}: {}", source, target, err),
    }
}

//Unmounting
fn unmount(target: &str, flags: mount::MntFlags) {
    match mount::umount2(target, flags) {
        Ok(_) => trace!("Unmounting {} done", target),
        Err(err) => error!("Unmounting {} failed :: {}", target, err),
    }
}

//Parse mount -a output
fn parse_mounts() -> io::Result<()> {
    let mut reader = ReaderBuilder::new()
        .delimiter(b' ')
        .has_headers(false)
        .from_path("/proc/mounts")?;

     //deserialize the mount data   
    for result in reader.deserialize::<Mount>() {
        let record = result?;
        trace!("{:?}", record);

        //Check filesystem
        match record.file_system.as_str() {
            "ext2" | "ext3" | "ext4" | "btrfs" | "xfs" | "vfat" | "msdos" | "overlay" => {
                unmount(record.mount_point.as_str(), mount::MntFlags::empty())
            }
            "nfs" | "nfs4" | "cifs" => {
                unmount(record.mount_point.as_str(), mount::MntFlags::MNT_DETACH)
            }
            _ => trace!("Nothing required"),
        }
    }

    Ok(())
}

//Unmounting FS
fn do_umounts() {
    match parse_mounts() {
        Ok(_) => trace!("Umount completed."),
        Err(err) => error!("Umount failed with {}", err),
    }
}

//Mounting FS
fn do_mounts() {
    // mount proc filesystem
    mount(
        "proc",
        "/proc",
        "proc",
        mount::MsFlags::MS_NODEV
            | mount::MsFlags::MS_NOSUID
            | mount::MsFlags::MS_NOEXEC
            | mount::MsFlags::MS_RELATIME,
        "",
    );

    // remount rootfs read only if it is not already
    // mount_silent(
    //     "",
    //     "/",
    //     "",
    //     mount::MsFlags::MS_REMOUNT | mount::MsFlags::MS_RDONLY,
    //     "",
    // );

    mount_silent(
        "",
        "/",
        "",
        mount::MsFlags::MS_REMOUNT,
        "",
    );

    // mount tmpfs for /tmp and /run
    mount(
        "tmpfs",
        "/run",
        "tmpfs",
        mount::MsFlags::MS_NODEV
            | mount::MsFlags::MS_NOSUID
            | mount::MsFlags::MS_NOEXEC
            | mount::MsFlags::MS_RELATIME,
        "size=10%,mode=755",
    );
    
    mount(
        "tmpfs",
        "/tmp",
        "tmpfs",
        mount::MsFlags::MS_NODEV
            | mount::MsFlags::MS_NOSUID
            | mount::MsFlags::MS_NOEXEC
            | mount::MsFlags::MS_RELATIME,
        "size=10%,mode=1777",
    );

    // mount tmpfs for /var. This may be overmounted with a persistent filesystem later
    mount(
        "tmpfs",
        "/var",
        "tmpfs",
        mount::MsFlags::MS_NODEV
            | mount::MsFlags::MS_NOSUID
            | mount::MsFlags::MS_NOEXEC
            | mount::MsFlags::MS_RELATIME,
        "size=50%,mode=755",
    );

    //Check if possibly all these can be created in rootfs
    let rwx_rx_rx = stat::Mode::S_IRWXU
        | stat::Mode::S_IROTH
        | stat::Mode::S_IXOTH
        | stat::Mode::S_IRGRP
        | stat::Mode::S_IXGRP;
    let rx_rx_rx = stat::Mode::S_IRUSR
        | stat::Mode::S_IXUSR
        | stat::Mode::S_IROTH
        | stat::Mode::S_IXOTH
        | stat::Mode::S_IRGRP
        | stat::Mode::S_IXGRP;
    let rwx_rwx_rwx = stat::Mode::S_IRWXU | stat::Mode::S_IRWXG | stat::Mode::S_IRWXO;
    mkdir("/var/cache", rwx_rx_rx);
    mkdir("/var/empty", rx_rx_rx);
    mkdir("/var/lib", rwx_rx_rx);
    mkdir("/var/local", rwx_rx_rx);
    mkdir("/var/lock", rwx_rx_rx);
    mkdir("/var/log", rwx_rx_rx);
    mkdir("/var/opt", rwx_rx_rx);
    mkdir("/var/spool", rwx_rx_rx);
    mkdir("/var/tmp", stat::Mode::S_ISVTX | rwx_rwx_rwx);
    // mount devfs
    mount(
        "dev",
        "/dev",
        "devtmpfs",
        mount::MsFlags::MS_NOSUID | mount::MsFlags::MS_NOEXEC | mount::MsFlags::MS_RELATIME,
        "size=10m,nr_inodes=248418,mode=755",
    );

    // make minimum necessary devices
    let rw__ = stat::Mode::S_IRUSR | stat::Mode::S_IWUSR;
    let rw_w_ = stat::Mode::S_IRUSR | stat::Mode::S_IWUSR | stat::Mode::S_IWGRP;
    let rw_rw_rw = stat::Mode::S_IRUSR
        | stat::Mode::S_IWUSR
        | stat::Mode::S_IRGRP
        | stat::Mode::S_IWGRP
        | stat::Mode::S_IROTH
        | stat::Mode::S_IWOTH;
    let rw_rw_ =
        stat::Mode::S_IRUSR | stat::Mode::S_IWUSR | stat::Mode::S_IRGRP | stat::Mode::S_IWGRP;

    mkchar("/dev/console", rw__, 5, 1);
    mkchar("/dev/tty1", rw_w_, 4, 1);
    mkchar("/dev/tty", rw_rw_rw, 5, 0);
    mkchar("/dev/null", rw_rw_rw, 1, 3);
    mkchar("/dev/kmsg", rw_rw_, 1, 11);
    // make standard symlinks
    symlink("/proc/self/fd", "/dev/fd");
    symlink("/proc/self/fd/0", "/dev/stdin");
    symlink("/proc/self/fd/1", "/dev/stdout");
    symlink("/proc/self/fd/2", "/dev/stderr");
    symlink("/proc/kcore", "/dev/kcore");
    // dev mountpoints
    mkdir("/dev/mqueue", stat::Mode::S_ISVTX | rwx_rwx_rwx);
    mkdir("/dev/shm", stat::Mode::S_ISVTX | rwx_rwx_rwx);
    mkdir("/dev/pts", rwx_rx_rx);
    // mounts on /dev
    mount(
        "mqueue",
        "/dev/mqueue",
        "mqueue",
        mount::MsFlags::MS_NODEV | mount::MsFlags::MS_NOSUID | mount::MsFlags::MS_NOEXEC,
        "",
    );
    mount(
        "shm",
        "/dev/shm",
        "tmpfs",
        mount::MsFlags::MS_NODEV | mount::MsFlags::MS_NOSUID | mount::MsFlags::MS_NOEXEC,
        "mode=1777",
    );
    mount(
        "devpts",
        "/dev/pts",
        "devpts",
        mount::MsFlags::MS_NODEV | mount::MsFlags::MS_NOSUID,
        "gid=5,mode=0620",
    );

    // sysfs
    mount(
        "sysfs",
        "/sys",
        "sysfs",
        mount::MsFlags::MS_NODEV | mount::MsFlags::MS_NOSUID | mount::MsFlags::MS_NOEXEC,
        "",
    );
    // some of the subsystems may not exist, so ignore errors
    mount_silent(
        "securityfs",
        "/sys/kernel/security",
        "securityfs",
        mount::MsFlags::MS_NODEV | mount::MsFlags::MS_NOSUID | mount::MsFlags::MS_NOEXEC,
        "",
    );
    mount_silent(
        "debugfs",
        "/sys/kernel/debug",
        "debugfs",
        mount::MsFlags::MS_NODEV | mount::MsFlags::MS_NOSUID | mount::MsFlags::MS_NOEXEC,
        "",
    );
    mount_silent(
        "configfs",
        "/sys/kernel/config",
        "configfs",
        mount::MsFlags::MS_NODEV | mount::MsFlags::MS_NOSUID | mount::MsFlags::MS_NOEXEC,
        "",
    );
    mount_silent(
        "fusectl",
        "/sys/fs/fuse/connections",
        "fusectl",
        mount::MsFlags::MS_NODEV | mount::MsFlags::MS_NOSUID | mount::MsFlags::MS_NOEXEC,
        "",
    );
    mount_silent(
        "selinuxfs",
        "/sys/fs/selinux",
        "selinuxfs",
        mount::MsFlags::MS_NODEV | mount::MsFlags::MS_NOSUID | mount::MsFlags::MS_NOEXEC,
        "",
    );
    mount_silent(
        "pstore",
        "/sys/fs/pstore",
        "pstore",
        mount::MsFlags::MS_NODEV | mount::MsFlags::MS_NOSUID | mount::MsFlags::MS_NOEXEC,
        "",
    );
    mount_silent("bpffs", "/sys/fs/bpf", "bpf", mount::MsFlags::MS_NOSUID, "");

    mount_silent(
        "efivarfs",
        "/sys/firmware/efi/efivars",
        "efivarfs",
        mount::MsFlags::MS_NODEV | mount::MsFlags::MS_NOSUID | mount::MsFlags::MS_NOEXEC,
        "",
    );

    // misc /proc mounted fs
    mount_silent(
        "binfmt_misc",
        "/proc/sys/fs/binfmt_misc",
        "binfmt_misc",
        mount::MsFlags::MS_NODEV | mount::MsFlags::MS_NOSUID | mount::MsFlags::MS_NOEXEC,
        "",
    );

    // mount cgroup root tmpfs
    mount(
        "cgroup_root",
        "/sys/fs/cgroup",
        "tmpfs",
        mount::MsFlags::MS_NODEV | mount::MsFlags::MS_NOSUID | mount::MsFlags::MS_NOEXEC,
        "mode=755,size=10m",
    );

    // mount cgroups filesystems for all enabled cgroups
    let cg_list = match cgroup_parse() {
        Ok(vec) => vec,
        Err(_) => Vec::new(), //Creating a empty vector
    };

    for cg in cg_list {
        let path = format!("{}{}", "/sys/fs/cgroup/", cg);
        mkdir(&path, rw_rw_rw);
        mount(
            &cg,
            &path,
            "cgroup",
            mount::MsFlags::MS_NODEV | mount::MsFlags::MS_NOSUID | mount::MsFlags::MS_NOEXEC,
            &cg,
        );
    }

    // use hierarchy for memory
    match write("/sys/fs/cgroup/memory/memory.use_hierarchy", "1") {
        Ok(_) => trace!("Write  Sucessfull on file."),
        Err(_) => error!("Error while writing to file."),
    }

    // many things assume systemd
    mkdir("/sys/fs/cgroup/systemd", rw_rw_rw);
    mount(
        "cgroup",
        "/sys/fs/cgroup/systemd",
        "cgroup",
        mount::MsFlags::empty(),
        "none,name=systemd",
    );

    // make / rshared
    mount(
        "",
        "/",
        "",
        mount::MsFlags::MS_REC | mount::MsFlags::MS_SHARED,
        "",
    );
}

//Hotplug for devices
fn do_hotplug() {
    let mdev = "/sbin/mdev";
    
    // start mdev for hotplug
    match write("/proc/sys/kernel/hotplug", mdev) {
        Ok(_) => debug!("Succesfully wrote to hotplug"),
        Err(_) => error!("Failed to write to hotplug."),
    }

    let devpath = "/sys/devices";

    let files = match readdir(devpath) {
        Ok(fvec) => fvec,
        Err(_) => return,
    };

    trace!("Device files present in {:?} are {:?}", devpath, files);
    let file_iter = files.iter();
    for f in file_iter {
        let ufile = match f.to_str() {
            Some(file) => file,
            None => continue,
        };

        let uevent = format!("{}{}{}", devpath, ufile, "uevent");
        if ufile.starts_with("usb") && exists(&uevent) {
            match write(&uevent, "add") {
                Ok(_) => trace!("Succesfully wrote to {}", uevent),
                Err(_) => error!("Failed to write to {}.", uevent),
            }
        }
    }

    let output = Command::new("mdev")
        .arg("-s")
        .output()
        .expect("failed to execute process");
    debug!("mdev -s status: {}", output.status);

    //support for cold plug
    for entry in glob("sys/bus/*/devices/*/modalias").expect("Failed to read glob for cold plug") {
        match entry {
            Ok(path) => {
                trace!(" Loading driver for glob pattern {:?}", path.display());
                let pstr = match path.to_str() {
                    Some(ps) => ps,
                    None => continue,
                };
                modalias(&pstr);
            }
            Err(err) => {
                error!("Error while matching pattern. Moving to next {:?}", err);
                continue;
            }
        }
    }

}

//Adjust HW clock
fn do_clock() {
    let output = Command::new("/sbin/hwclock")
        .arg("--hctosys")
        .arg("--utc")
        .output()
        .expect("failed to execute process");
    debug!("HwClock status: {}", output.status);
}

//Add loopback devices
fn do_loopback() {
    let output = Command::new("/sbin/ip")
        .arg("addr")
        .arg("add")
        .arg("127.0.0.1/8")
        .arg("dev")
        .arg("lo")
        .arg("brd")
        .arg("+")
        .arg("scope")
        .arg("host")
        .output()
        .expect("failed to execute process");
    debug!("sbin addr : {}", output.status);

    let output = Command::new("/sbin/ip")
        .arg("route")
        .arg("add")
        .arg("127.0.0.0/8")
        .arg("dev")
        .arg("lo")
        .arg("scope")
        .arg("host")
        .output()
        .expect("failed to execute process");
    debug!("sbin route : {}", output.status);

    let output = Command::new("/sbin/ip")
        .arg("link")
        .arg("set")
        .arg("lo")
        .arg("up")
        .output()
        .expect("failed to execute process");
    debug!("sbin route : {}", output.status);
}

//Set resource limits
fn rlimits(resource: Resource, soft: Rlim, hard: Rlim) {
    match setrlimit(resource, soft, hard) {
        Ok(_) => debug!("Limits set for {:?} {:?} {:?}.", resource, soft, hard),
        Err(err) => error!(
            "Limits set failed for {:?} {:?} {:?} with err {:?}.",
            resource, soft, hard, err
        ),
    }
}

//Request for setting resource limits.
fn do_limits() {
    rlimits(
        Resource::NOFILE,
        Rlim::from_raw(1048576),
        Rlim::from_raw(1048576),
    );
    rlimits(Resource::NPROC, Rlim::INFINITY, Rlim::INFINITY);
}

//Request to set hostname
fn do_hostname() {
    //Read hostname 
    let hname = match read("/etc/hostname") {
        Ok(data) => data,
        Err(_) => String::from(""),
    };

    //if /etc/hostname exist
    if hname != "" {
        match unistd::sethostname(&hname) {
            Ok(_) => println!("Hostname set to {}", hname),
            Err(err) => error!("Failed to set host name to {} :: {}", hname, err),
        }
    }

    //if no hostname exist
    let mut buf = [0u8; 64];
    let hname_cstr = unistd::gethostname(&mut buf).expect("Failed getting hostname");
    let hname = hname_cstr.to_str().expect("Hostname wasn't valid UTF-8");
    if hname != "" {
        debug!("Hostname: {}", hname);
        return;
    }

    //Read MAC
    let mac = match read("/sys/class/net/eth0/address") {
        Ok(data) => data,
        Err(_) => String::from(""),
    };
    if mac == "" {
        warn!("No mac address found.");
        return;
    }
    let mac = mac.replace(":", "");
    let hname = format!("{}{}", "ukama-", mac);
    
    //Set hostname with mac as a postfix to ukama
    match unistd::sethostname(&hname) {
        Ok(_) => info!("Hostname set to {}", hname),
        Err(err) => error!("Failed to set host name to {} :: {}", hname, err),
    }

}


//Resolvconf for network.
fn do_resolvconf() {
    let link = match fs::read_link("/etc/resolv.conf") {
        Ok(ln) => ln,
        Err(_) => return,
    };

    let parent = link.parent().unwrap();
    let dir = parent.to_str().unwrap();
    // {
    //     None => "",
    //     Some(str) => str,
    // };
    

    let link = link.to_str().unwrap();
    let rwx_rx_rx = stat::Mode::S_IRWXU
        | stat::Mode::S_IROTH
        | stat::Mode::S_IXOTH
        | stat::Mode::S_IRGRP
        | stat::Mode::S_IXGRP;
    mkdir(dir, rwx_rx_rx);

    match write(link, "") {
        Ok(_) => debug!("Write is success to {}.", link),
        Err(err) => error!("Write failed to {} :: {}", link, err),
    }
}


//Reaping the exitted child.
fn do_reap() {
    let mut status = 0;
    loop {
        unsafe {
            let pid = libc::waitpid(-1, &mut status, libc::WNOHANG);
            if pid < 0 {
                return;
            }
        }
    }
}

//start the process
fn run_init(path: &str) {
    let files = match readdir(path) {
        Ok(fvec) => fvec,
        Err(_) => return,
    };

    println!("All files present in {:?} are {:?}", path, files);
    let file_iter = files.iter();
    for f in file_iter {
        let file = Path::new(path).join(f);
        let file = match file.to_str() {
            Some(file) => file,
            None => continue,
        };

        let stat = match sys::stat::stat(file) {
            Ok(fs) => {
                println!("Found {} file for execution", file);
                fs
            }
            Err(err) => {
                println!("{} not found : {}", file, err);
                continue;
            }
        };

        println!("Stat: {:?}", stat);

        //Mode is a not regular file
        //continue;
        let metadata = match fs::metadata(file) {
            Ok(data) => data,
            Err(err) => {
                println!("Failed to get meta data for {}: {}", file, err);
                continue;
            }
        };

        if !metadata.is_file() {
            println!("File {} is not regular file, Stat: {:?} ", file, metadata);
            continue;
        }
        println!("Exec {}", file);
        Command::new(file)
            .spawn()
            .expect("failed to execute process");   
    }
}

//poweroff handler
fn do_shutdown(action: &str) {
    run_init("/etc/shutdown.d");
    unistd::sync();
    do_umounts();
    match action {
        "poweroff" => match shutdown() {
            Ok(_) => info!("Shutting down system...!"),
            Err(e) => error!("Failed to shutdown system: {}", e),
        },
        _ => match reboot() {
            Ok(_) => info!("Rebooting system...!"),
            Err(err) => error!("Failed to reboot system: {}", err),
        },
    }
    std::process::exit(0);
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

    trace!("Starting rc.init...!");
    let args: Vec<String> = env::args().collect();
    debug!("Argument is {}", args[0]);
    let arg = Path::new(&args[0]).file_name().unwrap();
    trace!("Unwrapped arg is {:#?}", arg);
    if arg == "rc.shutdown" {
        let action = "poweroff";
        //TODO:: Revisit later 
        // if args.len() > 1 {
        //     println!("Argument1 is {}", args[1]);
        //     action = &args[1];
        // }
        do_shutdown(action);
    }

    let userspace = exists("/proc/self");

    if userspace {
        set_subreaper(true);
    } else {
        do_mounts();
        do_hotplug();
        do_clock();
        do_loopback();
    }

    do_limits();
    do_hostname();
    do_resolvconf();

    //start the executables from /etc/init.d path
    run_init("/etc/init.d");

    if userspace {
        do_reap();
    }
}
