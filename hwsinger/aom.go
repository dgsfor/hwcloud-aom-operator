package hwsinger

import (
	"fmt"
	"hwcloud-aom-operator/utils"
)

type Response struct {
	Code  int    `json:"code"`
	Data  []byte `json:"data,omitempty"`
	Msg   string `json:"msg,omitempty"`
	Error string `json:"error,omitempty"`
}

type AomGroupResult struct {
	Config struct {
		CooldownTime int    `json:"cooldown_time"`
		ID           string `json:"id"`
		MaxInstances int    `json:"max_instances"`
		MinInstances int    `json:"min_instances"`
	} `json:"config,omitempty"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func GetAomGroupWithDeployment(namespace string, deployment string) Response {
	headerMap := make(map[string]string)
	headerMap["content-type"] = "application/json"
	headerMap["ResourceType"] = "app"
	headerMap["Cluster-Id"] = utils.GetConfig("hwcloud::cluster_id").String()
	headerMap["Namespace"] = namespace
	headerMap["Deployment-Name"] = deployment
	url := fmt.Sprintf("https://aom.%s.myhuaweicloud.com/v1/%s/pe/policy/config",
		utils.GetConfig("hwcloud::region").String(), utils.GetConfig("hwcloud::project_id").String())
	response, code, err := HwSinger("GET", url, headerMap, nil)
	if err != nil {
		return Response{
			Code:  code,
			Data:  response,
			Error: err.Error(),
		}
	}
	return Response{
		Code: code,
		Data: response,
	}
}
