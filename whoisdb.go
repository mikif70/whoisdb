// whoisdb
package main

import (
	"fmt"

	"gopkg.in/mgo.v2"
	//	"gopkg.in/mgo.v2/bson"
	"bufio"
	"os"
	"strings"
)

var (
	filename = "ripe.db.inetnum"
	filepath = "/usr/local/Backup"
	mongodb  = "mongodb://127.0.0.127018"
)

type Whois struct {
	Inetnum string `json:"inetnum" bson:"inetnum"`
	Start   string `json:"start" bson:"start"`
	Stop    string `json:"stop" bson:"stop`
	Netname string `json:"netname" bson:"netname"`
	Country string `json:"country" bson:"country"`
	Descr   string `json:"desc" bson:"desc"`
	Mnt     string `json:"mnt" bson:"mnt"`
}

func main() {

	session, err := mgo.Dial(mongodb)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	db := session.DB("whois").C("inetnum")

	file, err := os.Open(filepath + "/" + filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var whois = &Whois{}

	count := 0
	for scanner.Scan() {
		count += 1
		var line = strings.Split(scanner.Text(), ":")

		switch line[0] {
		case "inetnum":
			whois.Inetnum = strings.TrimSpace(line[1])
			var ips = strings.Split(line[1], "-")
			whois.Start = strings.TrimSpace(ips[0])
			whois.Stop = strings.TrimSpace(ips[1])
		case "netname":
			whois.Netname = strings.TrimSpace(line[1])
		case "descr":
			if whois.Descr == "" {
				whois.Descr = strings.TrimSpace(line[1])
			}
		case "country":
			whois.Country = strings.TrimSpace(line[1])
		case "mnt-by":
			whois.Mnt = strings.TrimSpace(line[1])
		default:
			if strings.Contains(line[0], "%") {
				if whois.Inetnum != "" {
					fmt.Println(whois)
					db.Insert(whois)
					whois = &Whois{}
				}
			}
		}
	}
}
