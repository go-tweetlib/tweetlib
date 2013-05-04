// tweetlib - A fully oauth-authenticated Go Twitter library
//
// Copyright 2011 The Tweetlib Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tweetlib

// Groups search functionality
type SearchGroup struct {
	*Client
}

// Returns a collection of relevant Tweets matching a specified query.
// See https://dev.twitter.com/docs/api/1.1/get/search/tweets
func (sg *SearchGroup) Tweets(q string, opts *Optionals) (tweets *TweetList, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	opts.Add("q", q)
	tweets = &TweetList{}
	err = sg.Call("GET", "search/tweets", opts, tweets)
	return
}
