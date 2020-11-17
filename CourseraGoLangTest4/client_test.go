package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

type UserXML struct {
	ID        int    `xml:"id"`
	Firstname string `xml:"first_name"`
	Lasttname string `xml:"last_name"`
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
}
type ListUser struct {
	RowUser []UserXML `xml:"row"`
}

type sortOderUser struct {
	sortUser []User
	nameOder string
	orderBy  int
}

func (s sortOderUser) Len() int { return len(s.sortUser) }

func (s sortOderUser) Less(i, j int) bool {

	chek := func() bool {
		if s.nameOder == "Id" {
			return s.sortUser[i].Id < s.sortUser[j].Id
		}

		if s.nameOder == "Name" {
			return s.sortUser[i].Name < s.sortUser[j].Name
		}

		if s.nameOder == "Age" {
			return s.sortUser[i].Age < s.sortUser[j].Age
		}
		return false
	}
	if s.orderBy == -1 {
		return !chek()
	}
	return chek()
}

func (s sortOderUser) Swap(i, j int) { s.sortUser[i], s.sortUser[j] = s.sortUser[j], s.sortUser[i] }

func SearchServer(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("AccessToken")
	if token != "GoodToken" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Add("Content-Type", "text/json")

	query := r.URL.Query().Get("query")
	orderBy, err := strconv.Atoi(r.URL.Query().Get("order_by"))
	if err != nil {
		serr := SearchErrorResponse{"ErrorBadOrderBy"}
		b, _ := json.Marshal(serr)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(b)
		return
	}

	ordersort := r.URL.Query().Get("order_field")
	if ordersort != "" && ordersort != "Name" && ordersort != "Age" && ordersort != "Id" {

		serr := SearchErrorResponse{"ErrorBadOrderField"}
		b, _ := json.Marshal(serr)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(b)
		return
	}

	f, err := ioutil.ReadFile("./dataset.xml")
	if err != nil {
		fmt.Println("File can't open")
	}
	var lu ListUser

	xml.Unmarshal(f, &lu)

	resUsers := make([]User, 0, len(lu.RowUser))
	for _, val := range lu.RowUser {
		if query == "" || (strings.Contains(val.Firstname+val.Lasttname, query) || strings.Contains(val.About, query)) {
			resUsers = append(resUsers, User{val.ID, val.Firstname + val.Lasttname, val.Age, val.About, val.Gender})
		}
	}

	sr := sortOderUser{resUsers, ordersort, orderBy}
	sort.Sort(sr)

	bjson, err := json.Marshal(sr.sortUser)
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(bjson))

}

type testCase struct {
	sc  SearchClient
	sr  SearchRequest
	res string // http body json
}

func Test_Server1(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	sc := SearchClient{"GoodToken", ts.URL}
	sr := SearchRequest{25, 1, "", "Name", 1}

	_, err := sc.FindUsers(sr)

	if err != nil {
		t.Errorf("Don,t work seach")
		return
	}

}

func Test_Server2(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	sc := SearchClient{"GoodToken", ts.URL}
	sr := SearchRequest{-1, 1, "", "Name", 1}

	_, err := sc.FindUsers(sr)

	if err == nil {
		t.Errorf("Don,t work seach")
		return
	}

}

func Test_Server3(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	sc := SearchClient{"GoodToken", ts.URL}
	sr := SearchRequest{25, -1, "", "Name", 1}

	_, err := sc.FindUsers(sr)

	if err == nil {
		t.Errorf("Don,t work seach")
		return
	}

}

func Test_Server4(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	sc := SearchClient{"GoodToken", ts.URL}
	sr := SearchRequest{30, 1, "", "Name", 1}

	_, err := sc.FindUsers(sr)

	if err != nil {
		t.Errorf("Don,t work seach")
		return
	}

}

func Test_Server5(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	sc := SearchClient{"BadToken", ts.URL}
	sr := SearchRequest{30, 1, "", "Name", 1}

	_, err := sc.FindUsers(sr)
	errfunc := true
	if err != nil {
		if err.Error() == "Bad AccessToken" {
			errfunc = false
		}
	}

	if errfunc {
		t.Errorf("Test_Server5 don't work!!!")
		return
	}

}

