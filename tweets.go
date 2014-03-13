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
	"reflect"
)

type TweetsService struct {
	*Client
}

// Holds a single tweet. Depending on the API call used, this
// struct may or may not be fully populated.
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

// A list of tweets
type TweetList []Tweet

// A media attached to a tweet. In practice, this represents
// an image file.
type TweetMedia struct {
	Filename string // Name for the file (e.g. image.png)
	Data     []byte // Raw file data
}

// Returns the 20 (by default) most recent tweets containing a users's
// @screen_name for the authenticating user.
// THis method can only return up to 800 tweets (via the "count" optional
// parameter.
// See https://dev.twitter.com/docs/api/1.1/get/statuses/mentions_timeline
func (tg *TweetsService) Mentions(opts *Optionals) (tweets *TweetList, err error) {
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
func (tg *TweetsService) UserTimeline(screenname string, opts *Optionals) (tweets *TweetList, err error) {
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
func (tg *TweetsService) HomeTimeline(opts *Optionals) (tweets *TweetList, err error) {
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
func (tg *TweetsService) RetweetsOfMe(opts *Optionals) (tweets *TweetList, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	tweets = new(TweetList)
	err = tg.Call("GET", "statuses/retweets_of_me", opts, tweets)
	return
}

// Update: posts a status update to Twitter
// See https://dev.twitter.com/docs/api/1.1/post/statuses/update
func (tg *TweetsService) Update(status string, opts *Optionals) (tweet *Tweet, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	opts.Add("status", status)
	tweet = &Tweet{}
	err = tg.Call("POST", "statuses/update", opts, tweet)
	return tweet, err
}

// Returns up to 100 of the first retweets of a given tweet Id
func (tg *TweetsService) Retweets(id int64, opts *Optionals) (tweets *TweetList, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	tweets = &TweetList{}
	err = tg.Call("GET", fmt.Sprintf("statuses/retweets/%d", id), opts, tweets)
	return
}

// Returns a single Tweet, specified by the id parameter.
// The Tweet's author will also be embedded within the tweet.
func (tg *TweetsService) Get(id int64, opts *Optionals) (tweet *Tweet, err error) {
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
func (tg *TweetsService) Destroy(id int64, opts *Optionals) (tweet *Tweet, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	opts.Add("id", id)
	tweet = &Tweet{}
	err = tg.Call("POST", fmt.Sprintf("statuses/destroy/%d", id), opts, tweet)
	return tweet, err
}

// Retweets a tweet. Returns the original tweet with retweet details embedded.
func (tg *TweetsService) Retweet(id int64, opts *Optionals) (tweet *Tweet, err error) {
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
func (tg *TweetsService) UpdateWithMedia(status string, media *TweetMedia, opts *Optionals) (tweet *Tweet, err error) {
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
	defer res.Body.Close()
	if err = checkResponse(res); err != nil {
		return
	}
	tweet = &Tweet{}
	if err = json.NewDecoder(res.Body).Decode(tweet); err != nil &&
		reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return
	}
	return tweet, nil

}

// Results of a search for tweets.
type TweetSearchResults struct {
	// Tweets matching the search query
	Results TweetList `json:"statuses"`
	// Search metadata
	Metadata TweetSearchMetadata `json:"search_metadata"`
}

// When searching, Twitter returns this metadata
// along with results
type TweetSearchMetadata struct {
	MaxId       int64   `json:"max_id"`
	SinceId     int64   `json:"since_id"`
	RefreshUrl  string  `json:"refresh_url"`
	NextResults string  `json:"next_results"`
	Count       int     `json:"count"`
	CompletedIn float64 `json:"completed_in"`
	SinceIdStr  string  `json:"since_id_str"`
	Query       string  `json:"query"`
	MaxIdStr    string  `json:"max_id_str"`
}

// Returns a collection of relevant Tweets matching a specified query.
// See https://dev.twitter.com/docs/api/1.1/get/search/tweets
// and also https://dev.twitter.com/docs/using-search
func (tg *TweetsService) Tweets(q string, opts *Optionals) (searchResults *TweetSearchResults, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	opts.Add("q", q)
	searchResults = &TweetSearchResults{}
	err = tg.Call("GET", "search/tweets", opts, searchResults)
	return
}
