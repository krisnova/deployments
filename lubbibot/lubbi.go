package main

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"

	"github.com/kris-nova/thenovadiary"

	"github.com/ChimeraCoder/anaconda"
	"github.com/kris-nova/logger"
	"github.com/kris-nova/novaarchive/bot"
	"github.com/kris-nova/photoprism-client-go"
)

func main() {
	albumID := os.Getenv("PHOTOPRISMALBUM")
	logger.BitwiseLevel = logger.LogEverything
	logger.Info("Starting bot...")
	robot := bot.NewTwitterBot(bot.NewTwitterBotCredentialsFromEnvironmentalVariables())
	robot.AddKey("/lubbi")
	logger.Info("Setting command /lubbi...")
	robot.SetBufferSizeGBytes(1)
	logger.Info("Setting buffer 1Gb...")
	robot.SetSendTweet(func(api *anaconda.TwitterApi, tweet anaconda.Tweet) error {
		logger.Always("Found tweet: %s", tweet.IdStr)
		// Photoprism Connection
		client := photoprism.New(os.Getenv("PHOTOPRISMCONN"))
		err := client.Auth(photoprism.NewClientAuthLogin(os.Getenv("PHOTOPRISMUSER"), os.Getenv("PHOTOPRISMPASS")))
		if err != nil {
			return fmt.Errorf("unable to auth with photoprism: %v", err)
		}

		// --- Photo ---
		photo, err := thenovadiary.FindNextPhotoInAlbum(client, albumID)
		if err != nil {
			return fmt.Errorf("Unable to FindNextPhotoInAlbum(): %v", err)
		}

		pBytes, err := client.V1().GetPhotoDownload(photo.PhotoUID)
		if err != nil {
			return fmt.Errorf("Unable to download photo: %v", err)
		}

		// --- Upload Photo ---
		b64str := string(b64e(pBytes))
		media, err := api.UploadMedia(b64str)
		if err != nil {
			return fmt.Errorf("Ynable to upload photo to twitter: %v", err)
		}

		// --- Send Tweet ---
		v := url.Values{}
		v.Set("media_ids", media.MediaIDString)
		v.Set("in_reply_to_status_id", tweet.IdStr)
		v.Set("auto_populate_reply_metadata", "true")
		v.Set("display_coordinates", "false") // TODO set 101
		sentTweet, err := api.PostTweet(getStatus(), v)
		if err != nil {
			return fmt.Errorf("unabble to send lubbi tweet: %v", err)
		}
		logger.Always("Sent tweet: https://twitter.com/%s/status/%s", sentTweet.User.ScreenName, sentTweet.IdStr)
		today := thenovadiary.TimeToday()
		data := thenovadiary.GetCustomData(*photo)
		data.LastTweet = &today
		err = thenovadiary.SetCustomData(data, photo)
		if err != nil {
			return err
		}
		_, err = client.V1().UpdatePhoto(*photo)
		if err != nil {
			return fmt.Errorf("unable to update photoprism photo: %v", err)
		}
		return nil
	})
	logger.Info("Setting SendTweet...")
	user, err := robot.Login()
	if err != nil {
		logger.Critical(err.Error())
		os.Exit(1)
	}
	logger.Info("Logged in as @%s (%s)", user.ScreenName, user.Name)
	err = robot.Run()
	if err != nil {
		logger.Critical(err.Error())
		os.Exit(1)
	}
	logger.Info("Running bot...")
	for {
		err := robot.NextError()
		logger.Critical(err.Error())
	}
}

func getStatus() string {
	// TODO get random lubbi string from liblubbi
	return "o hæ..."
}

func b64e(message []byte) []byte {
	b := make([]byte, base64.StdEncoding.EncodedLen(len(message)))
	base64.StdEncoding.Encode(b, message)
	return b
}

func b64d(message []byte) (b []byte, err error) {
	var l int
	b = make([]byte, base64.StdEncoding.DecodedLen(len(message)))
	l, err = base64.StdEncoding.Decode(b, message)
	if err != nil {
		return
	}
	return b[:l], nil
}
