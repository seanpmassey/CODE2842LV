package main

import (
	//"bytes"
	"encoding/json"
	"fmt"
	//"io"
	//"net/http"
	//"net/url"
	//"fmt"
	//"errors"
	//"flag"
	//"log"
	//"os"
	// "strings"
	// "time"
)

func staledesktops() {
	getdesktoppools("VIRTUAL_CENTER", "AUTOMATED", "DEDICATED")

	for _, val := range DesktopPoolList {
		//fmt.Println("Desktop Pools:")
		//fmt.Println(val.PoolDisplayName + "		" + val.PoolName)
		fmt.Println("Checking For Stale Sessions for " + val.PoolDisplayName + "	" + val.PoolID)

		var url string = "/inventory/v5/machines"
		var filter string = `{"type":"And","filters": [{"type":"Equals","name":"state","value":"DISCONNECTED"},{"type":"Equals","name":"desktop_pool_id","value":"` + val.PoolID + `"}]}`

		body := getREST(url, filter)
		//fmt.Println(body)
		//Uncomment to debug output
		//fmt.Println(json.Unmarshal([]byte(body), &DesktopDetailsList))
		if err := json.Unmarshal([]byte(body), &DesktopDetailsList); err != nil {
			panic(err)
		}
		fmt.Println("Stale Sessions for Pool " + val.PoolDisplayName + ":")
		for _, val := range DesktopDetailsList {
			//fmt.Println("Desktop Pools:")
			fmt.Println(val.DesktopName + "     " + val.DesktopState)
		}
	}
}

func getalldesktoppools() {
	var url string = "/inventory/v7/desktop-pools"
	var filter string = ""

	body := getREST(url, filter)
	//fmt.Println(body)
	if err := json.Unmarshal([]byte(body), &DesktopPoolList); err != nil {
		panic(err)
	}
	fmt.Println("Desktop Pools:")
	for _, val := range DesktopPoolList {
		//fmt.Println("Desktop Pools:")
		fmt.Println(val.PoolDisplayName + "		" + val.PoolName)
	}
}

func getdesktoppools(source string, desktoppooltype string, userassignment string) {
	var url string = "/inventory/v7/desktop-pools"
	var filter string = `{"type":"And","filters": [{"type":"Equals","name":"source","value":"` + source + `"},{"type":"Equals","name":"type","value":"` + desktoppooltype + `"},{"type":"Equals","name":"user_assignment","value":"` + userassignment + `"}]}`

	body := getREST(url, filter)
	//fmt.Println(body)
	if err := json.Unmarshal([]byte(body), &DesktopPoolList); err != nil {
		panic(err)
	}
	fmt.Println("Desktop Pools:")
	for _, val := range DesktopPoolList {
		//fmt.Println("Desktop Pools:")
		fmt.Println(val.PoolDisplayName + "     " + val.PoolName)
	}
}

//func logoutdesktopsession(){
//var url string = "/monitor/connection-servers"

//body := getREST(url,filter)
//fmt.Println(body)
//}
