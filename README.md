tweetlib
==========

A fully OAuth-authenticated library to interact with the new Twitter's
[REST API](https://dev.twitter.com/docs/api/)

Installing
----------

Best way to isntall this package is by running goinstall:

    goinstall github.com/robteix/tweetlib

And then you can import it in your code like this:

    import (
            ...
            "github.com/robteix/tweetlib"
            ...
        )

And of course, if you want to do it manually, just check the code,
cd to the project's root dir and do a:

    gomake install

General use
-----------

1. First create a `tweetlib.Transport` that will handle OAuth:

        config := &Config{
                    ConsumerKey: "your-consumer-key",
                    ConsumerSecret: "your-consumer-secret",
                    Callback: "http://www.my-app.com/my-callback"}
        token := &Token{
                    OAuthSecret: "your-oauth-secret",
                    OAuthToken: "your-oauth-token"}
        tr := &tweetlib.Transport{Config: config,
                              Token: token}

2. Create a `tweetlib` object:

        tl := tweetlib.New(tr.Client())

3. Post a tweet:

        tl.Tweets.Update("this is cool!").Do()

4. Get the 10 last tweets of your home timeline:

        if tweets, ok := tl.Timelines.HomeTimeline().Count(10).Do(); ok {
            for _, tweet := range tweets {
                fmt.Println(tweet.Message)
            }
        }

5. Get a specific tweet, get only basic User info

        tweet, err := tl.Tweets.Show("12345678").TrimUser(true).Do()

6. Retweet a tweet

        tl.Tweets.Retweet("12345").Do()

7. Delete a tweet

        tl.Tweets.Destroy("1234556").Do()

Check out http://dev.twitter.com/api/docs/ for more.

Obtaining a Twitter Access Token
--------------------------------

In order to autheticate with Twitter and perform actions on behalf of
an user, Twitter requires you to request an Access Token. The process
is often confusing to many people, so I thought I'd go through the
whole process here. You'll see that with tweetlib this will be quite
simple to do:

1. The first thing you need to do is to create a new Twitter
   application by going to: http://dev.twitter.com/apps/new

2. You will need to take note of the *Consumer Key* and *Consumer
   Secret* of your new Twitter application.

3. In order to take actions on behalf of a users, we need to start our
   OAuth dance with Twitter. The first thing to do is to create a new
   tweetlib Transport with an empty Token:

        config := &Config{
                ConsumerKey: "your-consumer-key",
                ConsumerSecret: "your-consumer-secret",
                Callback: "http://www.my-app.com/my-callback"}
        tr := &tweetlib.Transport{Config: config,
                              Token: &Token{}}

4. (Optional) tweetlib.Transport uses the http package to talk to
   Twitter.  If you need to (e.g., you're using Google App Engine),
   you can set a custom underlying transport to be used. E.g:

        tr.Transport = &urlfetch.Transport{Context: c}

5. Now it's time to request a temporary token. This is the actual
   beginning of the OAuth dance with Twitter

        tt, err := tr.TempToken()

6. With `tweetlib.TempToken` ready, you will not need your user to
   authorize your application to use Twitter on their behalf. This is
   done by redirecting your user to an authorization URL.

        authorizationURL := tt.AUthURL()
        // redirect the user to it...

7. After the user authorizes your application to take action on their
   behalf, Twitter will send a verification code to your
   callback. More on that step 7 below.


8. You now need to wait for Twitter to get back to you. They will do
   so by doing a GET on the callback specified in your
   `tweetlib.Config` (step 3 above.)  Twitter will send two parameters
   to your callback: `oauth_token` (same as your `TempToken.Token`)
   and an `oauth_verifier`. You're almost ready. Now you only need to
   exchange your `tweetlib.TempToken` for a permanent one.

        tok, err := tr.AccessToken(tt, oauth_verifier)

9. Note that you do not need to update your `tweetlib.Transport.Token`
   with the new token, since this is done automatically. You should
   save the token for future use, though, so you don't have to go
   through all this dance again.

License
-------

Copyright 2011 The Tweetlib Authors.  All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
