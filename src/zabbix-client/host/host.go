package host

import (
	"log"

	"github.com/cavaliercoder/go-zabbix"

	"zabbix-client/utils"
)

var isDebug bool

func init() {
	isDebug = true
}

/**
	GetHostList
		Parameters
			- filter (map[string]interface{}) : 필터 맵
			- groupIds ([]string) : 호스트가 속한 그룹의 ID 배열
			- hostIds ([]string) : 호스트 ID 배열
 */
func GetHostList(session *zabbix.Session, params map[string]interface{}) []zabbix.Host {
	var hostParams zabbix.HostGetParams

	filterMap, ok := params["filter"]
	if ok {
		hostParams.Filter = filterMap.(map[string]interface{})
	}

	groupIds, ok := params["groupIds"]
	if ok {
		hostParams.GroupIDs = groupIds.([]string)
	}


	hostParams.SelectItems = zabbix.SelectFields{"name", "lastvalue", "units", "itemid", "lastclock", "value_type"}
	hostParams.SelectInterfaces = zabbix.SelectFields{"ip"}
	hostParams.OutputFields = "extend"
	result, err := session.GetHosts(hostParams)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	if isDebug {
		utils.PrintJson(result)
	}

	return result
}
