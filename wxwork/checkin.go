package wxwork

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

const (
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
	api *API
}

func NewCAPI() *CAPI {
	api := NewAPI(os.Getenv("EXWECHAT_CORP_ID"), os.Getenv("EXWECHAT_CHECKIN_SECRET"))
	return &CAPI{api}
}

func (a *CAPI) ListCheckin(days int, userIDs ...string) (result *CheckInResult, err error) {
	if len(userIDs) == 0 {
		err = ErrEmptyArg
		return
	}
	result = new(CheckInResult)
	var token string
	token, err = a.api.c.GetAuthToken()
	if err != nil {
		return nil, err
	}
	if days == 0 {
		days = 7
	}
	if days > 30 {
		err = ErrOutofRange
		return
	}

	startTime := time.Now().Add(0 - time.Hour*24*time.Duration(days)).Unix()
	endTime := startTime + int64(time.Hour*24*7/time.Second)
	log.Printf("start %v, end %v", startTime, endTime)
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
	err = a.api.c.PostJSON(uri, data, result)
	return
}
