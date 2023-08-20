package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	//"fmt"

	//"errors"

	//"flag"
	"log"
	"os"
	"time"
	// "strings"
)

func HorizonAuthentication() {
	LoginUrl := ConnectionServerInfo.ServerURL + "/login"
	fmt.Println(LoginUrl)

	postBody, _ := json.Marshal(map[string]string{
		"username": ConnectionUserInfo.Username,
		"domain":   ConnectionUserInfo.Domain,
		"password": ConnectionUserInfo.Password,
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(LoginUrl, "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body
	r := resp.Status
	fmt.Println(r)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if r == "200 OK" {
		populateTokens(string(body))
		log.Println("Authentication Successful")

		//sb := string(body)
		//log.Printf(sb)

		//json.Unmarshal([]byte(body), &ConnectionTokenInfo)

		//log.Printf("Access Token:")
		//fmt.Println(ConnectionTokenInfo.AccessToken)
		//log.Printf("Refresh Token:")
		//fmt.Println(ConnectionTokenInfo.RefreshToken)
	} else {
		log.Println("Unable to authenticate. Please check your credentials, server URL, and service availability and try again.")
		os.Exit(1)
	}

}

func HorizonLogout() {
	LogoutUrl := ConnectionServerInfo.ServerURL + "/logout"

	postBody, _ := json.Marshal(map[string]string{
		"refresh_token": ConnectionTokenInfo.RefreshToken,
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(LogoutUrl, "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body
	r := resp.Status
	fmt.Println(r)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if r == "200 OK" {
		log.Println(ConnectionUserInfo.Username + " Logged Out Successful")
		//log.Print(body)
	} else {
		log.Println("Error during logout.")
		log.Print(body)
		os.Exit(1)
	}

}

func getREST(ServerURL string, filter string) (body []byte) {
	//GetURL := ConnectionServerInfo.ServerURL + ServerURL
	//fmt.Println(GetURL)
	//fmt.Println(filter)

	var queryURL string

	if filter != "" {
		//fmt.Println(url.QueryEscape(filter))
		queryURL = ConnectionServerInfo.ServerURL + ServerURL + "?filter=" + url.QueryEscape(filter)
	} else {
		queryURL = ConnectionServerInfo.ServerURL + ServerURL
	}

	//fmt.Println(queryURL)

	// Create a new request using http
	//client := http.Client{Timeout: 5 * time.Second}
	client := http.Client{}
	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	// add authorization header to the req
	var bearer = "Bearer " + ConnectionTokenInfo.AccessToken
	req.Header.Set("Authorization", bearer)

	resp, err := client.Do(req)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body
	r := resp.Status
	fmt.Println(r)
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if r == "200 OK" {
		//populateTokens(string(body))
		//fmt.Print(body)
		log.Println("REST Call Successful")
	} else {
		log.Println("REST Call Failed")
		HorizonLogout()
		//os.Exit(1)
	}

	return

}

func postREST(ServerURL string, postbody string) (body []byte) {
	//GetURL := ConnectionServerInfo.ServerURL + ServerURL
	//fmt.Println(GetURL)
	//fmt.Println(filter)

	var queryURL string

	//if filter != "" {
	//fmt.Println(url.QueryEscape(filter))
	//	queryURL = ConnectionServerInfo.ServerURL + ServerURL + "?filter=" + url.QueryEscape(filter)
	//} else {
	queryURL = ConnectionServerInfo.ServerURL + ServerURL
	//}

	//fmt.Println(queryURL)

	jsonBody := []byte(postbody)
	bodyReader := bytes.NewReader(jsonBody)

	// Create a new request using http
	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("POST", queryURL, bodyReader)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	// add authorization header to the req
	var bearer = "Bearer " + ConnectionTokenInfo.AccessToken
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body
	r := resp.Status
	fmt.Println(r)
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if r == "200 OK" {
		//populateTokens(string(body))
		//fmt.Print(body)
		log.Println("REST Call Successful")
	} else {
		log.Println("REST Call Failed")
		HorizonLogout()
		//os.Exit(1)
	}

	return

}

func populateTokens(jsonBody string) {
	json.Unmarshal([]byte(jsonBody), &ConnectionTokenInfo)
}
