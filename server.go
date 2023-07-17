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
	"os"
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

// execute templates from templatedir
func executeTemplate(w http.ResponseWriter, tmplFile string, data interface{}) {
	tmpl, err := template.New(tmplFile).ParseGlob("./templates/" + tmplFile)
	if err != nil {
		panic(err)
	}
	// Execute the template
	err = tmpl.ExecuteTemplate(w, tmplFile, data)
	if err != nil{
		panic(err)
	}
}

// parses three good things entry into struct
func parseDay(entry string, escape bool) Day {
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

			if escape {
				rest = template.HTMLEscapeString(rest)
			}

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

func parseDays(filename string, escape bool) []Day {
	// Read in a text file containing TGT
	content, err := ioutil.ReadFile(filename)
	if (err != nil) {
		panic(err);
	}
	text := string(content)

	// Split it into posts based on pagebreak elements ***
	var days []Day
	posts := strings.Split(text, "***")
	for _, post := range posts {
		days = append(days, parseDay(post, escape));
	}
	return days
}

// Handler for three good things posts
func handleThreeGoodThings(w http.ResponseWriter, r *http.Request) {
	// Split it into posts based on pagebreak elements ***
	days := parseDays("./static/three-good-things/index.txt", false)
	executeTemplate(w, "three-good-things.tmpl", days)
}

// Build atom feed for three good things
func handleThreeGoodThingsFeed(w http.ResponseWriter, r *http.Request) {
	days := parseDays("./static/three-good-things/index.txt", true)
	executeTemplate(w, "three-good-things-feed.tmpl", days)
}


// get the proper title name from the file path
func getTitleFromURL(url string) string {
	// Get list of all parts of the filepath
	ps := strings.Split(strings.Trim(url, "/"), "/")
	if len(ps) == 0 { return "" }

	// Get the last entry as a candidate
	c := ps[len(ps) - 1]

	if c == "index.html" {
		if len(ps) == 1 { return ""}
		return ps[len(ps) - 2]
	}

	return strings.TrimSuffix(c, ".html")
}

func readURL(url string) ([]byte, error) {
	info, err := os.Stat(url)
	if err != nil {
		panic(err)
	}

	if info.IsDir() {
		return ioutil.ReadFile(url + "index.html")
	} else {
		return ioutil.ReadFile(url)
	}

}

// Serve a standard html style page
func handlePage(w http.ResponseWriter, r *http.Request) {
	content, err := readURL("./static" + r.URL.Path)
	if err != nil {
		panic(err)
	}
	title := getTitleFromURL(r.URL.Path)
	if title != "" { 
		title = "HD-DN: " + strings.ToLower(title)
	} else { 
		title = "HD-DN" 
	} 

	page := Page{
		Title: title,
		Content: string(content),
	}

	executeTemplate(w, "page.tmpl", page)
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
	w.Header().Set("Content-Type", "text/html")
	executeTemplate(w, "color.tmpl", getDailyColor())
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
	executeTemplate(w, "runs.tmpl", runs)
}

func main() {
	fileServer := http.FileServer(http.Dir(staticDir))
	
	http.Handle("/css/", fileServer)
	http.Handle("/img/", fileServer)
	http.Handle("/favicon.ico", fileServer)
	
	http.HandleFunc("/", handlePage)


	// Handle Three Good Things separately coz she's special
	http.HandleFunc("/three-good-things/", handleThreeGoodThings)
	http.HandleFunc("/three-good-things/atom.atom", handleThreeGoodThingsFeed)
	
	// do colo(u)r of the day
	http.HandleFunc("/color/", handleColor)
	http.HandleFunc("/colour/", handleColor)

	http.HandleFunc("/runs/", handleRuns)


	fmt.Printf("Starting server at http://localhost:8040\n")
	if err := http.ListenAndServe(":8040", nil); err != nil {
		log.Fatal(err)
	}
}