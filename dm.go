// tweetlib - A fully oauth-authenticated Go Twitter library
//
// Copyright 2011 The Tweetlib Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tweetlib

type DMService struct {
	*Client
}

// Represents a direct message -- a message between
// two users who follow each other.
type DirectMessage struct {
	CreatedAt           string `json:"created_at"`
	SenderScreenName    string `json:"sender_screen_name"`
	Sender              *User  `json:"sender"`
	Text                string `json:"text"`
	RecipientScreenName string `json:"recipient_screen_name"`
	Id                  string `json:"id"`
	Recipient           *User  `json:"recipient"`
	RecipientId         string `json:"recipient_id"`
	SenderId            string `json:"sender_id"`
}

// A list of direct messages
type DirectMessageList []DirectMessage

// Returns the 20 most recent direct messages sent to the authenticating user.
// Includes detailed information about the sender and recipient user. You can
// request up to 200 direct messages per call, up to a maximum of 800 incoming DMs
// See https://dev.twitter.com/docs/api/1.1/get/direct_messages
func (dm *DMService) List(opts *Optionals) (messages *DirectMessageList, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	messages = &DirectMessageList{}
	err = dm.Call("GET", "direct_messages", opts, messages)
	return
}

// Returns the 20 most recent direct messages sent by the authenticating user.
// Includes detailed information about the sender and recipient user. You can
// request up to 200 direct messages per call, up to a maximum of 800 outgoing DMs.
// See https://dev.twitter.com/docs/api/1.1/get/direct_messages/sent
func (dm *DMService) Sent(opts *Optionals) (messages *DirectMessageList, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	messages = &DirectMessageList{}
	err = dm.Call("GET", "direct_messages/sent", opts, messages)
	return
}

// Returns a single direct message, specified by an id parameter.
// See https://dev.twitter.com/docs/api/1.1/get/direct_messages/show
func (dm *DMService) Get(id int64, opts *Optionals) (message *DirectMessage, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	opts.Add("id", id)
	message = &DirectMessage{}
	err = dm.Call("GET", "direct_messages/show", opts, message)
	return
}

// Destroys the direct message specified in the required ID parameter.
// The authenticating user must be the recipient of the specified direct
// message.
// See https://dev.twitter.com/docs/api/1.1/post/direct_messages/destroy
func (dm *DMService) Destroy(id int64, opts *Optionals) (message *DirectMessage, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	opts.Add("id", id)
	message = &DirectMessage{}
	err = dm.Call("POST", "direct_messages/show", opts, message)
	return
}

// Sends a new direct message to the specified user from the authenticating user.
// See https://dev.twitter.com/docs/api/1.1/post/direct_messages/new
func (dm *DMService) Send(screenname, text string, opts *Optionals) (message *DirectMessage, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	opts.Add("screen_name", screenname)
	opts.Add("text", text)
	message = &DirectMessage{}
	err = dm.Call("POST", "direct_messages/new", opts, message)
	return
}
