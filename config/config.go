package config

import (
	"github.com/gobuffalo/envy"
)

var DB_DSN = envy.Get("DB_DSN", "")
var DISCORD_TOKEN = envy.Get("DISCORD_TOKEN", "")
var ENV = envy.Get("GO_ENV", "development")
