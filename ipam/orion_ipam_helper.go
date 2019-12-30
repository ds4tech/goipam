package ipam

import (
    "encoding/json"
    "regexp"
    "log"
    "fmt"

    "github.com/mrxinu/gosolar"
)

// Node struct holds query results
type Node struct {
    IPNODEID int `json:"ipnodeid"` //not used
    SUBNETID int `json:"subnetid"`
    STATUS int `json:"status"`
    IPADDRESS string `json:"ipaddress"`
    COMMENTS string `json:"comments"`
    URI string `json:"uri"`
    DISPLAYNAME string `json:DisplayName"`
}

func getAllIpAddresses(client * gosolar.Client, vlan string)  {
    var response []*Node
    var status string

    subnetIpAddr := findIP(vlan)
    cidr := string(vlan[len(vlan)-2:]) //substring last 2 chars

    log.Printf("[DEBUG] vsphere_vlan: %s, subnetIpAddr: %s, cidr: %s", vlan, subnetIpAddr, cidr)

    querySubnet := fmt.Sprintf("SELECT SubnetId FROM IPAM.Subnet WHERE  Address='%s' AND CIDR=%s", subnetIpAddr, cidr)
    response = queryOrionServer(client, querySubnet)
    log.Printf("[DEBUG] SUBNETID: %v\n", response[0].SUBNETID)

    queryIpnode := fmt.Sprintf("SELECT IpNodeId,IPAddress,Comments,Status,Uri,DisplayName FROM IPAM.IPNode WHERE SubnetId='%d'", response[0].SUBNETID)
    response = queryOrionServer(client, queryIpnode)

    fmt.Println("DisplayName | IPAddress | Comments | Status")
    for _, el := range response {
    comment := el.COMMENTS
    if( len(comment) == 0 ){
  	   comment = "---"
    }
    if( 1 == el.STATUS ){
  	   status = "Used"
    } else if( 2 == el.STATUS ){
  	   status = "Available"
    } else if( 4 == el.STATUS ){
  	   status = "Reserved"
    } else if( 8 == el.STATUS ){
  	   status = "Transient"
    } else {
  	   status = string(el.STATUS)
    }
      fmt.Printf("%s | %s | %s | %s\n", el.DISPLAYNAME, el.IPADDRESS, comment, status)
    }
}

func reserveIpFromVlan(client * gosolar.Client, vlan string, comment string)  {
    var response []*Node
    subnetIpAddr := findIP(vlan)
    cidr := string(vlan[len(vlan)-2:]) //substring last 2 chars

    log.Printf("[DEBUG] vsphere_vlan: %s, subnetIpAddr: %s, cidr: %s", vlan, subnetIpAddr, cidr)

    querySubnet := fmt.Sprintf("SELECT SubnetId FROM IPAM.Subnet WHERE  Address='%s' AND CIDR=%s", subnetIpAddr, cidr)
    response = queryOrionServer(client, querySubnet)

    log.Printf("[DEBUG] SUBNETID: %v\n", response[0].SUBNETID)

    queryIpnode := fmt.Sprintf("SELECT TOP 1 IpNodeId,IPAddress,Comments,Status,Uri FROM IPAM.IPNode WHERE SubnetId='%d' and status=2 AND IPOrdinal BETWEEN 11 AND 254", response[0].SUBNETID)
    response = queryOrionServer(client, queryIpnode)

    updateIPNodeStatus(client, response[0].URI, "1", comment) // '1' == ip used
    fmt.Println("Reserved ip: ", response[0].IPADDRESS)
}

func reserveIpAddress(client * gosolar.Client, ipaddress string, comment string)  {
    var response []*Node

    querySubnet := fmt.Sprintf("SELECT IpNodeId,IPAddress,Comments,Status,Uri FROM IPAM.IPNode WHERE IPAddress='%s'", ipaddress)
    response = queryOrionServer(client, querySubnet)
    fmt.Printf("Reserve: IPADDRESS: %v, COMMENTS: %v, URI: %v\n", ipaddress, comment, response[0].URI)

    if(2 == response[0].STATUS) {
      updateIPNodeStatus(client, response[0].URI, "1", comment) // '1' == ip used
    } else {
      log.Fatalf("[Error] Ip address is not available. It has status %v.", response[0].STATUS)
    }
    fmt.Println("Reserved")
}

func changeIpAddressStatus(client * gosolar.Client, ipaddress string)  {}

func releaseIpAddress(client * gosolar.Client, ipaddress string)  {
    var response []*Node

    querySubnet := fmt.Sprintf("SELECT IpNodeId,IPAddress,Comments,Status,Uri FROM IPAM.IPNode WHERE IPAddress='%s'", ipaddress)
    response = queryOrionServer(client, querySubnet)
    fmt.Printf("Release: IPADDRESS: %v, COMMENTS: %v, URI: %v\n", ipaddress, response[0].COMMENTS, response[0].URI)

    if(2 != response[0].STATUS) {
      updateIPNodeStatus(client, response[0].URI, "2", "") // '2' == ip available
    } else {
      log.Fatalf("[Error] Ip address is already available. It has status %v.", response[0].STATUS)
    }
    fmt.Println("Released")
}

func queryOrionServer(client * gosolar.Client, query string) []*Node {
    var data []*Node
    ipNodeDetails, err := client.Query(query, nil)
    if err != nil {
            log.Fatal(err)
    }
    if err := json.Unmarshal(ipNodeDetails, &data); err != nil {
            log.Fatal(err)
    }
    return data
}

func updateIPNodeStatus(client * gosolar.Client, uri string, status string, comment string) {
    req := map[string]interface{}{
            "Status": status,
            "Comments": comment,
    }
    _, err := client.Update(uri, req)
    if err != nil {
    	log.Printf("[ERROR]\n")
    	log.Printf("[ERROR] It looks like you do not have privilege to modify this VLAN.\n" )
    	log.Printf("[ERROR]\n")
        log.Fatal(err)
    }
}

func findIP(input string) string {
   numBlock := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
   regexPattern := numBlock + "\\." + numBlock + "\\." + numBlock + "\\." + numBlock

   regEx := regexp.MustCompile(regexPattern)
   return regEx.FindString(input)
}
