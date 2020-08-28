package log


import (
  "fmt"
  "github.com/op/go-logging"
  "io"
  "os"
)

type LogLevel int
type LogOutput int
type LogState bool
type LogLevelStr string


// Log levels:

const (
  LL_NONE     LogLevel = 0
  LL_ERROR    LogLevel   = 1
  LL_WARNING  LogLevel  = 2
  LL_INFO     LogLevel   = 3
  LL_DEBUG    LogLevel   = 4
  LL_TRACE    LogLevel   = 5
  LL_ALL      LogLevel   = 6
  LL_NOTSET   LogLevel  = 10
)


//Log Levels String value corresponding to string
const (
  LL_NONE_STR     LogLevelStr = "none"
  LL_ERROR_STR    LogLevelStr = "error"
  LL_WARNING_STR  LogLevelStr = "warning"
  LL_INFO_STR     LogLevelStr = "info"
  LL_DEBUG_STR    LogLevelStr = "debug"
  LL_TRACE_STR    LogLevelStr = "trace"
  LL_ALL_STR      LogLevelStr = "all"
  LL_NOTSET_STR   LogLevelStr = "notset"
)

//Log output destination corresponding to LogOutput
const (
  LO_NONE       LogOutput = 0
  LO_SCREEN     LogOutput = 1
  LO_FILE       LogOutput = 2
  LO_BOTH       LogOutput = 3
  LO_ALl        LogOutput = 4
)

// Log state Corresponding to LogState
const (
  LOG_ENABLED LogState  = true
  LOG_DISABLED LogState = false
)


const (
  SECU_LOG_ENABLED  = true
  SECU_LOG_DISABLED = false
)

var loggingEnabled LogState = LOG_ENABLED
var defaultModule = "default"
var log = logging.MustGetLoger(defaultModule)
var AuthServerVersion = ""
var seculoggingEnabled = SECU_LOG_ENABLED

// Setting format for log format , basd on color, location, path, loglevel
var format = logging.MustStringFormatter(
  `%{level:.1s}%{color:reset}%{message}`,
)

// Secure is used to prevent seneitive infoirmation to be printed in log, like
//passwords
// Secure is just an example type which is being used to implement the
//Redactor Interface, it will return an interface which you can use
// for any types. here we use Redacted function which return *****

type Secure string

func (p Secure) Redacted() interface{} {
  if seculoggingEnabled {
    return logging.Redact(string(p))
  } else{
    return string(p)
  }
}


// Setup the authserver utility
// This function must be called only once by the client of this library
// Must be called before any other logging fucntion of this library is called
// It runs 'true' if setup is successful, else return 'false'
// The caller is expected to check the return value before procedding to other functions

func SetupLogger(logState LogState, level LogLevel, dirName string, processName string, logOutput LogOutput, secuenabled LogState) bool {
  loggingEnabled = logState
  seculoggingEnabled  = secuenabled
  if logOutput == LO_NONE || level == LL_NONE {
    loggingEnabled = false
  }
  if level == LL_ALL || level == LL_TRACE {
    level = LL_DEBUG
  }
  if loggingEnabled {
    // Logging Output destination handlers
    var stdoutBackendLeveled, fileBackendLeveled logging.LeveledBackend
  }
}
