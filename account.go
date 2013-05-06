// tweetlib - A fully oauth-authenticated Go Twitter library
//
// Copyright 2011 The Tweetlib Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tweetlib

import (
	"encoding/base64"
)

// Groups account-related functions
type AccountService struct {
	*Client
}

// Returns settings (including current trend, geo and sleep time information)
// for the authenticating user
// See https://dev.twitter.com/docs/api/1.1/get/account/settings
func (ag *AccountService) Settings() (settings *AccountSettings, err error) {
	settings = &AccountSettings{}
	err = ag.Call("GET", "account/settings", nil, settings)
	return
}

// Helper function to verify if credentials are valid. Returns the
// user object if they are.
// See https://dev.twitter.com/docs/api/1.1/get/account/verify_credentials
func (ag *AccountService) VerifyCredentials(opts *Optionals) (user *User, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	user = &User{}
	err = ag.Call("GET", "account/verify_credentials", opts, user)
	return
}

// Update authenticating user's settings.
// See https://dev.twitter.com/docs/api/1.1/post/account/settings
func (ag *AccountService) UpdateSettings(opts *Optionals) (newSettings *AccountSettings, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	newSettings = &AccountSettings{}
	err = ag.Call("POST", "account/settings", opts, newSettings)
	return
}

// Enables/disables SMS delivery
// See https://dev.twitter.com/docs/api/1.1/post/account/update_delivery_device
func (ag *AccountService) EnableSMS(enable bool) (err error) {
	opts := NewOptionals()
	if enable {
		opts.Add("device", "sms")
	} else {
		opts.Add("device", "none")
	}
	err = ag.Call("POST", "account/update_delivery_device", opts, nil)
	return
}

// Sets values that users are able to set under the "Account" tab of their
// settings page. Only the parameters specified will be updated.
// See https://dev.twitter.com/docs/api/1.1/post/account/update_profile
func (ag *AccountService) UpdateProfile(opts *Optionals) (user *User, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	user = &User{}
	err = ag.Call("POST", "account/update_profile", opts, user)
	return
}

// Updates the authenticating user's profile background image.
// Passing an empty []byte as image will disable the current
// background image.
// https://dev.twitter.com/docs/api/1.1/post/account/update_profile_background_image
func (ag *AccountService) UpdateProfileBackgroundImage(image []byte, opts *Optionals) (user *User, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	if len(image) > 0 {
		opts.Add("image", base64.StdEncoding.EncodeToString(image))
		opts.Add("use", true)
	} else {
		opts.Add("use", false)
	}
	user = &User{}
	err = ag.Call("POST", "account/update_profile_background_image", opts, user)
	return

}

// Sets one or more hex values that control the color scheme of the
// authenticating user's profile page on twitter.com. Each parameter's value
// must be a valid hexidecimal value, and may be either three or six characters
// (ex: #fff or #ffffff).
func (ag *AccountService) UpdateProfileColors(opts *Optionals) (user *User, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	user = &User{}
	err = ag.Call("POST", "account/update_profile_colors", opts, user)
	return
}

// Updates the authenticating user's profile image. The image parameter should
// be the raw data from the image file, not a path or URL
// See https://dev.twitter.com/docs/api/1.1/post/account/update_profile_image
func (ag *AccountService) UpdateProfileImage(image []byte, opts *Optionals) (user *User, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	opts.Add("image", base64.StdEncoding.EncodeToString(image))
	user = &User{}
	err = ag.Call("POST", "account/update_profile_image", opts, user)
	return
}

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
