package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// "browsers":["Lynx/2.8.5rel.1 libwww-FM/2.14 SSL-MM/1.4.1 GNUTLS/0.8.12","Mozilla/5.0 (iPad; CPU OS 10_0 like Mac OS X) AppleWebKit/601.1 (KHTML, like Gecko) CriOS/49.0.2623.109 Mobile/14A5335b Safari/601.1.46","Mozilla/4.0 (compatible; MSIE 6.0; Windows CE; IEMobile 7.11) Sprint:PPC6800","Mozilla/5.0 (WindowsCE 6.0; rv:2.0.1) Gecko/20100101 Firefox/4.0.1"],
// "company":"Jaxbean",
// "country":"Syria",
// "email":"optio_in_adipisci@Photobug.mil",
// "job":"Sales Associate",
// "name":"James Griffin",
// "phone":"1-310-576-97-47"}

type user struct {
	Browsers []string `json:"browsers"`
	Company  string   `json:"company"`
	Country  string   `json:"country"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
}

func changeEmail(e string) string {
	mstr := strings.Split(e, "@")
	if len(mstr) > 1 {
		return mstr[1] + " [at] " + mstr[0]
	}
	return ""
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	// SlowSearch(out)
	var u user
	var strb strings.Builder

	seenBrowsers := make(map[string]bool)
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	scan := bufio.NewScanner(file)
	if err != nil {
		panic(err)
	}

	for i := 0; scan.Scan(); i++ {

		isAndroid := false
		isMSIE := false
		str := scan.Bytes()
		json.Unmarshal(str, &u)
		for _, val := range u.Browsers {
			if strings.Contains(val, "Android") {
				seenBrowsers[val] = true
				isAndroid = true
			}
			if strings.Contains(val, "MSIE") {
				seenBrowsers[val] = true
				isMSIE = true
			}
		}

		if isAndroid && isMSIE {
			// resstr := fmt.Sprintf("[%d] %s <%s>\n", i, u.Name, changeEmail(u.Email))
			resstr := fmt.Sprintf("[%d] %s <%s>\n", i, u.Name, strings.Replace(u.Email, "@", " [at] ", 1))
			strb.WriteString(resstr)
		}

	}
	fmt.Fprintln(out, "found users:\n"+strb.String())
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
