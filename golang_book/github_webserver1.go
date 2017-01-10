package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const IssuesURL = "https://api.github.com/search/issues"
const MilestoneURL = "https://api.github.com/repos/golang/go/milestones"

var issueList = template.Must(template.New("issueslist").
	Funcs(template.FuncMap{"daysAgo": daysAgo}).
	Parse(`
<h1>{{.TotalCount}} issues</h1>
<h3><a href='/users'>Users List</a></h3>
<h3><a href='/milestones'>Milestones</a></h3>
<table>
<tr style='text-align: left'>
  <th>#</th>
  <th>State</th>
  <th>Days ago</th>
  <th>User</th>
  <th>Title</th>
  <th>Milestone</th>
</tr>
{{range .Items}}
<tr>
  <td><a href='{{.HTMLURL}}'>{{.Number}}</a></td>
  <td>{{.State}}</td>
  <td>{{.CreatedAt | daysAgo}}</td>
  <td><a href='{{.User.HTMLURL}}'>{{.User.Login}}</a></td>
  <td><a href='{{.HTMLURL}}'>{{.Title}}</a></td>
  <td>{{.Milestone.Title}}</td>
</tr>
{{end}}
</table>
`))

var usersList = template.Must(template.New("users").Parse(`
<h1>{{.TotalCount}} Users</h1>
<h3><a href='/reports'>Reports</a></h3>
<h3><a href='/milestones'>Milestones</a></h3>
<table>
<tr style='text-align: left'>
  <th>User</th>
  <th>Login</th>
</tr>
{{range .Items}}
<tr>
  <td>{{.User.Login}}</td>
  <td><a href='{{.User.HTMLURL}}'>{{.User.HTMLURL}}</a></td>
</tr>
{{end}}
</table>
`))

var milesList = template.Must(template.New("milestones").Parse(`
<h1>Milestones List</h1>
<h3><a href='/reports'>Reports</a></h3>
<h3><a href='/users'>Users</a></h3>
<table>
<tr style='text-align: left'>
  <th>Milestone</th>
  <th>Description</th>
  <th>URL</th>
  <th>Creator</th>
</tr>
{{range .}}
<tr>
  <td>{{.Title}}</td>
  <td>{{.Description}}</td>
  <td><a href='{{.HTMLURL}}'>{{.HTMLURL}}</a></td>
  <td>{{.Creator.Login}}</td>
</tr>
{{end}}
</table>
`))

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string    // in Markdown format
	Milestone *Milestones
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

type Milestones struct {
	URL         string
	HTMLURL     string `json:"html_url"`
	Number      int
	Title       string
	Description string
	Creator     *Creator
	CreatedAt   time.Time `json:"created_at"`
}

type Creator struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

var IssueRes *IssuesSearchResult
var MilestoneRes *[]Milestones

func init() {
	var args []string
	args = os.Args[1:]
	if len(args) == 0 {
		args = append(args, "repo:golang/go is:open json decoder")
	}
	res, err := SearchIssues(args)
	if err != nil {
		log.Fatal(err)
	}
	IssueRes = res

	res2, err := SearchMilestones()
	if err != nil {
		log.Fatal(err)
	}
	MilestoneRes = res2

	log.Println("Data from github was retrieved successfully")
}

func SearchMilestones() (*[]Milestones, error) {
	resp, err := http.Get(MilestoneURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}
	var result []Milestones
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}

func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := http.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}
	//!-
	// For long-term stability, instead of http.Get, use the
	// variant below which adds an HTTP request header indicating
	// that only version 3 of the GitHub API is acceptable.
	//
	//   req, err := http.NewRequest("GET", IssuesURL+"?q="+q, nil)
	//   if err != nil {
	//       return nil, err
	//   }
	//   req.Header.Set(
	//       "Accept", "application/vnd.github.v3.text-match+json")
	//   resp, err := http.DefaultClient.Do(req)
	//!+

	// We must close resp.Body on all execution paths.
	// (Chapter 5 presents 'defer', which makes this simpler.)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}

func daysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}

func reports(w http.ResponseWriter, r *http.Request) {
	if err := issueList.Execute(w, IssueRes); err != nil {
		log.Fatal(err)
	}
}

func milestones(w http.ResponseWriter, r *http.Request) {
	if err := milesList.Execute(w, MilestoneRes); err != nil {
		log.Fatal(err)
	}
}

func users(w http.ResponseWriter, r *http.Request) {
	if err := usersList.Execute(w, IssueRes); err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", reports)
	http.HandleFunc("/reports", reports)
	http.HandleFunc("/milestones", milestones)
	http.HandleFunc("/users", users)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
