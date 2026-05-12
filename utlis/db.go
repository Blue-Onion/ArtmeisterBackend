package utlis

import "github.com/Blue-Onion/RestApi-Go/config"

var conf *config.Config = config.LoadConfig()
var Db string = conf.DbUrl
