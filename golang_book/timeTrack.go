package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	_ "strings"
	"sync"
)

//func main() {
//	err := ui.Main(func() {
//		name := ui.NewEntry()
//		button := ui.NewButton("Timetrack Hours")
//		greeting := ui.NewLabel("")
//		box := ui.NewVerticalBox()
//		box.Append(ui.NewLabel("Enter your name:"), false)
//		box.Append(name, false)
//		box.Append(button, false)
//		box.Append(greeting, false)
//		window := ui.NewWindow("Hello", 200, 100, false)
//		window.SetChild(box)
//		button.OnClicked(func(*ui.Button) {
//			greeting.SetText("Hello, " + name.Text() + "!")
//		})
//		window.OnClosing(func(*ui.Window) bool {
//			ui.Quit()
//			return true
//		})
//		window.Show()
//	})
//	if err != nil {
//		panic(err)
//	}
//}

type Jar struct {
	sync.Mutex
	cookies map[string][]*http.Cookie
}

func NewJar() *Jar {
	jar := new(Jar)
	jar.cookies = make(map[string][]*http.Cookie)
	return jar
}

func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.Lock()
	if _, ok := jar.cookies[u.Host]; ok {
		for _, c := range cookies {
			log.Println(jar.cookies[u.Host], c)
			jar.cookies[u.Host] = append(jar.cookies[u.Host], c)
		}
	} else {
		jar.cookies[u.Host] = cookies
	}
	jar.Unlock()
}

func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies[u.Host]
}

func NewJarClient() *http.Client {
	return &http.Client{
		Jar: NewJar(),
	}
}

