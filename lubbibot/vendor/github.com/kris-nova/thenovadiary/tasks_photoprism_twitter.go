package thenovadiary

import (
	"fmt"
	"time"

	"github.com/ChimeraCoder/anaconda"

	"github.com/kris-nova/logger"
)

const (
	DailyTweetCacheKey = "dailytweet"
)

func (d *Diary) SendDailyTweet() error {
	// ------ [ Check if we need to tweet ] ------
	today := TimeToday()
	cachedRecord := d.cache.Get(DailyTweetCacheKey)
	if !cachedRecord.Found {
		cachedRecord.Value = TimeDeltaDaysFromNow(-100)
	}
	lastTweetTime, err := time.Parse(TimeLayoutTimeTime, fmt.Sprintf("%v", cachedRecord.Value))
	if err != nil {
		// Unable to pull time from record
		return fmt.Errorf("unable to read time from cache: %v %v", cachedRecord.Value, err)
	}
	lastTweetDeltaDays := TimeDeltaDays(today, lastTweetTime)

	// ------ [ Cron Logic] ------
	//
	// This is a daily tweet, so we check days = 1
	//
	if lastTweetDeltaDays < 1 {
		// This is the silent edge case
		// we will hit the majority of
		// the time
		return nil
	}
	logger.Debug("Days since last tweet: %d", lastTweetDeltaDays)

	// ------- [ Get Photo From Photo Search Algorithm ] -----
	pClient, err := d.NewPhotoprismClient()
	if err != nil {
		return fmt.Errorf("unable to build new photoprism client: %v", err)
	}
	photoPtr, err := FindNextPhotoInAlbum(pClient, d.config.PhotoprismAlbum)
	if err != nil {
		return err
	}
	photo := *photoPtr
	logger.Debug("Processing photo: %s", photo.UUID)

	// ------- [ Get Photo Content ] ------
	pBytes, err := pClient.V1().GetPhotoDownload(photo.PhotoUID)
	if err != nil {
		return fmt.Errorf("unable to download photo: %v", err)
	}
	logger.Debug("Downloaded photo %d kb", len(pBytes)/1024)

	// ------- [ Auth Twitter Client ] ------
	tURL := DefaultTwitterURL
	if !BypassSendDailyTweetTwitter {
		cfg := d.config
		tClient := anaconda.NewTwitterApiWithCredentials(cfg.TwitterToken, cfg.TwitterTokenSecret, cfg.TwitterConsumerKey, cfg.TwitterConsumerKeySecret)
		tURL, err = SendPhotoTweet(tClient, photo, pBytes)
		if err != nil {
			return fmt.Errorf("unable to send tweet: %v", err)
		}
	}

	// ------- [ Photo Found, Tweet Sent (lol), Update Cache ] ------
	if !BypassSendDailyTweetCache {
		logger.Debug("Syncing photo: %s", photo.PhotoUID)
		data := GetCustomData(photo)
		data.LastTweet = &today
		err = SetCustomData(data, &photo)
		if err != nil {
			return err
		}
		_, err = pClient.V1().UpdatePhoto(photo)
		if err != nil {
			return fmt.Errorf("unable to update photoprism photo: %v", err)
		}
		// ------- [ Update Records ] -------
		cachedRecord.Value = today.String()
		d.cache.Set(DailyTweetCacheKey, cachedRecord)
		d.cache.Persist()
		logger.Always("Successful Daily Tweet: %s", tURL)
	}
	return nil

}
