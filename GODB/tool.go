package GODB

import (
	"fmt"
	"io"
	"os"
)

func usage() {
	io.WriteString(os.Stderr, "Usage:")
	io.WriteString(os.Stderr, "\n tool DBNAME get KEY")
	io.WriteString(os.Stderr, "\n tool DBNAME set KEY VALUE")
	io.WriteString(os.Stderr, "\n tool DBNAME delete KEY\n")
}

func Run() {

	Nargs := len(os.Args)
	if Nargs < 4 || Nargs > 5 {
		usage()
		return
	}

	var dbname, verb, key, val string
	dbname = os.Args[1]
	verb = os.Args[2]
	key = os.Args[3]
	if Nargs > 4 {
		val = os.Args[4]
	}

	set := make(map[string]bool)
	set["get"] = true
	set["set"] = true
	set["delete"] = true

	// check verb valid
	if _, ok := set[verb]; ok != true {
		usage()
		return
	}

	fmt.Println(dbname, key, val)
	// 省略了connect , 封装进了实例化， 本质这个connect就是个DB实例化

	db := NewGODB(dbname)
	if db == nil {
		io.WriteString(os.Stderr, "DB connect failed!")
		return
	}

	switch verb {
	case "get":
		retV, ok := db.get(key)
		if ok {
			fmt.Println(retV)
		} else {
			io.WriteString(os.Stderr, "key not found!")
		}
	case "set":
		ok := db.set(key, val)
		if ok {
			if cok := db.commit(); cok == false {
				io.WriteString(os.Stderr, "commit error!")
			}
		} else {
			io.WriteString(os.Stderr, "set error!")
		}

	case "delete":
		ok := db.del(key)
		if ok {
			if cok := db.commit(); cok == false {
				io.WriteString(os.Stderr, "commit error!")
			}
		} else {
			io.WriteString(os.Stderr, "delet error!")
		}
	default:
		io.WriteString(os.Stderr, "verb not found!")
	}

	return
}
