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

type Configuration struct {
	CharactersReservedPerMedia int      `json:"characters_reserved_per_media"`
	MaxMediaPerUpload          int      `json:"max_media_per_upload"`
	NonUsernamePaths           []string `json:"non_username_paths"`
	PhotoSizeLimit             int      `json:"photo_size_limit"`
	PhotoSizes                 map[string]struct {
		Width  int    `json:"w"`
		Height int    `json:"h"`
		Resize string `json:"resize"`
	} `json:"photo_sizes"`
	ShortUrlLengthHttps int `json:"short_url_length_https"`
	ShortUrlLength      int `json:"short_url_length"`
}

type LimitStatus struct {
	RemainingHits      int    `json:"remaining_hits"`
	ResetTimeInSeconds int    `json:"reset_time_in_secods"`
	HourlyLimit        int    `json:"hourly_limit"`
	ResetTime          string `json:"reset_time"`
	Photos             struct {
		RemainingHits      int    `json:"remaining_hits"`
		ResetTimeInSeconds int    `json:"reset_time_in_secods"`
		ResetTime          string `json:"reset_time"`
		DailyLimit         int    `json:"daily_limit"`
	} `json:"photos"`
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
