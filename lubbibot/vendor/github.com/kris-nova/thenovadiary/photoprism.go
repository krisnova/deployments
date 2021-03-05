package thenovadiary

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/kris-nova/logger"
	"github.com/kris-nova/photoprism-client-go"
	"github.com/kris-nova/photoprism-client-go/api/v1"
)

type CustomData struct {
	LastTweet   *time.Time
	NoteStrings []string
	KeyValue    map[string]string
	Description string
}

func SetCustomData(d *CustomData, photo *api.Photo) error {
	jBytes, err := json.Marshal(&d)
	if err != nil {
		return fmt.Errorf("unable to set custom photo data: %v", err)
	}
	photo.PhotoDescription = string(jBytes)
	return nil
}

func GetCustomData(photo api.Photo) *CustomData {
	d := &CustomData{}
	if photo.PhotoDescription == "" {
		return d
	}
	noteStr := photo.PhotoDescription
	err := json.Unmarshal([]byte(noteStr), &d)
	if err != nil {
		logger.Warning("INVALID JSON in Notes: %v", err)
	}
	return d
}

func FindNextPhotoInAlbum(client *photoprism.Client, albumID string) (*api.Photo, error) {
	// TODO Nóva add pagination
	photos, err := client.V1().GetPhotos(&api.PhotoOptions{
		AlbumUID: albumID,
		Count:    500, // TODO This is an  enormous number holy shit
	})
	if err != nil {
		return nil, fmt.Errorf("unable to list photos: %v", err)
	}
	logger.Debug("Searching photo album: %s", albumID)
	logger.Debug("Found %d photos to process", len(photos))
	return FindPhoto(photos)
}

// FindNextPhoto will list all photos and
// linear search for a photo that was last
// updated with the greatest delta in days
// ago based on the timestamp found in
//
//
// The function will return unprocessed
// photos first by design.
func FindPhoto(photos []api.Photo) (*api.Photo, error) {

	// --------------------------------------------------------------------------
	//
	// Photo processing rules
	//   - Rules are processed linearly (slow) using a linear search for each
	//   - The first rule to match will win the search
	//
	//
	// ------ [ Shuffle Photos ] ------
	photos = shufflePhotos(photos)
	//
	// ------ [ Photo Pointer ] ------
	var photoToTweet *api.Photo
	//
	// ------ [ Favorites ] ------
	if photoToTweet == nil {
		photoToTweet = findFirstFavoritePhoto(photos)
	}
	// ------ [ New ] ------
	if photoToTweet == nil {
		photoToTweet = findFirstNewPhoto(photos)
	}
	// ------ [ Oldest ] ------
	if photoToTweet == nil {
		photoToTweet = findOldestPhoto(photos)
	}
	if photoToTweet == nil {
		return nil, fmt.Errorf(
			"unable to find photo in album check for correct albumid and photos exist")
	}
	return photoToTweet, nil
	//
	//
	// --------------------------------------------------------------------------
}

// The old Fisher–Yates shuffle
func shufflePhotos(photos []api.Photo) []api.Photo {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < len(photos); i++ {
		j := rand.Intn(i + 1)
		photos[i], photos[j] = photos[j], photos[i]
	}
	return photos
}

func findOldestPhoto(photos []api.Photo) *api.Photo {
	today := TimeToday()
	var oldestPhoto *api.Photo
	delta := 0
	for _, photo := range photos {
		data := GetCustomData(photo)
		if data.LastTweet != nil {
			timeFromData := *data.LastTweet
			pDelta := TimeDeltaDays(today, timeFromData)
			logger.Info("[%s] delta(%d) pDelta(%d)", photo.PhotoTitle, delta, pDelta)
			if pDelta > delta {
				oldestPhoto = &photo
			}
		}
	}
	if delta > 0 {
		logger.Debug("Found oldest photo %s with delta %d", oldestPhoto.PhotoTitle, delta)
	}
	return oldestPhoto
}

func findFirstNewPhoto(photos []api.Photo) *api.Photo {
	for _, photo := range photos {
		data := GetCustomData(photo)
		if data.LastTweet == nil {
			logger.Debug("Found first new photo: %s", photo.PhotoTitle)
			return &photo
		}
	}
	return nil
}

func findFirstFavoritePhoto(photos []api.Photo) *api.Photo {
	for _, photo := range photos {
		if photo.PhotoFavorite {
			logger.Debug("Found first favorite photo: %s", photo.PhotoTitle)
			photo.PhotoFavorite = false
			return &photo
		}
	}
	return nil
}

// NewPhotoprismClient will always return a new client with a fresh token
func (d *Diary) NewPhotoprismClient() (*photoprism.Client, error) {
	mask := ""
	for i := 0; i < len(d.config.PhotoprismPass); i++ {
		mask = fmt.Sprintf("%s%s", mask, "*")
	}
	logger.Debug("Photoprism client (%s) %s: %s", d.config.PhotoprismConn, d.config.PhotoprismUser, mask)
	client := photoprism.New(d.config.PhotoprismConn)
	err := client.Auth(photoprism.NewClientAuthLogin(d.config.PhotoprismUser, d.config.PhotoprismPass))
	if err != nil {
		return nil, fmt.Errorf("unable to auth with photoprism client: %v", err)
	}
	return client, nil
}
