package main

import (
	"html"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"unicode"
)

type Lark struct {
	Articles []Article
}

type Article struct {
	Sections []Section
}

type Section struct {
	Blocks []Block
}

type Block struct {
	Glyph    string
	Contents []string
}

func getGlyph(line string) string {
	// Handle empty lines
	if len(line) == 0 {
		return "paragraph"
	}

	// Handle dividers
	if len(line) > 2 {
		switch line[0:3] {
		case "```":
			return "pre"
		case "'''":
			return "pre"
		case "***":
			return "article"
		case "---":
			return "section"
		}

	}

	char := line[0]
	glyphs := map[byte]string{
		'=': "header",
		'-': "subheader",
		'_': "subsubheader",
		'+': "date",
		'~': "author",
		'@': "link",
		'!': "image",
		'>': "blockquote",
		'*': "ulist",
		':': "olist",
	}
	for k, v := range glyphs {
		if char == k {
			return v
		}
	}
	return "paragraph"
}

func parseLine(line string, glyph string) Block {
	// Get the glyph and store it in a block
	var content string
	block := Block{Glyph: glyph}

	if glyph != "paragraph" {
		content = line[1:]
	} else {
		content = line
	}

	// Personal opinion
	content = strings.Replace(content, "---", "—", -1)
	content = strings.Replace(content, "--", "—", -1)

	block.Contents = []string{strings.TrimLeftFunc(content, unicode.IsSpace)}
	return block
}

func isPre(line string, pre bool) bool {
	if pre && getGlyph(line) != "pre" {
		return true
	}
	if !pre && getGlyph(line) == "pre" {
		return true
	}
	return false
}

func encodeLark(lines []string) Lark {
	lark := Lark{}
	article := Article{}
	section := Section{}
	block := Block{}
	preblock := Block{Glyph: "pre"}
	ublock := Block{Glyph: "ulist"}
	oblock := Block{Glyph: "olist"}

	pre := false

	for _, line := range lines {
		glyph := getGlyph(line)

		pre = isPre(line, pre)
		if pre {
			if glyph != "pre" {
				preblock.Contents = append(preblock.Contents, line)
			}
		} else {
			if glyph == "article" {
				if section.Blocks != nil {
					article.Sections = append(article.Sections, section)
					section.Blocks = nil
				}
				if article.Sections != nil {
					lark.Articles = append(lark.Articles, article)
					article.Sections = nil
				}
			} else if glyph == "section" {
				if section.Blocks != nil {
					article.Sections = append(article.Sections, section)
					section.Blocks = nil
				}
			} else {
				// Handle pre block printing now
				if preblock.Contents != nil {
					section.Blocks = append(section.Blocks, preblock)
					preblock.Contents = nil
				}
				if glyph != "pre" {
					if glyph == "ulist" {
						ublock.Contents = append(ublock.Contents, parseLine(line, glyph).Contents[0])
					} else if glyph == "olist" {
						oblock.Contents = append(oblock.Contents, parseLine(line, glyph).Contents[0])
					} else {
						if ublock.Contents != nil {
							section.Blocks = append(section.Blocks, ublock)
							ublock.Contents = nil
						}

						if oblock.Contents != nil {
							section.Blocks = append(section.Blocks, oblock)
							oblock.Contents = nil
						}

						block = parseLine(line, glyph)
						section.Blocks = append(section.Blocks, block)
					}
				}
			}
		}
	}
	if preblock.Contents != nil {
		section.Blocks = append(section.Blocks, preblock)
	}
	if ublock.Contents != nil {
		section.Blocks = append(section.Blocks, ublock)
	}
	if oblock.Contents != nil {
		section.Blocks = append(section.Blocks, oblock)
	}

	if section.Blocks != nil {
		article.Sections = append(article.Sections, section)
	}
	if article.Sections != nil {
		lark.Articles = append(lark.Articles, article)
	}

	return lark
}

func (b *Block) GetHTMLTags() string {
	tags := map[string]string{
		"header":       "h1",
		"subheader":    "h2",
		"subsubheader": "h3",
		"date":         "h2 class='date'",
		"author":       "h2 class='author'",
		"link":         "a",
		"image":        "a",
		"blockquote":   "blockquote",
		"ulist":        "ul",
		"olist":        "ol",
		"pre":          "pre",
	}
	for k, v := range tags {
		if b.Glyph == k {
			return v
		}
	}
	return "p"
}

