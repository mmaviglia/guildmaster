package config

import (
	"fmt"
	"strconv"

	"github.com/gobuffalo/envy"
)

var DB_DSN = envy.Get("DB_DSN", "")
var DISCORD_TOKEN = envy.Get("DISCORD_TOKEN", "")
var ENV = envy.Get("GO_ENV", "development")
var BOT_EMBED_THUMBNAIL_URL = envy.Get("BOT_EMBED_THUMBNAIL_URL", "")
var BOT_EMBED_COLOR = mustParseInt(envy.Get("BOT_EMBED_COLOR", ""))
var BOT_URL = envy.Get("BOT_URL", "")

func mustParseInt(s string) int {
	d, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		panic(fmt.Errorf("parse int: %w", err))
	}
	return int(d)
}
