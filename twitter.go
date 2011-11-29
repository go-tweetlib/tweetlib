// tweetlib - A fully oauth-authenticated Go Twitter library
//
// Copyright 2011 The Tweetlib Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tweetlib

import (
	"fmt"
	"http"
	"io"
	"io/ioutil"
	"json"
	"os"
	"url"
	"bytes"
)

const (
	postURL = "http://api.twitter.com/1/statuses/update.json"
	apiURL  = "https://api.twitter.com/1"
)

var (
	ErrOAuth = os.NewError("OAuth failure")
)

type errorReply struct {
	Error   string `json:"error"`
	Request string `json:"request"`
}

type errorsReply struct {
	Errors string
}

func checkResponse(res *http.Response) os.Error {
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}
	slurp, err := ioutil.ReadAll(res.Body)
	fmt.Printf("%s\n", slurp)
	if err == nil {
		jerr := new(errorReply)
		err = json.Unmarshal(slurp, jerr)
		if err == nil && jerr.Error != "" {
			return os.NewError(jerr.Error)
		}
		errs := new(errorsReply)
		err = json.Unmarshal(slurp, errs)
		if err == nil && errs.Errors != "" {
			return os.NewError(errs.Errors)
		}

	}
	return fmt.Errorf("googleapi: got HTTP response code %d and error reading body: %v",
		res.StatusCode, err)
}

func New(client *http.Client) (*Service, os.Error) {
	if client == nil {
		return nil, os.NewError("client is nil")
	}
	s := &Service{client: client}
	s.Timelines = &TimelinesService{s: s}
	s.Tweets = &TweetsService{s: s}
	s.Help = &HelpService{s: s}
	return s, nil
}

type Service struct {
	client *http.Client

	Timelines *TimelinesService
	Tweets    *TweetsService
	Help      *HelpService
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

func (c *TimelinesListCall) Do() (*TweetList, os.Error) {
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
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return ret, err
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

func (c *TimelinesHomeTimelineCall) Do() (*TweetList, os.Error) {
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
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return ret, err
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

func (c *TimelinesMentionsCall) Do() (*TweetList, os.Error) {
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
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return ret, err
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

func (c *TimelinesPublicTimelineCall) Do() (*TweetList, os.Error) {
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
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return ret, err
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

func (c *TimelinesRetweetedByMeCall) Do() (*TweetList, os.Error) {
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
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return ret, err
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

func (c *TimelinesRetweetedToMeCall) Do() (*TweetList, os.Error) {
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
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return ret, err
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

func (c *TimelinesRetweetsOfMeCall) Do() (*TweetList, os.Error) {
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
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return ret, err
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

func (c *TimelinesRetweetedToUserCall) Do() (*TweetList, os.Error) {
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
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return ret, err
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

func (c *TimelinesRetweetedByUserCall) Do() (*TweetList, os.Error) {
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
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return ret, err
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

func (c *TweetsRetweetedByCall) Do() (*UserList, os.Error) {
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
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return ret, err
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

func (c *TweetsRetweetedByIdsCall) Do() (*[]string, os.Error) {
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
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return ret, err
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

func (c *TweetsRetweetsCall) Do() (*TweetList, os.Error) {
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
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return ret, err
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

func (c *TweetsShowCall) Do() (*Tweet, os.Error) {
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
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return ret, err
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

func (c *TweetsDestroyCall) Do() (*Tweet, os.Error) {
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
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return ret, err
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

func (c *TweetsRetweetCall) Do() (*Tweet, os.Error) {
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
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return ret, err
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

func (c *TweetsUpdateCall) Do() (*Tweet, os.Error) {
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
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return ret, err
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

func (c *HelpConfigurationCall) Do() (*Configuration, os.Error) {
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
	ret := new(Configuration)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return ret, err
	}
	return ret, nil
}
