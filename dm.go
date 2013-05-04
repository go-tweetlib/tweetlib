// tweetlib - A fully oauth-authenticated Go Twitter library
//
// Copyright 2011 The Tweetlib Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tweetlib

type DMGroup struct {
	*Client
}

// Returns the 20 most recent direct messages sent to the authenticating user.
// Includes detailed information about the sender and recipient user. You can
// request up to 200 direct messages per call, up to a maximum of 800 incoming DMs
// See https://dev.twitter.com/docs/api/1.1/get/direct_messages
func (dm *DMGroup) List(opts *Optionals) (messages *MessageList, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	messages = &MessageList{}
	err = dm.Call("GET", "direct_messages", opts, messages)
	return
}

// Returns the 20 most recent direct messages sent by the authenticating user.
// Includes detailed information about the sender and recipient user. You can
// request up to 200 direct messages per call, up to a maximum of 800 outgoing DMs.
// See https://dev.twitter.com/docs/api/1.1/get/direct_messages/sent
func (dm *DMGroup) Sent(opts *Optionals) (messages *MessageList, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	messages = &MessageList{}
	err = dm.Call("GET", "direct_messages/sent", opts, messages)
	return
}

// Returns a single direct message, specified by an id parameter.
// See https://dev.twitter.com/docs/api/1.1/get/direct_messages/show
func (dm *DMGroup) Get(id int64, opts *Optionals) (message *Message, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	opts.Add("id", id)
	message = &Message{}
	err = dm.Call("GET", "direct_messages/show", opts, message)
	return
}

// Destroys the direct message specified in the required ID parameter.
// The authenticating user must be the recipient of the specified direct
// message.
// See https://dev.twitter.com/docs/api/1.1/post/direct_messages/destroy
func (dm *DMGroup) Destroy(id int64, opts *Optionals) (message *Message, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	opts.Add("id", id)
	message = &Message{}
	err = dm.Call("POST", "direct_messages/show", opts, message)
	return
}

// Sends a new direct message to the specified user from the authenticating user.
// See https://dev.twitter.com/docs/api/1.1/post/direct_messages/new
func (dm *DMGroup) Send(screenname, text string, opts *Optionals) (message *Message, err error) {
	if opts == nil {
		opts = NewOptionals()
	}
	opts.Add("screen_name", screenname)
	opts.Add("text", text)
	message = &Message{}
	err = dm.Call("POST", "direct_messages/new", opts, message)
	return
}
