package main
import (
	"testing"
)

func TestSingleHash(t *testing.T){
	testres:= "4108050209~502633748"
	in := make(chan interface{})
	out := make(chan interface{})
	go SingleHash(in,out)
	in <- 0
	res:= <- out
	if res!=testres{
		t.Errorf("Result: %s mismatch \n taget %s ",res,testres)
	}

}

func TestMultiHash(t *testing.T){
	testres:= "29568666068035183841425683795340791879727309630931025356555"
	in := make(chan interface{})
	out := make(chan interface{})
	go MultiHash(in,out)
	in <- "4108050209~502633748"
	res:= <- out
	if res!=testres{
		t.Errorf("Result: %s mismatch \n taget %s ",res,testres)
	}

}
