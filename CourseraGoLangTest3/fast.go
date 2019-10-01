package main

import (
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

	r := regexp.MustCompile("@")
	seenBrowsers := []string{}
	uniqueBrowsers := 0
	foundUsers := ""

	var strBuilder strings.Builder

	type strData struct {
		browsers []string
		email    string
		name     string
	}

	regexp1, err1 := regexp.Compile("Android")
	regexp2, err2 := regexp.Compile("MSIE")
	//user := make(map[string]interface{})

	jsonReader := json.NewDecoder(file)

	t, err := jsonReader.Token()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%T: %v\n", t, t)

	var strtest interface{}
	
	err = jsonReader.Decode(&strtest)
	fmt.Printf("%v", strtest)

	if err != nil {
		panic(err)
	}

	jsonReader.More()

	t, err = jsonReader.Token()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%T: %v\n", t, t)

	//fileScanText := bufio.NewScanner(file)

	for i := 0; jsonReader.More(); i++ {


		var dataline strData

		err = jsonReader.Decode(&dataline)
		fmt.Printf("%v", dataline)

		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		//		browsers, ok := user["browsers"].([]interface{})
		browsers := dataline.browsers

		// if !ok {
		// 	// log.Println("cant cast browsers")
		// 	continue
		// }

		for _, browserRaw := range browsers {

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

		for _, browserRaw := range browsers {

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
		//email := r.ReplaceAllString(dataline.email, " [at] ")
		email := r.ReplaceAllString("   ", " [at] ")
		//		email := r.ReplaceAllString(user["email"].(string), " [at] ")
		//		strBuilder.WriteString(fmt.Sprintf("[%d] %s <%s>\n", i, dataline.name, email))
		strBuilder.WriteString(fmt.Sprintf("[%d] %s <%s>\n", i, "dataline.name", email))
		//foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user["name"], email)
	}

	foundUsers = strBuilder.String()

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
