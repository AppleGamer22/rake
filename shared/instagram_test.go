package shared_test

import (
	"net/url"
	"regexp"
	"testing"

	"github.com/AppleGamer22/rake/shared"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var filePathRegularExpression = regexp.MustCompile(`\.(jpg)|(webp)|(mp4)|(webm)`)
var instagramDomainRegularExpression = regexp.MustCompile(`(cdninstagram\.com)|(fbcdn\.net)`)

func init() {
	viper.SetConfigName(".rake")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("..")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func TestInstagramSingleImage(t *testing.T) {
	instagram := shared.NewInstagram(viper.GetString("fbsr"), viper.GetString("session"), viper.GetString("app"))
	URLs, username, err := instagram.Post("CbgDyqkFBdj")
	assert.NoError(t, err)
	assert.Equal(t, "wikipedia", username)
	assert.Len(t, URLs, 1)
	URL, err := url.Parse(URLs[0])
	assert.NoError(t, err)
	assert.Equal(t, "https", URL.Scheme)
	assert.Regexp(t, instagramDomainRegularExpression, URL.Host, URLs[0])
	assert.Regexp(t, filePathRegularExpression, URL.Path)
}

func TestInstagramSingleVideo(t *testing.T) {
	instagram := shared.NewInstagram(viper.GetString("fbsr"), viper.GetString("session"), viper.GetString("app"))
	URLs, username, err := instagram.Post("BKyN0E2AApX")
	assert.NoError(t, err)
	assert.Equal(t, "wikipedia", username)
	assert.Len(t, URLs, 1)
	URL, err := url.Parse(URLs[0])
	assert.NoError(t, err)
	assert.Equal(t, "https", URL.Scheme)
	assert.Regexp(t, instagramDomainRegularExpression, URL.Host, URLs[0])
	assert.Regexp(t, filePathRegularExpression, URL.Path)
}

func TestInstagramBundleImages(t *testing.T) {
	instagram := shared.NewInstagram(viper.GetString("fbsr"), viper.GetString("session"), viper.GetString("app"))
	URLs, username, err := instagram.Post("CZNJeAil1BC")
	assert.NoError(t, err)
	assert.Equal(t, "wikipedia", username)
	assert.Len(t, URLs, 2)
	for _, urlString := range URLs {
		URL, err := url.Parse(urlString)
		assert.NoError(t, err)
		assert.Equal(t, "https", URL.Scheme)
		assert.Regexp(t, instagramDomainRegularExpression, URL.Host, urlString)
		assert.Regexp(t, filePathRegularExpression, URL.Path)
	}
}
