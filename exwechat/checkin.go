package exwechat

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var (
	urlCheckinData = "https://qyapi.weixin.qq.com/cgi-bin/checkin/getcheckindata"
)

type CheckInReq struct {
	OpenCheckInDataType int      `json:"opencheckindatatype,omitempty"`
	StartTime           int64    `json:"starttime,omitempty"`
	EndTime             int64    `json:"endtime,omitempty"`
	UserIdList          []string `json:"useridlist,omitempty"`
}

type CheckInResult struct {
	ErrCode     int           `json:"errcode,omitempty"`
	ErrMsg      string        `json:"errmsg,omitempty"`
	CheckInData []CheckInData `json:"checkindata,omitempty"`
}

type CheckInData struct {
	UserID         string   `json:"userid,omitempty"`
	GroupName      string   `json:"groupname,omitempty"`
	CheckInType    string   `json:"checkin_type,omitempty"`
	ExceptionType  string   `json:"exception_type,omitempty"`
	CheckInTime    int64    `json:"checkin_time,omitempty"`
	LocationTitle  string   `json:"location_title,omitempty"`
	LocationDetail string   `json:"location_detail,omitempty"`
	WifiName       string   `json:"wifiname,omitempty"`
	Notes          string   `json:"notes,omitempty"`
	WifiMac        string   `json:"wifimac,omitempty"`
	Mediaids       []string `json:"mediaids,omitempty"`
}

type CAPI struct {
	*API
}

func NewCAPI() *CAPI {
	return &CAPI{API: New(os.Getenv("EXWECHAT_CORP_ID"), os.Getenv("EXWECHAT_CHECKIN_SECRET"))}
}

func (a *CAPI) GetCheckInData(userIDs []string, startTime int64) (result *CheckInResult, err error) {
	result = new(CheckInResult)
	var token string
	token, err = a.c.GetAuthToken()
	if err != nil {
		return nil, err
	}

	endTime := time.Unix(startTime, 0).AddDate(0, 1, 0).Unix()
	req := CheckInReq{
		OpenCheckInDataType: 3,
		StartTime:           startTime,
		EndTime:             endTime,
		UserIdList:          userIDs,
	}
	var data []byte
	data, err = json.Marshal(req)
	if err != nil {
		return
	}

	uri := fmt.Sprintf("%s?access_token=%s", urlCheckinData, token)
	err = a.c.PostJSON(uri, data, result)
	return
}
