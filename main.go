package main

import (
	"fmt"

	//"errors"

	"flag"
	"log"
	"os"
	// "strings"
)

type LoginUser struct {
	Username string `json:"username"`
	Domain   string `json:"domain"`
	Password string `json:"password"`
}

type LoginToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type EnvDetails struct {
	Server    string
	Protocol  string
	ServerURL string
}

type CSDetails struct {
	Name     string `json:"name"`
	CSStatus string `json:"status"`
	Details  struct {
		Version string `json:"version"`
	} `json:"details"`
}

type DesktopPoolDetails struct {
	PoolName        string `json:"name"`
	PoolDisplayName string `json:"display_name"`
	PoolID          string `json:"id"`
}

type DesktopDetails struct {
	DesktopID        string `json:"id"`
	DesktopName      string `json:"name"`
	DesktopDNSName   string `json:"dns_name"`
	DesktopPoolID    string `json:"desktop_pool_id"`
	DesktopState     string `json:"state"`
	DesktopSessionID string `json:"session_id"`
}

type SessionDetails struct {
	SessionID             string `json:"id"`
	SessionUserID         string `json:"user_id"`
	SessionMachineID      string `json:"machine_id"`
	SessionDesktopPoolID  string `json:"desktop_pool_id"`
	SessionType           string `json:"session_type"`
	SessionState          string `json:"session_state"`
	SessionStartTime      int64  `json:"start_time"`
	SessionDisconenctTime int64  `json:"disconnected_time"`
	SessionDuration       int64  `json:"last_session_duration_ms"`
}

var ConnectionUserInfo LoginUser
var ConnectionServerInfo EnvDetails
var ConnectionTokenInfo LoginToken
var CSList []CSDetails
var DesktopPoolList []DesktopPoolDetails
var DesktopDetailsList []DesktopDetails
var SessionDetailsList []SessionDetails

func init() {
	PopulateConfig("config.yml")
	populateStructs()

	if ConnectionServerInfo.Protocol == "" {
		log.Println("Info: No Connection Server Protocol Set. Defaulting to HTTPS.")
		ConnectionServerInfo.Protocol = "https://"
	}

	if ConnectionServerInfo.Server == "" {
		log.Println("Error: No Connection Server URL set. Please add a Connection Server URL to your config.yml file.")
		os.Exit(1)
	}

	ConnectionServerInfo.ServerURL = ConnectionServerInfo.Protocol + ConnectionServerInfo.Server + "/rest"
	fmt.Println("Connection Server URL:>", ConnectionServerInfo.ServerURL)

}

func main() {
	log.SetPrefix("Horizon-Tool: ")
	log.SetFlags(0)

	adminuser := flag.String("adminuser", "", "-adminuser <username> [Username for Server Connection. Used with -adminpwd]")
	adminpwd := flag.String("adminpwd", "", "-adminpwd <username> [Password for Server Connection. Used with -adminuser. Warning: Passwords are entered in plain text and may show up on STOUT or other logs.]")
	admindomain := flag.String("admindomain", "", "-admindomain <Active Directory Domain> [Domain for Server Connection. Used with -adminuser and -adminpwd.]")
	listcs := flag.Bool("listcs", false, "-listcs")
	listdesktoppools := flag.Bool("listdesktoppools", false, "-listdesktoppools")
	logoutstaledesktops := flag.Bool("logoutstaledesktops", false, "-logoutstaledesktops")

	flag.Parse()

	if *adminuser == "" {
		log.Println("Error: No Admin User entered.  Cannot proceed")
		os.Exit(1)
	} else {
		//ConnectionUserInfo.Username = *adminuser
		ConnectionUserInfo.Username = *adminuser
	}

	if *adminpwd == "" {
		log.Println("Error: No Admin Password entered.  Cannot proceed")
		os.Exit(1)
	} else {
		//ConnectionUserInfo.Password = *adminpwd
		ConnectionUserInfo.Password = *adminpwd
	}

	if *admindomain == "" {
		log.Println("Error: No Admin Domain entered.  Cannot proceed")
		os.Exit(1)
	} else {
		//ConnectionUserInfo.Domain = *admindomain
		ConnectionUserInfo.Domain = *admindomain
	}

	HorizonAuthentication()
	//fmt.Println(ConnectionTokenInfo.AccessToken)

	if *listcs {
		log.Println("Listing Connection Servers in Pod")
		cslist()
	}

	if *listdesktoppools {
		log.Println("Displaying All Desktop Pools")
		getalldesktoppools()
	}

	if *logoutstaledesktops {
		log.Println("Logging Out Stale Desktops")
		staledesktops()
	}

	HorizonLogout()

}

func populateStructs() {
	ConnectionUserInfo.Domain = ""
	ConnectionUserInfo.Password = ""
	ConnectionUserInfo.Username = ""
}
