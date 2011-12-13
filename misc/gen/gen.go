package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	service = flag.String("service", "", "Name of the service")
	required = flag.String("required", "", "Comma-separated list of required parameters (name:type)")
	options = flag.String("options", "", "Comma-separated list of options (name:type)")
	call = flag.String("call", "", "Call name")
	method = flag.String("method", "GET", "HTTP method (GET or POST)")
	ret = flag.String("ret", "interface{}", "Return type")
	endpoint = flag.String("endpoint", "", "Twitter RESP API endpoint")
)

func camel(s string) string {
	var news string
	for _, w := range strings.Split(s, "_") {
		news += strings.ToUpper(w[0:1]) + w[1:]
	}
	return news
}

func printComment() {
	fmt.Println("// Automatically generated")
	fmt.Printf("// %s\n", strings.Join(os.Args, " "))
}

func optSnippets() (funcs string, assigns string) {
	if *options == "" {
		return "", ""
	}
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
	c := &`+*service+*call+`Call{s: r.s, opt_: make(map[string]interface{})}`)
	fmt.Printf("%s\treturn c\n}\n", assigns)
}

func printDo(assigns string) {
	fmt.Print(`
func (c *`+*service+*call+`Call) Do() (*`+*ret+`, os.Error) {
	var body io.Reader = nil
	params := make(url.Values)
`)

	if *required != "" {
		for _, r := range strings.Split(*required, ",") {
			rName := strings.Split(r, ":")[0]
			fmt.Printf("\tparams.Set(\"%s\", fmt.Sprintf(\"%%v\", c.%s))\n", rName, rName)
		}
	}

	fmt.Print(assigns)

	if *method == "POST" {
		fmt.Println(`
	urls := fmt.Sprintf("%s/%s.json", apiURL, "`+*endpoint+`")
	urls += "?" + params.Encode()
	body = bytes.NewBuffer([]byte(params.Encode()))
	ctype := "application/x-www-form-urlencoded"
	req, _ := http.NewRequest("POST", urls, body)
	req.Header.Set("Content-Type", ctype)
	res, err := c.s.client.Do(req)`)
	} else {
		fmt.Println(`
	urls := fmt.Sprintf("%s/%s.json", apiURL, "`+*endpoint+`")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	res, err := c.s.client.Do(req)`)
	}
	fmt.Println(`
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := new(`+*ret+`)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil && reflect.TypeOf(err) != reflect.TypeOf(&json.UnmarshalTypeError{}){
		return ret, err
	}
	return ret, nil
}`)

}

func errorReq(s string) {
	fmt.Printf("Parameter '%s' is required.\n", s)
	flag.PrintDefaults()
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

	printComment()
	printType()
	printCall()
	fmt.Println(funcs)
	printDo(assigns)
}
