// tweetlib - A fully oauth-authenticated Go Twitter library
//
// Copyright 2011 The Tweetlib Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tweetlib

import "strconv"

type FriendsService struct {
	*Client
}

type Cursor struct {
	Next        int64   `json:"next_cursor"`
	NextStr     string  `json:"next_cursor_str"`
	Previous    int64   `json:"previous_cursor"`
	PreviousStr string  `json:"previous_cursor_str"`
	IDs         []int64 `json:"ids"`
}

// IDs returns a cursored collection of user IDs.
// See https://dev.twitter.com/docs/api/1.1/get/friends/ids
func (ls *FriendsService) IDs(screenName string, userID int64, cursor int64, opts *Optionals) (IDs *Cursor, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	switch {
	case screenName != "":
		opts.Add("screen_name", screenName)
	default:
		opts.Add("user_id", userID)
	}
	if cursor != 0 {
		opts.Add("cursor", strconv.FormatInt(cursor, 10))
	}
	IDs = &Cursor{}
	err = ls.Call("GET", "friends/ids", opts, IDs)
	return
}

type FollowersService struct {
	*Client
}

// IDs returns a cursored collection of user IDs.
// See https://dev.twitter.com/docs/api/1.1/get/followers/ids
func (ls *FollowersService) IDs(screenName string, userID int64, cursor int64, opts *Optionals) (IDs *Cursor, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	switch {
	case screenName != "":
		opts.Add("screen_name", screenName)
	default:
		opts.Add("user_id", userID)
	}
	if cursor != 0 {
		opts.Add("cursor", strconv.FormatInt(cursor, 10))
	}
	IDs = &Cursor{}
	err = ls.Call("GET", "followers/ids", opts, IDs)
	return
}

// https://dev.twitter.com/docs/api/1.1/get/followers/ids
