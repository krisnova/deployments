package thenovadiary

import "github.com/kris-nova/logger"

const (

	// BypassSendDailyTweetTwitter
	//
	// Bypass hitting the twitter API for the SendDailyTweet task
	BypassSendDailyTweetTwitter = false

	// BBypassSendDailyTweetCache
	//
	// Bypass saving the cache with an updated last tweet time
	BypassSendDailyTweetCache = false
)

func DebugConfig() {
	// ----------------------------------------------------------------------------------------
	//
	//
	logger.Info("[CONFIG SendDailyTweet] BypassTwitter: %t", BypassSendDailyTweetTwitter)
	logger.Info("[CONFIG SendDailyTweet] BypassCacheSave: %t", BypassSendDailyTweetCache)
	if BypassSendDailyTweetTwitter != BypassSendDailyTweetCache {
		logger.Critical("Unusual configuration for SendDailyTweet!!")
	}
	//
	//
	// ----------------------------------------------------------------------------------------
}
