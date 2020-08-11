package core

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

func ListJobStats() ([]*JobInfo, error) {
	url := "/schedule/job/{ids}"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("err:%s", resp.Status)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ret := make([]*JobInfo, 0)
	err = jsoniter.ConfigDefault.Unmarshal(content, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
