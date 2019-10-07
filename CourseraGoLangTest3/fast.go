package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	seenBrowsers := []string{}
	uniqueBrowsers := 0

	var strBuilder strings.Builder
	strBuilder.WriteString("found users:\n")

	type StrData struct {
		Browsers []string "json:browsers"
		Email    string   "json:email"
		Name     string   "json:name"
	}

	r := regexp.MustCompile("@")
	regexp1, err1 := regexp.Compile("Android")
	regexp2, err2 := regexp.Compile("MSIE")
	//user := make(map[string]interface{})

	fileScanText := bufio.NewScanner(file)
	var dataline StrData

	for i := 0; fileScanText.Scan(); i++ {

		tmpbyte := fileScanText.Bytes()
		err = json.Unmarshal(tmpbyte, &dataline)

		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false
		notSeenBefore := true

		//		browsers, ok := user["browsers"].([]interface{})
		browsers := dataline.Browsers

		// if !ok {
		// 	// log.Println("cant cast browsers")
		// 	continue
		// }

		for _, browserRaw := range browsers {

			//			if ok, err := regexp.MatchString("Android", browser); ok && err == nil {
			if ok := regexp1.MatchString(browserRaw); ok && err1 == nil {
				isAndroid = true
				notSeenBefore = true
				for _, item := range seenBrowsers {
					if item == browserRaw {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browserRaw)
					uniqueBrowsers++
				}
			}

			if ok := regexp2.MatchString(browserRaw); ok && err2 == nil {
				isMSIE = true
				notSeenBefore = true
				for _, item := range seenBrowsers {
					if item == browserRaw {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browserRaw)
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		email := r.ReplaceAllString(dataline.Email, " [at] ")
		strBuilder.WriteString(fmt.Sprintf("[%d] %s <%s>\n", i, dataline.Name, email))
	}

	fmt.Fprintln(out, strBuilder.String())
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))

}
