include $(GOROOT)/src/Make.inc

TARG=github.com/robteix/tweetlib
GOFILES=\
	twitter.go\
	types.go\
	twitter_oauth.go\

include $(GOROOT)/src/Make.pkg
