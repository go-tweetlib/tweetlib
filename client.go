// tweetlib - A fully oauth-authenticated Go Twitter library
//
// Copyright 2011 The Tweetlib Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tweetlib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

var _ = reflect.TypeOf

const (
	// URL to post tweets
	postURL = "http://api.twitter.com/1.1/statuses/update.json"
	// General URL for API calls
	apiURL = "https://api.twitter.com/1.1"
)

// Checks whether the response is an error
func checkResponse(res *http.Response) (err error) {
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}
	slurp, err := ioutil.ReadAll(res.Body)
	fmt.Printf("%s\n", slurp)
	if err != nil {
		return err
	}
	var jerr TwitterErrorReply
	if err = json.Unmarshal(slurp, &jerr); err != nil {
		return
	}
	return errors.New(jerr.String())
}

// Client: Twitter API client provides access to the various
// API services
type Client struct {
	client *http.Client

	// Twitter's 'statuses' function group
	Tweets *TweetsService

	// Direct Messages function group
	DM *DMService

	// Twitter's Help function group
	Help *HelpService

	// Account functiona group
	Account *AccountService

	// Search functionality group
	Search *SearchService

	// User services
	User *UserService

	// List services
	Lists *ListService

	// Friend services
	Friends *FriendsService

	// Followers services
	Followers *FollowersService

	// API base endpoint. This is the base endpoing URL for API calls. This
	// can be overwritten by an application that needs to use a different
	// version of the library or maybe a mock.
	Endpoint string

	// The token for twitter application we are using. If it is set to "" then
	// client will assume that we are not making application-only API calls and
	// are instead making calls using user authenticated APIs
	ApplicationToken string
}

// Creates a new twitter client for user authenticated API calls
func New(oauthClient *http.Client) (*Client, error) {
	if oauthClient == nil {
		return nil, errors.New("oauthClient is nil")
	}
	return constructClient(oauthClient, ""), nil
}

// Creates a new twitter client for application-only API calls
func NewApplicationClient(httpClient *http.Client, bearerToken string) (*Client, error) {
	if httpClient == nil {
		return nil, errors.New("httpClient is nil")
	}
	if bearerToken == "" {
		return nil, errors.New("The Bearer Token must be a valid and non-empty")
	}
	return constructClient(httpClient, bearerToken), nil
}

func constructClient(httpClient *http.Client, bearerToken string) *Client {
	c := &Client{client: httpClient}
	c.Help = &HelpService{c}
	c.DM = &DMService{c}
	c.Tweets = &TweetsService{c}
	c.Account = &AccountService{c}
	c.Search = &SearchService{c}
	c.User = &UserService{c}
	c.Lists = &ListService{c}
	c.Friends = &FriendsService{c}
	c.Followers = &FollowersService{c}
	c.Endpoint = "https://api.twitter.com/1.1"
	c.ApplicationToken = bearerToken
	return c
}

// Performs an arbitrary API call and returns the response JSON if successful.
// This is generally used internally by other functions but it can also
// be used to perform API calls not directly supported by tweetlib.
//
// For example
//
//   opts := NewOptionals()
//   opts.Add("status", "Hello, world")
//   rawJSON, _ := client.CallJSON("POST", "statuses/update_status", opts)
//   var tweet Tweet
//   err := json.Unmarshal(rawJSON, &tweet)
//
// is the same as
//
//   tweet, err := client.UpdateStatus("Hello, world", nil)
func (c *Client) CallJSON(method, endpoint string, opts *Optionals) (rawJSON []byte, err error) {
	if method != "GET" && method != "POST" {
		err = fmt.Errorf("Invalid method '%s'. Must be either GET or POST.", method)
		return
	}
	if opts == nil {
		opts = NewOptionals()
	}
	endpoint = fmt.Sprintf("%s/%s.json?%s", apiURL, endpoint, opts.Values.Encode())
	fmt.Println(endpoint)
	var req *http.Request
	if method == "POST" {
		body := bytes.NewBuffer([]byte(opts.Values.Encode()))
		req, _ = http.NewRequest(method, endpoint, body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req, _ = http.NewRequest(method, endpoint, nil)
	}
	if c.ApplicationToken != "" {
		req.Header.Add("Authorization", "Bearer "+c.ApplicationToken)
	}
	res, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if err = checkResponse(res); err != nil {
		return
	}
	rawJSON, err = ioutil.ReadAll(res.Body)
	return
}

// Performs an arbitrary API call and tries to unmarshal the result into
// 'resp' on success. This is generally used internally by the other functions
// but it could be used to perform unsupported API calls.
//
// Example usage:
//
//     var tweet Tweet
//     opts := NewOptionals()
//     opts.Add("status", "Hello, world")
//     err := client.Call("POST", "statuses/update_status", opts, &tweet)
//
// is equivalent to
//
//     tweet, err := client.UpdateStatus("Hello, world", nil)
func (c *Client) Call(method, endpoint string, opts *Optionals, resp interface{}) (err error) {
	rawJSON, err := c.CallJSON(method, endpoint, opts)
	if err != nil {
		return
	}
	fmt.Printf("Response: %s\n", rawJSON)
	if resp != nil {
		if err = json.Unmarshal(rawJSON, resp); err != nil &&
			reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
			return err
		}
	}
	return nil
}