func (b *Block) IsLink() bool {
	return b.Glyph == "link"
}

func (b *Block) IsImage() bool {
	return b.Glyph == "image"
}

func (b *Block) IsList() bool {
	return b.Glyph == "ulist" || b.Glyph == "olist"
}

func (b *Block) IsPre() bool {
	return b.Glyph == "pre" || b.Glyph == "code"
}

func (b *Block) EncodeImage() string {
	sep := strings.SplitN(b.Contents[0], " ", 2)
	link := sep[0]
	desc := sep[0]
	if len(sep) == 2 {
		desc = sep[1]
	}

	return "<img src='" + link + "' alt='" + desc + "' loading='lazy' />\n"
}

func (b *Block) EncodeLink() string {
	sep := strings.SplitN(b.Contents[0], " ", 2)
	link := sep[0]
	desc := sep[0]
	if len(sep) == 2 {
		desc = sep[1]
	}

	if link[0] != '/' {
		return "<a href='" + link + "' target='_blank'>" + desc + "</a>\n"
	}

	return "<a href='" + link + "'>" + desc + "</a>\n"
}

func (b *Block) EncodeList() string {
	output := "<" + b.GetHTMLTags() + ">\n"
	for _, content := range b.Contents {
		output += "\t<li>" + content + "</li>\n"
	}
	output += "</" + b.GetHTMLTags() + ">\n"

	return output
}

func (b *Block) EncodePre() string {
	output := "<" + b.GetHTMLTags() + ">\n"
	for _, content := range b.Contents {
		output += content + "\n"
	}
	output += "</" + b.GetHTMLTags() + ">\n"

	return output
}

func (b *Block) EscapeString(s string) string {
	return html.EscapeString(s)
}

func (a *Article) GetID() string {
	for _, s := range a.Sections {
		for _, b := range s.Blocks {
			if b.Glyph == "header" {
				str := b.Contents[0]
				str = regexp.MustCompile(`[^a-zA-Z0-9 \-_]+`).ReplaceAllString(str, "")
				str = strings.ToLower(strings.Replace(str, " ", "-", -1))
				return str
			}
			if b.Glyph == "subheader" {
				str := b.Contents[0]
				str = regexp.MustCompile(`[^a-zA-Z0-9 \-_]+`).ReplaceAllString(str, "")
				str = strings.ToLower(strings.Replace(str, " ", "-", -1))
				return str
			}
		}
	}

	return ""
}

func parseLarkToHTML(lark Lark) string {
	output := ""

	for _, article := range lark.Articles {
		output += "<article>\n"

		for _, section := range article.Sections {
			if len(article.Sections) > 1 {
				output += "<section>\n"
			}

			for _, block := range section.Blocks {
				// Handle images and links
				if block.Glyph == "link" {
					sep := strings.SplitN(block.Contents[0], " ", 2)
					link := sep[0]
					desc := sep[0]
					if len(sep) == 2 {
						desc = sep[1]
					}
					output += "<p><a href='" + link + "'>" + desc + "</a></p>\n"
				} else if block.Glyph == "image" {
					sep := strings.SplitN(block.Contents[0], " ", 2)
					link := sep[0]
					desc := sep[0]
					if len(sep) == 2 {
						desc = sep[1]
					}

					output += "<img src='" + link + "' alt='" + desc + "' />\n"
				} else {
					output += "<" + block.GetHTMLTags() + ">"

					for _, content := range block.Contents {
						if block.Glyph == "pre" {
							output += html.EscapeString(content) + "\n"
						} else if block.Glyph == "ulist" || block.Glyph == "olist" {
							output += "\n<li>" + content + "</li>\n"
						} else {
							if content != "" {
								output += "\n" + content + "\n"
							}
						}
					}
					output += "</" + strings.Split(block.GetHTMLTags(), " ")[0] + ">\n"
				}

			}

			if len(article.Sections) > 1 {
				output += "</section>\n"
			}
		}

		output += "</article>\n"
	}

	return output
}

func handleLark(w http.ResponseWriter, r *http.Request) {
	// Read in a text file containing TGT
	content, err := ioutil.ReadFile("./static/lark/test.lark")
	if err != nil {
		panic(err)
	}
	text := string(content)

	lines := strings.Split(text, "\n")

	lark := encodeLark(lines)

	executeTemplate(w, "three-good-things.tmpl", lark)
}
