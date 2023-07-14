package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"io/ioutil"
	"strings"
	"hash/fnv"
	"time"
)

type Page struct {
	Title		string
	Content	string
}

// Define a Day struct to hold TGT posts
type Day struct {
	Title		string
	First		string
	Second	string
	Third		string
}

// Color struct for daily colors
type Color struct {
	ColorName	string
	ColorHex	string
}

type Run struct {
	Title			string
	Distance	string
	Notes			[]string
}

// const staticDir = "/home/debian/hd-dn/static"
const staticDir = "./static"

// parses three good things entry into struct
func parseDay(entry string) Day {
	var day = Day{
		Title: "",
		First: "",
		Second: "",
		Third: "",
	}
	ls := strings.Split(entry, "\n")

	for _, l := range ls {
		if len(l) > 2 {
			check := l[0:2]
			rest := l[2:len(l)]

			rest = strings.TrimLeft(rest, " ")

			// t - signifies a title
			// 1., 2., 3. - first, second, third entry
			if check == "t " {
				day.Title = rest
			}
			if check == "1." {
				day.First = rest
			}
			if check == "2." {
				day.Second = rest
			}
			if check == "3." {
				day.Third = rest
			}
		}
	}
	return day
}

// Handler for three good things posts
func handleThreeGoodThings(w http.ResponseWriter, r *http.Request) {
		// Read in a text file containing TGT
		content, err := ioutil.ReadFile("./static/three-good-things/index.txt")
		if (err != nil) {
			panic(err);
		}
		text := string(content)
	
		
		// Split it into posts based on pagebreak elements ***
		var days []Day
		posts := strings.Split(text, "***")
		for _, post := range posts {
			days = append(days, parseDay(post));
		}
		
		// Load the TGT template
		var tmplFile = "three-good-things.tmpl"
		tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
		if err != nil {
			panic(err)
		}
		// Execute the template
		err = tmpl.Execute(w, days)
		if err != nil{
			panic(err)
		}
}

// hash function for color handling
func hash(s string) int {
	h := fnv.New32()
	h.Write([]byte(s))
	return int(h.Sum32())
}

func getDailyColor() Color {
	// Stub for now, finish writing this!
	colors := []Color{
		Color{"red!", "dc143c"},
		Color{"orange!", "ff8c00"},
		Color{"yellow!", "ffff00"},
		Color{"green!", "3cb371"},
		Color{"blue!", "00bfff"},
		Color{"indigo!", "4b0082"},
		Color{"viole(n)t!", "e582ee"},
		Color{"pink!", "ff69b4"},
		Color{"black!", "000000"},
		Color{"white!", "ffffff"},
	}

	sToday := time.Now().Format("20060102150405")[0:8]
	h := hash(sToday)
	return colors[h % len(colors)]
}

// Handler for color of the day
func handleColor(w http.ResponseWriter, r *http.Request) {
	c := getDailyColor()
	w.Header().Set("Content-Type", "text/html")

	// Load color template
	var tmplFile = "color.tmpl"
	tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}
	// Execute the template
	err = tmpl.Execute(w, c)
	if err != nil{
		panic(err)
	}
}

// Running journal
// Parse a run entry
func parseRun(entry string) Run {
	var run = Run {
		Title: "",
		Distance: "",
		Notes: []string{},
	}
	ls := strings.Split(entry, "\n")

	for _, l := range ls {
		if len(l) > 2 {
			check := l[0:2]
			rest := l[2:len(l)]

			rest = strings.TrimLeft(rest, " ")

			// t - signifies a title
			// 1., 2., 3. - first, second, third entry
			if check == "t " {
				run.Title = rest
			} else if check == "d " {
				run.Distance = rest
			} else {
				run.Notes = append(run.Notes, l)
			}
		}
	}
	return run
}



func handleRuns(w http.ResponseWriter, r *http.Request) {
// Read in a text file containing TGT
content, err := ioutil.ReadFile("./static/runs/index.txt")
if (err != nil) {
	panic(err);
}
text := string(content)


// Split it into posts based on pagebreak elements ***
var runs []Run
posts := strings.Split(text, "***")
for _, post := range posts {
	runs = append(runs, parseRun(post));
}

// Load the TGT template
var tmplFile = "runs.tmpl"
tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
if err != nil {
	panic(err)
}
// Execute the template
err = tmpl.Execute(w, runs)
if err != nil{
	panic(err)
}
}


func main() {
	fileServer := http.FileServer(http.Dir(staticDir))

	http.Handle("/", fileServer)

	// Handle Three Good Things separately coz she's special
	http.HandleFunc("/three-good-things/", handleThreeGoodThings)
	http.Handle("/three-good-things/atom.atom", fileServer)
	
	// do colo(u)r of the day
	http.HandleFunc("/color/", handleColor)
	http.HandleFunc("/colour/", handleColor)

	http.HandleFunc("/runs/", handleRuns)

	fmt.Printf("Starting server at port 8040\n")
	if err := http.ListenAndServe(":8040", nil); err != nil {
		log.Fatal(err)
	}
}