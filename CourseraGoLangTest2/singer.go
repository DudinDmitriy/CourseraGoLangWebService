package main

import (
	"fmt"
	"sort"
	"sync"
)

// сюда писать код

func executeJob(in chan interface{}, out chan interface{}, f job, wg *sync.WaitGroup) {
	f(in, out)
	close(out)
	wg.Done()
}

// ExecutePipeline execute massiv function
func ExecutePipeline(f ...job) {

	var wg sync.WaitGroup
	in := make(chan interface{})

	for _, val := range f {
		out := make(chan interface{})
		wg.Add(1)
		// fmt.Println(val)
		go executeJob(in, out, val, &wg)
		in = out
	}
	wg.Wait()
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

func mySingleHash(data string, wg *sync.WaitGroup, out chan interface{}) {

	defer wg.Done()
	res1 := make(chan string)
	res2 := make(chan string)
	res3 := make(chan string)
	go myDataSignerCrc32(data, res1)
	go myDataSignerMd5(data, res2)
	go myDataSignerCrc32(<-res2, res3)
	out <- (<-res1 + "~" + <-res3)

}

func SingleHash(in, out chan interface{}) {

	var wg sync.WaitGroup

	for val := range in {
		dataint := val.(int)
		data := fmt.Sprintf("%d", dataint)
		wg.Add(1)
		go mySingleHash(data, &wg, out)
	}
	wg.Wait()
}

func myMultiHash(data string, wg *sync.WaitGroup, out chan interface{}) {

	defer wg.Done()
	chm := make([]chan string, 6)
	for i := 0; i < 6; i++ {
		ch := make(chan string)
		pref := fmt.Sprintf("%d", i)
		go myDataSignerCrc32(pref+data, ch)
		chm[i] = ch
	}
	res := ""
	for _, val := range chm {
		res = res + <-val
	}
	fmt.Println("Data MulriHash", res)
	out <- res
}

func MultiHash(in, out chan interface{}) {

	var wg sync.WaitGroup
	for val := range in {
		data := val.(string)
		wg.Add(1)
		go myMultiHash(data, &wg, out)
	}
	wg.Wait()

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
		if res != "" {
			res += "_"
		}
		res += val
	}
	out <- res

}
