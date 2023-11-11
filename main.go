package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

type member struct {
	BirthName  string
	Birthday   string
	Birthplace string
	Revealed   string
}

func main() {

	isins := []string{"black-pink", "twice", "girls-generation-snsd"}
	// Create a new Colly collector
	c := colly.NewCollector()

	// Define the URL you want to scrape

	var members []member
	var m member
	groupInfos := make([][]member, 0, 1)
	textToPrintArr := []string{}
	var textToPrint string
	groupName := ""
	groupNames := []string{}

	c.OnHTML("title", func(e *colly.HTMLElement) {
		groupName = strings.TrimSuffix(e.Text, " Members Profile (Updated!)")
		groupNames = append(groupNames, groupName)
	})
	c.OnHTML("meta[property='og:description']", func(e *colly.HTMLElement) {
		textToPrint = ""
		descrip := e.Attr("content")
		if strings.Contains(descrip, "consists of ") || strings.Contains(descrip, "consisting of ") {
			lines := strings.Split(descrip, "Ideal Types ")
			rightDescrip := lines[1]
			finalDescrip := strings.Split(rightDescrip, ".")
			for i := 0; i < 2; i++ {
				textToPrint += string(finalDescrip[i]) + "."
			}
		}
	})

	// Set up a callback function to handle scraping events
	c.OnHTML("body", func(e *colly.HTMLElement) {
		// Extract the text from the span following "Birth Name:"
		text := e.Text
		c := strings.Contains
		if textToPrint == "" {
			if c(text, "consists of ") || c(text, "consisting of ") {
				lines := strings.Split(text, "Members Profile:")
				rightDescrip := lines[1]
				finalDescrip := strings.Split(rightDescrip, ".")
				for i := 0; i < 2; i++ {
					textToPrint += string(finalDescrip[i]) + "."
				}
			}
		}
		textToPrintArr = append(textToPrintArr, textToPrint)
		lines := strings.Split(text, "\n")

		for _, line := range lines {
			if c(line, "Birth Name:") {
				m.BirthName = line
			}
			if c(line, "Birthday:") || c(line, "Birth Date:") {
				m.Birthday = line
			}
			if c(line, "Birthplace:") || c(line, "Birth Place:") {
				m.Birthplace = line
			} else if m.Birthplace == "" && c(line, "was born in") {
				m.Birthplace = line
			}
			if c(line, "member to be revealed") || c(line, "member to be confirmed") {
				m.Revealed = line
				members = append(members, m)
				m = member{} // Create a new member struct
			} else if c(line, "Show more") && c(line, "fun facts") {
				if m != (member{}) {
					members = append(members, m)
					m = member{}
				}
			}
		}
	})

	// Visit the URL and start scraping

	c.OnScraped(func(r *colly.Response) {
		groupInfos = append(groupInfos, members)
		members = []member{}
	})
	for _, isin := range isins {
		err := c.Visit(scrapeUrl(isin))
		if err != nil {
			log.Fatal(err)
		}
	}

	for i, group := range groupInfos {
		s := "The Kpop group"
		w := 110 // or whatever
		fmt.Printf("%[1]*s", -w, fmt.Sprintf("%[1]*s", (w+len(s))/2, s))
		fmt.Println()
		fmt.Println("Group name: " + groupNames[i])
		fmt.Println(textToPrintArr[i] + "\n")
		fmt.Println("Members info:")
		for i, m := range group {
			fmt.Println(m.BirthName)
			fmt.Println(m.Birthday)
			fmt.Println(m.Birthplace)
			fmt.Println(m.Revealed)
			if i != len(group)-1 {
				fmt.Println("*************************************")
			}
		}
	}

	// fmt.Println("\nMembers info")

	// for i, m := range members {
	// 	fmt.Println(m.BirthName)
	// 	fmt.Println(m.Birthday)
	// 	fmt.Println(m.Birthplace)
	// 	fmt.Println(m.Revealed)
	// 	if i != len(members)-1 {
	// 		fmt.Println("*************************************")
	// 	}
	// }
}

func scrapeUrl(isin string) string {
	return "https://kprofiles.com/" + isin + "-members-profile/"
}
