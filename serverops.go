package main

import (
	//"bytes"
	"encoding/json"
	"fmt"
	//"io"
	//"net/http"
	//"errors"
	//"flag"
	//"log"
	//"os"
	// "strings"
)

func cslist() {
	var url string = "/monitor/connection-servers"
	var filter string = ""

	body := getREST(url, filter)
	if err := json.Unmarshal([]byte(body), &CSList); err != nil {
		panic(err)
	}
	for _, val := range CSList {
		fmt.Println("Server: " + val.Name + "     " + val.CSStatus + "     " + val.Details.Version)
	}
	//fmt.Println(body)

}