func main() {
	client := NewJarClient()
	req, _ := http.NewRequest("GET", "https://assistant.edvantis.com/Periods.aspx",
		nil)
	//req.SetBasicAuth("rostyslav.kovtun", "Enter192prise")
	// create the client
	//req.Header.Set("Content-Type", "text/html; charset=utf-8")
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:52.0) Gecko/20100101 Firefox/52.0")
	//req.Header.Set("Host","assistant.edvantis.com")
	//req.Header.Set("Connection", "keep-alive")
	//req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	for _, cookie := range client.Jar.Cookies(req.URL) {
		fmt.Printf("  %s: %s\n", cookie.Name, cookie.Value)
	}
	keys := make(map[string]string)
	dates := make(map[string]string)
	var date string

	keys, dates = dataSearch(keys, resp, dates, false)

	fmt.Println(keys, dates)
	//post on the login form.
	fmt.Println("+++++++++++++++++++++++++++LOGIN FORM+++++++++++++++++++++++++++++++++++++++++++++++++++++")
	resp, err = client.PostForm("https://assistant.edvantis.com/Common/Login.aspx?ReturnUrl=%2fPeriods.aspx", url.Values{
		"__EVENTARGUMENT":                       {""},
		"__EVENTTARGET":                         {""},
		"__EVENTVALIDATION":                     {keys["__EVENTVALIDATION"]},
		"__LASTFOCUS":                           {""},
		"__VIEWSTATE":                           {keys["__VIEWSTATE"]},
		"__VIEWSTATEGENERATOR":                  {keys["__VIEWSTATEGENERATOR"]},
		"ctl00$PlaceLoginContent$txtLogin":      {"rostyslav.kovtun"},
		"ctl00$PlaceLoginContent$txtPassword":   {"Enter192prise"},
		"ctl00$PlaceLoginContent$ckbRememberMe": {"on"},
		"ctl00$PlaceLoginContent$btnLogon":      {"Sign in"},
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	for _, cookie := range client.Jar.Cookies(req.URL) {
		fmt.Printf("  %s: %s\n", cookie.Name, cookie.Value)
	}
	keys, dates = dataSearch(keys, resp, dates, true)

	fmt.Println("+++++++++++++++++++++++++++PROJECT SELECT+++++++++++++++++++++++++++++++++++++++++++++++++++++")
	resp, err = client.PostForm("https://assistant.edvantis.com/Periods.aspx", url.Values{
		"__EVENTARGUMENT":                            {""},
		"__EVENTTARGET":                              {"ctl00$PlaceTimeTrackContent$ddlProjects|ctl00$PlaceTimeTrackContent$ddlTasks"},
		"__EVENTVALIDATION":                          {keys["__EVENTVALIDATION"]},
		"__LASTFOCUS":                                {""},
		"__VIEWSTATE":                                {keys["__VIEWSTATE"]},
		"__VIEWSTATEGENERATOR":                       {keys["__VIEWSTATEGENERATOR"]},
		"ctl00$PlaceTimeTrackContent$ddlEmployees":   {"rostyslav.kovtun"},
		"__AjaxControlToolkitCalendarCssLoaded":      {""},
		"__ASYNCPOST":                                {"false"},
		"ctl00$PlaceTimeTrackContent$ddlHours":       {""},
		"ctl00$PlaceTimeTrackContent$ddlMinutes":     {""},
		"ctl00$PlaceTimeTrackContent$ddlProjects":    {"102"},
		"ctl00$PlaceTimeTrackContent$ddlTasks":       {"1269"},
		"ctl00$PlaceTimeTrackContent$txtDate":        {dates["Today"]},
		"ctl00$PlaceTimeTrackContent$txtDescription": {""},
		"ctl00$ScriptManager1":                       {"ctl00$PlaceTimeTrackContent$ctl00|ctl00$PlaceTimeTrackContent$ddlProjects|ctl00$PlaceTimeTrackContent$ddlTasks"},
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	keys, dates = dataSearch(keys, resp, dates, false)

	var counter string
	for i := 5; i <= 10; i++ {
		counter = "0" + strconv.Itoa(i)
		switch i {
		case 5:
			date = dates["Mon"]
			counter = "0" + strconv.Itoa(6)
		case 6:
			date = dates["Mon"]
		case 7:
			date = dates["Tue"]
		case 8:
			date = dates["Wed"]
		case 9:
			date = dates["Thu"]
		case 10:
			date = dates["Fri"]
			counter = strconv.Itoa(i)
		case 11:
			date = dates["Sat"]
			counter = strconv.Itoa(i)
		case 12:
			date = dates["Sun"]
			counter = strconv.Itoa(i)
		}
		/*Select record*/
		fmt.Println("+++++++++++++++++++++++++++RECORD SELECT+++++++++++++++++++++++++++++++++++++++++++++++++++++" + counter)
		resp, err = client.PostForm("https://assistant.edvantis.com/Periods.aspx", url.Values{
			"__EVENTARGUMENT":                            {""},
			"__EVENTTARGET":                              {"ctl00$PlaceTimeTrackContent$PeriodTable$ctl03$ctl" + counter},
			"__EVENTVALIDATION":                          {keys["__EVENTVALIDATION"]},
			"__LASTFOCUS":                                {""},
			"__VIEWSTATE":                                {keys["__VIEWSTATE"]},
			"__VIEWSTATEGENERATOR":                       {keys["__VIEWSTATEGENERATOR"]},
			"ctl00$PlaceTimeTrackContent$ddlEmployees":   {"rostyslav.kovtun"},
			"__AjaxControlToolkitCalendarCssLoaded":      {""},
			"__ASYNCPOST":                                {"false"},
			"ctl00$PlaceTimeTrackContent$ddlHours":       {""},
			"ctl00$PlaceTimeTrackContent$ddlMinutes":     {""},
			"ctl00$PlaceTimeTrackContent$ddlProjects":    {"102"},
			"ctl00$PlaceTimeTrackContent$ddlTasks":       {"1269"},
			"ctl00$PlaceTimeTrackContent$txtDate":        {date},
			"ctl00$PlaceTimeTrackContent$txtDescription": {""},
			"ctl00$ScriptManager1":                       {"ctl00$PlaceTimeTrackContent$ctl00|ctl00$PlaceTimeTrackContent$PeriodTable$ctl03$ctl" + counter},
		})

		if err != nil {
			log.Fatal(err.Error())
		}

		/*Update record*/
		fmt.Println("+++++++++++++++++++++++++++RECORD UPDATE+++++++++++++++++++++++++++++++++++++++++++++++++++++" + counter)
		resp, err = client.PostForm("https://assistant.edvantis.com/Periods.aspx", url.Values{
			"__EVENTARGUMENT":                            {""},
			"__EVENTTARGET":                              {""},
			"__EVENTVALIDATION":                          {keys["__EVENTVALIDATION"]},
			"__LASTFOCUS":                                {""},
			"__VIEWSTATE":                                {keys["__VIEWSTATE"]},
			"__VIEWSTATEGENERATOR":                       {keys["__VIEWSTATEGENERATOR"]},
			"ctl00$PlaceTimeTrackContent$ddlEmployees":   {"rostyslav.kovtun"},
			"__AjaxControlToolkitCalendarCssLoaded":      {""},
			"__ASYNCPOST":                                {"false"},
			"ctl00$PlaceTimeTrackContent$ddlHours":       {"8"},
			"ctl00$PlaceTimeTrackContent$ddlMinutes":     {"0"},
			"ctl00$PlaceTimeTrackContent$ddlProjects":    {"102"},
			"ctl00$PlaceTimeTrackContent$ddlTasks":       {"1269"},
			"ctl00$PlaceTimeTrackContent$txtDate":        {date},
			"ctl00$PlaceTimeTrackContent$txtDescription": {"internal work"},
			"ctl00$ScriptManager1":                       {"ctl00$PlaceTimeTrackContent$ctl00|ctl00$PlaceTimeTrackContent$btnInsert"},
			"ctl00$PlaceTimeTrackContent$btnInsert":      {"Insert"},
		})
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	for _, cookie := range client.Jar.Cookies(req.URL) {
		fmt.Printf("  %s: %s\n", cookie.Name, cookie.Value)
	}
	log.Println(resp)
	//doc, err := goquery.NewDocumentFromResponse(resp)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//html, _ := doc.Html()
	//log.Println(string(html))
}

func dataSearch(keys map[string]string, resp *http.Response, dates map[string]string, single bool) (map[string]string, map[string]string) {
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatal(err)
	}
	dates = dateGrabber(dates, doc, single)
	types := doc.Find("input")
	for node := range types.Nodes {
		singleThing := types.Eq(node)

		hidden_input, _ := singleThing.Attr("type")

		if hidden_input == "hidden" {
			id, _ := singleThing.Attr("id")
			value, _ := singleThing.Attr("value")
			keys[id] = value
		}
	}
	return keys, dates
}

func dateGrabber(dates map[string]string, doc *goquery.Document, single bool) map[string]string {
	singleDate := doc.Find("input#ctl00_PlaceTimeTrackContent_txtDate")
	types := doc.Find("table#ctl00_MainGrid th nobr span")
	if single {
		for node := range singleDate.Nodes {
			singleThing := types.Eq(node)
			value, _ := singleThing.Attr("value")
			dates["Today"] = value
		}
		return dates
	}
	for node := range types.Nodes {
		singleThing := types.Eq(node)
		text := singleThing.Text()
		values := strings.SplitN(text, ",", 2)
		if values[0] != "Total" {
			dates[values[0]] = strings.TrimSpace(values[1])
		}
	}
	return dates
}
