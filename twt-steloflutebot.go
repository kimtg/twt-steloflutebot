package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func slurp(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	res.Body.Close()
	return string(body), nil
}

// return all specific groups
func re_groups(re *regexp.Regexp, text string, group int) []string {
	result := []string{}
	found := re.FindAllStringSubmatch(text, -1)
	for _, v := range found {
		result = append(result, v[group])
	}
	return result
}

func list_naver() []string {
	s, err := slurp("http://www.naver.com")
	if err != nil {
		return nil
	}
	return re_groups(
		regexp.MustCompile("<option value=\".+\">.+: (.+)</option>"),
		s,
		1)
}

func list_daum() []string {
	s, err := slurp("http://www.daum.net")
	if err != nil {
		return nil
	}
	return re_groups(
		regexp.MustCompile("<span class=\"txt_issue\">\n.+tabindex.+\n(<.+>)?(.+?)(<.+>)?\n"),
		s,
		2)
}

func print_and_twt(text string) {
	if len(text) > 140 {
		text = text[:140]
	}
	fmt.Println(text)
	_, _, err := client.Statuses.Update(text, nil)
	if err != nil {
		fmt.Println(err)
	}
}

var client *twitter.Client

func main() {
	consumerKey := ""
	consumerSecret := ""
	accessToken := ""
	accessSecret := ""

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client = twitter.NewClient(httpClient)

	const interval = 120
	fmt.Println("Refreshes every", interval, "minutes.")
	for {
		fmt.Println(time.Now())
		print_and_twt("Naver:" + strings.Join(list_naver(), ","))
		print_and_twt("Daum:" + strings.Join(list_daum(), ","))
		time.Sleep(interval * time.Minute)
	}
}
