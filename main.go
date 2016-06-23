package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"log/syslog"
	"net/http"
	"regexp"

	"github.com/BurntSushi/toml"

	"minoris.se/rabbitmq/camq"
)

type Trigger struct {
	Id      string
	Queue   string
	Message string

	amq *camq.Amq
}

type Config struct {
	AMQ camq.AMQConfig
}

type State struct {
	Triggers []*Trigger
}

var config Config

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func (state *State) trigger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	re := regexp.MustCompile("^/trigger/([^/]*)")
	matches := re.FindString(r.URL.String())[9:]

	if len(matches) > 0 {
		for _, trigger := range state.Triggers {
			if trigger.Id == matches {
				mac := hmac.New(sha1.New, []byte("secret"))
				data, _ := ioutil.ReadAll(r.Body)
				mac.Write(data)
				expectedMAC := mac.Sum(nil)

				headerMac, _ := hex.DecodeString(r.Header.Get("X-Hub-Signature"))
				if hmac.Equal(headerMac, expectedMAC) {
					fmt.Fprintf(w, "Triggered<br/>", trigger.Message)
					trigger.amq.Publish(trigger.Message)
				}
				fmt.Fprintf(w, "%x\n%+v", expectedMAC, r)

			}
		}
	}
}

func (trigger *Trigger) Init(config Config) {
	trigger.amq = camq.GetAMQChannel(config.AMQ)
	trigger.amq.DeclareExchange()
}

func (state *State) addTrigger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	if r.Method == "POST" {
		var trigger Trigger
		trigger.Id = r.FormValue("id")
		trigger.Message = r.FormValue("message")
		trigger.Queue = r.FormValue("queue")

		trigger.Init(config)

		state.Triggers = append(state.Triggers, &trigger)
		SaveState(*state)
	} else {

		fmt.Fprintf(w, "<form method='post'>")
		fmt.Fprintf(w, "Id: <input name='id'></input><br/>")
		fmt.Fprintf(w, "Queue: <input name='queue'></input><br/>")
		fmt.Fprintf(w, "Message: <input name='message'></input><br/>")
		fmt.Fprintf(w, "<input type='submit'></input>")
		fmt.Fprintf(w, "</form>")
	}
}

func main() {
	logwriter, err := syslog.Dial("tcp", "syslog:5556", syslog.LOG_NOTICE, "trigger")
	failOnError(err, "Unable to setup syslog logger")
	log.SetOutput(logwriter)

	toml.DecodeFile("/opt/trigger/config.toml", &config)
	state := LoadState(config)

	http.HandleFunc("/trigger/", state.trigger)
	http.HandleFunc("/addtrigger/", state.addTrigger)
	http.ListenAndServe(":8080", nil)
}
