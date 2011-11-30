package main

import (
	"flag"
	"fmt"
	"strings"
)

var (
	service = flag.String("service", "", "Name of the service")
	required = flag.String("required", "", "Comma-separated list of required parameters (name:type)")
	options = flag.String("options", "", "Comma-separated list of options (name:type)")
	call = flag.String("call", "", "Call name")
	method = flag.String("method", "", "HTTP method (GET or POST)")
)

func camel(s string) string {
	var news string
	for _, w := range strings.Split(s, "_") {
		news += strings.ToUpper(w[0:1]) + w[1:]
	}
	return news
}

func optSnippets() (funcs string, assigns string) {
	for _, opt := range strings.Split(*options, ",") {
		optName := strings.Split(opt, ":")[0]
		optType := strings.Split(opt, ":")[1]
		funcs += `
func (c *`+*service+*call+`Call) `+camel(optName)+`(`+optName+` `+optType+`) *`+*service+*call+`Call {
	c.opt_["`+optName+`"] = `+optName+`
	return c
}
`
		assigns += `
	if v, ok := c.opt_["`+optName+`"]; ok {
		params.Set("`+optName+`", fmt.Sprintf("%v", v))
	}
`

	}
	return
}

func printType() {

	fmt.Println(`
type `+*service+*call+`Call struct {
	s    *Service`)
	if *required != "" {
	for _, r := range strings.Split(*required, ",") {
		rName := strings.Split(r, ":")[0]
		rType := strings.Split(r, ":")[1]
		fmt.Printf("\t%s %s\n", rName, rType)
	}}
	fmt.Println("\topt_ map[string]interface{}\n}\n\n")
}

func printCall() {
	var assigns string
	if *required != "" {
	for _, r := range strings.Split(*required, ",") {
		rName := strings.Split(r, ":")[0]
		assigns += fmt.Sprintf("\tc.%s = %s\n", rName, rName)
	}}
	pars := strings.Replace(*required, ",", " ", -1)
	fmt.Println(`func (r *`+*service+`Service) `+*call+`(`+pars+`) *`+*service+*call+`Call {
	c := &`+*service+*call+`Call{{s: r.s, opt_: make(map[string]interface{{}})}}`)
	fmt.Printf("%s\treturn c\n}\n", assigns)
}

func printDo(assigns string) {
	fmt.Print(`
func (c *`+*service+*call+`Call) Do() (*Tweet, os.Error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("status", fmt.Sprintf("%v", c.status))`)
	fmt.Print(assigns)

	if *method == "POST" {
		fmt.Println(`
	urls := fmt.Sprintf("%s/%s.json", apiURL, "statuses/update")
	urls += "?" + params.Encode()
	body = bytes.NewBuffer([]byte(params.Encode()))
	ctype := "application/x-www-form-urlencoded"
	req, _ := http.NewRequest("POST", urls, body)
	req.Header.Set("Content-Type", ctype)
	res, err := c.s.client.Do(req)`)
	} else {
		fmt.Println(`
	urls := fmt.Sprintf("%s/%s.json", apiURL, fmt.Sprintf("statuses/%s/retweeted_by/ids", c.id))
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)`)

	fmt.Println(`
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
}`)
}

}

func errorReq(s string) {
	fmt.Printf("Parameter '%s' is required.\n", s)
}

func main() {

	flag.Parse()

	if *service == "" {
		errorReq("service")
		return
	}
	if *call == "" {
		errorReq("call")
		return
	}
	if *method != "GET" && *method != "POST" {
		errorReq("method (GET or POST)")
		return
	}

	funcs, assigns := optSnippets()

	printType()
	printCall()
	fmt.Println(funcs)
	printDo(assigns)
}
