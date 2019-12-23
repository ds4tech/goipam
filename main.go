package ipam

import (
  "log"
  "os"
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
  reserve     bool
  release     bool
  client * gosolar.Client
  )

func main() {
  flag.Parse()

  if(len(ipaddress) > 0) {
      log.Printf("[DEBUG] Checks DISABLED")
      if( checkDhcp(ipaddress) ) {
        log.Printf("[DEBUG] Ipaddress is assigned by DHCP")
      } else {
        log.Printf("[DEBUG] Ipaddress is NOT assigned by DHCP")
      }
  /*
      if( checkConnectivity(ipaddress) ) {
        log.Printf("[DEBUG] CheckConnectivity finished with SUCCESS")
      } else {
        log.Printf("[DEBUG] Checks NOT PASSED")
        //log.Fatalf("[Error] Checks NOT PASSED")
      }
  */
  } else {
    log.Fatalf("[Error] Ipaddress not provided. Please specify it using ip param. Example: -ip=127.0.0.1 ")
  }


  if(len(username)==0 || len(password)==0 || len(host)==0) {
      log.Fatalf("Failed to conntect to orion. Connection details not provided.")
  }
  log.Printf("[TRACE] Connecting to Orion with credentials: username: %s, password: %s*****, orion_host: %s", username, password[0:len(password)-5], host)
  log.Printf("[DEBUG] reserve: %v, release %v", reserve, release)

  if(len(vlan)>0) {
    getAllIpAddresses(client, vlan)
  } else if(reserve && len(ipaddress)>0) {
    reserveIpAddress(client, ipaddress, comment)
  } else if(release && len(ipaddress)>0) {
    releaseIpAddress(client, ipaddress)
  }else {}

}

func init() {
  flag.IntVar(&status, "status", 4, "4 - in transition")
  flag.StringVar(&ipaddress, "ip", "", "ipaddress of which status should be changed")
  flag.StringVar(&vlan, "vlan", "", "VLAN name, eg: VLAN100_10.141.16.0m24")
  flag.StringVar(&comment, "comment", "", "any comment that will be added to ipaddress record in IPAM")
  flag.BoolVar(&reserve, "reserve", false, "a bool")
  flag.BoolVar(&release, "release", false, "a bool")

  username = os.Getenv("ORION_USER")
  password = os.Getenv("ORION_PASSWORD")
  host     = os.Getenv("ORION_IP")
  client = gosolar.NewClient(host, username, password, true)
}
