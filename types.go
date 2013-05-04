// tweetlib - A fully oauth-authenticated Go Twitter library
//
// Copyright 2011 The Tweetlib Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package tweetlib

import (
	"fmt"
	"net/url"
)

type Tweet struct {
	Contributors        string `json:"contributors"`
	User                *User  `json:"user"`
	Truncated           bool   `json:"truncated"`
	Text                string `json:"text"`
	InReplyToScreenName string `json:"in_reply_to_screen_name"`
	RetweetCount        int64  `json:"retweet_count"`
	Entities            struct {
		Urls []struct {
			DisplayUrl  string  `json:"display_url"`
			Indices     []int64 `json:"indices"`
			Url         string  `json:"url"`
			ExpandedUrl string  `json:"expanded_url"`
		} `json:"urls"`
		Hashtags []struct {
			Text    string  `json:"text"`
			Indices []int64 `json:"indices"`
		} `json:"hashtags"`
		UserMentions []struct {
			ScreenName string  `json:"screen_name"`
			Name       string  `json:"name"`
			Indices    []int64 `json:"indices"`
			IdStr      string  `json:"id_str"`
			Id         int64   `json:"id"`
		} `json:"user_mentions"`
	} `json:"entities"`
	Geo                  string `json:"geo"`
	InReplyToUserId      int64  `json:"in_reply_to_user_id"`
	IdStr                string `json:"id_str"`
	CreatedAt            string `json:"created_at"`
	Source               string `json:"source"`
	Id                   int64  `json:"id"`
	InReplyToStatusId    string `json:"in_reply_to_status_id"`
	PossiblySensitive    bool   `json:"possibly_sensitive"`
	Retweeted            bool   `json:"retweeted"`
	InReplyToUserIdStr   string `json:"in_reply_to_user_id_str"`
	Coordinates          string `json:"coordinates"`
	Favorited            bool   `json:"favorited"`
	Place                string `json:"place"`
	InReplyToStatusIdStr string `json:"in_reply_to_status_id_str"`
}

type TweetList []Tweet

type User struct {
	ScreenName                string `json:"screen_name"`
	ListedCount               int64  `json:"listed_count"`
	FollowersCount            int64  `json:"followers_count"`
	Location                  string `json:"location"`
	ProfileBackgroundImageUrl string `json:"profile_background_image_url"`
	Name                      string `json:"name"`
	Notifications             bool   `json:"notifications"`
	Protected                 bool   `json:"protected"`
	IdStr                     string `json:"id_str"`
	ProfileBackgroundColor    string `json:"profile_background_color"`
	CreatedAt                 string `json:"created_at"`
	Url                       string `json:"url"`
	TimeZone                  string `json:"time_zone"`
	Id                        int64  `json:"id"`
	Verified                  bool   `json:"verified"`
	ProfileLinkColor          string `json:"profile_link_color"`
	ProfileImageUrl           string `json:"profile_image_url"`
	Status                    *Tweet `json:"status"`
	ProfileUseBackgroundImage bool   `json:"profile_use_background_image"`
	FavouritesCount           int64  `json:"favourites_count"`
	ProfileSidebarFillColor   string `json:"profile_sidebar_fill_color"`
	UtcOffset                 int64  `json:"utc_offset"`
	IsTranslator              bool   `json:"is_translator"`
	FollowRequestSent         bool   `json:"follow_request_sent"`
	Following                 bool   `json:"following"`
	ProfileBackgroundTile     bool   `json:"profile_background_tile"`
	ShowAllInlineMedia        bool   `json:"show_all_inline_media"`
	ProfileTextColor          string `json:"profile_text_color"`
	Lang                      string `json:"lang"`
	StatusesCount             int64  `json:"statuses_count"`
	ContributorsEnabled       bool   `json:"contributors_enabled"`
	FriendsCount              int64  `json:"friends_count"`
	GeoEnabled                bool   `json:"geo_enabled"`
	Description               string `json:"description"`
	ProfileSidebarBorderColor string `json:"profile_sidebar_border_color"`
}

type UserList []User

type MediaEntity struct {
	Id_str          string
	Media_url       string
	Media_url_https string
}

type UrlEntity struct {
	Expanded_url string
}

type Entities struct {
	Media *[]MediaEntity
	Urls  *[]UrlEntity
}

type Configuration struct {
	CharactersReservedPerMedia int                  `json:"characters_reserved_per_media"`
	MaxMediaPerUpload          int                  `json:"max_media_per_upload"`
	NonUsernamePaths           []string             `json:"non_username_paths"`
	PhotoSizeLimit             int                  `json:"photo_size_limit"`
	PhotoSizes                 map[string]PhotoSize `json:"photo_sizes"`
	ShortUrlLengthHttps        int                  `json:"short_url_length_https"`
	ShortUrlLength             int                  `json:"short_url_length"`
}

type PhotoSize struct {
	Width  int    `json:"w"`
	Height int    `json:"h"`
	Resize string `json:"resize"`
}

