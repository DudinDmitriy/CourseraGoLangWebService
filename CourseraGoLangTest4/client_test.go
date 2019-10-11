package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sort"
	"strconv"
	"testing"
)

const (
	orderIDAgeName = "Id_Age_Name"
	nameValidToken = "validtoken"
)

type TestCase struct {
	Query      string
	Orderfield string
	OrderBy    int
}

//Функция заполняет данными массив в зависимости от параметров
//В функцию передают пустой массив, данные дополнятся в массив
func GetData(query string, orderField string, orderBy int, limit int, offset int, result []User) (err error) {

	if orderField != "" && orderField != orderIDAgeName {
		err = fmt.Errorf("ErrorBadOrderField")
		return
	}

	if orderBy != orderDesc && orderBy != OrderByAsIs && orderBy != OrderByAsc {
		err = fmt.Errorf("ErrorBadOrderBy")
		return
	}

	f, err := os.Open("./dataset.xml")
	if err != nil {
		err = fmt.Errorf("can`t open file: %s", err)
		return
	}

	users := make([]User,0)

	fileCon, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	type DataRow struct {
		Id        int    `xml:"id"`
		FirstName string `xml:"first_name"`
		LastName  string `xml:"last_name"`
		Age       int    `xml:"age"`
		About     string `xml:"about"`
		Gender    string `xml:"gender"`
	}
	type Structdataset struct {
		Rows []DataRow `xml:"row"`
	}

	dataset := &Structdataset{}

	xml.Unmarshal(fileCon, dataset)

	for i := 0; i < len(dataset.Rows); i++ {
		d := dataset.Rows[i]
		if query == "" || query == (d.FirstName+d.LastName+d.About) {

			usr := User{}
			usr.Name = d.FirstName + d.LastName
			usr.Id = d.Id
			usr.Age = d.Age
			usr.About = d.About
			usr.Gender = d.Gender

			users = append(users, usr)
		}
	}

	if orderBy == OrderByAsc || orderBy == OrderByDesc {

		type sortfunc func(i, j int) bool
		var sf sortfunc

		sf = func(i, j int) bool {
			var k, l int

			if orderBy == OrderByAsc {
				k = i
				l = j
			} else {
				k = j
				l = l
			}

			if users[k].Id < users[l].Id {
				return true
			}

			if users[k].Id == users[l].Id && users[k].Age < users[l].Age {
				return true
			}

			if users[k].Id == users[l].Id && users[k].Age == users[l].Age && users[k].Name < users[l].Name {
				return true
			}
			return false
		}
		sort.Slice(users, sf)
	}

	for i := limit * (offset - 1); i >= 0 && i < limit*offset && i < len(users); i++ {
		result = append(result, users[i])
	}
	return nil
}

func GetDataSearchRequest(sr SearchRequest, result []User) (err error) {
	return GetData(sr.Query, sr.OrderField, sr.OrderBy, sr.Limit, sr.Offset, result)
}

// код писать тут
func SearchServer(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("AccessToken")
	if token != nameValidToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//Получаем параметры
	query := r.FormValue("query")
	orderFiled := r.FormValue("order_field")
	orderBy, _ := strconv.Atoi(r.FormValue("order_by"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	offset, _ := strconv.Atoi(r.FormValue("offset"))

	//Получим данные
	users := make([]User, 0)
	err := GetData(query, orderFiled, orderBy, limit, offset, users)
	if err != nil {
		if err.Error() == "ErrorBadOrderField" || err.Error() == "ErrorBadOrderBy" {
			w.WriteHeader(http.StatusBadRequest)
			errRes := SearchErrorResponse{Error: err.Error()}
			bytejson, _ := json.Marshal(errRes)
			io.WriteString(w, string(bytejson))
			return
		}
	}

	bytejson, _ := json.Marshal(users)
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(bytejson))

}

func TestSearchServer(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	sc := SearchClient{
		AccessToken: nameValidToken,
		URL:         ts.URL,
	}

	TestCases := []SearchRequest{
		SearchRequest{
			Limit:      5,
			Offset:     3,
			Query:      "",
			OrderField: orderIDAgeName,
			OrderBy:    OrderByDesc,
		},
		SearchRequest{
			Limit:      25,
			Offset:     1,
			Query:      "",
			OrderField: orderIDAgeName,
			OrderBy:    orderAsc,
		},
	}
	
	for _, tt := range TestCases {
		
		res, err := sc.FindUsers(tt)
		userseq := make([]User, 0)
		
		if err != nil {
			t.Errorf("Ошибка: %s", err.Error())
			continue
		}
		
		GetDataSearchRequest(tt, userseq)
		if !reflect.DeepEqual(res.Users, userseq) {
			t.Errorf("Для структуры запросы %+v данные различаются", tt)
		}
	}
}

func TestSearchServerToken(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	sc := SearchClient{
		AccessToken: nameValidToken,
		URL:         ts.URL,
	}

	TestCases := []SearchRequest{
		SearchRequest{
			Limit:      1,
			Offset:     1,
			Query:      "",
			OrderField: "NotLegal",
			OrderBy:    orderAsc,
		},
	}
	
	for _, tt := range TestCases {
		
		res, err := sc.FindUsers(tt)
		
		if res != nil || err.Error()!= "OrderFeld NotLegal invalid" {
			t.Errorf("Ошибка: неккорктно отрабатывает ошибку заполнения поля  <OrderField>")
		}
	}
}
