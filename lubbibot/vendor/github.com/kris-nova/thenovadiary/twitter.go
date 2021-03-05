package thenovadiary

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/kris-nova/photoprism-client-go/api/v1"
)

const (

	// Default Twitter Values
	// These are constants and if
	// we every change our app account
	// we should update here
	DefaultTwitterHandle = "thenovadiary"
	DefaultTwitterURL    = "https://twitter.com/thenovadiary"

	// The twitter API now accepts tweets with 280 characters
	TwitterLengthLimit = 280
)

// SendPhotoTweet will return the URL of the sent tweet and/or an
// error.
func SendPhotoTweet(twitter *anaconda.TwitterApi, photo api.Photo, pBytes []byte) (string, error) {
	if BypassSendDailyTweetTwitter == true {
		return "", nil
	}
	// TODO We should be memory conscience here
	b64str := string(b64e(pBytes))

	// ------ [ Upload Photo ] ------
	media, err := twitter.UploadMedia(b64str)
	if err != nil {
		return "", fmt.Errorf("unable to upload photo to twitter: %s", err)
	}
	// ------ [ Send Tweet ] ------
	v := url.Values{}
	v.Set("media_ids", media.MediaIDString)
	tweet, err := twitter.PostTweet(GetStatus(photo), v)
	if err != nil {
		return DefaultTwitterURL, fmt.Errorf("unable to tweet: %v", err)
	}
	url := fmt.Sprintf("https://twitter.com/%s/status/%s", DefaultTwitterHandle, tweet.IdStr)
	return url, nil
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

func GetStatus(photo api.Photo) string {
	var metaFields []string
	if photo.PhotoFocalLength != 0 {
		metaFields = append(metaFields, fmt.Sprintf("FocalLength(%dmm)", photo.PhotoFocalLength))
	}
	if photo.PhotoIso != 0 {
		metaFields = append(metaFields, fmt.Sprintf("ISO(%d)", photo.PhotoIso))
	}
	if photo.PhotoExposure != "" {
		metaFields = append(metaFields, fmt.Sprintf("Exposure(%s)", photo.PhotoExposure))
	}
	if photo.PhotoFNumber != 0 {
		metaFields = append(metaFields, fmt.Sprintf("F(%v)", int(photo.PhotoFNumber)))
	}
	var status string
	meta := strings.Join(metaFields, " | ")
	if meta == "" {
		status = photo.PhotoTitle
	} else {
		status = fmt.Sprintf("%s\n\n[%s]", photo.PhotoTitle, meta)
	}
	if len(status) > TwitterLengthLimit {
		if len(photo.PhotoTitle) < TwitterLengthLimit {
			return photo.PhotoTitle
		}

		// Limit string by length
		truncated := ""
		for i := 0; i < TwitterLengthLimit-3; i++ {
			charRune := status[i]
			truncated = fmt.Sprintf("%s%c", truncated, charRune)
		}
		return fmt.Sprintf("%s...", truncated)
	}
	return status
}
