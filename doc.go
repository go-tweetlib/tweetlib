// tweetlib - A fully oauth-authenticated Go Twitter library
//
// Copyright 2011 The Tweetlib Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
A fully OAuth-authenticated implementation of Twitter's REST API v1.1

See https:dev.twitter.com/docs/api/


Usage example:

  import "robteix.com/v2/tweetlib"

  config := &tweetlib.Config{
  	ConsumerKey: "your-consumer-key",
  	ConsumerSecret: "your-consumer-secret",
	Callback: "http:www.my-app.com/my-callback",
  }
  token := &tweetlib.Token{
        OAuthSecret: "your-oauth-secret",
        OAuthToken: "your-oauth-token",
  }
  tr := &tweetlib.Transport{
	Config: config,
        Token: token,
  }

  client, _ := tweetlib.New(tr.Client())
  client.Tweets.Update("Hello, world")

Authentication

Twitter uses OAUTH 1.0a for authentication and tweetlib supports the
3-legged authorization method.

See https://dev.twitter.com/docs/auth/3-legged-authorization

1. Setup a Config structure with the Consuler Key and Secret found on your
application's page. You can create new applications by visiting
http://dev.twitter.com/apps/new. Twitter also requires a callback URL
that will receive Twitter's token

    config := &tweetlib.Config{
        ConsumerKey: "your-consumer-key",
	ConsumerSecret: "your-consumer-secret",
	Callback: "http://www.my-app.com/my-callback",
    }

2. Populate a tweetlibTransport structure with the Config from the previous
step and, for now, an empty Token

    tr := &tweetlib.Transport{
	    Config: config,
            Token: &tweetlib.Token{}
    }

2a. (Optional) tweetlib.Transport uses the http package to talk to Twitter
by default. This may not be possible or desirable. For example, if the
Client is to be used in a Google Appengine app, it becomes necessary to
change the underlying transport to be used. E.g.:

	tr.Transport = &urlfetch.Transport{Context: c}

3. Not it's possible to request the temporary token. This will start the
little Oauth dance with Twitter

	tt, err := tr.TempToken()

4. With the tweetlib.TempToken ready, now it's time to request the user to
authorize your application. This is done by redirecting the user to the URL
returned by tweetlib.TempToken.AuthURL(). (Note that you must save the
temporary token as it will be necessary later to request the permanent token)

	authorizationURL := tt.AUthURL()
	// save tt so it is accessible from the callback later
	// redirect the user to authorizationURL...

5. The user will be promted by Twitter to authorize your application. If they
authorize it, Twitter will call your callback as set in step 1. Twitter will
issue a GET request to your callback with two parameters:

    oauth_token      same as your TempToken.Token
    oauth_verifier   A code that will be used to verify that your call back
                     was valid and received information.


6. Finally, you'll request the permanent token from Twitter

    tok, err := tr.AccessToken(tt, oauth_verifier)

Note that you do not need to update your tweetlib.Transport.Token with the new
token, as this is done automatically, meaning you can immediatly start making
API calls with the same transport.

That said, you must save the token for future use so you don't have to go
through all this dance again each time. Next time you need to make calls
on behalf of a user you already have a token for, you simply set the
Transport with the saved token.

    tr := &tweetlib.Transport{
	    Config: config,
	    Token: savedUserToken,
    }

Making API calls

Making an API call is trivial once authentication is set up. It all starts
with getting an API Client object:

   tr := &tweetlib.Transport{
	   Config: config,
           Token: token
   }
   client, err := tweetlib.New(tr)
   if err != nil {
	   panic(err)
   }

Once you have the client, you can make API calls easily. For example,
to post a tweet as the authenticating user

   tweet, err := client.Tweets.Update("Hello, world", nil)

The vast majority of API calls to the Twitter REST API takes one or two
required parameters along with any number of optional ones. For example,
when updating the status, it is possible to attached geographical
coordinates to it via the 'lat' and 'long' optional parameters.

To provide optional parameters, use tweetlib.Optionals

    opts := tweetlib.NewOptionals()
    opts.Add("lat", 37.7821120598956)
    opts.Add("long", -122.400612831116)
    tweet, err := client.Tweets.Update("Hello, world", opts)


There's also two ways of making arbitrary API calls. This is useful
when you need to call a new API that is not directly supported
by tweetlib's utility functions or maybe you want better control of
the response objects.

The first way is using Client.Call like this:

    var user User
    opts := NewOptionals()
    opts.Add("screen_name", "sometwitteruser")
    err := client.Call("GET", "users/show", opts, &user)

Client.Call will try to unmarshal the response returned from Twitter. If
however you wish to do it yourself or maybe not use the types defined
by tweetlib (User, Tweet, etc), you can use CallJSON instead:

    rawJSON, err := client.CallJSON("GET", "users/show", opts)
    // rawJSON now has the JSON response from Twitter

These two functions are usually internally by the many helper functions
defined in tweetlib and also add flexibility to


*/
package tweetlib