type Message struct {
	CreatedAt           string `json:"created_at"`
	SenderScreenName    string `json:"sender_screen_name"`
	Sender              *User  `json:"sender"`
	Text                string `json:"text"`
	RecipientScreenName string `json:"recipient_screen_name"`
	Id                  string `json:"id"`
	Recipient           *User  `json:"recipient"`
	RecipientId         string `json:"recipient_id"`
	SenderId            string `json:"sender_id"`
}

type MessageList []Message

type Category struct {
	Name  string    `json:"name"`
	Slug  string    `json:"slug"`
	Size  int       `json:"size"`
	Users *UserList `json:"users"`
}

type CategoryList []Category

type LimitStatus struct {
	RemainingHits      int          `json:"remaining_hits"`
	ResetTimeInSeconds int          `json:"reset_time_in_secods"`
	HourlyLimit        int          `json:"hourly_limit"`
	ResetTime          string       `json:"reset_time"`
	Photos             *PhotoLimits `json:"photos"`
}

type PhotoLimits struct {
	RemainingHits      int    `json:"remaining_hits"`
	ResetTimeInSeconds int    `json:"reset_time_in_secods"`
	ResetTime          string `json:"reset_time"`
	DailyLimit         int    `json:"daily_limit"`
}

type List struct {
	User            *User  `json:"user"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	Mode            string `json:"mode"`
	IdStr           string `json:"id_str"`
	Uri             string `json:"uri"`
	Id              int    `json:"id"`
	MemberCount     int    `json:"member_count"`
	Following       bool   `json:"following"`
	FullName        string `json:"full_name"`
	SubscriberCount int    `json:"subscriber_count"`
	Description     string `json:"description"`
}

type ListList []List

type TrendLocation struct {
	Woeid       int64  `json:"woeid"`
	Name        string `json:"name"`
	CountryCode string `json:"countryCode"`
	Country     string `json:"country"`
	Url         string `json:"url"`
	PlaceType   struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
	} `json:"placeType"`
}

type TrendLocationList []TrendLocation

// Settings for the authenticated account
type AccountSettings struct {
	AlwaysUseHttps      bool   `json:"always_use_https"`
	DiscoverableByEmail bool   `json:"discoverable_by_email"`
	GeoEnabled          bool   `json:"geo_enabled"`
	Language            string `json:"language"`
	Protected           bool   `json:"protected"`
	ScreenName          string `json:"screen_name"`
	ShowAllInlineMedia  bool   `json:"show_all_inline_media"`
	SleepTime           struct {
		Enabled   bool `json:"enabled"`
		EndTime   int
		StartTime int
	} `json:"sleep_time"`
	TimeZone struct {
		Name       string `json:"name"`
		TzinfoName string `json:"tzinfo_name"`
		UtcOffset  int64
	} `json:"time_zone"`
	TrendLocation []struct {
		Country     string `json:"country"`
		CountryCode string `json:"countryCode"`
		Name        string `json:"name"`
		ParentId    int64  `json:"parentid"`
		PlaceType   struct {
			Code int64  `json:"code"`
			Name string `json:"name"`
		}
		Url   string `json:"url"`
		WoeId int64  `json:"woeid"`
	} `json:"trend_location"`
	UseCookiePersonalization bool `json:"use_cookie_personalization"`
}

// Represents Twitter's current resource limits
// See https://dev.twitter.com/docs/rate-limiting/1.1/limits
// Usage:
//	limits, _ := c.Help.Limits()
//	fmt.Printf("App has %d user_timeline calls remaining\n",
//		limits["statuses"]["/statuses/user_timeline"].Remaining)
type Limits struct {
	// For which context these limits are
	// For tweetlib this will always be the user token
	Context struct {
		AccessToken string `json:"access_token"`
	} `json:"rate_limit_context"`

	// Resrouce families are "accounts", "help", etc
	ResourceFamilies map[string]map[string]struct {
		// How many calls remaining for this resource
		Remaining int `json:"remaining"`
		// When the limit will reset (epoch time)
		Reset int64 `json:"reset"`
		// Total number of calls allowed
		Limit int `json:"limit"`
	} `json:"resources"`
}

type LimitResourceFamily map[string]struct {
	Remaining int   `json:"remaining"`
	Reset     int64 `json:"reset"`
	Limit     int   `json:"limit"`
}

// A media attached to a tweet. In practice, this represents
// an image file.
type TweetMedia struct {
	Filename string // Name for the file (e.g. image.png)
	Data     []byte // Raw file data
}

// Optionals: used to provide optional arguments to
// API calls
//
// Usage:
//
//   opts := NewOptionals()
//   opts.Add("count", 10)
//   opts.Add("lat", -37.102013)
//   opts.Add("user_id", "twitteruser")
type Optionals struct {
	Values url.Values
}

// NewOptionals returns a new instance of Optionals
func NewOptionals() *Optionals {
	return &Optionals{make(url.Values)}
}

// Add: adds a new optional parameter to be used in
// an API request. The value needs to be "stringified"
func (o *Optionals) Add(name string, value interface{}) {
	o.Values.Add(name, fmt.Sprintf("%v", value))
}
