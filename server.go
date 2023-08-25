package main

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
)

type Page struct {
	Title   string
	Content string
}

type LarkPage struct {
	Title string
	Lark  Lark
}

// Define a Day struct to hold TGT posts
type Day struct {
	Title  string
	First  string
	Second string
	Third  string
}

// Color struct for daily colors
type Color struct {
	ColorName string
	ColorHex  string
}

type Photo struct {
	Path    string
	AltText string
}

type PhotoEntry struct {
	Date      string
	FilmStock string
	Notes     []string
	Photos    []Photo
}

// Running journal entries
type Run struct {
	Title    string
	Distance string
	Notes    []string
}

// const staticDir = "/home/debian/hd-dn/static"
const staticDir = "./static"

// execute templates from templatedir
func executeTemplate(w http.ResponseWriter, tmplFile string, data interface{}) {
	templates := template.Must(template.ParseGlob("./templates/*"))
	// Execute the template
	err := templates.ExecuteTemplate(w, tmplFile, data)
	if err != nil {
		panic(err)
	}
}

func parseDays(filename string, escape bool) Lark {
	// Read in a text file containing TGT
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	text := string(content)
	lines := strings.Split(text, "\n")
	lark := encodeLark(lines)
	return lark
}

// Handler for three good things posts
func handleThreeGoodThings(w http.ResponseWriter, r *http.Request) {
	// Cat the files together
	days := parseDays("./static/three-good-things/index.lark", false)
	executeTemplate(w, "three-good-things.tmpl", days)
}

// Build atom feed for three good things
func handleThreeGoodThingsFeed(w http.ResponseWriter, r *http.Request) {
	// Split it into posts based on pagebreak elements ***
	days := parseDays("./static/three-good-things/index.lark", false)

	executeTemplate(w, "three-good-things-feed.tmpl", days)
}

// get the proper title name from the file path
func getTitleFromURL(url string, suffix string) string {
	// Get list of all parts of the filepath
	ps := strings.Split(strings.Trim(url, "/"), "/")
	if len(ps) == 0 {
		return ""
	}

	// Get the last entry as a candidate
	c := ps[len(ps)-1]

	if c == "index"+suffix {
		if len(ps) == 1 {
			return ""
		}
		return ps[len(ps)-2]
	}

	return strings.TrimSuffix(c, suffix)
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
	splits := strings.Split(r.URL.Path, ".")
	suffix := splits[len(splits)-1]
	if suffix == r.URL.Path {
		suffix = "html"
		r.URL.Path += "index.html"
	}
	rest := strings.TrimSuffix(r.URL.Path, suffix)

	title := getTitleFromURL(r.URL.Path, "."+suffix)
	if title != "" {
		title = "HD-DN: " + strings.ToLower(title)
	} else {
		title = "HD-DN"
	}

	if suffix == "html" {
		// Check to see if there's a lark file
		if _, err := os.Stat("./static" + rest + "lark"); os.IsNotExist(err) {
			// If not, read in the HTML file
			content, err := readURL("./static" + r.URL.Path)
			if err != nil {
				panic(err)
			}

			text := string(content)
			page := Page{
				Title:   title,
				Content: text,
			}

			executeTemplate(w, "page.tmpl", page)
		} else {
			// If there is, use it
			content, err := readURL("./static" + rest + "lark")
			if err != nil {
				panic(err)
			}

			text := string(content)
			page := LarkPage{
				Title: title,
				Lark:  encodeLark(strings.Split(text, "\n")),
			}

			executeTemplate(w, "larkpage.tmpl", page)
		}
	}

	if suffix == "lark" {
		content, err := readURL("./static" + r.URL.Path)
		if err != nil {
			panic(err)
		}

		text := string(content)
		page := LarkPage{
			Title: title,
			Lark:  encodeLark(strings.Split(text, "\n")),
		}

		executeTemplate(w, "larkpage.tmpl", page)
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
	return colors[h%len(colors)]
}

// Handler for color of the day
func handleColor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	executeTemplate(w, "color.tmpl", getDailyColor())
}

