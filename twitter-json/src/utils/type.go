package utils

import twitterscraper "github.com/n0madic/twitter-scraper"

type (
	// DbProfile type.
	DbProfile struct {
		Username       string
		Profile        twitterscraper.Profile
		LastTweetTime  int64
		LastUpdateTime int64
	}
	// Twitter type.
	Twitter struct {
		DbProfile
		Tweets []twitterscraper.Tweet
	}
	// Settings type default value is used by dev.
	Settings struct {
		PicBed       string `env:"PIC_BED" envDefault:"github"`
		UseWebhook   bool   `env:"USE_WEBHOOK" envDefault:"false"`
		DbInit       bool   `env:"DB_INIT" envDefault:"false"`
		SettingsPath string `env:"SETTINGS_PATH" envDefault:"../settings.json"` // 与 main.go 的相对位置
	}
	// DbImage image info in db
	DbImage struct {
		Key  string
		User string
	}
)
