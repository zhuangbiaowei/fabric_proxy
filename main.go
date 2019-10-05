package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	. "./util"
	"strconv"
	"strings"
	"os"
	"time"
)

type RestapiReq struct {
	ChannelId        string    `json:"channelId"`
	ChaincodeId      string    `json:"chaincodeId"`
	ChaincodeVersion string    `json:"chaincodeVersion"`
	UserId           string    `json:"userId"`
	OrgId            string    `json:"orgId"`
	OrgPeers         string    `json:"orgPeers"`
	Opmethod         string    `json:"opmethod"`  //invoke / query
	Args             string    `json:"args"`      //[“invoke”,”a”,”b”,”20”]
	Timestamp        string    `json:"timestamp"` //时间戳格式："2006-01-02 15:04:05"
	Cert             string    `json:"cert"`      //用户证书
}

type ReqBody struct {
	SignGzipSwitch   string    `yaml:"SignGzipSwitch"`
	ChannelId        string    `yaml:"ChannelId"`
	ChaincodeId      string    `yaml:"ChaincodeId"`
	ChaincodeVersion string    `yaml:"ChaincodeVersion"`
	UserId           string    `yaml:"UserId"`
	OrgId            string    `yaml:"OrgId"`
	OrgPeers         string    `yaml:"OrgPeers"`
	Opmethod         string    `yaml:"Opmethod"`
	Args             string    `yaml:"Args"`
}

type TransactionID struct {
	ID    string
	Nonce []byte
}

type OrgPeer struct {
	OrgId          string
	PeerDomainName string
}


func main() {
	InitConfig("./config/config.yaml")
	signCert := getCert(GlobalConfig.SignCert)
	req := ReqBody{}
	req.SignGzipSwitch = "1"
	req.ChaincodeId = GlobalConfig.ChaincodeId
	req.ChaincodeVersion = GlobalConfig.ChaincodeVersion
	req.ChannelId = GlobalConfig.ChannelId
	req.UserId = GlobalConfig.UserId
	req.OrgId = GlobalConfig.OrgId
	if (len(os.Args)>1) {
		req.Opmethod = os.Args[1]
		if(req.Opmethod == "invoke"){
			orgPeers := []OrgPeer{}
			for _, v := range GlobalConfig.Peers {
				orgPeer := OrgPeer{
					OrgId: GlobalConfig.OrgId,
					PeerDomainName: v,
				}
				orgPeers = append(orgPeers, orgPeer)
			}
			orgPeersByte, _ := json.Marshal(orgPeers)
			req.OrgPeers = string(orgPeersByte)
		}
	}
	if (len(os.Args)>2) {
		args := os.Args[2:]
		req.Args = "[\"" + strings.Join(args, "\",\"") + "\"]"
	}
	if(req.Args!=""){
		result := DoReq(string(signCert), GlobalConfig.PrvKey, req)
		fmt.Println(string(result))
	}
}

func getCert(path string) []byte {
	certPEMBlock, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("cert path err" + err.Error())
	}
	return certPEMBlock
}

func DoReq(signCert, userPrvKeyPath string, reqPara ReqBody) (respBody []byte) {
	tempReqInvoke := RestapiReq{
		ChannelId:        reqPara.ChannelId,
		ChaincodeId:      reqPara.ChaincodeId,
		ChaincodeVersion: reqPara.ChaincodeVersion,
		UserId:           reqPara.UserId,
		OrgId:            reqPara.OrgId,
		OrgPeers:         reqPara.OrgPeers,
		Opmethod:         reqPara.Opmethod,
		Args:             reqPara.Args,
		Timestamp:        time.Now().Format(time.RFC3339),
		Cert:             signCert,
	}
	reqBody, err := json.Marshal(tempReqInvoke)
	//fmt.Println("user post request is :", string(tempReqInvokeBody))
	signGzipSwitch := reqPara.SignGzipSwitch
	//fmt.Println("signGzipSwitch:",signGzipSwitch)
	if signGzipSwitch != "1" && signGzipSwitch != "0" {
		fmt.Println("The Header x-bcs-signature-sign-gzip does not match the rule.signGzipSwitch:", signGzipSwitch)
		return
	}

	signresult, err := EcdsaSignWithSha256Hex(reqBody, userPrvKeyPath, signGzipSwitch)
	if err != nil {
		fmt.Println("Signing reqeuset body failed " + err.Error())
	}
	// fmt.Println("the encode result is:", signresult)
	headers := make(map[string]string)

	headers["x-bcs-signature-sign"] = signresult
	headers["x-bcs-signature-method"] = GlobalConfig.CryptoMethod
	headers["x-bcs-signature-sign-gzip"] = signGzipSwitch

	var statusCode int
	resp, err := DoHTTPRequest("POST", GlobalConfig.Endpoint, GlobalConfig.Path, headers, nil, reqBody)
	defer ReleaseBody(resp)
	if err != nil {
		fmt.Println("POST request to " + GlobalConfig.Endpoint + " failed: " + err.Error())
	} else if !IsResponseStatusOk(resp) {
		if respBadStatusBody, e := CopyResponseBody(resp); e == nil {
			fmt.Println("Post status not ok: ", string(respBadStatusBody))
		}
		statusCode = int(resp.StatusCode)
		fmt.Println(GlobalConfig.Endpoint + " response state not OK, code: " + strconv.Itoa(statusCode))
	} else if respBody, err = CopyResponseBody(resp); err != nil {
		fmt.Println("copy response body failed: " + err.Error())
	}

	return
}