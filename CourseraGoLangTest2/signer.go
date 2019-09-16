package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
	//"time"
)

// сюда писать код

var Md5 sync.Mutex

func main() {

}

func ExecutePipeline(argJobs ...job) {

	var wg sync.WaitGroup
	in := make(chan interface{})
	for jj := range argJobs {
		out := make(chan interface{})
		wg.Add(1)
		go runjob(in, out, argJobs[jj], &wg)
		in = out
	}
	wg.Wait()
}

func runjob(in, out chan interface{}, f job, wg *sync.WaitGroup) {
	defer wg.Done()
	f(in, out)
	close(out)
}

func SingleHash(in, out chan interface{}) {
	//	defer close(out)
	var wg sync.WaitGroup
	for val := range in {
		data := strconv.Itoa(val.(int))
		wg.Add(1)
		go func(outres chan interface{}, datares string,wg *sync.WaitGroup) {
			defer wg.Done()
			var wg1 sync.WaitGroup
			var s1, s2 string
			f1 := func() {
				s1 = DataSignerCrc32(datares)
				wg1.Done()
			}
			f2 := func() {
				var wg2 sync.WaitGroup
				var s3 string
				defer wg1.Done()

				f3 := func() {
					defer wg2.Done()
					s3 = DataSignerMd5Mutex(datares)
				}
				wg2.Add(1)
				go f3()
				wg2.Wait()

				f4 := func() {
					defer wg2.Done()
					s2 = DataSignerCrc32(s3)
				}
				wg2.Add(1)
				go f4()
				wg2.Wait()
			}
			wg1.Add(1)
			go f1()
			wg1.Add(1)
			go f2()
			wg1.Wait()
			res := s1 + "~" + s2
			outres <- res
		}(out, data, &wg)
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {

    var wg sync.WaitGroup
	for i := range in {

		wg.Add(1)
		go func(outres chan interface{},dataIn string,wg *sync.WaitGroup) {
			defer wg.Done()
			ss := make([]string, 6)

			var wg1 sync.WaitGroup
			for j := 0; j < 6; j++ {
				wg1.Add(1)
				go func(ind int, data string, resss []string) {
					defer wg1.Done()
					resss[ind] = DataSignerCrc32(data)
				}(j, (strconv.Itoa(j) + dataIn), ss)
				
			}
			wg1.Wait()
			out <- strings.Join(ss, "")
		}(out,i.(string),&wg)
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {

	//defer close(out)
	res := make([]string, 0)
	for val := range in {
		res = append(res, val.(string))
	}
	sort.Strings(res)
	out <- strings.Join(res, "_")
	//println(res)
}

func DataSignerMd5Mutex(data string) (res string) {
	Md5.Lock()
	res = DataSignerMd5(data)
	Md5.Unlock()
	return
}
