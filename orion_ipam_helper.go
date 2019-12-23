package main

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

    subnetIpAddr := findIP(vlan)
    cidr := string(vlan[len(vlan)-2:]) //substring last 2 chars

    log.Printf("[TRACE] vsphere_vlan: %s, subnetIpAddr: %s, cidr: %s", vlan, subnetIpAddr, cidr)

    querySubnet := fmt.Sprintf("SELECT SubnetId FROM IPAM.Subnet WHERE  Address='%s' AND CIDR=%s", subnetIpAddr, cidr)
    response = queryOrionServer(client, querySubnet)
    log.Printf("[DEBUG] SUBNETID: %v\n", response[0].SUBNETID)

    queryIpnode := fmt.Sprintf("SELECT IpNodeId,IPAddress,Comments,Status,Uri,DisplayName FROM IPAM.IPNode WHERE SubnetId='%d' and status=2", response[0].SUBNETID)
    response = queryOrionServer(client, queryIpnode)

    log.Println("[DEBUG] DisplayName, IPAddress, comments")
    for _, el := range response {
      log.Printf("[DEBUG] %s, %s, %s\n", el.DISPLAYNAME, el.IPADDRESS, el.COMMENTS)
    }
}

func reserveIpAddress(client * gosolar.Client, ipaddress string, comment string)  {
    var response []*Node

    querySubnet := fmt.Sprintf("SELECT IpNodeId,IPAddress,Comments,Status,Uri FROM IPAM.IPNode WHERE IPAddress='%s'", ipaddress)
    response = queryOrionServer(client, querySubnet)
    log.Printf("[DEBUG] Reserve: IPADDRESS: %v, COMMENTS: %v\n", response[0].IPADDRESS, response[0].COMMENTS)

    if(2 == response[0].STATUS) {
      updateIPNodeStatus(client, response[0].URI, "1", comment) // '1' == ip used
    } else {
      log.Fatalf("[DEBUG] Ip address is not available. It has status %v.", response[0].STATUS)
    }
    log.Printf("[DEBUG] Reserved")
}

func changeIpAddressStatus(client * gosolar.Client, ipaddress string)  {}

func releaseIpAddress(client * gosolar.Client, ipaddress string)  {
    var response []*Node

    querySubnet := fmt.Sprintf("SELECT IpNodeId,IPAddress,Comments,Status,Uri FROM IPAM.IPNode WHERE IPAddress='%s'", ipaddress)
    response = queryOrionServer(client, querySubnet)
    log.Printf("[DEBUG] Release: IPADDRESS: %v, COMMENTS: %v\n", response[0].IPADDRESS, response[0].COMMENTS)

    if(2 != response[0].STATUS) {
      updateIPNodeStatus(client, response[0].URI, "2", "") // '2' == ip available
    } else {
      log.Fatalf("[DEBUG] Ip address is already available. It has status %v.", response[0].STATUS)
    }
    log.Printf("[DEBUG] Released")
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
            log.Fatal(err)
    }
}

func findIP(input string) string {
   numBlock := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
   regexPattern := numBlock + "\\." + numBlock + "\\." + numBlock + "\\." + numBlock

   regEx := regexp.MustCompile(regexPattern)
   return regEx.FindString(input)
}
