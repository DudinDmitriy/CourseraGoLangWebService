package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
)

type fileSort []os.FileInfo

func (fs fileSort) Len() int {
	return len(fs)
}

func (fs fileSort) Swap(i, j int) { fs[i], fs[j] = fs[j], fs[i] }

func (fs fileSort) Less(i, j int) bool { return (fs[i].Name() < fs[j].Name()) }

func dirTree(f io.Writer, dir string, p bool) error {
	return dirTreePref(f, dir, p, "")
}

func nameFile(f os.FileInfo) string {
	if f.IsDir() {
		return f.Name()
	}
	strsize := "empty"
	if f.Size() > 0 {
		strsize = fmt.Sprintf("%db", f.Size())
	}
	return fmt.Sprintf("%s (%s)", f.Name(), strsize)

}

func dirTreePref(f io.Writer, dir string, p bool, pref string) error {

	var fs1 []os.FileInfo

	strprefc := "├───"
	strprefend := "└───"
	prefnext := "│\t"
	prefend := "\t"

	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("Dir can't read: " + dir)
	}

	//уберем файлы из списка
	if !p {
		for _, val := range fs {
			if val.IsDir() {
				fs1 = append(fs1, val)
			}
		}
	} else {
		fs1 = fs
	}

	sort.Sort(fileSort(fs1))

	for id, val := range fs1 {

		strname := pref
		prefn := pref + prefnext
		if (id + 1) == len(fs1) {
			strname = strname + strprefend
			prefn = pref + prefend
		} else {
			strname = strname + strprefc
		}
		strname = strname + nameFile(val) + "\n"

		fmt.Fprintf(f, strname)

		if val.IsDir() {
			dirTreePref(f, dir+"/"+val.Name(), p, prefn)
		}

	}

	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
