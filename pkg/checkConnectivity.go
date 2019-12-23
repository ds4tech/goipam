package ipam

import (
  "fmt"
  "log"
  "net"
  "net/http"
  "crypto/tls"
  "time"
  "io/ioutil"
  "github.com/sparrc/go-ping"
)

func checkDhcp(ipaddress string) bool {
  var response []*Node

  queryIpnode := fmt.Sprintf("SELECT IpNodeId,SubnetId,Uri FROM IPAM.IPNode WHERE IPAddress='%s'", ipaddress)
  response = queryOrionServer(client, queryIpnode)
  log.Println("[DEBUG] IPAM.IPNode", response[0])

  queryDhcp := fmt.Sprintf("SELECT NodeId,ScopeId FROM IPAM.DhcpScope WHERE SubnetId='%d'", response[0].SUBNETID)
  response = queryOrionServer(client, queryDhcp )
  log.Println("[DEBUG] IPAM.DhcpScope", queryDhcp)

  if( len(response)>0 ){
  log.Println("[DEBUG] IPAM.DhcpScope", response[0])
  return true 
  } else {
  return false
  }
}

func checkConnectivity(ipaddress string) bool {
  ports := []string{"22","80","443","3389", "17778"}
  if( pingTest(ipaddress) ){
    log.Println("[DEBUG] Ping test SUCCESS")
    if( ncPortTest(ipaddress, ports) ){
      log.Println("[DEBUG] NCPort test SUCCESS")
        if( curlPortTest(ipaddress, "80") ) {
          log.Println("[DEBUG] CURL test SUCCESS")
          return true
        } else {
          log.Println("[DEBUG] CURL test FAILED")
          return false
        }
    } else {
      log.Println("[DEBUG] NC Port test FAILED")
      return false
    }
  } else {
    log.Println("[DEBUG] Ping test FAILED")
    return false
  }
}

func pingTest(host string) bool {
  var packetLoss float64
  pinger, err := ping.NewPinger(host)
  if err != nil {
      log.Fatalf("[ERROR] %s\n", err.Error()) //przerywa działanie programu
      return false
  }

  pinger.OnRecv = func(pkt *ping.Packet) {
        log.Printf("[DEBUG] %d bytes from %s: icmp_seq=%d time=%v\n",
                            pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
  }
  pinger.OnFinish = func(stats *ping.Statistics) {
        log.Printf("[DEBUG] --- %s ping statistics ---\n", stats.Addr)
        log.Printf("[DEBUG] %d packets transmitted, %d packets received, %v%% packet loss\n",
                                stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
        log.Printf("[DEBUG] round-trip min/avg/max/stddev = %v/%v/%v/%v\n", stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
        packetLoss = stats.PacketLoss
  }

  log.Printf("[DEBUG] PING %s:\n", pinger.IPAddr())
  pinger.Count = 3
  pinger.Timeout = 3 * time.Second
  pinger.Run()

  if(packetLoss > 50.0) { // 50 reflects 50 % of loss
    return false
  } else {
    return true
  }
}

func ncPortTest(host string, ports []string) bool {
  log.Println("[DEBUG] Trying to connect on ports: ", ports)
  for _, port := range ports {
      timeout := time.Second
      conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
      if err != nil {
      	  log.Fatalf("[ERROR] Connecting error: %s\n", err) //przerywa działanie programu
          return false
      }
      if conn != nil {
          defer conn.Close()
          log.Println("[DEBUG] Opened", net.JoinHostPort(host, port))
      }
  }
  return true
}

func curlPortTest(host string, port string) bool {
  url := fmt.Sprintf("https://%s:%s", host, port)
  log.Printf("[DEBUG] CURL to %s", url)

  transCfg := &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
  }

  var netClient = &http.Client{
    Transport: transCfg,
    Timeout: time.Second * 5,
  }

  resp, err := netClient.Get(url)
  if err != nil {
  	log.Println(err)
  	return false
  }

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
  	log.Fatalln(err)
  }

  log.Println(string(body))
  return true
}
