package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/simon3z/golang-bugzilla" // cspell:disable-line
)

var cmdFlags = struct {
	BugzillaURL      string
	BugzillaUsername string
	BugzillaPassword string
}{}

func main() {
	scan := bufio.NewScanner(os.Stdin)

	cmdFlags.BugzillaURL = os.Getenv("BZURL")
	cmdFlags.BugzillaUsername = os.Getenv("BZUSER")
	cmdFlags.BugzillaPassword = os.Getenv("BZPASS")

	bz, err := bugzilla.NewClient(cmdFlags.BugzillaURL, cmdFlags.BugzillaUsername, cmdFlags.BugzillaPassword)

	if err != nil {
		log.Fatalln(err)
	}

	w := csv.NewWriter(os.Stdout)
	w.Comma = '\t'

	w.Write([]string{
		"ID",
		"Product",
		"Component",
		"Assignee",
		"Status",
		"Summary",
		"Pillar",
		"Severity",
		"Target Release",
		"Cases",
	})

	for scan.Scan() {
		line := scan.Text()

		id, err := strconv.Atoi(line)

		if err != nil {
			log.Println("cannot parse:", err)
			continue
		}

		bugInfo, err := bz.BugInfo(id, &bugzilla.BugInfoOptions{IncludeFields: []string{"_default", "external_bugs"}})

		if err != nil {
			log.Println("cannot get bug info:", err)
			continue
		}

		//fmt.Println(bugInfo)

		w.Write([]string{
			fmt.Sprintf("=HYPERLINK(\"%s/%d\", \"%d\")", cmdFlags.BugzillaURL, id, id),
			shortString(bugInfo["product"].(string), 30),
			bugInfo["component"].([]interface{})[0].(string),
			bugInfo["assigned_to"].(string),
			bugInfo["status"].(string),
			shortString(bugInfo["summary"].(string), 50),
			bugInfo["cf_internal_whiteboard"].(string),
			bugInfo["severity"].(string),
			bugInfo["target_release"].([]interface{})[0].(string),
			// fmt.Sprintf("%#v", bugInfo["external_bugs"].([]interface{})),
			countExternalBugs(bugInfo["external_bugs"]),
		})

		w.Flush()
	}

	if scan.Err() != nil {
		log.Fatalln(scan.Err())
	}
}

func countExternalBugs(extbz interface{}) string {
	cases := []string{}

	for _, i := range extbz.([]interface{}) {
		ii := i.(map[string]interface{})
		typ := ii["type"].(map[string]interface{})
		if typ["type"].(string) == "SFDC" {
			cases = append(cases, fmt.Sprintf("=HYPERLINK(\"%s\", \"%s\")",
				strings.ReplaceAll(typ["full_url"].(string), "%id%", ii["ext_bz_bug_id"].(string)),
				ii["ext_bz_bug_id"].(string)))
		}
	}

	if len(cases) == 1 {
		return cases[0]
	} else if len(cases) > 1 {
		return fmt.Sprintf("#%d cases", len(cases))
	}

	return ""
}

func shortString(s string, n int) string {
	if len(s) > n {
		return s[:n-1] + "…"
	}

	return s
}
