package hwsinger

import (
	"encoding/json"
	"hwcloud-aom-operator/hwsinger/core"
	"hwcloud-aom-operator/utils"
	"io/ioutil"
	"net/http"
	"strings"
)

func HwSinger(method string, url string, headerMap map[string]string, postBody map[string]interface{}) ([]byte, int, error) {
	s := core.Signer{
		Key:    utils.GetConfig("hwcloud::key").String(),
		Secret: utils.GetConfig("hwcloud::secret").String(),
	}

	data, err := json.Marshal(postBody)
	if err != nil {
		return nil, 400, err
	}
	r, _ := http.NewRequest(method, url, strings.NewReader(string(data)))

	// 循环添加请求头
	for k, v := range headerMap {
		r.Header.Add(k, v)
	}

	_ = s.Sign(r)
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, 400, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}
