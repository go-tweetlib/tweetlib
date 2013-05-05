// tweetlib - A fully oauth-authenticated Go Twitter library
//
// Copyright 2011 The Tweetlib Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tweetlib

type ListService struct {
	*Client
}

// Returns all lists the authenticating or specified user subscribes to,
// including their own.
// See https://dev.twitter.com/docs/api/1.1/get/lists/list
func (ls *ListService) GetAll(screenName string, opts *Optionals) (lists *ListList, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	opts.Add("screen_name", screenName)
	lists = &ListList{}
	err = ls.Call("GET", "lists/list", opts, lists)
	return
}

// https://dev.twitter.com/docs/api/1.1/get/lists/statuses
