from bs4 import BeautifulSoup as soup
import xml.etree.ElementTree as et

clean = lambda x: x.text.strip()

def html_to_soups(html):
  with open(html, "r") as f:
    page_soup = soup(f, "html.parser")
    title = clean(page_soup.find("h1"))
    articles = page_soup.find_all("article")

    # Handle relative pathing of image names
    for article in articles:
      imgs = article.find_all("img")
      for img in imgs:
        img["src"] = "https://hd-dn.com{}".format(img["src"])

  return title, articles

def articles_to_xml_entries(title, articles):
  entries = []
  format_date = lambda x: x + "T00:00:00+00:00"
  for article in articles:
    sub_element = lambda x: et.SubElement(entry, x)
    pid = clean(article.find("h2"))
    ds = format_date(pid)

    # Create an entry
    entry = et.Element("entry")

    # Set the title
    el_title = sub_element("title")
    el_title.text = pid

    # Set the post id
    el_pid = sub_element("id")
    el_pid.text = "https://hd-dn.com/three-good-things#{}".format(pid)

    # Set the published and updated dates
    el_pub = sub_element("published")
    el_upd = sub_element("updated")
    el_pub.text = ds
    el_upd.text = ds

    # Set the content
    el_content = sub_element("content")
    el_content.set("type", "html")
    el_content.text = str(article.find("ol")).replace("\n", "")

    # Insert into entries list
    entries.append(entry)
  return entries

def entries_to_xml(entries):
  out = ""
  out += ("""<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
    <title>Three Good Things</title>
    <id>https://hd-dn.com/three-good-things</id>
    <link rel="alternate" href="https://hd-dn.com/three-good-things"></link>
    <updated>{}T00:00:00+00:00</updated>
    <author>
        <name>Three Good Things</name>
        <uri>https://hd-dn.com/three-good-things</uri>
    </author>""".format(entries[0][1].text))

  for entry in entries:
    out += "{}".format(et.tostring(entry).decode())

  out += """    </entry>
</feed>"""
    
  return out

def write_xml(xmlstr, filename):
  xml_data = soup(xmlstr, "xml")
  with open(filename, "w") as f:
    f.write(xml_data.prettify())

if __name__ == "__main__":
  title, articles = html_to_soups("index.html")
  entries = articles_to_xml_entries(title, articles)
  write_xml(entries_to_xml(entries), "atom.atom")
  