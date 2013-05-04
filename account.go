// tweetlib - A fully oauth-authenticated Go Twitter library
//
// Copyright 2011 The Tweetlib Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tweetlib

import (
	"encoding/base64"
)

type AccountGroup struct {
	*Client
}

// Returns settings (including current trend, geo and sleep time information)
// for the authenticating user
// See https://dev.twitter.com/docs/api/1.1/get/account/settings
func (ag *AccountGroup) AccountSettings() (settings *AccountSettings, err error) {
	settings = &AccountSettings{}
	err = ag.Call("GET", "account/settings", nil, settings)
	return
}

// Helper function to verify if credentials are valid. Returns the
// user object if they are.
// See https://dev.twitter.com/docs/api/1.1/get/account/verify_credentials
func (ag *AccountGroup) VerifyCredentials(opts *Optionals) (user *User, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	user = &User{}
	err = ag.Call("GET", "account/verify_credentials", opts, user)
	return
}

// Update authenticating user's settings.
// See https://dev.twitter.com/docs/api/1.1/post/account/settings
func (ag *AccountGroup) UpdateSettings(opts *Optionals) (newSettings *AccountSettings, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	newSettings = &AccountSettings{}
	err = ag.Call("POST", "account/settings", opts, newSettings)
	return
}

// Enables/disables SMS delivery
// See https://dev.twitter.com/docs/api/1.1/post/account/update_delivery_device
func (ag *AccountGroup) EnableSMS(enable bool) (err error) {
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
func (ag *AccountGroup) UpdateProfile(opts *Optionals) (user *User, err error) {
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
func (ag *AccountGroup) UpdateProfileBackgroundImage(image []byte, opts *Optionals) (user *User, err error) {
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
