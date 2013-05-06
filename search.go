// tweetlib - A fully oauth-authenticated Go Twitter library
//
// Copyright 2011 The Tweetlib Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tweetlib

// Groups search functionality
type SearchService struct {
	*Client
}

// Results of a search for tweets.
type SearchResults struct {
	// Tweets matching the search query
	Results TweetList `json:"statuses"`
	// Search metadata
	Metadata SearchMetadata `json:"search_metadata"`
}

// When searching, Twitter returns this metadata
// along with results
type SearchMetadata struct {
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
func (sg *SearchService) Tweets(q string, opts *Optionals) (searchResults *SearchResults, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	opts.Add("q", q)
	searchResults = &SearchResults{}
	err = sg.Call("GET", "search/tweets", opts, searchResults)
	return
}
