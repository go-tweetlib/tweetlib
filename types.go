// tweetlib - A fully oauth-authenticated Go Twitter library
//
// Copyright 2011 The Tweetlib Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tweetlib

type Tweet struct {
	Id_str                    string    `json:"id_str,omitempty"`
	Text                      string    `json:"text,omitempty"`
	In_reply_to_screen_name   string    `json:"in_reply_to_screen_name,omitempty"`
	In_reply_to_status_id_str string    `json:"in_reply_to_status_id_str,,omitempty"`
	In_reply_to_user_id_str   int64     `json:"in_reply_to_user_id_str,omitempty"`
	Created_at                string    `json:"created_at,omitempty"`
	User                      *User     `json:"user,omitempty"`
	Entities                  *Entities `json:"entities,omitempty"`
}

type TweetList []Tweet

type User struct {
	Id_str                  string
	Screen_name             string
	Name                    string
	Profile_image_url       string
	Profile_image_url_https string
}

type UserList []User

type MediaEntity struct {
	Id_str          string
	Media_url       string
	Media_url_https string
}

type UrlEntity struct {
	Expanded_url string
}

type Entities struct {
	Media *[]MediaEntity
	Urls  *[]UrlEntity
}

type Configuration struct {
	CharactersReservedPerMedia int                  `json:"characters_reserved_per_media,omitempty"`
	MaxMediaPerUpload          int                  `json:"max_media_per_upload,omitempty"`
	NonUsernamePaths           []string             `json:"non_username_paths,omitempty"`
	PhotoSizeLimit             int                  `json:"photo_size_limit,omitempty"`
	PhotoSizes                 map[string]PhotoSize `json:"photo_sizes,omitempty"`
	ShortUrlLengthHttps        int                  `json:"short_url_length_https,omitempty"`
	ShortUrlLength             int                  `json:"short_url_length,omitempty"`
}

type PhotoSize struct {
	Width  int    `json:"w,omitempty"`
	Height int    `json:"h,omitempty"`
	Resize string `json:"resize,omitempty"`
}

type Message struct {
	CreatedAt string `json:"created_at,omitempty"`
	SenderScreenName string `json:"sender_screen_name,omitempty"`
	Sender *User `json:"sender,omitempty"`
	Text string `json:"text,omitempty"`
	RecipientScreenName string `json:"recipient_screen_name,omitempty"`
	Id string `json:"id,omitempty"`
	Recipient *User `json:"recipient,omitempty"`
	RecipientId string `json:"recipient_id,omitempty"`
	SenderId string `json:"sender_id,omitempty"`
}

type MessageList []Message