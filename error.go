package tweetlib

import (
    "net/http"
    "net/http/httputil"
    "strconv"
)


type TwitterError struct {
    errorString string    
}

func (e *TwitterError) Error() string {
    return e.errorString
}

func NoBearerTokenError() *TwitterError {
    return &TwitterError{"No bearer token was attached with your token response"}
}

func StatusCodeError(resp *http.Response) *TwitterError {
    dump, _ := httputil.DumpResponse(resp, true)
    errorString := "A "+strconv.Itoa(resp.StatusCode)+" code was returned. Full body response:\n"+string(dump)
    return &TwitterError{errorString}
}
