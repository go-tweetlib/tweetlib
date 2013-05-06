// tweetlib - A fully oauth-authenticated Go Twitter library
//
// Copyright 2011 The Tweetlib Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tweetlib

// Groups help functions
type HelpService struct {
	*Client
}

// Returns the current configuration used by Twitter including twitter.com
// slugs which are not usernames, maximum photo resolutions, and t.co URL
// lengths.
// See https://dev.twitter.com/docs/api/1.1/get/help/configuration
func (hs *HelpService) Configuration() (configuration *Configuration, err error) {
	configuration = &Configuration{}
	err = hs.Call("GET", "help/configuration", nil, configuration)
	return
}

// Returns Twitter's Privacy Policy
// Seehttps://dev.twitter.com/docs/api/1.1/get/help/privacy
func (hs *HelpService) PrivacyPolicy() (privacyPolicy string, err error) {
	type pp struct {
		Text string `json:"privacy"`
	}
	ret := &pp{}
	err = hs.Call("GET", "help/privacy", nil, ret)
	privacyPolicy = ret.Text
	return
}

// Returns Twitter's terms of service
// See https://dev.twitter.com/docs/api/1.1/get/help/tos
func (hs *HelpService) Tos() (string, error) {
	type tos struct {
		Text string `json:"tos"`
	}
	ret := &tos{}
	err := hs.Call("GET", "help/tos", nil, ret)
	return ret.Text, err
}

// Returns current Twitter's rate limits
// See https://dev.twitter.com/docs/api/1.1/get/application/rate_limit_status
func (hs *HelpService) Limits() (limits *Limits, err error) {
	limits = &Limits{}
	err = hs.Call("GET", "application/rate_limit_status", nil, limits)
	return
}
