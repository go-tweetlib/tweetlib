// tweetlib - A fully oauth-authenticated Go Twitter library
//
// Copyright 2011 The Tweetlib Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tweetlib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
)

type TweetsGroup struct {
	*Client
}

// Returns the 20 (by default) most recent tweets containing a users's
// @screen_name for the authenticating user.
// THis method can only return up to 800 tweets (via the "count" optional
// parameter.
// See https://dev.twitter.com/docs/api/1.1/get/statuses/mentions_timeline
func (tg *TweetsGroup) Mentions(opts *Optionals) (tweets *TweetList, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	tweets = &TweetList{}
	err = tg.Call("GET", "statuses/mentions_timeline", opts, tweets)
	return
}

// Returns a collection of the most recent Tweets posted by the user indicated
// by the screen_name.
// See https://dev.twitter.com/docs/api/1.1/get/statuses/user_timeline
func (tg *TweetsGroup) UserTimeline(screenname string, opts *Optionals) (tweets *TweetList, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	opts.Add("screen_name", screenname)
	tweets = new(TweetList)
	err = tg.Call("GET", "statuses/user_timeline", opts, tweets)
	return
}

// Returns a collection of the most recent Tweets and retweets posted by
// the authenticating user and the users they follow.
// See https://dev.twitter.com/docs/api/1.1/get/statuses/home_timeline
func (tg *TweetsGroup) HomeTimeline(opts *Optionals) (tweets *TweetList, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	tweets = new(TweetList)
	err = tg.Call("GET", "statuses/home_timeline", opts, tweets)
	return
}

// Returns a collection of the  most recent tweets authored by the
// authenticating user that have been retweeted by others.
// See https://dev.twitter.com/docs/api/1.1/get/statuses/retweets_of_me
func (tg *TweetsGroup) RetweetsOfMe(opts *Optionals) (tweets *TweetList, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	tweets = new(TweetList)
	err = tg.Call("GET", "statuses/retweets_of_me", opts, tweets)
	return
}

// Update: posts a status update to Twitter
// See https://dev.twitter.com/docs/api/1.1/post/statuses/update
func (tg *TweetsGroup) Update(status string, opts *Optionals) (tweet *Tweet, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	opts.Add("status", status)
	tweet = &Tweet{}
	err = tg.Call("POST", "statuses/update", opts, tweet)
	return tweet, err
}

// Returns up to 100 of the first retweets of a given tweet Id
func (tg *TweetsGroup) Retweets(id int64, opts *Optionals) (tweets *TweetList, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	tweets = &TweetList{}
	err = tg.Call("GET", fmt.Sprintf("statuses/retweets/%d", id), opts, tweets)
	return
}

// Returns a single Tweet, specified by the id parameter.
// The Tweet's author will also be embedded within the tweet.
func (tg *TweetsGroup) Get(id int64, opts *Optionals) (tweet *Tweet, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	opts.Add("id", id)
	tweet = &Tweet{}
	err = tg.Call("GET", "statuses/show", opts, tweet)
	return
}

// Destroys the status specified by the required ID parameter.
// The authenticating user must be the author of the specified
// status. returns the destroyed tweet if successful
func (tg *TweetsGroup) Destroy(id int64, opts *Optionals) (tweet *Tweet, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	opts.Add("id", id)
	tweet = &Tweet{}
	err = tg.Call("POST", fmt.Sprintf("statuses/destroy/%d", id), opts, tweet)
	return tweet, err
}

// Retweets a tweet. Returns the original tweet with retweet details embedded.
func (tg *TweetsGroup) Retweet(id int64, opts *Optionals) (tweet *Tweet, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	opts.Add("id", id)
	tweet = &Tweet{}
	err = tg.Call("POST", fmt.Sprintf("statuses/retweet/%d", id), opts, tweet)
	return tweet, err
}

// Updates the authenticating user's current status and attaches media for
// upload. In other words, it creates a Tweet with a picture attached.
func (tg *TweetsGroup) UpdateWithMedia(status string, media *TweetMedia, opts *Optionals) (tweet *Tweet, err error) {
	if opts == nil {
		opts = NewOptionals()
	}

	body := bytes.NewBufferString("")
	mp := multipart.NewWriter(body)
	mp.WriteField("status", status)
	for n, v := range opts.Values {
		mp.WriteField(n, v[0])
	}
	writer, err := mp.CreateFormFile("media[]", media.Filename)
	if err != nil {
		return nil, err
	}
	writer.Write(media.Data)
	header := fmt.Sprintf("multipart/form-data;boundary=%v", mp.Boundary())
	mp.Close()

	endpoint := fmt.Sprintf("%s/statuses/update_with_media.json?%s", apiURL, opts.Values.Encode())
	req, _ := http.NewRequest("POST", endpoint, body)
	req.Header.Set("Content-Type", header)
	res, err := tg.client.Do(req)
	if err != nil {
		return
	}
	if err = checkResponse(res); err != nil {
		return
	}
	if err = json.NewDecoder(res.Body).Decode(tweet); err != nil {
		return
	}
	return

}
