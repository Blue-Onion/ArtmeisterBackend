package utlis

import "github.com/Blue-Onion/ArtmeisterBackend/config"

var conf *config.Config = config.GetConfig()
var Db string = conf.DbUrl
