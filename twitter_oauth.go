// tweetlib - A fully oauth-authenticated Go Twitter library
//
// Copyright 2011 The Tweetlib Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tweetlib

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	tokenRequestURL = "https://api.twitter.com/oauth/request_token" // request token endpoint
	authURL         = "https://api.twitter.com/oauth/authorize"     // user authorization endpoint
	accessTokenURL  = "https://api.twitter.com/oauth/access_token"  // access token endpoint
)

type Config struct {
	ConsumerKey    string
	ConsumerSecret string
	Callback       string
}

type Token struct {
	OAuthSecret string
	OAuthToken  string
}
type TempToken struct {
	Token  string
	Secret string
}

func (tt *TempToken) AuthURL() string {
	return fmt.Sprintf("%s?oauth_token=%s", authURL, tt.Token)
}

func (t *Transport) nonce() string {
	s := time.Now()
	return strconv.FormatInt(s.Unix(), 10)
}

func (c *Config) callback() string {
	if c.Callback != "" {
		return c.Callback
	}
	return "oob"
}

type Transport struct {
	*Config
	*Token

	// Transport is the HTTP transport to use when making requests.
	// It will default to http.DefaultTransport if nil.
	// (It should never be an oauth.Transport.)
	Transport http.RoundTripper
}

// Client returns an *http.Client that makes OAuth-authenticated requests.
func (t *Transport) Client() *http.Client {
	return &http.Client{Transport: t}
}

func (t *Transport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Config == nil {
		return nil, errors.New("no Config supplied")
	}
	if t.Token == nil {
		return nil, errors.New("no Token supplied")
	}

	// Refresh the Token if it has expired.
	//if t.Expired() {
	//	if err := t.Refresh(); err != nil {
	//		return nil, err
	//	}
	//}
	// Make the HTTP request.
	t.sign(req)
	return t.transport().RoundTrip(req)
}

// Twitter requires that all authenticated requests be
// signed with HMAC-SHA1
//
// https://dev.twitter.com/docs/auth/oauth
//
// The base string is a special combination
// of parameters:
//
//      httpMethod + "&" +
//      url_encode(  base_uri ) + "&" +
//      sorted_query_params.each  { | k, v |
//          url_encode ( k ) + "%3D" +
//          url_encode ( v )
//      }.join("%26")
//
// And then you sign this with HMAC-SHA1 with the key:
//
//    consumer_secret&oauth_token_secret
//
func (t *Transport) sign(req *http.Request) error {
	u, _ := url.ParseQuery(req.URL.RawQuery) //"status": {"testing..."}}
	u.Set("oauth_signature_method", "HMAC-SHA1")
	u.Set("oauth_timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	u.Set("oauth_nonce", t.nonce())
	u.Set("oauth_version", "1.0")
	u.Set("oauth_consumer_key", t.ConsumerKey)
	if t.OAuthToken != "" {
		u.Set("oauth_token", t.OAuthToken)
	}

	// url-encode and sort all parameters as required by twitter
	var pairs, oauthPairs []string
	for k, v := range u {
		pairs = append(pairs, t.percentEncode(k)+"="+t.percentEncode((v[0])))
		if len(k) > 4 && k[:5] == "oauth" {
			oauthPairs = append(oauthPairs, t.percentEncode((k))+"=\""+t.percentEncode((v[0]))+"\"")
		}
	}
	sort.Strings(pairs)

	// create the base string
	parameters := strings.Join(pairs, "&")
	urlForBase := strings.Split(req.URL.String(), "?")[0]
	if req.Method == "POST" {
		req.URL, _ = url.Parse(urlForBase)
	}
	base := req.Method + "&" +
		t.percentEncode(urlForBase) + "&" + t.percentEncode((parameters))
	// sign the base string with the consumer secret and aouth token string
	signature := t.createSignature(base)

	// Create the Authentication header
	authHeader := "OAuth " + strings.Join(oauthPairs, ", ")
	authHeader += ", oauth_signature=\"" + t.percentEncode((signature)) + "\""
	req.Header.Set("Authorization", authHeader)
	return nil
}

func (t *Transport) createSignature(base string) string {
	key := t.percentEncode((t.ConsumerSecret)) + "&" + t.percentEncode((t.OAuthSecret))
	hash := hmac.New(sha1.New, []byte(key))
	hash.Write([]byte(base))
	var sha1Hash bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &sha1Hash)
	encoder.Write(hash.Sum(nil))
	encoder.Close()
	return sha1Hash.String()
}

func (t *Transport) AccessToken(tempToken *TempToken, oauthVerifier string) (*Token, error) {

	u := &url.Values{"oauth_token": {tempToken.Token},
		"oauth_verifier": {oauthVerifier}}
	var body io.Reader
	body = bytes.NewBuffer([]byte(""))
	urls := fmt.Sprintf("%s?%s", accessTokenURL, u.Encode())
	req, err := http.NewRequest("POST", urls, body)
	if err != nil {
		return nil, err
	}
	resp, err := t.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	//if resp.StatusCode != 200 {
	//	return os.NewError("Authentication Error")
	//}
	defer resp.Body.Close()

	var respBody []byte
	respBody, _ = ioutil.ReadAll(resp.Body)
	data, err := url.ParseQuery(string(respBody))
	if err != nil {
		return nil, err
	}
	t.OAuthToken = data.Get("oauth_token")
	t.OAuthSecret = data.Get("oauth_token_secret")
	return &Token{OAuthToken: t.OAuthToken, OAuthSecret: t.OAuthSecret}, nil
}

func (t *Transport) TempToken() (*TempToken, error) {
	var body io.Reader
	body = bytes.NewBuffer([]byte(""))
	req, err := http.NewRequest("POST", tokenRequestURL+"?oauth_callback="+t.percentEncode(t.callback()), body)
	if err != nil {
		return nil, err
	}
	resp, err := t.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	var respBody []byte
	respBody, _ = ioutil.ReadAll(resp.Body)
	//f resp.StatusCode != 200 {
	//return nil, tc.parseError(respBody)
	//

	data, err := url.ParseQuery(string(respBody))
	if err != nil {
		return nil, err
	}

	confirmed, _ := strconv.ParseBool(data.Get("oauth_callback_confirmed"))
	if !confirmed {
		return nil, errors.New("Rejected Callback")
	}

	return &TempToken{Token: data.Get("oauth_token"),
		Secret: data.Get("oauth_token_secret")}, nil
}

func (t *Transport) shouldEscape(c byte) bool {
	switch {
	case c >= 0x41 && c <= 0x5A:
		return false
	case c >= 0x61 && c <= 0x7A:
		return false
	case c >= 0x30 && c <= 0x39:
		return false
	case c == '-', c == '.', c == '_', c == '~':
		return false
	}
	return true
}

func (tr *Transport) percentEncode(s string) string {
	spaceCount, hexCount := 0, 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if tr.shouldEscape(c) {
			hexCount++
		}
	}

	if spaceCount == 0 && hexCount == 0 {
		return s
	}

	t := make([]byte, len(s)+2*hexCount)
	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case tr.shouldEscape(c):
			t[j] = '%'
			t[j+1] = "0123456789ABCDEF"[c>>4]
			t[j+2] = "0123456789ABCDEF"[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}