func Test_Server6(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	sc := SearchClient{"GoodToken", ts.URL}
	sr := SearchRequest{30, 1, "", "BadName", 1}

	_, err := sc.FindUsers(sr)
	errfunc := true
	if err != nil {
		if err.Error() == "OrderFeld BadName invalid" {
			errfunc = false
		}
	}

	if errfunc {
		t.Errorf("Test_Server6 don't work!!!")
		return
	}

}

func Test_Server7(t *testing.T) {

	ftest := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 2)
		w.WriteHeader(http.StatusOK)

	}

	ts := httptest.NewServer(http.HandlerFunc(ftest))

	sc := SearchClient{"GoodToken", ts.URL}
	sr := SearchRequest{30, 1, "", "Name", 1}

	_, err := sc.FindUsers(sr)
	errfunc := true
	if err != nil {
		if strings.Contains(err.Error(), "timeout for") {
			errfunc = false
		}
	}

	if errfunc {
		t.Errorf("Test_Server7 don't work!!!")
		return
	}

}

func Test_Server8(t *testing.T) {

	ftest := func(w http.ResponseWriter, r *http.Request) {
		serr := SearchErrorResponse{"ErrorBadOrderField"}
		b, _ := json.Marshal(serr)
		w.Write(b)

	}

	ts := httptest.NewServer(http.HandlerFunc(ftest))

	sc := SearchClient{"GoodToken", ts.URL}
	sr := SearchRequest{30, 1, "", "Name", 1}

	_, err := sc.FindUsers(sr)
	errfunc := true
	if err != nil {
		if strings.Contains(err.Error(), "ant unpack result json") {
			errfunc = false
		}
	}

	if errfunc {
		t.Errorf("Test_Server8 don't work!!!")
		return
	}

}

func Test_Server9(t *testing.T) {

	sc := SearchClient{"GoodToken", "BadURL"}
	sr := SearchRequest{30, 1, "", "Name", 1}

	_, err := sc.FindUsers(sr)
	errfunc := true
	if err != nil {
		if strings.Contains(err.Error(), "unknown error") {
			errfunc = false
		}
	}

	if errfunc {
		t.Errorf("Test_Server9 don't work!!!")
		return
	}

}

func Test_Server10(t *testing.T) {

	ftest := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)

	}

	ts := httptest.NewServer(http.HandlerFunc(ftest))

	sc := SearchClient{"GoodToken", ts.URL}
	sr := SearchRequest{30, 1, "", "Name", 1}

	_, err := sc.FindUsers(sr)
	errfunc := true
	if err != nil {
		if strings.Contains(err.Error(), "cant unpack error json") {
			errfunc = false
		}
	}

	if errfunc {
		t.Errorf("Test_Server10 don't work!!!")
		return
	}

}

func Test_Server11(t *testing.T) {

	ftest := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		serr := SearchClient{"GoodToken", "BadURL"}
		b, _ := json.Marshal(serr)
		w.Write(b)

	}

	ts := httptest.NewServer(http.HandlerFunc(ftest))

	sc := SearchClient{"GoodToken", ts.URL}
	sr := SearchRequest{30, 1, "", "Name", 1}

	_, err := sc.FindUsers(sr)
	errfunc := true
	if err != nil {
		if strings.Contains(err.Error(), "unknown bad request error:") {
			errfunc = false
		}
	}

	if errfunc {
		t.Errorf("Test_Server11 don't work!!!")
		return
	}

}

func Test_Server12(t *testing.T) {

	ftest := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)

	}

	ts := httptest.NewServer(http.HandlerFunc(ftest))

	sc := SearchClient{"GoodToken", ts.URL}
	sr := SearchRequest{30, 1, "", "Name", 1}

	_, err := sc.FindUsers(sr)
	errfunc := true
	if err != nil {
		if strings.Contains(err.Error(), "SearchServer fatal error") {
			errfunc = false
		}
	}

	if errfunc {
		t.Errorf("Test_Server12 don't work!!!")
		return
	}

}

func Test_Server13(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	sc := SearchClient{"GoodToken", ts.URL}
	sr := SearchRequest{0, 1, "Whitney", "Name", 1}

	res, err := sc.FindUsers(sr)
	fmt.Println(len(res.Users))
	if err != nil {
		t.Errorf("Test_Server13 don,t work seach")
		return
	}

}
