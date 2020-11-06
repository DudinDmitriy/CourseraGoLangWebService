package main

import (
	"fmt"
	"sort"
	"sync"
)

// сюда писать код

// ExecutePipeline execute massiv function
func ExecutePipeline(f ...job) {

	in := make(chan interface{})
	for _, val := range f {
		out := make(chan interface{})
		go val(in, out)
		in = out
	}

}

var mutexMD5 sync.Mutex

func myDataSignerMd5(data string, ch chan string) {
	mutexMD5.Lock()
	ch <- DataSignerMd5(data)
	mutexMD5.Unlock()

}

func myDataSignerCrc32(data string, ch chan string) {
	ch <- DataSignerCrc32(data)
}

func mySingleHash(data string, out chan interface{}) {

	res1 := make(chan string)
	res2 := make(chan string)
	res3 := make(chan string)
	go myDataSignerCrc32(data, res1)
	go myDataSignerMd5(data, res2)
	go myDataSignerCrc32(<-res2, res3)
	out <- (<-res1 + "~" + <-res3)

}

func SingleHash(in, out chan interface{}) {
	for val := range in {
		data := val.(string)
		go mySingleHash(data, out)
	}
	close(out)
}

func myMultiHash(data string, out chan interface{}) {

	chm := make([]chan string, 6)
	for i := 0; i < 6; i++ {
		ch := make(chan string)
		pref := fmt.Sprintf("%02d", i)
		myDataSignerCrc32(pref+data, ch)
		chm[i] = ch
	}
	res := ""
	for _, val := range chm {
		res = res + <-val
	}
	out <- res
}

func MultiHash(in, out chan interface{}) {

	for val := range in {
		data := val.(string)
		go myMultiHash(data, out)
	}
	close(out)

}

type ressort []string

func (rs ressort) Len() int           { return len(rs) }
func (rs ressort) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
func (rs ressort) Less(i, j int) bool { return rs[i] < rs[j] }

func CombineResults(in, out chan interface{}) {

	resm := make([]string, 0)
	for val := range in {
		data := val.(string)
		resm = append(resm, data)
	}
	sort.Sort(ressort(resm))
	res := ""
	for _, val := range resm {
		res = res + val
	}
	out <- res
	close(out)

}

func main() {

}
