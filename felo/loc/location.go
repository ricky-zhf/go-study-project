package loc

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
)

func GetLoc() {
	geoip()
	//var dbPath = "./ip2region.xdb"
	//searcher, err := xdb.NewWithFileOnly(dbPath)
	//if err != nil {
	//	fmt.Printf("failed to create searcher: %s\n", err.Error())
	//	return
	//}
	//
	//defer searcher.Close()
	//
	//// do the search
	//var ip = "1.2.3.4"
	//var tStart = time.Now()
	//region, err := searcher.SearchByStr(ip)
	//if err != nil {
	//	fmt.Printf("failed to SearchIP(%s): %s\n", ip, err)
	//	return
	//}
	//
	//fmt.Printf("{region: %s, took: %s}\n", region, time.Since(tStart))
}

func geoip() {
	db, err := geoip2.Open("/Users/zhouhuaifeng/GoWorkspace/src/StudyProject/felo/loc/GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// If you are using strings that may be invalid, check that ip is not nil
	//ip := net.ParseIP("81.2.69.142")
	ip := net.ParseIP("32.241.60.192")
	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Portuguese (BR) city name: %v\n", record.City.Names["pt-BR"])
	if len(record.Subdivisions) > 0 {
		fmt.Printf("English subdivision name: %v\n", record.Subdivisions[0].Names["en"])
	}
	fmt.Printf("Russian country name: %v\n", record.Country.Names["ru"])
	fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)
	fmt.Printf("Time zone: %v\n", record.Location.TimeZone)
	fmt.Printf("Coordinates: %v, %v\n", record.Location.Latitude, record.Location.Longitude)
}
