package utlis

import "github.com/Blue-Onion/ArtmeisterBackend/config"

var conf *config.Config = config.LoadConfig()
var Db string = conf.DbUrl
