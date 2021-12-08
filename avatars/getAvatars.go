package avatars

import (
	"encoding/xml"
	"fmt"
	"github.com/MrWaggel/gosteamconv"
	"io/ioutil"
	"log"
	"net/http"
)

type Profile struct {
	XMLName    xml.Name `xml:"profile"`
	AvatarIcon string   `xml:"avatarIcon"`
}

func GetAvatar(id uint32) (string, error) {
	steamString, err := gosteamconv.SteamInt32ToString(int32(id))
	if err != nil {
		return "", err
	}
	steamID64, err := gosteamconv.SteamStringToInt64(steamString)
	if err != nil {
		return "", err
	}

	res, err := http.Get(fmt.Sprintf("https://steamcommunity.com/profiles/%d/?xml=1", steamID64))
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var profile Profile

	err = xml.Unmarshal(content, &profile)
	if err != nil {
		return "", err
	}

	if len(profile.AvatarIcon) < 45 {
		return "", fmt.Errorf("wrong link: %s for user: %d", profile.AvatarIcon, id)
	}
	hash := profile.AvatarIcon[len(profile.AvatarIcon)-44:]

	return hash[:40], nil
}
