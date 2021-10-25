package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"

	_ "github.com/orrche/trigger/triggerkasa"

	"github.com/gorilla/pat"
)

func letitbedark(subid int, payload int) {
	conn, err := net.Dial("tcp", "192.168.2.56:5003")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Fprintf(conn, "12;%d;1;0;2;%d\n", subid, payload)
}

func main() {
	p := pat.New()
	p.Post("/ttft", func(w http.ResponseWriter, req *http.Request) {
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("DATA: %s\n", string(data))

		if string(data) == "dark" {
			go letitbedark(3, 0)
			go letitbedark(4, 0)
			exec.Command("amixer", "set", "Master", "mute").Output()
			exec.Command("xset", "dpms", "force", "suspend").Output()
		}
		if string(data) == "light" {
			go letitbedark(3, 1)
			go letitbedark(4, 1)
			exec.Command("amixer", "set", "Master", "unmute").Output()
		}

	})
	p.Get("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("Got a get request\n")
		fmt.Fprintf(w, "Hello World\n")
	})

	http.ListenAndServe(":8080", p)
}
