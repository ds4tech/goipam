package main

import (
  "log"
  "os"
  "fmt"
  "flag"
	"github.com/mrxinu/gosolar"
)

var (
  username    string
  password    string
  host        string
  vlan        string
  ipaddress   string
  comment     string
  status      int
  list        bool
  reserve     bool
  release     bool
  client * gosolar.Client
  )

func main() {
  flag.Parse()

  if(len(ipaddress) > 0) {
      log.Printf("[DEBUG] Checks DISABLED")
      if( checkDhcp(ipaddress) ) {
        fmt.Println("Ipaddress is assigned by DHCP")
      } else {
        fmt.Println("Ipaddress is NOT assigned by DHCP")
      }
  /*
      if( checkConnectivity(ipaddress) ) {
        log.Printf("[DEBUG] CheckConnectivity finished with SUCCESS")
      } else {
        log.Printf("[DEBUG] Checks NOT PASSED")
        //log.Fatalf("[Error] Checks NOT PASSED")
      }
  */
  } else if(len(vlan) > 0) {
  }else {
      log.Fatalf("[Error] Nor VLAN or Ipaddress is provided. Please specify it using params. \nExample: \n\tgoipam -ip=127.0.0.1 \n\tgoipam -vlan=VLAN_141810.14.18.0m24 -list")
  }


  if(len(username)==0 || len(password)==0 || len(host)==0) {
      log.Fatalf("Failed to conntect to orion. Connection details not provided.")
  }
  fmt.Println("Connecting to Orion with credentials: username: %s, password: %s*****, orion_host: %s", username, password[0:len(password)-5], host)
  log.Printf("[DEBUG] list: %v, reserve: %v, release %v", list, reserve, release)

  if( list && len(vlan)>0 ) { //list all ipaddress from provided vlan
      getAllIpAddresses(client, vlan)
  } else if( reserve && len(vlan)>0 ) { //reserve first available ipaddress from provided vlan
      reserveIpFromVlan(client, vlan, comment)
  } else if( reserve && len(ipaddress)>0 ) {  //reserve provided ipaddress if available
      reserveIpAddress(client, ipaddress, comment)
  } else if( release && len(ipaddress)>0 ) {  //release provided ipaddress
      releaseIpAddress(client, ipaddress)
  } else {
      log.Printf("Please, provide one of parameters: list, reserve or release")
  }

}

func init() {
  flag.IntVar(&status, "status", 4, "4 - in transition")
  flag.StringVar(&ipaddress, "ip", "", "ipaddress of which status should be changed")
  flag.StringVar(&vlan, "vlan", "", "VLAN name, eg: VLAN100_10.141.16.0m24")
  flag.StringVar(&comment, "comment", "", "any comment that will be added to ipaddress record in IPAM")
  flag.BoolVar(&list , "list", false, "a bool")
  flag.BoolVar(&reserve, "reserve", false, "a bool")
  flag.BoolVar(&release, "release", false, "a bool")

  username = os.Getenv("ORION_USER")
  password = os.Getenv("ORION_PASSWORD")
  host     = os.Getenv("ORION_IP")
  client = gosolar.NewClient(host, username, password, true)
}
