#!/usr/bin/python3
from bs4 import BeautifulSoup as soup
import xml.etree.ElementTree as et

# SET THESE

filename = "index.html"
title = "HD-DN Three Good Things"
url = "https://hd-dn.com/three-good-things/"

clean = lambda x: x.text.strip()

def html_to_soups(html):
  with open(html, "r") as f:
    page_soup = soup(f, "html.parser")
    # title = clean(page_soup.find("h1"))
    articles = page_soup.find_all("article")

  return articles

def articles_to_xml_entries(title, articles, url):
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

    # Set the published and updated dates
    el_pub = sub_element("published")
    el_upd = sub_element("updated")
    el_pub.text = ds
    el_upd.text = ds

    # Set the content
    article.find("h2").decompose()
    el_content = sub_element("content")
    el_content.set("type", "html")
    el_content.text = "<h2>{}</h2>".format(el_title.text)
    el_content.text += str(article).replace("\n", "")[9:-10]

    # Insert into entries list
    entries.append(entry)
  return entries

def entries_to_xml(entries, url, title):
  out = ""
  out += ("""<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
    <title>{title}</title>
    <id>{url}</id>
    <link rel="alternate" href="{url}"></link>
    <updated>{last}T00:00:00+00:00</updated>
    <author>
        <name>{title}</name>
        <uri>{url}</uri>
    </author>""".format(url=url, title=title, last=entries[0][1].text))

  for entry in entries:
    out += "{}".format(et.tostring(entry).decode())

  out += """    </entry>
</feed>"""
    
  return out

def write_xml(xmlstr, outfile):
  xml_data = soup(xmlstr, "xml")
  with open(outfile, "w") as f:
    f.write(xml_data.prettify())

if __name__ == "__main__":
  articles = html_to_soups(filename)
  entries = articles_to_xml_entries(title, articles, url)
  write_xml(entries_to_xml(entries, url, title), "atom.atom")
  