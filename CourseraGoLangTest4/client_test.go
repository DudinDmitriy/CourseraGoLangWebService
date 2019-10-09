package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"testing"
)

//Функция заполняет данными массив в зависимости от параметров
//В функцию передают пустой массив, данные дополнятся в массив
func GetData(query string, orderFiled string, orderBy int, result []User) {

	f, err := os.Open("./dataset.xml")
	if err != nil {
		fmt.Errorf("can`t open file: %s", err)
	}

	users := []User{}
	fileContext, err := ioutil.ReadAll(f)

	err = xml.Unmarshal(fileContext, users)
	if query != "" {
		sort.Slice(users, func(i, j int) bool { return (users[i].Name + users[i].About) < (users[j].Name + users[j].About) })
		in := sort.Search(len(users), func(i int) bool { return (users[i].Name + users[i].About) >= query })
	}

}

// код писать тут
func SearchServer(w http.ResponseWriter, r *http.Request) {

	//Получаем параметры
	query := r.FormValue("query")
	order_filed := r.FormValue("order_field")
	order_by := r.FormValue("order_by")

	//Получим данные
	f, err := os.Open("./dataset.xml")
	if err != nil {
		fmt.Errorf("can`t open file: %s", err)
	}

	users := []User{}
	fileContext, err := ioutil.ReadAll(f)
	err = xml.Unmarshal(fileContext, users)

	//Поиск данных

}

func Testreqwest1(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	sc := SearchClient{
		AccessToken: "valodtoken",
		URL:         ts.URL}

}
