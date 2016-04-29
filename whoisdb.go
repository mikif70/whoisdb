// whoisdb
package main

import (
	"fmt"

	"gopkg.in/mgo.v2"
	//	"gopkg.in/mgo.v2/bson"
	"bufio"
	"flag"
	"math"
	"os"
	"strconv"
	"strings"
)

var (
	filename = "./ripe.db.inetnum"
	mongodb  = "mongodb://127.0.0.127018"
)

type Whois struct {
	Inetnum string `json:"inetnum" bson:"inetnum"`
	Start   int    `json:"start" bson:"start"`
	Stop    int    `json:"stop" bson:"stop`
	Netname string `json:"netname" bson:"netname"`
	Country string `json:"country" bson:"country"`
	Descr   string `json:"desc" bson:"desc"`
	Mnt     string `json:"mnt" bson:"mnt"`
}

func init() {
	flag.StringVar(&mongodb, "m", mongodb, "Mongodb")
	flag.StringVar(&filename, "f", filename, "ripe inetnum filename")
}

func ip2dec(ip string) int {
	ips := strings.Split(ip, ".")
	dec := 0
	//	dec = ips[0]*(256^3) + ips[1]*(256^2) + ips[2]*(256)+ ips[3]
	pot := 3.0
	fmt.Printf("IP: %s => \n", ip)
	for i := range ips {
		num, err := strconv.Atoi(ips[i])
		if err != nil {
			fmt.Printf("Number error: %v\n", err)
			return 0
		}
		res := num * int(math.Pow(256, pot))
		fmt.Printf("\tnum(%d): %d * ( 256 ^ %f) == %d\n", i, num, pot, res)
		dec += res
		pot -= 1
	}

	return dec
}

func main() {

	flag.Parse()

	session, err := mgo.Dial(mongodb)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	db := session.DB("whois").C("inetnum")

	file, err := os.Open(filename)
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
			whois.Start = ip2dec(strings.TrimSpace(ips[0]))
			whois.Stop = ip2dec(strings.TrimSpace(ips[1]))
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
