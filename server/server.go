package server

import (
  "log"
  "github.com/gorilla/mux"
  "net/http"
  "time"

)
var Authurl             string
var UserHandler         string
var MaintainenceHandler string
var AdminHandler        string
