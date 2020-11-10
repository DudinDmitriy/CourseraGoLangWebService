package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// "browsers":["Lynx/2.8.5rel.1 libwww-FM/2.14 SSL-MM/1.4.1 GNUTLS/0.8.12","Mozilla/5.0 (iPad; CPU OS 10_0 like Mac OS X) AppleWebKit/601.1 (KHTML, like Gecko) CriOS/49.0.2623.109 Mobile/14A5335b Safari/601.1.46","Mozilla/4.0 (compatible; MSIE 6.0; Windows CE; IEMobile 7.11) Sprint:PPC6800","Mozilla/5.0 (WindowsCE 6.0; rv:2.0.1) Gecko/20100101 Firefox/4.0.1"],
// "company":"Jaxbean",
// "country":"Syria",
// "email":"optio_in_adipisci@Photobug.mil",
// "job":"Sales Associate",
// "name":"James Griffin",
// "phone":"1-310-576-97-47"}

// type user struct {
// 	Browsers []string `json:"browsers"`
// 	Company  string   `json:"company"`
// 	Country  string   `json:"country"`
// 	Email    string   `json:"email"`
// 	Name     string   `json:"name"`
// 	Phone    string   `json:"phone"`
// }

func changeEmail(e string) string {
	mstr := strings.Split(e, "@")
	if len(mstr) > 1 {
		return mstr[1] + " [at] " + mstr[0]
	}
	return ""
}

var poolBytes = sync.Pool{
	New: func() interface{} { return []byte{} },
}

// func unmarshallTxt(str string) (res user) {
// 	res = user{}
// 	return res
// }

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	// SlowSearch(out)
	var strName, strEmail string
	// var str string

	seenBrowsers := make(map[string]bool)
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	isAndroid := false
	isMSIE := false
	isAndroidb := false
	isMSIEb := false

	byte1 := byte('"')

	fmt.Fprintln(out, "found users:")

	// scan := bufio.NewScanner(file)
	r := bufio.NewReader(file)
	next := true
	status := 0
	b := make([]byte, 0, 100)
	for i := 0; next; i++ {

		b, _, err = r.ReadLine()
		if err != nil {
			next = false
		}
		iBBrow := 0
		iEBrow := 0
		iBEm1 := 0
		iEEm1 := 0
		iBEm2 := 0
		iEEm2 := 0
		iBN := 0
		iEN := 0

		for ind, oneb := range b {
			if status == 32 && oneb == byte('{') {
				status = 0
			} else if status == 0 && oneb == byte1 {
				status = 1
			} else if status == 1 && oneb == byte1 {
				status = 2 // прочтали имя brower
			} else if status == 1 && oneb == byte1 {
				status = 2 // прочтали имя brower
			} else if status == 2 && oneb == byte(':') {
			} else if status == 2 && oneb == byte('[') {
			} else if status == 2 && oneb == byte1 {
				status = 3 // читаем строку broweser
				iBBrow = ind + 1
			} else if status == 3 && oneb == byte1 {
				iEBrow = ind - 1
				status = 4 // прочитали browser
			} else if status == 4 && oneb == byte(',') {
				status = 2 // еще раз читаем bro
			} else if status == 4 && oneb == byte(']') {
				status = 5
			} else if status == 5 && oneb == byte(',') {

			} else if status == 5 && oneb == byte1 { //begin company
				status = 6
			} else if status == 6 && oneb == byte1 { //end  company
				status = 7
			} else if status == 7 && oneb == byte(':') {
			} else if status == 7 && oneb == byte1 { //begin company value
				status = 8
			} else if status == 8 && oneb == byte1 { //end company value
				status = 9
			} else if status == 9 && oneb == byte(',') {
			} else if status == 9 && oneb == byte1 { //begin country
				status = 10
			} else if status == 10 && oneb == byte1 { //end country
				status = 11
			} else if status == 10 && oneb == byte(':') {
			} else if status == 11 && oneb == byte1 { //begin conutry value
				status = 12
			} else if status == 12 && oneb == byte1 { //end country value
				status = 13
			} else if status == 13 && oneb == byte(',') {
			} else if status == 13 && oneb == byte1 { //begin email
				status = 14
			} else if status == 14 && oneb == byte1 { //end email
				status = 15
			} else if status == 15 && oneb == byte(':') {
			} else if status == 15 && oneb == byte1 { //begin email value
				status = 17
				iBEm1 = ind + 1
			} else if status == 17 && oneb == byte('@') { //end email
				iEEm1 = ind
				iBEm2 = ind + 1
			} else if status == 17 && oneb == byte1 { //end email value
				iEEm2 = ind
				status = 18
			} else if status == 18 && oneb == byte(',') {
			} else if status == 18 && oneb == byte1 { //begin job
				status = 19
			} else if status == 19 && oneb == byte1 { //end job
				status = 20
			} else if status == 20 && oneb == byte(':') {
			} else if status == 20 && oneb == byte1 { //begin job value
				status = 21
			} else if status == 21 && oneb == byte1 { //end job value
				status = 22
			} else if status == 22 && oneb == byte(',') {
			} else if status == 22 && oneb == byte1 { //begin name
				status = 23
			} else if status == 23 && oneb == byte1 { //end name
				status = 24
			} else if status == 24 && oneb == byte(':') {
			} else if status == 24 && oneb == byte1 { //begin name value
				status = 25
				iBN = ind + 1
			} else if status == 25 && oneb == byte1 { //end name value
				status = 26
				iEN = ind
			} else if status == 26 && oneb == byte(',') {
			} else if status == 26 && oneb == byte1 { //begin phone
				status = 27
			} else if status == 27 && oneb == byte1 { //end phone
				status = 29
			} else if status == 29 && oneb == byte(':') {
			} else if status == 29 && oneb == byte1 { //begin phone value
				status = 30
			} else if status == 30 && oneb == byte1 { //end phone value
				status = 31
			} else if status == 31 && oneb == byte('}') { //end
				status = 32
			}
			if status == 0 {
				isAndroid = false
				isMSIE = false
				isAndroidb = false
				isMSIEb = false

			}

			if status == 4 {
				valstr := string(b[iBBrow:iEBrow])
				isAndroidb = false
				isMSIEb = false
				if strings.Contains((valstr), "Android") {
					isAndroid = true
					isAndroidb = true
				}
				if strings.Contains((valstr), "MSIE") {
					isMSIE = true
					isMSIEb = true
				}
				if isAndroidb || isMSIEb {
					seenBrowsers[valstr] = true
				}
			}
			if status == 18 && oneb == byte1 {
				strEmail = string(b[iBEm1:iEEm1]) + " [at] " + string(b[iBEm2:iEEm2])
			}
			if status == 26 && oneb == byte1 {
				strName = string(b[iBN:iEN])
			}
			if status == 32 {
				if isAndroid && isMSIE {
					fmt.Fprintf(out, "[%d] %s <%s>\n", i, strName, strEmail)
				}
			}
		}
		b = b[:0]

	}
	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}
