include $(GOROOT)/src/Make.inc

TARG=tweetlib
GOFILES=\
	twitter.go\
	types.go\
	twitter_oauth.go\

include $(GOROOT)/src/Make.pkg