// Photography
// Parse a photography entry
func parsePhotoEntry(path string, entry string) PhotoEntry {
	pe := PhotoEntry{}
	photos := false

	path = "/img/" + strings.Trim(path, ".html") + "/"

	ls := strings.Split(entry, "\n")

	for _, l := range ls {
		if len(l) > 1 {
			check := l[0:2]
			rest := l[2:len(l)]
			rest = strings.TrimLeft(rest, " ")
			if check == ":p" && !photos {
				photos = true
			}

			if !photos {
				if check == "d " {
					pe.Date = rest
				} else if check == "f " {
					pe.FilmStock = rest
				} else {
					pe.Notes = append(pe.Notes, l)
				}
			} else if check != ":p" {
				s := strings.Split(l, "|")
				photo := Photo{Path: path + strings.Trim(s[0], " ") + ".jpg", AltText: strings.Trim(s[1], " ")}
				pe.Photos = append(pe.Photos, photo)
			}
		}
	}
	return pe
}

func handlePhotos(w http.ResponseWriter, r *http.Request) {
	// Read in a text file containing photos
	path := strings.TrimSuffix(r.URL.Path, ".html") + ".txt"
	ss := strings.Split(r.URL.Path, "/")
	fp := ss[len(ss)-1]
	if fp == "index.html" || fp == "" {
		handlePage(w, r)
		return
	}

	content, err := readURL("./static" + path)
	if err != nil {
		panic(err)
	}

	text := string(content)
	executeTemplate(w, "photos.tmpl", parsePhotoEntry(fp, text))
}

// Running journal
// Parse a run entry
func parseRun(entry string) Run {
	var run = Run{
		Title:    "",
		Distance: "",
		Notes:    []string{},
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
	if err != nil {
		panic(err)
	}
	text := string(content)

	// Split it into posts based on pagebreak elements ***
	var runs []Run
	posts := strings.Split(text, "***")
	for _, post := range posts {
		runs = append(runs, parseRun(post))
	}
	executeTemplate(w, "runs.tmpl", runs)
}

// Seasonal CSS
func setSeason() string {
	now := time.Now()
	current_year := now.Year()

	spring_start := time.Date(current_year, 3, 1, 0, 0, 0, 0, time.Local)
	summer_start := time.Date(current_year, 6, 1, 0, 0, 0, 0, time.Local)
	autumn_start := time.Date(current_year, 9, 1, 0, 0, 0, 0, time.Local)
	winter_start := time.Date(current_year, 12, 1, 0, 0, 0, 0, time.Local)

	// Spring
	if now.After(spring_start) && now.Before(summer_start) {
		return "spring"
	}
	if now.After(summer_start) && now.Before(autumn_start) {
		return "summer"
	}
	if now.After(autumn_start) && now.Before(winter_start) {
		return "autumn"
	}
	if now.After(winter_start) || now.Before(spring_start) {
		return "winter"
	}

	return "style"
}

func handleSeasonalCSS(w http.ResponseWriter, r *http.Request) {
	content, _ := ioutil.ReadFile("./static/css/" + setSeason() + ".css")
	w.Header().Set("Content-Type", "text/css")
	w.Write(content)
}

// Main server function
func main() {
	fileServer := http.FileServer(http.Dir(staticDir))

	http.Handle("/css/", fileServer)
	http.HandleFunc("/css/style.css", handleSeasonalCSS)
	http.Handle("/img/", fileServer)
	http.Handle("/favicon.ico", fileServer)
	http.Handle("/three-good-things/atom.atom", fileServer)
	http.HandleFunc("/", handlePage)

	// Breathe handling
	http.Handle("/breathe/", fileServer)

	// do colo(u)r of the day
	http.HandleFunc("/color/", handleColor)
	http.HandleFunc("/colour/", handleColor)

	http.HandleFunc("/runs/", handleRuns)
	http.HandleFunc("/photography/film/", handlePhotos)

	fmt.Printf("Starting server at http://localhost:8040\n")
	if err := http.ListenAndServe(":8040", nil); err != nil {
		log.Fatal(err)
	}
}
