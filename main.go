package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type Task struct {
	Des       string // description
	Start     int64  // start time
	DStart    string
	Expires   int64 // expires time
	DExpires  string
	Iterval   int64 // interval time
	Itervaled bool  // interval ?
	Status    bool  // on off
}

var (
	index_tpl *template.Template
	tasks     []*Task
)

func init() {
	index_tpl, _ = template.New("index.html").ParseFiles("index.html")
	tasks = make([]*Task, 0, 100)
}
func DisplayTime(t int64) string {
	// time.ParseInLocation("layout", value, loc)
	loc, err := time.LoadLocation("Asia/Shanghai") //"Asia/Shanghai"
	if nil != err {
		return "location nil"
	}
	return time.Unix(t, 0).In(loc).Format("2006-01-02 15:04:05")
}

func New(des string, start, expires, interval int64) *Task {
	t := &Task{Des: des, Start: start, DStart: DisplayTime(start), Expires: expires, DExpires: DisplayTime(expires)}
	if interval > 0 {
		t.Itervaled = true
		t.Iterval = interval
	}
	t.Status = true
	var now time.Time
	ticker := time.NewTicker(time.Millisecond * 1000)
	go func() {
		for {
			<-ticker.C
			if !t.Status {
				continue
			}
			now = time.Now()
			if now.Unix() > start {
				fmt.Println("\n")
				fmt.Print(t.Des, "*")
			}
			if now.Unix() > expires {
				fmt.Print(t.Des, "END\n")
				t.Status = false
				break
			}
		}
	}()
	return t
}

func main() {
	now := time.Now()
	thr := now.Add(time.Second * 3).Unix()
	t := New("des", now.Unix(), thr, 0)
	tasks = append(tasks, t)
	tasks = append(tasks, New("sec", now.Unix(), thr, 0))
	http.HandleFunc("/", index)
	http.ListenAndServe(":80", nil)
	t.Status = false
	time.Sleep(2e9)
	t.Status = true
	time.Sleep(10e9)
}

func index(rw http.ResponseWriter, req *http.Request) {
	index_tpl.Execute(rw, tasks)
}