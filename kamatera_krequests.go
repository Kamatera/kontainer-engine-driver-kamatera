package main

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var cloudcliBaseUrl = "https://cloudcli.cloudwm.com"
var cloudcliDebug = ""

type KTaskIdResponse struct {
	TaskId string `json:"task_id"`
}

type KTaskResponse struct {
	TaskName *string      `json:"task_name"`
	State    string       `json:"state"`
	Result   *interface{} `json:"result"`
	Error    *string      `json:"error"`
	Meta     *interface{} `json:"meta"`
}

type KubeConfig struct {
	ApiVersion     string `json:"apiVersion"`
	Kind           string `json:"kind"`
	CurrentContext string `json:"current-context"`
	Preferences    interface{}
	Clusters       []struct {
		Name    string `json:"name"`
		Cluster struct {
			CertificateAuthorityData string `json:"certificate-authority-data"`
			Server                   string `json:"server"`
		} `json:"cluster"`
	}
	Contexts []struct {
		Name    string `json:"name"`
		Context struct {
			Cluster string `json:"cluster"`
			User    string `json:"user"`
		} `json:"context"`
	}
	Users []struct {
		Name string `json:"name"`
		User struct {
			ClientCertificateData string `json:"client-certificate-data"`
			ClientKeyData         string `json:"client-key-data"`
		}
	}
}

func KPost(apiClientid string, apiSecret string, path string, data url.Values) ([]byte, error) {
	if cloudcliDebug == "true" {
		logrus.Infof("POST %s/k8s/%s", cloudcliBaseUrl, path)
		logrus.Infof("data: %s", data)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/k8s/%s", cloudcliBaseUrl, path), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("AuthClientId", apiClientid)
	req.Header.Set("AuthSecret", apiSecret)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if cloudcliDebug == "true" {
		logrus.Infof("response: %s", body)
	}
	return body, nil
}

func KCreateCluster(apiClientid string, apiSecret string, config KConfig) (*KTaskResponse, error) {
	return KClusterTask(apiClientid, apiSecret, config, "create_cluster")
}

func KGetKubeconfig(apiClientid string, apiSecret string, config KConfig) (*KubeConfig, error) {
	taskResponse, err := KClusterTask(apiClientid, apiSecret, config, "kubeconfig?json_format=true")
	if err != nil {
		return nil, err
	}
	if taskResponse.State != "SUCCESS" {
		return nil, fmt.Errorf("failed to get kubeconfig: %s", *taskResponse.Error)
	}
	if taskResponse.Result == nil {
		return nil, fmt.Errorf("failed to get kubeconfig: result is nil")
	}
	kubeconfigJson := (*taskResponse.Result).(string)
	var kubeconfig KubeConfig
	err = json.Unmarshal([]byte(kubeconfigJson), &kubeconfig)
	if err != nil {
		return nil, err
	}
	return &kubeconfig, nil
}

func KClusterTask(apiClientid string, apiSecret string, config KConfig, taskName string) (*KTaskResponse, error) {
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set("kconfig", string(jsonBytes))
	body, err := KPost(apiClientid, apiSecret, taskName, data)
	if err != nil {
		return nil, err
	}
	var response KTaskIdResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return KWaitTask(apiClientid, apiSecret, response.TaskId)
}

func KWaitTask(apiClientId string, apiSecret string, taskId string) (*KTaskResponse, error) {
	for {
		data := url.Values{}
		data.Set("task_id", taskId)
		body, err := KPost(apiClientId, apiSecret, "task_status", data)
		if err != nil {
			return nil, err
		}
		var response KTaskResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}
		if response.State != "PENDING" {
			return &response, nil
		}
		time.Sleep(5 * time.Second)
	}
}
