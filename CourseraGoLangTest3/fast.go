package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
)

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	r := regexp.MustCompile("@")
	seenBrowsers := []string{}
	uniqueBrowsers := 0

	var strBuilder strings.Builder
	strBuilder.WriteString("found users:\n")


	type Strdata struct {
		Browsers []string "json:browsers"
		Email    string "json:email"
		Name     string "json:name"
	}

	var dataPool = sync.Pool{
		New: func() interface{} {
			return &Strdata{}
		},
	}

	regexp1, err1 := regexp.Compile("Android")
	regexp2, err2 := regexp.Compile("MSIE")

    scanfile := bufio.NewScanner(file)

	for i:=0;scanfile.Scan();i++{

		dataline := dataPool.Get().(*Strdata)
		

		err := json.Unmarshal(scanfile.Bytes(),dataline)

		if err != nil {
			panic(err)
		}

//		fmt.Printf("%v \n", dataline)


		isAndroid := false
		isMSIE := false

		//		browsers, ok := user["browsers"].([]interface{})
		//browsers := dataline.Browsers

		// if !ok {
		// 	// log.Println("cant cast browsers")
		// 	continue
		// }

		for _, browserRaw := range dataline.Browsers {

			browser := browserRaw

			//browser, ok := browserRaw
			// if !ok {
			// 	// log.Println("cant cast browser to string")
			// 	continue
			// }

			//			if ok, err := regexp.MatchString("Android", browser); ok && err == nil {
			if ok := regexp1.MatchString(browser); ok && err1 == nil {
				isAndroid = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		for _, browserRaw := range dataline.Browsers {

			browser := browserRaw
			// browser, ok := browserRaw.(string)
			// if !ok {
			// 	// log.Println("cant cast browser to string")
			// 	continue
			// }
			//			if ok, err := regexp.MatchString("MSIE", browser); ok && err == nil {

			if ok := regexp2.MatchString(browser); ok && err2 == nil {
				isMSIE = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
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
		//foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user["name"], email)
		
		dataPool.Put(dataline)
	}

	fmt.Fprintln(out, strBuilder.String())
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
