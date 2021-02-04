package utils

import twitterscraper "github.com/n0madic/twitter-scraper"

type (
	// UserInfo user info type.
	UserInfo struct {
		Avatar         string
		Username       string
		Name           string
		LastTweetTime  int64
		LastUpdateTime int64
		TweetsCount    int64
		Pages          int
	}
	// DbProfile type.
	DbProfile struct {
		UserInfo
		Profile twitterscraper.Profile
	}
	// Twitter type.
	Twitter struct {
		DbProfile
		Tweets []twitterscraper.Tweet
	}
	// Settings type default value is used by dev.
	Settings struct {
		PicBed       string `env:"PIC_BED" envDefault:"github"`                 // 图床: github, qiniu, smms
		UseWebhook   bool   `env:"USE_WEBHOOK" envDefault:"false"`              // 是否用 webhook 触发 actions
		DbInit       bool   `env:"DB_INIT" envDefault:"false"`                  // 数据库是否已初始化
		PageSize     int64  `env:"PAGE_SIZE" envDefault:"10"`                   // 每页推特数量
		SettingsPath string `env:"SETTINGS_PATH" envDefault:"../settings.json"` // 与 main.go 的相对位置
	}
	// DbImage image info in db
	DbImage struct {
		PicBed   string
		FileName string
		Key      string
		User     string
		URL      string
		Type     string // photo video or profile
		Idx      int    // photo/video: index of photos/videos  profile: 0 -> avatar 1 ->banner
		Status   string // upload status: success or fail
		TweetID  string // if Type is tweet use tweetid for upload fail
	}
	// UpdateInfo update info
	UpdateInfo struct {
		UpdateTime int64
		IsUpdate   bool
		Users      []string
	}
	// UserChangeInfo users change info
	UserChangeInfo struct {
		Type  string `env:"EVENT_TYPE" envDefault:""`
		Users string `env:"USERS_DATA" envDefault:""`
	}
)
