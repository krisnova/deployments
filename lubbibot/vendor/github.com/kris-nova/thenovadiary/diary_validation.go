package thenovadiary

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

type validationCheck func(cfg *DiaryConfig) string

var (
	// --------------------------------------
	// Define the hard coded checks for the
	// diary configuration here
	//
	// TODO Use reflection to pre-populate all
	// functions that start with "v_"
	validationChecks = []validationCheck{
		v_emptyPhotoprismAlbum,
		v_emptyPhotoprismUser,
		v_emptyPhotoprismPass,
		v_emptyPhotoprismConn,
		v_emptyTwitterConsumerKey,
		v_emptyTwitterConsumerKeySecret,
		v_emptyTwitterToken,
		v_emptyTwitterTokenSecret,
	}
)

func ValidateConfig(cfg *DiaryConfig) error {
	var strerr string
	for _, v := range validationChecks {
		e := v(cfg)
		if e != "" {
			strerr = fmt.Sprintf("%s%s", strerr, e)
		}
	}
	if strerr != "" {
		return fmt.Errorf("Errors during validation: %s", strerr)
	}
	// We know we have a valid config
	cfg.validated = true
	return nil
}

func v_emptyTwitterToken(cfg *DiaryConfig) string {
	if cfg.TwitterToken == "" {
		cfg.TwitterToken = os.Getenv("DIARY_TWITTERTOKEN")
	}
	if cfg.TwitterToken == "" {
		return ferr("Empty TwitterToken")
	}
	if len(cfg.TwitterToken) < 3 {
		return ferr("TwitterToken < 3 chars")
	}
	return ""
}

func v_emptyTwitterTokenSecret(cfg *DiaryConfig) string {
	if cfg.TwitterTokenSecret == "" {
		cfg.TwitterTokenSecret = os.Getenv("DIARY_TWITTERTOKENSECRET")
	}
	if cfg.TwitterTokenSecret == "" {
		return ferr("Empty TwitterTokenSecret")
	}
	if len(cfg.TwitterTokenSecret) < 3 {
		return ferr("TwitterTokenSecret < 3 chars")
	}
	return ""
}

func v_emptyTwitterConsumerKey(cfg *DiaryConfig) string {
	if cfg.TwitterConsumerKey == "" {
		cfg.TwitterConsumerKey = os.Getenv("DIARY_TWITTERCONSUMERKEY")
	}
	if cfg.TwitterConsumerKey == "" {
		return ferr("Empty TwitterConsumerKey")
	}
	if len(cfg.TwitterConsumerKey) < 3 {
		return ferr("TwitterConsumerKey < 3 chars")
	}
	return ""
}

func v_emptyTwitterConsumerKeySecret(cfg *DiaryConfig) string {
	if cfg.TwitterConsumerKeySecret == "" {
		cfg.TwitterConsumerKeySecret = os.Getenv("DIARY_TWITTERCONSUMERKEYSECRET")
	}
	if cfg.TwitterConsumerKeySecret == "" {
		return ferr("Empty TwitterConsumerKeySecret")
	}
	if len(cfg.TwitterConsumerKeySecret) < 3 {
		return ferr("TwitterConsumerKeySecret < 3 chars")
	}
	return ""
}

func v_emptyPhotoprismAlbum(cfg *DiaryConfig) string {
	if cfg.PhotoprismAlbum == "" {
		cfg.PhotoprismAlbum = os.Getenv("DIARY_PHOTOPRISMALBUM")
	}
	if cfg.PhotoprismAlbum == "" {
		return ferr("Empty PhotoprismAlbum")
	}
	if len(cfg.PhotoprismAlbum) < 3 {
		return ferr("PhotoprismAlbum < 3 chars")
	}
	return ""
}

func v_emptyPhotoprismUser(cfg *DiaryConfig) string {
	if cfg.PhotoprismUser == "" {
		cfg.PhotoprismUser = os.Getenv("DIARY_PHOTOPRISMUSER")
	}
	if cfg.PhotoprismUser == "" {
		return ferr("Empty PhotoprismUser")
	}
	if len(cfg.PhotoprismUser) < 3 {
		return ferr("PhotoprismUser < 3 chars")
	}
	return ""
}

func v_emptyPhotoprismPass(cfg *DiaryConfig) string {
	if cfg.PhotoprismPass == "" {
		cfg.PhotoprismPass = os.Getenv("DIARY_PHOTOPRISMPASS")
	}
	if cfg.PhotoprismPass == "" {
		return ferr("Empty PhotoprismPass")
	}
	if len(cfg.PhotoprismPass) < 3 {
		return ferr("PhotoprismPass < 3 chars")
	}
	return ""
}

func v_emptyPhotoprismConn(cfg *DiaryConfig) string {
	if cfg.PhotoprismConn == "" {
		cfg.PhotoprismConn = os.Getenv("DIARY_PHOTOPRISMCONN")
	}
	if cfg.PhotoprismConn == "" {
		return ferr("Empty PhotoprismConn")
	}
	if len(cfg.PhotoprismConn) < 3 {
		return ferr("PhotoprismConn < 3 chars")
	}
	return ""
}

func ferr(str string) string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next() // Ignoring next bool (This is error safe)
	fqn := frame.Func.Name()
	spl := strings.Split(fqn, ".")
	if len(spl) < 3 {
		return fmt.Sprintf("[%s] %s ", frame.Func.Name(), str)
	}
	return fmt.Sprintf("[%s] %s ", spl[2], str)
}
