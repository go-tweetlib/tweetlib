// tweetlib - A fully oauth-authenticated Go Twitter library
//
// Copyright 2011 The Tweetlib Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tweetlib

import (
	"errors"
	"strconv"
	"strings"
)

type UserService struct {
	*Client
}

// Represents a Twitter user. This struct is
// often part of a tweet. Depending on the API
// call used, this struct may not be fully filled.
// Some calls use only the Id field.
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

// A list of users
type UserList []User

// Provides a simple, relevance-based search interface to public user accounts
// on Twitter. Try querying by topical interest, full name, company name,
// location, or other criteria. Exact match searches are not supported.
// See https://dev.twitter.com/docs/api/1.1/get/users/search
func (us *UserService) Search(q string, opts *Optionals) (users *UserList, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	opts.Add("q", q)
	users = &UserList{}
	err = us.Call("GET", "users/search", opts, users)
	return
}

// See https://dev.twitter.com/docs/api/1.1/get/users/show
func (us *UserService) Show(screenName string, opts *Optionals) (user *User, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	if screenName != "" {
		opts.Add("screen_name", screenName)
	}
	user = &User{}
	err = us.Call("GET", "users/show", opts, user)
	return
}

// See https://dev.twitter.com/docs/api/1.1/get/users/lookup
func (us *UserService) Lookup(screenNames []string, userIDs []int64, opts *Optionals) (users *UserList, err error) {
	if opts == nil {
		opts = NewOptionals()
	}

	var n int
	switch {
	case screenNames != nil && len(screenNames) <= 100:
		opts.Add("screen_name", strings.Join(screenNames, ","))
		n = len(screenNames)
	case userIDs != nil && len(userIDs) <= 100:
		n = len(userIDs)
		ids := make([]string, n)
		for k, v := range userIDs {
			ids[k] = strconv.FormatInt(v, 10)
		}
		opts.Add("user_id", strings.Join(ids, ","))
	default:
		return nil, errors.New("Invalid request.")
	}

	users = &UserList{}
	err = us.Call("POST", "users/lookup", opts, users)
	return
}
