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
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

const (
	postURL = "http://api.twitter.com/1/statuses/update.json"
	apiURL  = "https://api.twitter.com/1"
)

var (
	ErrOAuth = errors.New("OAuth failure")
)

type errorReply struct {
	Error   string `json:"error"`
	Request string `json:"request"`
}

type errorsReply struct {
	Errors string
}

func checkResponse(res *http.Response) error {
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}
	slurp, err := ioutil.ReadAll(res.Body)
	fmt.Printf("%s\n", slurp)
	if err == nil {
		jerr := new(errorReply)
		err = json.Unmarshal(slurp, jerr)
		if err == nil && jerr.Error != "" {
			return errors.New(jerr.Error)
		}
		errs := new(errorsReply)
		err = json.Unmarshal(slurp, errs)
		if err == nil && errs.Errors != "" {
			return errors.New(errs.Errors)
		}

	}
	return fmt.Errorf("googleapi: got HTTP response code %d and error reading body: %v",
		res.StatusCode, err)
}

func New(client *http.Client) (*Service, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	s := &Service{client: client}
	s.Timelines = &TimelinesService{s: s}
	s.Tweets = &TweetsService{s: s}
	s.Help = &HelpService{s: s}
	s.DM = &DMService{s: s}
	s.Users = &UsersService{s: s}
	s.Account = &AccountService{s: s}
	s.Lists = &ListsService{s: s}
	return s, nil
}

type Service struct {
	client *http.Client

	Timelines *TimelinesService
	Tweets    *TweetsService
	Help      *HelpService
	Search    *SearchService
	DM        *DMService
	Users     *UsersService
	Account   *AccountService
	Lists     *ListsService
}

type TimelinesService struct {
	s *Service
}

type TweetsService struct {
	s *Service
}

type HelpService struct {
	s *Service
}

type SearchService struct {
	s *Service
}

type DMService struct {
	s *Service
}

type UsersService struct {
	s *Service
}

type AccountService struct {
	s *Service
}

type ListsService struct {
	s *Service
}

type TimelinesListCall struct {
	s      *Service
	userid string
	opt_   map[string]interface{}
}

// Search: Search all public profiles.
func (r *TimelinesService) List(userid string) *TimelinesListCall {
	c := &TimelinesListCall{s: r.s, opt_: make(map[string]interface{})}
	c.userid = userid
	return c
}

// Language sets the optional parameter "language": Specify the
// preferred language to search with. See search language codes for
// available values.
func (c *TimelinesListCall) Count(count int) *TimelinesListCall {
	c.opt_["count"] = count
	return c
}

