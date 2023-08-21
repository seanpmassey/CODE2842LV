package main

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"time"
	//"io"
	//"net/http"
	//"net/url"
	//"fmt"
	//"errors"
	//"flag"
	//"log"
	//"os"
	// "strings"
)

func staledesktops() {
	getdesktoppools("VIRTUAL_CENTER", "AUTOMATED", "DEDICATED")

	evalTime := time.Now()
	disconnectThresholdTime := evalTime.Add(time.Duration(-2) * time.Hour)

	for _, val := range DesktopPoolList {
		//fmt.Println("Desktop Pools:")
		//fmt.Println(val.PoolDisplayName + "		" + val.PoolName)
		fmt.Println("Checking For Stale Sessions for " + val.PoolDisplayName + "	" + val.PoolID)

		getdisconnectedmachines(val.PoolID)
		fmt.Println("Stale Sessions for Pool " + val.PoolDisplayName + ":")
		for _, val := range DesktopDetailsList {
			//fmt.Println("Desktop Pools:")
			fmt.Println(val.DesktopName + "     " + val.DesktopState)
		}
		getdisconnectedsessions(val.PoolID)
		for _, val := range SessionDetailsList {
			//fmt.Println("Desktop Pools:")
			//var sessionDisconnectTimeString string
			//timestamp, err := strconv.ParseInt(sessionDisconnectTimeString, 10, 64)
			//if err != nil {
			//	panic(err)
			//}

			sessionDisconnectTime := time.UnixMicro(val.SessionDisconenctTime)
			if disconnectThresholdTime.Before(sessionDisconnectTime) {
				fmt.Println("Session has not passed threshold")
			} else {
				fmt.Println("Logging Out Session")
				postLogoutSession(val.SessionID)
			}
			//fmt.Println(val.SessionID + "     " + sessionDisconnectTime.Local().GoString())
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

func getdisconnectedmachines(desktoppoolid string) {
	var url string = "/inventory/v5/machines"
	var filter string = `{"type":"And","filters": [{"type":"Equals","name":"state","value":"DISCONNECTED"},{"type":"Equals","name":"desktop_pool_id","value":"` + desktoppoolid + `"}]}`

	body := getREST(url, filter)
	//fmt.Println(body)
	//Uncomment to debug output
	//fmt.Println(json.Unmarshal([]byte(body), &DesktopDetailsList))
	if err := json.Unmarshal([]byte(body), &DesktopDetailsList); err != nil {
		panic(err)
	}
}

func getdisconnectedsessions(desktoppoolid string) {
	var url string = "/inventory/v1/sessions"
	var filter string = `{"type":"And","filters": [{"type":"Equals","name":"session_state","value":"DISCONNECTED"},{"type":"Equals","name":"desktop_pool_id","value":"` + desktoppoolid + `"}]}`

	body := getREST(url, filter)
	if err := json.Unmarshal([]byte(body), &SessionDetailsList); err != nil {
		panic(err)
	}
}

func postLogoutSession(sessionID string) {
	var url string = "/inventory/v1/sessions/action/logoff"
	var postbody string

	postbody = `["` + sessionID + `"]`
	fmt.Println(postbody)

	postREST(url, postbody)
	//if err := json.Unmarshal([]byte(body), &SessionDetailsList); err != nil {
	//	panic(err)
	//}
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
