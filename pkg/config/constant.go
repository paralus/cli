package config

import (
	"path/filepath"

	"github.com/RafaySystems/rcloud-cli/pkg/utils"
)

var RAFAY_DIR_DEFAULT_LOCATION = filepath.Join(utils.GetUserHome(), ".rafay")
var CLI_DIR_DEFAULT_LOCATION = filepath.Join(RAFAY_DIR_DEFAULT_LOCATION, "cli")
var CONFIG_FILE_DEFAULT_NAME = "config"
var LOG_FILE_DEFAULT_NAME = "cli.log"
var CONFIG_FILE_DEFAULT_LOCATION = filepath.Join(CLI_DIR_DEFAULT_LOCATION, CONFIG_FILE_DEFAULT_NAME)
var LOG_FILE_DEFAULT_LOCATION = filepath.Join(CLI_DIR_DEFAULT_LOCATION, LOG_FILE_DEFAULT_NAME)

var VERBOSE_FLAG_NAME = "verbose"
var CONFIG_FLAG_NAME = "config"
var PROFILE_FLAG_NAME = "profile"

var CONFIG_API_VERSION = "1.0"
var WP_API_VERSION = "1.0"
var CLI_VERSION = "1.0"
var CLI_BUILD_NUMBER = "NA"
var CLI_ARCH = "NA"
var CLI_BUILD_TIME = "NA"

var GENERIC_ERROR_MESSAGE = "CLI faced an issue while running the command %s. Please use -v flag to see debug logs."
var AUTH_FAILURE_ERROR_MESSAGE = "CLI faced an authentication failure. Please check your config file or use \"init\" command to configure cli with the appropriate credentials."
var CLI_MISMATCH_ERROR_MESSAGE = "Your CLI version is incompatible with the current API version. Please visit console to download the most uptodate CLI version."

var WORKLOAD_CREATED = 201
var REQUEST_SUCCESSFUL = 200
var AUTHENTICATION_FAILURE = 401
