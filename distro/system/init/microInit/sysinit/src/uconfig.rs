// src/uconfig.rs
use log::*;

use serde::Deserialize;
use std::{fs, io};

//pub static CONFIG_FILE: &str = "/etc/uInit.yml";
pub static CONFIG_FILE: &str = "/etc/microInit.toml";

//microInit.config
#[derive(Debug, Deserialize, Clone)]
pub struct Config {
  init: Vec<AppConfig>,
  init_bundle: Vec<AppConfig>,
  shutdown_bundle: Vec<AppConfig>,
}

//AppConfig
#[derive(Debug, Deserialize, Clone)]
pub struct AppConfig {
  pub name: String,
  pub image: Option<String>,
  pub version: String,
  pub path: Option<String>,
}

lazy_static! {
  static ref CONFIG: Config = Config::init();
}

impl Config {
  fn init() -> Self {
    let contents = Self::read_config_file().unwrap();
    trace!("Config contents are :: {}", contents);
    return toml::from_str(&contents).unwrap();
  }

  pub fn get() -> Self {
    CONFIG.to_owned()
  }

  //Read init config file
  pub fn read_config_file() -> Result<String, io::Error> {
    trace!("Reading init config {}", CONFIG_FILE);
    fs::read_to_string(CONFIG_FILE)
  }

  //Read init config
  pub fn get_init_config() -> &'static Vec<AppConfig> {
    &CONFIG.init
  }
  //Read init-bundles config
  pub fn get_init_bundle_config() -> &'static Vec<AppConfig> {
    &CONFIG.init_bundle
  }

  //Read shutdown-bundles config
  pub fn get_shutdown_bundle_config() -> &'static Vec<AppConfig> {
    &CONFIG.shutdown_bundle
  }

  //Read OCI runtime
  pub fn get_oci_runtime() -> Option<&'static AppConfig> {
    let init = Config::get_init_config();
    let list = init.iter();
    for app in list {
      if app.name == "oci-runtime" {
        return Some(&app);
      } else {
        continue;
      }
    }
    None
  }
}
