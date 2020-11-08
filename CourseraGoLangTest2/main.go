package main

func main() {
	inputData := []int{0, 1}
	// inputData := []int{0,1}

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
	}

	ExecutePipeline(hashSignJobs...)
}
