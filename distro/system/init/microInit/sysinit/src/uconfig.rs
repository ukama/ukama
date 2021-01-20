// src/uconfig.rs
use log::*;

use serde::Deserialize;
use std::{fs, io};

pub static CONFIG_FILE: &str = "/etc/uInit.yml";

//uInit.config
#[derive(Debug, Deserialize, Clone)]
pub struct Config {
  init: Vec<AppConfig>,
  onboot: Vec<AppConfig>,
  onshutdown: Vec<AppConfig>,
  sysservice: Vec<AppConfig>,
  userservice: Vec<AppConfig>,
}

//AppConfig
#[derive(Debug, Deserialize, Clone)]
pub struct AppConfig {
  pub name: String,
  pub image: Option<String>,
  pub version: String,
  pub path: Option <String>,
}

#[derive(Debug, Deserialize, Clone)]
pub struct OnBootConfig {
  onboot: Vec<AppConfig>,
}

lazy_static! {
  static ref CONFIG: Config = Config::init();
}

impl Config {
  fn init() -> Self {
    let contents = Self::read_config_file().unwrap();
    trace!("Data read is :: {}", contents);
    return serde_yaml::from_str(&contents).unwrap();
  }

  pub fn get() -> Self {
    CONFIG.to_owned()
  }

  pub fn read_config_file() -> Result<String, io::Error> {
    trace!("Trying reading {}",CONFIG_FILE);
    fs::read_to_string(CONFIG_FILE)
  }

  pub fn get_init_config() -> &'static Vec<AppConfig>  {
    //let config = CONFIG.to_owned();
    //println!("{:#?}", config.onboot);
    &CONFIG.init
  }

  pub fn get_onboot_config() -> &'static Vec<AppConfig>  {
    //let config = CONFIG.to_owned();
    //println!("{:#?}", config.onboot);
    &CONFIG.onboot
  }

  pub fn get_onshutdown_config() -> &'static Vec<AppConfig> {
    //let config = CONFIG.to_owned();
    //println!("{:#?}", config.onboot);
    &CONFIG.onshutdown
  }
  #[allow(dead_code)]
  pub fn get_sysservice_config() -> &'static Vec<AppConfig> {
    //let config = CONFIG.to_owned();
    //println!("{:#?}", config.onboot);
    &CONFIG.sysservice
  }
  #[allow(dead_code)]
  pub fn get_usrservice_config() -> &'static Vec<AppConfig> {
    //let config = CONFIG.to_owned();
    //println!("{:#?}", config.onboot);
    &CONFIG.sysservice
  }
}