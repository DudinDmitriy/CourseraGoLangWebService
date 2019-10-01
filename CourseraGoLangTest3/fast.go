package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

	r := regexp.MustCompile("@")
	seenBrowsers := []string{}
	uniqueBrowsers := 0
	foundUsers := ""

	regexp1, err1 := regexp.Compile("Android")
	regexp2, err2 := regexp.Compile("MSIE")
	user := make(map[string]interface{})
	
	fileScanText := bufio.NewScanner(file)
	for fileScanText.Scan() {

		strByte := fileScanText.Bytes()
		err = json.Unmarshal(strByte, &user)
		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		browsers, ok := user["browsers"].([]interface{})
		if !ok {
			// log.Println("cant cast browsers")
			continue
		}

		for _, browserRaw := range browsers {
			browser, ok := browserRaw.(string)
			if !ok {
				// log.Println("cant cast browser to string")
				continue
			}

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

		for _, browserRaw := range browsers {
			browser, ok := browserRaw.(string)
			if !ok {
				// log.Println("cant cast browser to string")
				continue
			}
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
		email := r.ReplaceAllString(user["email"].(string), " [at] ")
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user["name"], email)
	}

	//нужно читать файл построчно
	// fileContents, err := ioutil.ReadAll(file)
	// if err != nil {
	// 	panic(err)
	// }
    //затем каждую строку обрабатывать
	// r := regexp.MustCompile("@")
	// seenBrowsers := []string{}
	// uniqueBrowsers := 0
	// foundUsers := ""

	// regexp1, err1 := regexp.Compile("Android")
	// regexp2, err2 := regexp.Compile("MSIE")

	// lines := strings.Split(string(fileContents), "\n")

	//users := make([]map[string]interface{}, 0,)

	// users := make([]map[string]interface{}, 0, len(lines))

	// for _, line := range lines {
	// 	user := make(map[string]interface{})
	// 	// fmt.Printf("%v %v\n", err, line)
	// 	err := json.Unmarshal([]byte(line), &user)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	isAndroid := false
	// 	isMSIE := false

	// 	browsers, ok := user["browsers"].([]interface{})
	// 	if !ok {
	// 		// log.Println("cant cast browsers")
	// 		continue
	// 	}

	// 	for _, browserRaw := range browsers {
	// 		browser, ok := browserRaw.(string)
	// 		if !ok {
	// 			// log.Println("cant cast browser to string")
	// 			continue
	// 		}

	// 		//			if ok, err := regexp.MatchString("Android", browser); ok && err == nil {
	// 		if ok := regexp1.MatchString(browser); ok && err1 == nil {
	// 			isAndroid = true
	// 			notSeenBefore := true
	// 			for _, item := range seenBrowsers {
	// 				if item == browser {
	// 					notSeenBefore = false
	// 				}
	// 			}
	// 			if notSeenBefore {
	// 				// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
	// 				seenBrowsers = append(seenBrowsers, browser)
	// 				uniqueBrowsers++
	// 			}
	// 		}
	// 	}

	// 	for _, browserRaw := range browsers {
	// 		browser, ok := browserRaw.(string)
	// 		if !ok {
	// 			// log.Println("cant cast browser to string")
	// 			continue
	// 		}
	// 		//			if ok, err := regexp.MatchString("MSIE", browser); ok && err == nil {

	// 		if ok := regexp2.MatchString(browser); ok && err2 == nil {
	// 			isMSIE = true
	// 			notSeenBefore := true
	// 			for _, item := range seenBrowsers {
	// 				if item == browser {
	// 					notSeenBefore = false
	// 				}
	// 			}
	// 			if notSeenBefore {
	// 				// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
	// 				seenBrowsers = append(seenBrowsers, browser)
	// 				uniqueBrowsers++
	// 			}
	// 		}
	// 	}

	// 	if !(isAndroid && isMSIE) {
	// 		continue
	// 	}

	// 	// log.Println("Android and MSIE user:", user["name"], user["email"])
	// 	email := r.ReplaceAllString(user["email"].(string), " [at] ")
	// 	foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user["name"], email)
	// }

	// for i, user := range users {

	// 	isAndroid := false
	// 	isMSIE := false

	// 	browsers, ok := user["browsers"].([]interface{})
	// 	if !ok {
	// 		// log.Println("cant cast browsers")
	// 		continue
	// 	}

	// 	for _, browserRaw := range browsers {
	// 		browser, ok := browserRaw.(string)
	// 		if !ok {
	// 			// log.Println("cant cast browser to string")
	// 			continue
	// 		}

	// 		//			if ok, err := regexp.MatchString("Android", browser); ok && err == nil {
	// 		if ok := regexp1.MatchString(browser); ok && err1 == nil {
	// 			isAndroid = true
	// 			notSeenBefore := true
	// 			for _, item := range seenBrowsers {
	// 				if item == browser {
	// 					notSeenBefore = false
	// 				}
	// 			}
	// 			if notSeenBefore {
	// 				// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
	// 				seenBrowsers = append(seenBrowsers, browser)
	// 				uniqueBrowsers++
	// 			}
	// 		}
	// 	}

	// 	for _, browserRaw := range browsers {
	// 		browser, ok := browserRaw.(string)
	// 		if !ok {
	// 			// log.Println("cant cast browser to string")
	// 			continue
	// 		}
	// 		//			if ok, err := regexp.MatchString("MSIE", browser); ok && err == nil {

	// 		if ok := regexp2.MatchString(browser); ok && err2 == nil {
	// 			isMSIE = true
	// 			notSeenBefore := true
	// 			for _, item := range seenBrowsers {
	// 				if item == browser {
	// 					notSeenBefore = false
	// 				}
	// 			}
	// 			if notSeenBefore {
	// 				// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
	// 				seenBrowsers = append(seenBrowsers, browser)
	// 				uniqueBrowsers++
	// 			}
	// 		}
	// 	}

	// 	if !(isAndroid && isMSIE) {
	// 		continue
	// 	}

	// 	// log.Println("Android and MSIE user:", user["name"], user["email"])
	// 	email := r.ReplaceAllString(user["email"].(string), " [at] ")
	// 	foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user["name"], email)
	// }

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
