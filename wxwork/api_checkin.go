package wxwork

import (
	"time"

	"daxv.cn/gopak/tencent-api-go/client"
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

func (a *API) ListCheckin(days int, userIDs ...string) (result *CheckInResult, err error) {
	if len(userIDs) == 0 {
		err = ErrEmptyArg
		return
	}
	result = new(CheckInResult)

	if days == 0 {
		days = 7
	}
	if days > 30 {
		err = ErrOutofRange
		return
	}

	startTime := time.Now().Add(0 - time.Hour*24*time.Duration(days)).Unix()
	endTime := startTime + int64(time.Hour*24*7/time.Second)
	req := CheckInReq{
		OpenCheckInDataType: 3,
		StartTime:           startTime,
		EndTime:             endTime,
		UserIdList:          userIDs,
	}

	err = a.c.PostJSON(UriPrefix+"/checkin/getcheckindata", client.MustMarshal(&req), result)
	return
}
