package main

import (
	"fmt"
	// "github.com/limetext/lime/backend/util"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"sort"
	"strings"
)

var (
	re  = regexp.MustCompile("(build/target/)?(lib/[^\n]+)")
	re2 = regexp.MustCompile(`(0[0-9a-f]+)?\s+(\w\s+[^\n]+)\n`)
)

type (
	DB map[string][]string
)

func (d *DB) Add(so, symbols string) {
	log.Println(so)
	syms := re2.FindAllStringSubmatch(symbols, -1)
	s2 := strings.Split(strings.TrimSpace(symbols), "\n")
	if len(s2) != len(syms) {
		for i, s := range syms {
			log.Println(s[2], s2[i])
		}
		log.Fatalln(len(syms), len(s2), s2[len(s2)-1])
	}
	ss := make([]string, len(syms))
	for i, s := range syms {
		ss[i] = s[2]
	}
	sort.Strings(ss)
	(*d)[so] = ss
}

func (d *DB) ParseFile(fn string) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Fatalln(err)
	}
	dat := strings.TrimSpace(string(data))
	matches := re.FindAllStringSubmatchIndex(dat, -1)

	for i := 1; i < len(matches); i++ {
		m0 := matches[i-1]
		m1 := matches[i]
		so := dat[m0[4]:m0[5]]
		a, b := m0[5]+2, m1[0]
		if a >= b {
			continue
		}
		entries := dat[a:b]
		d.Add(so, entries)
	}
}
func main() {
	db := make(DB)
	db.ParseFile("out.txt")
	db2 := make(DB)
	db2.ParseFile("out2.txt")

	for k, v := range db {
		a := strings.Join(v, "\n")
		da, ok := db2[k]
		if !ok {
			fmt.Println("Only in orig:", k)
			continue
		}
		b := strings.Join(da, "\n")
		ioutil.WriteFile(k+".1", []byte(a), 0644)
		ioutil.WriteFile(k+".2", []byte(b), 0644)
		o, _ := exec.Command("diff", "-y", k+".1", k+".2").Output()
		spa := strings.TrimSpace(string(o))
		if spa == "" {
			continue
		}
		fmt.Println(k)
		fmt.Println(spa)
	}

}