func (c *TimelinesListCall) Do() (*TweetList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	params.Set("user_id", fmt.Sprintf("%v", c.userid))
	if v, ok := c.opt_["count"]; ok {
		params.Set("count", fmt.Sprintf("%v", v))
	}
	urls := fmt.Sprintf("%s/%s.json", apiURL, "statuses/user_timeline")
	//urls += "?" + params.Encode()
	fmt.Printf("URL: %s\n", urls)
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(TweetList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

type TimelinesHomeTimelineCall struct {
	s    *Service
	opt_ map[string]interface{}
}

func (r *TimelinesService) HomeTimeline() *TimelinesHomeTimelineCall {
	c := &TimelinesHomeTimelineCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *TimelinesHomeTimelineCall) Count(count int) *TimelinesHomeTimelineCall {
	c.opt_["count"] = count
	return c
}

func (c *TimelinesHomeTimelineCall) SinceId(since_id string) *TimelinesHomeTimelineCall {
	c.opt_["since_id"] = since_id
	return c
}

func (c *TimelinesHomeTimelineCall) MaxId(max_id string) *TimelinesHomeTimelineCall {
	c.opt_["max_id"] = max_id
	return c
}

func (c *TimelinesHomeTimelineCall) Page(page int) *TimelinesHomeTimelineCall {
	c.opt_["page"] = page
	return c
}

func (c *TimelinesHomeTimelineCall) IncludeRTS(include_rts bool) *TimelinesHomeTimelineCall {
	c.opt_["include_rts"] = include_rts
	return c
}

func (c *TimelinesHomeTimelineCall) TrimUser(trim_user bool) *TimelinesHomeTimelineCall {
	c.opt_["trim_user"] = trim_user
	return c
}

func (c *TimelinesHomeTimelineCall) IncludeEntities(include_entities bool) *TimelinesHomeTimelineCall {
	c.opt_["include_entities"] = include_entities
	return c
}

func (c *TimelinesHomeTimelineCall) ExcludeReplies(exclude_replies bool) *TimelinesHomeTimelineCall {
	c.opt_["exclude_replies"] = exclude_replies
	return c
}
func (c *TimelinesHomeTimelineCall) ContributorDetails(contributor_details bool) *TimelinesHomeTimelineCall {
	c.opt_["contributor_details"] = contributor_details
	return c
}

func (c *TimelinesHomeTimelineCall) Do() (*TweetList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	if v, ok := c.opt_["count"]; ok {
		params.Set("count", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["since_id"]; ok {
		params.Set("since_id", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["max_id"]; ok {
		params.Set("max_id", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["page"]; ok {
		params.Set("page", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["include_rts"]; ok {
		params.Set("include_rts", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["trim_user"]; ok {
		params.Set("trim_user", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["include_entities"]; ok {
		params.Set("include_entities", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["exclude_replies"]; ok {
		params.Set("exclude_replies", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["contributor_details"]; ok {
		params.Set("contributor_details", fmt.Sprintf("%v", v))
	}
	urls := fmt.Sprintf("%s/%s.json", apiURL, "statuses/user_timeline")
	urls += "?" + params.Encode()
	fmt.Printf("URL: %s\n", urls)
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(TweetList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

type TimelinesMentionsCall struct {
	s    *Service
	opt_ map[string]interface{}
}

func (r *TimelinesService) Mentions() *TimelinesMentionsCall {
	c := &TimelinesMentionsCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *TimelinesMentionsCall) Count(count int) *TimelinesMentionsCall {
	c.opt_["count"] = count
	return c
}

func (c *TimelinesMentionsCall) SinceId(since_id string) *TimelinesMentionsCall {
	c.opt_["since_id"] = since_id
	return c
}

func (c *TimelinesMentionsCall) MaxId(max_id string) *TimelinesMentionsCall {
	c.opt_["max_id"] = max_id
	return c
}

func (c *TimelinesMentionsCall) Page(page int) *TimelinesMentionsCall {
	c.opt_["page"] = page
	return c
}

func (c *TimelinesMentionsCall) IncludeRTS(include_rts bool) *TimelinesMentionsCall {
	c.opt_["include_rts"] = include_rts
	return c
}

func (c *TimelinesMentionsCall) TrimUser(trim_user bool) *TimelinesMentionsCall {
	c.opt_["trim_user"] = trim_user
	return c
}

func (c *TimelinesMentionsCall) IncludeEntities(include_entities bool) *TimelinesMentionsCall {
	c.opt_["include_entities"] = include_entities
	return c
}

func (c *TimelinesMentionsCall) ExcludeReplies(exclude_replies bool) *TimelinesMentionsCall {
	c.opt_["exclude_replies"] = exclude_replies
	return c
}
func (c *TimelinesMentionsCall) ContributorDetails(contributor_details bool) *TimelinesMentionsCall {
	c.opt_["contributor_details"] = contributor_details
	return c
}

func (c *TimelinesMentionsCall) Do() (*TweetList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	if v, ok := c.opt_["count"]; ok {
		params.Set("count", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["since_id"]; ok {
		params.Set("since_id", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["max_id"]; ok {
		params.Set("max_id", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["page"]; ok {
		params.Set("page", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["include_rts"]; ok {
		params.Set("include_rts", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["trim_user"]; ok {
		params.Set("trim_user", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["include_entities"]; ok {
		params.Set("include_entities", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["exclude_replies"]; ok {
		params.Set("exclude_replies", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["contributor_details"]; ok {
		params.Set("contributor_details", fmt.Sprintf("%v", v))
	}
	urls := fmt.Sprintf("%s/%s.json", apiURL, "statuses/mentions")
	urls += "?" + params.Encode()
	fmt.Printf("URL: %s\n", urls)
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(TweetList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

type TimelinesPublicTimelineCall struct {
	s    *Service
	opt_ map[string]interface{}
}

func (r *TimelinesService) PublicTimeline() *TimelinesPublicTimelineCall {
	c := &TimelinesPublicTimelineCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *TimelinesPublicTimelineCall) TrimUser(trim_user bool) *TimelinesPublicTimelineCall {
	c.opt_["trim_user"] = trim_user
	return c
}

func (c *TimelinesPublicTimelineCall) IncludeEntities(include_entities bool) *TimelinesPublicTimelineCall {
	c.opt_["include_entities"] = include_entities
	return c
}

func (c *TimelinesPublicTimelineCall) Do() (*TweetList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	if v, ok := c.opt_["trim_user"]; ok {
		params.Set("trim_user", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["include_entities"]; ok {
		params.Set("include_entities", fmt.Sprintf("%v", v))
	}
	urls := fmt.Sprintf("%s/%s.json", apiURL, "statuses/public_timeline")
	urls += "?" + params.Encode()
	fmt.Printf("URL: %s\n", urls)
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(TweetList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

type TimelinesRetweetedByMeCall struct {
	s    *Service
	opt_ map[string]interface{}
}

func (r *TimelinesService) RetweetedByMe() *TimelinesRetweetedByMeCall {
	c := &TimelinesRetweetedByMeCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *TimelinesRetweetedByMeCall) Count(count int) *TimelinesRetweetedByMeCall {
	c.opt_["count"] = count
	return c
}

func (c *TimelinesRetweetedByMeCall) SinceId(since_id string) *TimelinesRetweetedByMeCall {
	c.opt_["since_id"] = since_id
	return c
}

func (c *TimelinesRetweetedByMeCall) MaxId(max_id string) *TimelinesRetweetedByMeCall {
	c.opt_["max_id"] = max_id
	return c
}

func (c *TimelinesRetweetedByMeCall) Page(page int) *TimelinesRetweetedByMeCall {
	c.opt_["page"] = page
	return c
}

func (c *TimelinesRetweetedByMeCall) IncludeRTS(include_rts bool) *TimelinesRetweetedByMeCall {
	c.opt_["include_rts"] = include_rts
	return c
}

func (c *TimelinesRetweetedByMeCall) TrimUser(trim_user bool) *TimelinesRetweetedByMeCall {
	c.opt_["trim_user"] = trim_user
	return c
}

func (c *TimelinesRetweetedByMeCall) IncludeEntities(include_entities bool) *TimelinesRetweetedByMeCall {
	c.opt_["include_entities"] = include_entities
	return c
}

func (c *TimelinesRetweetedByMeCall) ExcludeReplies(exclude_replies bool) *TimelinesRetweetedByMeCall {
	c.opt_["exclude_replies"] = exclude_replies
	return c
}
func (c *TimelinesRetweetedByMeCall) ContributorDetails(contributor_details bool) *TimelinesRetweetedByMeCall {
	c.opt_["contributor_details"] = contributor_details
	return c
}

func (c *TimelinesRetweetedByMeCall) Do() (*TweetList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	if v, ok := c.opt_["count"]; ok {
		params.Set("count", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["since_id"]; ok {
		params.Set("since_id", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["max_id"]; ok {
		params.Set("max_id", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["page"]; ok {
		params.Set("page", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["include_rts"]; ok {
		params.Set("include_rts", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["trim_user"]; ok {
		params.Set("trim_user", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["include_entities"]; ok {
		params.Set("include_entities", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["exclude_replies"]; ok {
		params.Set("exclude_replies", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["contributor_details"]; ok {
		params.Set("contributor_details", fmt.Sprintf("%v", v))
	}
	urls := fmt.Sprintf("%s/%s.json", apiURL, "statuses/retweeted_by_me")
	urls += "?" + params.Encode()
	fmt.Printf("URL: %s\n", urls)
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(TweetList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

type TimelinesRetweetedToMeCall struct {
	s    *Service
	opt_ map[string]interface{}
}

func (r *TimelinesService) RetweetedToMe() *TimelinesRetweetedToMeCall {
	c := &TimelinesRetweetedToMeCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *TimelinesRetweetedToMeCall) Count(count int) *TimelinesRetweetedToMeCall {
	c.opt_["count"] = count
	return c
}

func (c *TimelinesRetweetedToMeCall) SinceId(since_id string) *TimelinesRetweetedToMeCall {
	c.opt_["since_id"] = since_id
	return c
}

func (c *TimelinesRetweetedToMeCall) MaxId(max_id string) *TimelinesRetweetedToMeCall {
	c.opt_["max_id"] = max_id
	return c
}

func (c *TimelinesRetweetedToMeCall) Page(page int) *TimelinesRetweetedToMeCall {
	c.opt_["page"] = page
	return c
}

func (c *TimelinesRetweetedToMeCall) IncludeRTS(include_rts bool) *TimelinesRetweetedToMeCall {
	c.opt_["include_rts"] = include_rts
	return c
}

func (c *TimelinesRetweetedToMeCall) TrimUser(trim_user bool) *TimelinesRetweetedToMeCall {
	c.opt_["trim_user"] = trim_user
	return c
}

func (c *TimelinesRetweetedToMeCall) IncludeEntities(include_entities bool) *TimelinesRetweetedToMeCall {
	c.opt_["include_entities"] = include_entities
	return c
}

func (c *TimelinesRetweetedToMeCall) ExcludeReplies(exclude_replies bool) *TimelinesRetweetedToMeCall {
	c.opt_["exclude_replies"] = exclude_replies
	return c
}
func (c *TimelinesRetweetedToMeCall) ContributorDetails(contributor_details bool) *TimelinesRetweetedToMeCall {
	c.opt_["contributor_details"] = contributor_details
	return c
}

func (c *TimelinesRetweetedToMeCall) Do() (*TweetList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	if v, ok := c.opt_["count"]; ok {
		params.Set("count", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["since_id"]; ok {
		params.Set("since_id", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["max_id"]; ok {
		params.Set("max_id", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["page"]; ok {
		params.Set("page", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["include_rts"]; ok {
		params.Set("include_rts", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["trim_user"]; ok {
		params.Set("trim_user", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["include_entities"]; ok {
		params.Set("include_entities", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["exclude_replies"]; ok {
		params.Set("exclude_replies", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["contributor_details"]; ok {
		params.Set("contributor_details", fmt.Sprintf("%v", v))
	}
	urls := fmt.Sprintf("%s/%s.json", apiURL, "statuses/retweeted_to_me")
	urls += "?" + params.Encode()
	fmt.Printf("URL: %s\n", urls)
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(TweetList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

type TimelinesRetweetsOfMeCall struct {
	s    *Service
	opt_ map[string]interface{}
}

func (r *TimelinesService) RetweetsOfMe() *TimelinesRetweetsOfMeCall {
	c := &TimelinesRetweetsOfMeCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *TimelinesRetweetsOfMeCall) Count(count int) *TimelinesRetweetsOfMeCall {
	c.opt_["count"] = count
	return c
}

func (c *TimelinesRetweetsOfMeCall) SinceId(since_id string) *TimelinesRetweetsOfMeCall {
	c.opt_["since_id"] = since_id
	return c
}

func (c *TimelinesRetweetsOfMeCall) MaxId(max_id string) *TimelinesRetweetsOfMeCall {
	c.opt_["max_id"] = max_id
	return c
}

func (c *TimelinesRetweetsOfMeCall) Page(page int) *TimelinesRetweetsOfMeCall {
	c.opt_["page"] = page
	return c
}

func (c *TimelinesRetweetsOfMeCall) TrimUser(trim_user bool) *TimelinesRetweetsOfMeCall {
	c.opt_["trim_user"] = trim_user
	return c
}

func (c *TimelinesRetweetsOfMeCall) IncludeEntities(include_entities bool) *TimelinesRetweetsOfMeCall {
	c.opt_["include_entities"] = include_entities
	return c
}

func (c *TimelinesRetweetsOfMeCall) Do() (*TweetList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	if v, ok := c.opt_["count"]; ok {
		params.Set("count", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["since_id"]; ok {
		params.Set("since_id", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["max_id"]; ok {
		params.Set("max_id", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["page"]; ok {
		params.Set("page", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["trim_user"]; ok {
		params.Set("trim_user", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["include_entities"]; ok {
		params.Set("include_entities", fmt.Sprintf("%v", v))
	}
	urls := fmt.Sprintf("%s/%s.json", apiURL, "statuses/retweets_of_me")
	urls += "?" + params.Encode()
	fmt.Printf("URL: %s\n", urls)
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(TweetList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

type TimelinesRetweetedToUserCall struct {
	s    *Service
	id   string
	opt_ map[string]interface{}
}

func (r *TimelinesService) RetweetedToUser(id string) *TimelinesRetweetedToUserCall {
	c := &TimelinesRetweetedToUserCall{s: r.s, opt_: make(map[string]interface{})}
	c.id = id
	return c
}

func (c *TimelinesRetweetedToUserCall) Count(count int) *TimelinesRetweetedToUserCall {
	c.opt_["count"] = count
	return c
}

func (c *TimelinesRetweetedToUserCall) SinceId(since_id string) *TimelinesRetweetedToUserCall {
	c.opt_["since_id"] = since_id
	return c
}

func (c *TimelinesRetweetedToUserCall) MaxId(max_id string) *TimelinesRetweetedToUserCall {
	c.opt_["max_id"] = max_id
	return c
}

func (c *TimelinesRetweetedToUserCall) Page(page int) *TimelinesRetweetedToUserCall {
	c.opt_["page"] = page
	return c
}

func (c *TimelinesRetweetedToUserCall) TrimUser(trim_user bool) *TimelinesRetweetedToUserCall {
	c.opt_["trim_user"] = trim_user
	return c
}

func (c *TimelinesRetweetedToUserCall) IncludeEntities(include_entities bool) *TimelinesRetweetedToUserCall {
	c.opt_["include_entities"] = include_entities
	return c
}

func (c *TimelinesRetweetedToUserCall) Do() (*TweetList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("id", c.id)
	if v, ok := c.opt_["count"]; ok {
		params.Set("count", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["since_id"]; ok {
		params.Set("since_id", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["max_id"]; ok {
		params.Set("max_id", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["page"]; ok {
		params.Set("page", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["trim_user"]; ok {
		params.Set("trim_user", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["include_entities"]; ok {
		params.Set("include_entities", fmt.Sprintf("%v", v))
	}
	urls := fmt.Sprintf("%s/%s.json", apiURL, "statuses/retweeted_to_user")
	urls += "?" + params.Encode()
	fmt.Printf("URL: %s\n", urls)
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(TweetList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

type TimelinesRetweetedByUserCall struct {
	s    *Service
	id   string
	opt_ map[string]interface{}
}

func (r *TimelinesService) RetweetedByUser(id string) *TimelinesRetweetedByUserCall {
	c := &TimelinesRetweetedByUserCall{s: r.s, opt_: make(map[string]interface{})}
	c.id = id
	return c
}

func (c *TimelinesRetweetedByUserCall) Count(count int) *TimelinesRetweetedByUserCall {
	c.opt_["count"] = count
	return c
}

func (c *TimelinesRetweetedByUserCall) SinceId(since_id string) *TimelinesRetweetedByUserCall {
	c.opt_["since_id"] = since_id
	return c
}

func (c *TimelinesRetweetedByUserCall) MaxId(max_id string) *TimelinesRetweetedByUserCall {
	c.opt_["max_id"] = max_id
	return c
}

func (c *TimelinesRetweetedByUserCall) Page(page int) *TimelinesRetweetedByUserCall {
	c.opt_["page"] = page
	return c
}

func (c *TimelinesRetweetedByUserCall) TrimUser(trim_user bool) *TimelinesRetweetedByUserCall {
	c.opt_["trim_user"] = trim_user
	return c
}

func (c *TimelinesRetweetedByUserCall) IncludeEntities(include_entities bool) *TimelinesRetweetedByUserCall {
	c.opt_["include_entities"] = include_entities
	return c
}

func (c *TimelinesRetweetedByUserCall) Do() (*TweetList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("id", c.id)
	if v, ok := c.opt_["count"]; ok {
		params.Set("count", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["since_id"]; ok {
		params.Set("since_id", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["max_id"]; ok {
		params.Set("max_id", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["page"]; ok {
		params.Set("page", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["trim_user"]; ok {
		params.Set("trim_user", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["include_entities"]; ok {
		params.Set("include_entities", fmt.Sprintf("%v", v))
	}
	urls := fmt.Sprintf("%s/%s.json", apiURL, "statuses/retweets_by_user")
	urls += "?" + params.Encode()
	fmt.Printf("URL: %s\n", urls)
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(TweetList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

type TweetsRetweetedByCall struct {
	s    *Service
	id   string
	opt_ map[string]interface{}
}

func (r *TweetsService) RetweetedBy(id string) *TweetsRetweetedByCall {
	c := &TweetsRetweetedByCall{s: r.s, opt_: make(map[string]interface{})}
	c.id = id
	return c
}

func (c *TweetsRetweetedByCall) Count(count int) *TweetsRetweetedByCall {
	c.opt_["count"] = count
	return c
}

func (c *TweetsRetweetedByCall) Page(page int) *TweetsRetweetedByCall {
	c.opt_["page"] = page
	return c
}

func (c *TweetsRetweetedByCall) Do() (*UserList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("id", c.id)
	if v, ok := c.opt_["count"]; ok {
		params.Set("count", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["page"]; ok {
		params.Set("page", fmt.Sprintf("%v", v))
	}
	urls := fmt.Sprintf("%s/%s.json", apiURL, fmt.Sprintf("statuses/%s/retweeted_by", c.id))
	urls += "?" + params.Encode()
	fmt.Printf("URL: %s\n", urls)
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(UserList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

type TweetsRetweetedByIdsCall struct {
	s    *Service
	id   string
	opt_ map[string]interface{}
}

func (r *TweetsService) RetweetedByIds(id string) *TweetsRetweetedByIdsCall {
	c := &TweetsRetweetedByIdsCall{s: r.s, opt_: make(map[string]interface{})}
	c.id = id
	return c
}

func (c *TweetsRetweetedByIdsCall) Count(count int) *TweetsRetweetedByIdsCall {
	c.opt_["count"] = count
	return c
}

func (c *TweetsRetweetedByIdsCall) Page(page int) *TweetsRetweetedByIdsCall {
	c.opt_["page"] = page
	return c
}

func (c *TweetsRetweetedByIdsCall) Do() (*[]string, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("id", c.id)
	if v, ok := c.opt_["count"]; ok {
		params.Set("count", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["page"]; ok {
		params.Set("page", fmt.Sprintf("%v", v))
	}
	urls := fmt.Sprintf("%s/%s.json", apiURL, fmt.Sprintf("statuses/%s/retweeted_by/ids", c.id))
	urls += "?" + params.Encode()
	fmt.Printf("URL: %s\n", urls)
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new([]string)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

type TweetsRetweetsCall struct {
	s    *Service
	id   string
	opt_ map[string]interface{}
}

func (r *TweetsService) Retweets(id string) *TweetsRetweetsCall {
	c := &TweetsRetweetsCall{s: r.s, opt_: make(map[string]interface{})}
	c.id = id
	return c
}

func (c *TweetsRetweetsCall) Count(count int) *TweetsRetweetsCall {
	c.opt_["count"] = count
	return c
}

func (c *TweetsRetweetsCall) Page(page int) *TweetsRetweetsCall {
	c.opt_["page"] = page
	return c
}
func (c *TweetsRetweetsCall) TrimUser(trim_user bool) *TweetsRetweetsCall {
	c.opt_["trim_user"] = trim_user
	return c
}

func (c *TweetsRetweetsCall) IncludeEntities(include_entities bool) *TweetsRetweetsCall {
	c.opt_["include_entities"] = include_entities
	return c
}

func (c *TweetsRetweetsCall) Do() (*TweetList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("id", c.id)
	if v, ok := c.opt_["count"]; ok {
		params.Set("count", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["page"]; ok {
		params.Set("page", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["include_entities"]; ok {
		params.Set("include_entities", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["trim_users"]; ok {
		params.Set("trim_users", fmt.Sprintf("%v", v))
	}
	urls := fmt.Sprintf("%s/%s.json", apiURL, fmt.Sprintf("statuses/retweets/%s", c.id))
	urls += "?" + params.Encode()
	fmt.Printf("URL: %s\n", urls)
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(TweetList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

type TweetsShowCall struct {
	s    *Service
	id   string
	opt_ map[string]interface{}
}

func (r *TweetsService) Show(id string) *TweetsShowCall {
	c := &TweetsShowCall{s: r.s, opt_: make(map[string]interface{})}
	c.id = id
	return c
}

func (c *TweetsShowCall) TrimUser(trim_user bool) *TweetsShowCall {
	c.opt_["trim_user"] = trim_user
	return c
}

func (c *TweetsShowCall) IncludeEntities(include_entities bool) *TweetsShowCall {
	c.opt_["include_entities"] = include_entities
	return c
}

func (c *TweetsShowCall) Do() (*Tweet, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("id", c.id)
	if v, ok := c.opt_["include_entities"]; ok {
		params.Set("include_entities", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["trim_users"]; ok {
		params.Set("trim_users", fmt.Sprintf("%v", v))
	}
	urls := fmt.Sprintf("%s/%s.json", apiURL, fmt.Sprintf("statuses/show/%s", c.id))
	urls += "?" + params.Encode()
	fmt.Printf("URL: %s\n", urls)
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(Tweet)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

type TweetsDestroyCall struct {
	s    *Service
	id   string
	opt_ map[string]interface{}
}

func (r *TweetsService) Destroy(id string) *TweetsDestroyCall {
	c := &TweetsDestroyCall{s: r.s, opt_: make(map[string]interface{})}
	c.id = id
	return c
}

func (c *TweetsDestroyCall) TrimUser(trim_user bool) *TweetsDestroyCall {
	c.opt_["trim_user"] = trim_user
	return c
}

func (c *TweetsDestroyCall) IncludeEntities(include_entities bool) *TweetsDestroyCall {
	c.opt_["include_entities"] = include_entities
	return c
}

func (c *TweetsDestroyCall) Do() (*Tweet, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("id", c.id)
	if v, ok := c.opt_["include_entities"]; ok {
		params.Set("include_entities", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["trim_users"]; ok {
		params.Set("trim_users", fmt.Sprintf("%v", v))
	}
	urls := fmt.Sprintf("%s/%s.json", apiURL, fmt.Sprintf("statuses/destroy/%s", c.id))
	urls += "?" + params.Encode()
	body = bytes.NewBuffer([]byte(params.Encode()))
	ctype := "application/x-www-form-urlencoded"
	fmt.Printf("URL: %s\n", urls)
	req, _ := http.NewRequest("POST", urls, body)
	req.Header.Set("Content-Type", ctype)
	res, err := c.s.client.Do(req)

	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(Tweet)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

type TweetsRetweetCall struct {
	s    *Service
	id   string
	opt_ map[string]interface{}
}

func (r *TweetsService) Retweet(id string) *TweetsRetweetCall {
	c := &TweetsRetweetCall{s: r.s, opt_: make(map[string]interface{})}
	c.id = id
	return c
}

func (c *TweetsRetweetCall) TrimUser(trim_user bool) *TweetsRetweetCall {
	c.opt_["trim_user"] = trim_user
	return c
}

func (c *TweetsRetweetCall) IncludeEntities(include_entities bool) *TweetsRetweetCall {
	c.opt_["include_entities"] = include_entities
	return c
}

func (c *TweetsRetweetCall) Do() (*Tweet, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("id", c.id)
	if v, ok := c.opt_["include_entities"]; ok {
		params.Set("include_entities", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["trim_users"]; ok {
		params.Set("trim_users", fmt.Sprintf("%v", v))
	}
	urls := fmt.Sprintf("%s/%s.json", apiURL, fmt.Sprintf("statuses/retweet/%s", c.id))
	urls += "?" + params.Encode()
	body = bytes.NewBuffer([]byte(params.Encode()))
	ctype := "application/x-www-form-urlencoded"
	fmt.Printf("URL: %s\n", urls)
	req, _ := http.NewRequest("POST", urls, body)
	req.Header.Set("Content-Type", ctype)
	res, err := c.s.client.Do(req)

	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(Tweet)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

type TweetsUpdateCall struct {
	s      *Service
	status string
	opt_   map[string]interface{}
}

// Search: Search all public profiles.
func (r *TweetsService) Update(status string) *TweetsUpdateCall {
	c := &TweetsUpdateCall{s: r.s, opt_: make(map[string]interface{})}
	c.status = status
	return c
}

func (c *TweetsUpdateCall) Do() (*Tweet, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("status", fmt.Sprintf("%v", c.status))
	urls := fmt.Sprintf("%s/%s.json", apiURL, "statuses/update")
	urls += "?" + params.Encode()
	body = bytes.NewBuffer([]byte(params.Encode()))
	ctype := "application/x-www-form-urlencoded"
	fmt.Printf("URL: %s\n", urls)
	req, _ := http.NewRequest("POST", urls, body)
	req.Header.Set("Content-Type", ctype)
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(Tweet)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

// Help calls --------------------------------------------------------------

type HelpConfigurationCall struct {
	s    *Service
	opt_ map[string]interface{}
}

func (r *HelpService) Configuration() *HelpConfigurationCall {
	c := &HelpConfigurationCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *HelpConfigurationCall) Do() (*Configuration, error) {
	var body io.Reader = nil
	urls := fmt.Sprintf("%s/%s.json", apiURL, "help/configuration")
	fmt.Printf("URL: %s\n", urls)
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	buf, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("sucks: %s\n", buf)
	ret := new(Configuration)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

// Search calls --------------------------------------------------------------

type SearchSearchCall struct {
	s    *Service
	q    string
	opt_ map[string]interface{}
}

func (r *SearchService) Search(q string) *SearchSearchCall {
	c := &SearchSearchCall{s: r.s, opt_: make(map[string]interface{})}
	c.q = q
	return c
}

func (c *SearchSearchCall) Callback(callback string) *SearchSearchCall {
	c.opt_["callback"] = callback
	return c
}

func (c *SearchSearchCall) Geocode(geocode string) *SearchSearchCall {
	c.opt_["geocode"] = geocode
	return c
}

func (c *SearchSearchCall) Lang(lang string) *SearchSearchCall {
	c.opt_["lang"] = lang
	return c
}

func (c *SearchSearchCall) Locale(locale string) *SearchSearchCall {
	c.opt_["locale"] = locale
	return c
}

func (c *SearchSearchCall) Page(page int) *SearchSearchCall {
	c.opt_["page"] = page
	return c
}

func (c *SearchSearchCall) ResultType(result_type string) *SearchSearchCall {
	c.opt_["result_type"] = result_type
	return c
}

func (c *SearchSearchCall) Rpp(rpp int) *SearchSearchCall {
	c.opt_["rpp"] = rpp
	return c
}

func (c *SearchSearchCall) ShowUser(show_user bool) *SearchSearchCall {
	c.opt_["show_user"] = show_user
	return c
}

func (c *SearchSearchCall) Until(until string) *SearchSearchCall {
	c.opt_["until"] = until
	return c
}

func (c *SearchSearchCall) SinceId(since_id string) *SearchSearchCall {
	c.opt_["since_id"] = since_id
	return c
}

func (c *SearchSearchCall) IncludeEntities(include_entities bool) *SearchSearchCall {
	c.opt_["include_entities"] = include_entities
	return c
}

func (c *SearchSearchCall) Do() (*TweetList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("q", c.q)
	if v, ok := c.opt_["callback"]; ok {
		params.Set("callback", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["geocode"]; ok {
		params.Set("geocode", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["lang"]; ok {
		params.Set("lang", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["locale"]; ok {
		params.Set("locale", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["page"]; ok {
		params.Set("page", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["result_type"]; ok {
		params.Set("result_type", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["rpp"]; ok {
		params.Set("rpp", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["show_user"]; ok {
		params.Set("show_user", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["until"]; ok {
		params.Set("until", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["since_id"]; ok {
		params.Set("since_id", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["include_entities"]; ok {
		params.Set("include_entities", fmt.Sprintf("%v", v))
	}
	urls := "http://search.twitter.com/search.json"
	urls += "?" + params.Encode()
	fmt.Printf("URL: %s\n", urls)
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(TweetList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

// Direct Messages -----------------------------------------------------------

type DMMessagesCall struct {
	s    *Service
	opt_ map[string]interface{}
}

func (r *DMService) Messages() *DMMessagesCall {
	c := &DMMessagesCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *DMMessagesCall) SinceId(since_id string) *DMMessagesCall {
	c.opt_["since_id"] = since_id
	return c
}

func (c *DMMessagesCall) MaxId(max_id string) *DMMessagesCall {
	c.opt_["max_id"] = max_id
	return c
}

func (c *DMMessagesCall) Count(count int) *DMMessagesCall {
	c.opt_["count"] = count
	return c
}

func (c *DMMessagesCall) Page(page int) *DMMessagesCall {
	c.opt_["page"] = page
	return c
}

func (c *DMMessagesCall) IncludeEntities(include_entities bool) *DMMessagesCall {
	c.opt_["include_entities"] = include_entities
	return c
}

func (c *DMMessagesCall) SkipStatus(skip_status bool) *DMMessagesCall {
	c.opt_["skip_status"] = skip_status
	return c
}

func (c *DMMessagesCall) Do() (*MessageList, error) {
	var body io.Reader = nil
	params := make(url.Values)

	if v, ok := c.opt_["since_id"]; ok {
		params.Set("since_id", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["max_id"]; ok {
		params.Set("max_id", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["count"]; ok {
		params.Set("count", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["page"]; ok {
		params.Set("page", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["include_entities"]; ok {
		params.Set("include_entities", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["skip_status"]; ok {
		params.Set("skip_status", fmt.Sprintf("%v", v))
	}

	urls := fmt.Sprintf("%s/%s.json", apiURL, "direct_messages")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)

	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(MessageList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

type DMMessagesSentCall struct {
	s    *Service
	opt_ map[string]interface{}
}

func (r *DMService) MessagesSent() *DMMessagesSentCall {
	c := &DMMessagesSentCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *DMMessagesSentCall) SinceId(since_id string) *DMMessagesSentCall {
	c.opt_["since_id"] = since_id
	return c
}

func (c *DMMessagesSentCall) MaxId(max_id string) *DMMessagesSentCall {
	c.opt_["max_id"] = max_id
	return c
}

func (c *DMMessagesSentCall) Count(count int) *DMMessagesSentCall {
	c.opt_["count"] = count
	return c
}

func (c *DMMessagesSentCall) Page(page int) *DMMessagesSentCall {
	c.opt_["page"] = page
	return c
}

func (c *DMMessagesSentCall) IncludeEntities(include_entities bool) *DMMessagesSentCall {
	c.opt_["include_entities"] = include_entities
	return c
}

func (c *DMMessagesSentCall) Do() (*MessageList, error) {
	var body io.Reader = nil
	params := make(url.Values)

	if v, ok := c.opt_["since_id"]; ok {
		params.Set("since_id", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["max_id"]; ok {
		params.Set("max_id", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["count"]; ok {
		params.Set("count", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["page"]; ok {
		params.Set("page", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["include_entities"]; ok {
		params.Set("include_entities", fmt.Sprintf("%v", v))
	}

	urls := fmt.Sprintf("%s/%s.json", apiURL, "direct_messages/sent")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)

	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(MessageList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

type DMDestroyCall struct {
	s    *Service
	id   string
	opt_ map[string]interface{}
}

func (r *DMService) Destroy(id string) *DMDestroyCall {
	c := &DMDestroyCall{s: r.s, opt_: make(map[string]interface{})}
	c.id = id
	return c
}

func (c *DMDestroyCall) Id(id string) *DMDestroyCall {
	c.opt_["id"] = id
	return c
}

func (c *DMDestroyCall) IncludeEntities(include_entities bool) *DMDestroyCall {
	c.opt_["include_entities"] = include_entities
	return c
}

func (c *DMDestroyCall) Do() (*Message, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("id", c.id)

	if v, ok := c.opt_["id"]; ok {
		params.Set("id", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["include_entities"]; ok {
		params.Set("include_entities", fmt.Sprintf("%v", v))
	}

	urls := fmt.Sprintf("%s/%s.json", apiURL, fmt.Sprintf("direct_messages/destroy/%s", c.id))
	urls += "?" + params.Encode()
	body = bytes.NewBuffer([]byte(params.Encode()))
	ctype := "application/x-www-form-urlencoded"
	req, _ := http.NewRequest("POST", urls, body)
	req.Header.Set("Content-Type", ctype)
	res, err := c.s.client.Do(req)

	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(Message)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

type DMNewCall struct {
	s    *Service
	text string
	opt_ map[string]interface{}
}

func (r *DMService) New(text string) *DMNewCall {
	c := &DMNewCall{s: r.s, opt_: make(map[string]interface{})}
	c.text = text
	return c
}

func (c *DMNewCall) UserId(user_id string) *DMNewCall {
	c.opt_["user_id"] = user_id
	return c
}

func (c *DMNewCall) ScreenName(screen_name string) *DMNewCall {
	c.opt_["screen_name"] = screen_name
	return c
}

func (c *DMNewCall) WrapLinks(wrap_links bool) *DMNewCall {
	c.opt_["wrap_links"] = wrap_links
	return c
}

func (c *DMNewCall) Do() (*Tweet, error) {
	var body io.Reader = nil
	params := make(url.Values)

	if v, ok := c.opt_["user_id"]; ok {
		params.Set("user_id", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["screen_name"]; ok {
		params.Set("screen_name", fmt.Sprintf("%v", v))
	}

	params.Set("text", fmt.Sprintf("%v", c.text))

	if v, ok := c.opt_["wrap_links"]; ok {
		params.Set("wrap_links", fmt.Sprintf("%v", v))
	}

	urls := fmt.Sprintf("%s/%s.json", apiURL, "direct_messages/new")
	urls += "?" + params.Encode()
	body = bytes.NewBuffer([]byte(params.Encode()))
	ctype := "application/x-www-form-urlencoded"
	req, _ := http.NewRequest("POST", urls, body)
	req.Header.Set("Content-Type", ctype)
	res, err := c.s.client.Do(req)

	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(Tweet)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

type DMShowCall struct {
	s    *Service
	id   string
	opt_ map[string]interface{}
}

func (r *DMService) Show(id string) *DMShowCall {
	c := &DMShowCall{s: r.s, opt_: make(map[string]interface{})}
	c.id = id
	return c
}

func (c *DMShowCall) Do() (*Message, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("id", fmt.Sprintf("%v", c.id))

	urls := fmt.Sprintf("%s/%s.json", apiURL, fmt.Sprintf("direct_messages/new/%v", c.id))
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)

	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(Message)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

// Users calls ---------------------------------------------------------------

type UsersLookUpCall struct {
	s    *Service
	opt_ map[string]interface{}
}

func (r *UsersService) LookUp() *UsersLookUpCall {
	c := &UsersLookUpCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *UsersLookUpCall) UserId(user_id string) *UsersLookUpCall {
	c.opt_["user_id"] = user_id
	return c
}

func (c *UsersLookUpCall) ScreenName(screen_name string) *UsersLookUpCall {
	c.opt_["screen_name"] = screen_name
	return c
}

func (c *UsersLookUpCall) IncludeEntities(include_entities bool) *UsersLookUpCall {
	c.opt_["include_entities"] = include_entities
	return c
}

func (c *UsersLookUpCall) Do() (*UserList, error) {
	var body io.Reader = nil
	params := make(url.Values)

	if v, ok := c.opt_["user_id"]; ok {
		params.Set("user_id", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["screen_name"]; ok {
		params.Set("screen_name", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["include_entities"]; ok {
		params.Set("include_entities", fmt.Sprintf("%v", v))
	}

	urls := fmt.Sprintf("%s/%s.json", apiURL, "users/lookup")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)

	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(UserList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

type UsersSearchCall struct {
	s    *Service
	q    string
	opt_ map[string]interface{}
}

func (r *UsersService) Search(q string) *UsersSearchCall {
	c := &UsersSearchCall{s: r.s, opt_: make(map[string]interface{})}
	c.q = q
	return c
}

func (c *UsersSearchCall) Page(page int) *UsersSearchCall {
	c.opt_["page"] = page
	return c
}

func (c *UsersSearchCall) PerPage(per_page int) *UsersSearchCall {
	c.opt_["per_page"] = per_page
	return c
}

func (c *UsersSearchCall) IncludeEntities(include_entities bool) *UsersSearchCall {
	c.opt_["include_entities"] = include_entities
	return c
}

func (c *UsersSearchCall) Do() (*UserList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("q", fmt.Sprintf("%v", c.q))

	if v, ok := c.opt_["page"]; ok {
		params.Set("page", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["per_page"]; ok {
		params.Set("per_page", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["include_entities"]; ok {
		params.Set("include_entities", fmt.Sprintf("%v", v))
	}

	urls := fmt.Sprintf("%s/%s.json", apiURL, "users/search")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)

	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(UserList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

type UsersShowCall struct {
	s    *Service
	opt_ map[string]interface{}
}

func (r *UsersService) Show() *UsersShowCall {
	c := &UsersShowCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *UsersShowCall) UserId(user_id string) *UsersShowCall {
	c.opt_["user_id"] = user_id
	return c
}

func (c *UsersShowCall) ScreenName(screen_name string) *UsersShowCall {
	c.opt_["screen_name"] = screen_name
	return c
}

func (c *UsersShowCall) IncludeEntities(include_entities bool) *UsersShowCall {
	c.opt_["include_entities"] = include_entities
	return c
}

func (c *UsersShowCall) Do() (*User, error) {
	var body io.Reader = nil
	params := make(url.Values)

	if v, ok := c.opt_["user_id"]; ok {
		params.Set("user_id", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["screen_name"]; ok {
		params.Set("screen_name", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["include_entities"]; ok {
		params.Set("include_entities", fmt.Sprintf("%v", v))
	}

	urls := fmt.Sprintf("%s/%s.json", apiURL, "users/show")
	urls += "?" + params.Encode()
	fmt.Printf("urls = %s\n", urls)
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)

	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(User)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

type UsersContributeesCall struct {
	s    *Service
	opt_ map[string]interface{}
}

func (r *UsersService) Contributees() *UsersContributeesCall {
	c := &UsersContributeesCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *UsersContributeesCall) UserId(user_id string) *UsersContributeesCall {
	c.opt_["user_id"] = user_id
	return c
}

func (c *UsersContributeesCall) ScreenName(screen_name string) *UsersContributeesCall {
	c.opt_["screen_name"] = screen_name
	return c
}

func (c *UsersContributeesCall) IncludeEntities(include_entities bool) *UsersContributeesCall {
	c.opt_["include_entities"] = include_entities
	return c
}

func (c *UsersContributeesCall) SkipStatus(skip_status bool) *UsersContributeesCall {
	c.opt_["skip_status"] = skip_status
	return c
}

func (c *UsersContributeesCall) Do() (*UserList, error) {
	var body io.Reader = nil
	params := make(url.Values)

	if v, ok := c.opt_["user_id"]; ok {
		params.Set("user_id", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["screen_name"]; ok {
		params.Set("screen_name", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["include_entities"]; ok {
		params.Set("include_entities", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["skip_status"]; ok {
		params.Set("skip_status", fmt.Sprintf("%v", v))
	}

	urls := fmt.Sprintf("%s/%s.json", apiURL, "users/contributees")

	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)

	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(UserList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

type UsersContributorsCall struct {
	s    *Service
	opt_ map[string]interface{}
}

func (r *UsersService) Contributors() *UsersContributorsCall {
	c := &UsersContributorsCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *UsersContributorsCall) UserId(user_id string) *UsersContributorsCall {
	c.opt_["user_id"] = user_id
	return c
}

func (c *UsersContributorsCall) ScreenName(screen_name string) *UsersContributorsCall {
	c.opt_["screen_name"] = screen_name
	return c
}

func (c *UsersContributorsCall) IncludeEntities(include_entities bool) *UsersContributorsCall {
	c.opt_["include_entities"] = include_entities
	return c
}

func (c *UsersContributorsCall) SkipStatus(skip_status bool) *UsersContributorsCall {
	c.opt_["skip_status"] = skip_status
	return c
}

func (c *UsersContributorsCall) Do() (*Tweet, error) {
	var body io.Reader = nil
	params := make(url.Values)

	if v, ok := c.opt_["user_id"]; ok {
		params.Set("user_id", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["screen_name"]; ok {
		params.Set("screen_name", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["include_entities"]; ok {
		params.Set("include_entities", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["skip_status"]; ok {
		params.Set("skip_status", fmt.Sprintf("%v", v))
	}

	urls := fmt.Sprintf("%s/%s.json", apiURL, "users/contributors")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)

	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(Tweet)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

// Automatically generated
// misc/gen --service Users --call SuggestedCategories --method GET --options lang:string

type UsersSuggestedCategoriesCall struct {
	s    *Service
	opt_ map[string]interface{}
}

func (r *UsersService) SuggestedCategories() *UsersSuggestedCategoriesCall {
	c := &UsersSuggestedCategoriesCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *UsersSuggestedCategoriesCall) Lang(lang string) *UsersSuggestedCategoriesCall {
	c.opt_["lang"] = lang
	return c
}

func (c *UsersSuggestedCategoriesCall) Do() (*CategoryList, error) {
	var body io.Reader = nil
	params := make(url.Values)

	if v, ok := c.opt_["lang"]; ok {
		params.Set("lang", fmt.Sprintf("%v", v))
	}

	urls := fmt.Sprintf("%s/%s.json", apiURL, "users/suggestions")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)

	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(CategoryList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

// Automatically generated
// ./misc/gen --service Users --call SuggestedUsers --required
// slug:string --options lang:string --ret Category --endpoint
// users/suggestions/slug

type UsersSuggestedUsersCall struct {
	s    *Service
	slug string
	opt_ map[string]interface{}
}

func (r *UsersService) SuggestedUsers(slug string) *UsersSuggestedUsersCall {
	c := &UsersSuggestedUsersCall{s: r.s, opt_: make(map[string]interface{})}
	c.slug = slug
	return c
}

func (c *UsersSuggestedUsersCall) Lang(lang string) *UsersSuggestedUsersCall {
	c.opt_["lang"] = lang
	return c
}

func (c *UsersSuggestedUsersCall) Do() (*Category, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("slug", fmt.Sprintf("%v", c.slug))

	if v, ok := c.opt_["lang"]; ok {
		params.Set("lang", fmt.Sprintf("%v", v))
	}

	urls := fmt.Sprintf("%s/%s.json", apiURL, fmt.Sprintf("users/suggestions/%s", c.slug))
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)

	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(Category)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

// Acount calls --------------------------------------------------------------

// Automatically generated
// ./misc/gen --service Account --call UpdateProfile --options
// name:string,url:string,location:string,description:string,include_entities:bool,skip_status:bool
// --method POST --ret User --endpoint account/update_profile

type AccountUpdateProfileCall struct {
	s    *Service
	opt_ map[string]interface{}
}

func (r *AccountService) UpdateProfile() *AccountUpdateProfileCall {
	c := &AccountUpdateProfileCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *AccountUpdateProfileCall) Name(name string) *AccountUpdateProfileCall {
	c.opt_["name"] = name
	return c
}

func (c *AccountUpdateProfileCall) Url(url string) *AccountUpdateProfileCall {
	c.opt_["url"] = url
	return c
}

func (c *AccountUpdateProfileCall) Location(location string) *AccountUpdateProfileCall {
	c.opt_["location"] = location
	return c
}

func (c *AccountUpdateProfileCall) Description(description string) *AccountUpdateProfileCall {
	c.opt_["description"] = description
	return c
}

func (c *AccountUpdateProfileCall) IncludeEntities(include_entities bool) *AccountUpdateProfileCall {
	c.opt_["include_entities"] = include_entities
	return c
}

func (c *AccountUpdateProfileCall) SkipStatus(skip_status bool) *AccountUpdateProfileCall {
	c.opt_["skip_status"] = skip_status
	return c
}

func (c *AccountUpdateProfileCall) Do() (*User, error) {
	var body io.Reader = nil
	params := make(url.Values)

	if v, ok := c.opt_["name"]; ok {
		params.Set("name", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["url"]; ok {
		params.Set("url", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["location"]; ok {
		params.Set("location", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["description"]; ok {
		params.Set("description", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["include_entities"]; ok {
		params.Set("include_entities", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["skip_status"]; ok {
		params.Set("skip_status", fmt.Sprintf("%v", v))
	}

	urls := fmt.Sprintf("%s/%s.json", apiURL, "account/update_profile")
	urls += "?" + params.Encode()
	body = bytes.NewBuffer([]byte(params.Encode()))
	ctype := "application/x-www-form-urlencoded"
	req, _ := http.NewRequest("POST", urls, body)
	req.Header.Set("Content-Type", ctype)
	res, err := c.s.client.Do(req)

	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(User)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

// Automatically generated
// ./misc/gen --service Account --call RateLimitStatus --ret
// LimitStatus --endpoint account/rate_limit_status

type AccountRateLimitStatusCall struct {
	s    *Service
	opt_ map[string]interface{}
}

func (r *AccountService) RateLimitStatus() *AccountRateLimitStatusCall {
	c := &AccountRateLimitStatusCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *AccountRateLimitStatusCall) Do() (*LimitStatus, error) {
	var body io.Reader = nil
	params := make(url.Values)

	urls := fmt.Sprintf("%s/%s.json", apiURL, "account/rate_limit_status")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)

	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(LimitStatus)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

// Automatically generated
// ./misc/gen/gen --service Account --call VerifyCredentials --options
// include_entities:bool,skip_status:bool --ret User --endpoint
// account/verify_credentials

type AccountVerifyCredentialsCall struct {
	s    *Service
	opt_ map[string]interface{}
}

func (r *AccountService) VerifyCredentials() *AccountVerifyCredentialsCall {
	c := &AccountVerifyCredentialsCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *AccountVerifyCredentialsCall) IncludeEntities(include_entities bool) *AccountVerifyCredentialsCall {
	c.opt_["include_entities"] = include_entities
	return c
}

func (c *AccountVerifyCredentialsCall) SkipStatus(skip_status bool) *AccountVerifyCredentialsCall {
	c.opt_["skip_status"] = skip_status
	return c
}

func (c *AccountVerifyCredentialsCall) Do() (*User, error) {
	var body io.Reader = nil
	params := make(url.Values)

	if v, ok := c.opt_["include_entities"]; ok {
		params.Set("include_entities", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["skip_status"]; ok {
		params.Set("skip_status", fmt.Sprintf("%v", v))
	}

	urls := fmt.Sprintf("%s/%s.json", apiURL, "account/verify_credentials")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)

	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(User)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}

// Automatically generated
// ./misc/gen --service Lists --call All --options
// user_id:string,screen_name:string --ret ListList --endpoint
// lists/all

type ListsAllCall struct {
	s    *Service
	opt_ map[string]interface{}
}

func (r *ListsService) All() *ListsAllCall {
	c := &ListsAllCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *ListsAllCall) UserId(user_id string) *ListsAllCall {
	c.opt_["user_id"] = user_id
	return c
}

func (c *ListsAllCall) ScreenName(screen_name string) *ListsAllCall {
	c.opt_["screen_name"] = screen_name
	return c
}

func (c *ListsAllCall) Do() (*ListList, error) {
	var body io.Reader = nil
	params := make(url.Values)

	if v, ok := c.opt_["user_id"]; ok {
		params.Set("user_id", fmt.Sprintf("%v", v))
	}

	if v, ok := c.opt_["screen_name"]; ok {
		params.Set("screen_name", fmt.Sprintf("%v", v))
	}

	urls := fmt.Sprintf("%s/%s.json", apiURL, "lists/all")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)

	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(ListList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}) {
		return nil, err
	}
	return ret, nil
}
