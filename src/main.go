package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strconv"
	zabbix_client "zabbix-client"
	"zabbix-client/host"
	"zabbix-client/hostgroup"
	"zabbix-client/utils"

	"github.com/cavaliercoder/go-zabbix"

	"zabbix-client/common"
	"zabbix-client/item"
)

var apiHost = "http://203.255.255.101:8080/zabbix/api_jsonrpc.php"
var zabbixId = "Admin"
var zabbixPw = "zabbix"

var session *zabbix.Session


func createSession() {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	cache := zabbix.NewSessionFileCache().SetFilePath("./zabbix_session")
	_session, err := zabbix.CreateClient(apiHost).
		WithCache(cache).
		WithHTTPClient(client).
		WithCredentials("Admin", "zabbix").Connect()
	session = _session
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	version, err := session.GetVersion()
	fmt.Printf("Connected to Zabbix API v%s\n", version)
}






func main() {
	createSession()

	// 호스트 그룹 ID를 얻기 위해 호스트 그룹 정보를 조회
	groupParams := make(map[string]interface{}, 0)
	groupParams["name"] = "Discovered hosts"
	hostgroupList := hostgroup.GetHostgroup(session, groupParams)
	groupId := hostgroupList[0].GroupID


	// 호스트 그룹에 속해있는 hostList 목록을 조회
	hostParams := make(map[string]interface{}, 0)
		groupIds := make([]string, 1)
		groupIds[0] = groupId
	hostParams["groupIds"] = groupIds
	hostList := host.GetHostList(session, hostParams)
	fmt.Printf("result length : %d\n", len(hostList))
	
	hostIdArr := make([]string, len(hostList))
	for idx, host := range hostList {
		hostIdArr[idx] = host.HostID
	}

	// 호스트 그룹에 속해있는 특정 호스트 정보 조회
	hostParams = make(map[string]interface{}, 0)
		filterMap := make(map[string]interface{}, 0)
		filterMap["hostid"] = hostIdArr[0]
	hostParams["groupIds"] = groupIds
	hostParams["filter"] = filterMap
	hostList = host.GetHostList(session, hostParams)
	fmt.Printf("result length : %d\n", len(hostList))

	// 호스트의 Item 정보를 조회 - CPU 사용률
	itemParams := make(map[string]interface{}, 0)
		keywordArr := make([]string, 1)
		keywordArr[0] = common.SYSTEM_CPU_UTIL
	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIdArr
	item.GetItemList(session, itemParams)

	// 호스트의 Item 정보를 조회 - CPU 코어 개수
	itemParams = make(map[string]interface{}, 0)
		keywordArr = make([]string, 1)
		keywordArr[0] = common.SYSTEM_CPU_NUM
	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIdArr
	item.GetItemList(session, itemParams)

	// 호스트의 Item 정보를 조회 - 메모리 사용률
	itemParams = make(map[string]interface{}, 0)
		keywordArr = make([]string, 1)
		keywordArr[0] = common.VM_MEMORY_UTILIZATION
	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIdArr
	item.GetItemList(session, itemParams)

	// 호스트의 Item 정보를 조회 - 네트워크 인바운드 패킷
	itemParams = make(map[string]interface{}, 0)
	keywordArr = make([]string, 1)
		keywordArr[0] = "net.if.in[*"
	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIdArr[0:2]
	result := item.GetItemList(session, itemParams)
	var sum int = 0
	for _, item := range result {
		value, _ := strconv.Atoi(item.LastValue)
		sum = sum + value
	}
	fmt.Printf("Network inbound packet sum : %d\n", sum)

	// 호스트의 Item 정보를 조회 - 네트워크 아웃바운드 패킷
	itemParams = make(map[string]interface{}, 0)
		keywordArr = make([]string, 1)
		keywordArr[0] = "net.if.out[*"
	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIdArr[0:2]
	result = item.GetItemList(session, itemParams)
	sum = 0
	for _, item := range result {
		value, _ := strconv.Atoi(item.LastValue)
		sum = sum + value
	}
	fmt.Printf("Network outbound packet sum : %d\n", sum)

	// 호스트의 Item 정보를 조회 - 네트워크 아웃바운드 드랍 패킷
	itemParams = make(map[string]interface{}, 0)
		keywordArr = make([]string, 2)
		keywordArr[0] = common.NETWORK_OUTPUT_DROPPED_PACKET
	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIdArr[1:3]
	item.GetItemList(session, itemParams)


	// 호스트의 Item 정보를 조회 - 네트워크 아웃바운드 에러 패킷
	itemParams = make(map[string]interface{}, 0)
	keywordArr = make([]string, 2)
	keywordArr[0] = common.NETWORK_OUTPUT_ERROR_PACKET
	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIdArr[0:2]
	result = item.GetItemList(session, itemParams)
	sum = 0
	for _, item := range result {
		value, _ := strconv.Atoi(item.LastValue)
		sum = sum + value
	}
	fmt.Printf("sum : %d\n", sum)


	// IP주소로 해당 호스트의 CPU 차트 데이터, Network IO 패킷, Disk IO rate 데이터 추출
	hostInfo, err := zabbix_client.GetHostInfo(session, "10.37.0.140")
	if err == nil {
		cpuHistory := zabbix_client.GetHistory(session, common.SYSTEM_CPU_UTIL, hostInfo.HostID)
		utils.PrintJson(cpuHistory)

		networkInputHistory := zabbix_client.GetHistory(session, common.NETWORK_INPUT_PACKET, hostInfo.HostID)
		utils.PrintJson(networkInputHistory)

		networkOutputHistory := zabbix_client.GetHistory(session, common.NETWORK_OUTPUT_PACKET, hostInfo.HostID)
		utils.PrintJson(networkOutputHistory)

		diskReadRateHistory := zabbix_client.GetHistory(session, common.DISK_READ_RATE, hostInfo.HostID)
		utils.PrintJson(diskReadRateHistory)

		diskWriteRateHistory := zabbix_client.GetHistory(session, common.DISK_WRITE_RATE, hostInfo.HostID)
		utils.PrintJson(diskWriteRateHistory)
	}
}

