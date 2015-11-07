package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"time"

	"github.com/everfore/gotest/mail"
)

type Excutor interface {
	Excute()
}

type MailExcutor struct {
}

func (e *MailExcutor) Excute() {
	fmt.Println("\n****************************")
	fmt.Println("Time is over!")
	fmt.Println("****************************")
	mail.SendMail("Time is over")
}

type Task struct {
	Des       string // description
	Seconds   int64
	Start     int64 // start time
	DStart    string
	Expires   int64 // expires time
	DExpires  string
	Iterval   int64   // interval time
	Itervaled bool    // interval ?
	Status    bool    // on off
	Exc       Excutor // excutor
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
	t := &Task{Des: des, Seconds: (expires - start), Start: start, DStart: DisplayTime(start), Expires: expires, DExpires: DisplayTime(expires), Exc: &MailExcutor{}}
	if interval > 0 {
		t.Itervaled = true
		t.Iterval = interval
	}
	t.Status = true
	var now time.Time
	ticker := time.NewTicker(time.Millisecond * 1000)
	go func() {
		i := t.Seconds
		for {
			<-ticker.C
			if !t.Status {
				continue
			}
			now = time.Now()
			if now.Unix() > start {
				fmt.Printf("%s(%d)-", t.Des, i)
			}
			if now.Unix() > expires {
				// fmt.Print(t.Des, "END\n")
				t.Exc.Excute()
				t.Status = false
				break
			}
			i--
		}
	}()
	return t
}

func main() {
	now := time.Now()
	thr := now.Add(time.Second * 3).Unix()
	tasks = append(tasks, New("first", now.Unix(), thr, 0))
	http.HandleFunc("/", index)
	http.HandleFunc("/add", add)
	http.HandleFunc("/clear", clear)
	http.ListenAndServe(":80", nil)
}

func index(rw http.ResponseWriter, req *http.Request) {
	new_tasks := make([]*Task, 0, cap(tasks))
	for _, it := range tasks {
		if it.Expires > time.Now().Unix() {
			new_tasks = append(new_tasks, it)
		}
	}
	tasks = new_tasks
	index_tpl.Execute(rw, tasks)
}

func add(rw http.ResponseWriter, req *http.Request) {
	i := rand.Intn(60)
	now := time.Now()
	d := (time.Duration)(int64(i) * 1e9)
	tasks = append(tasks, New(fmt.Sprintf("%d", i), now.Unix(), now.Add(d).Unix(), 0))
	http.Redirect(rw, req, "/", 302)
}

func clear(rw http.ResponseWriter, req *http.Request) {
	for _, it := range tasks {
		it.Status = false
	}
	tasks = make([]*Task, 0, 100)
	index_tpl.Execute(rw, tasks)
}
