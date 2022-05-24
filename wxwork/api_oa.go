package wxwork

import (
	"strconv"
	"time"

	"daxv.cn/gopak/tencent-api-go/client"
	"daxv.cn/gopak/tencent-api-go/models/oa"
)

type reqOAGetTemplateDetail struct {
	TemplateID string `json:"template_id"`
}

type respOAGetTemplateDetail struct {
	client.Error

	oa.TemplateDetail
}

type reqOAApplyEvent struct {
	oa.ApplyEvent
}

type respOAApplyEvent struct {
	client.Error

	// SpNo 表单提交成功后，返回的表单编号
	SpNo string `json:"sp_no"`
}

type reqOAGetApprovalInfo struct {
	StartTime string                  `json:"starttime"`
	EndTime   string                  `json:"endtime"`
	Cursor    int                     `json:"cursor"`
	Size      uint32                  `json:"size"`
	Filters   []oa.ApprovalInfoFilter `json:"filters"`
}

type respOAGetApprovalInfo struct {
	client.Error

	// SpNoList 审批单号列表，包含满足条件的审批申请
	SpNoList []string `json:"sp_no_list"`
}

type reqOAGetApprovalDetail struct {
	// SpNo 审批单编号。
	SpNo string `json:"sp_no"`
}

type respOAGetApprovalDetail struct {
	client.Error

	// Info 审批申请详情
	Info oa.ApprovalDetail `json:"info"`
}

// GetOAApprovalInfoReq 批量获取审批单号请求
type GetOAApprovalInfoReq struct {
	// StartTime 审批单提交的时间范围，开始时间，UNix时间戳
	StartTime time.Time
	// EndTime 审批单提交的时间范围，结束时间，Unix时间戳
	EndTime time.Time
	// Cursor 分页查询游标，默认为0，后续使用返回的next_cursor进行分页拉取
	Cursor int
	// Size 一次请求拉取审批单数量，默认值为100，上限值为100
	Size uint32
	// Filters 筛选条件，可对批量拉取的审批申请设置约束条件，支持设置多个条件
	Filters []oa.ApprovalInfoFilter
}

// GetOATemplateDetail 获取审批模板详情
func (a *API) GetOATemplateDetail(templateID string) (*oa.TemplateDetail, error) {
	var resp respOAGetTemplateDetail
	err := a.c.PostJSON(UriPrefix+"/oa/gettemplatedetail", client.MustMarshal(&reqOAGetTemplateDetail{
		TemplateID: templateID,
	}), &resp)
	if err != nil {
		return nil, err
	}

	return &resp.TemplateDetail, nil
}

// ApplyOAEvent 提交审批申请
func (a *API) ApplyOAEvent(applyInfo oa.ApplyEvent) (string, error) {
	var resp respOAApplyEvent
	err := a.c.PostJSON(UriPrefix+"/oa/applyevent", client.MustMarshal(&reqOAApplyEvent{
		ApplyEvent: applyInfo,
	}), &resp)
	if err != nil {
		return "", err
	}
	return resp.SpNo, nil
}

// GetOAApprovalInfo 批量获取审批单号
func (a *API) GetOAApprovalInfo(req GetOAApprovalInfoReq) ([]string, error) {
	var resp respOAGetApprovalInfo
	err := a.c.PostJSON(UriPrefix+"/oa/getapprovalinfo", client.MustMarshal(&reqOAGetApprovalInfo{
		StartTime: strconv.FormatInt(req.StartTime.Unix(), 10),
		EndTime:   strconv.FormatInt(req.EndTime.Unix(), 10),
		Cursor:    req.Cursor,
		Size:      req.Size,
		Filters:   req.Filters,
	}), &resp)
	if err != nil {
		return nil, err
	}
	return resp.SpNoList, nil
}

// GetOAApprovalDetail 提交审批申请
func (a *API) GetOAApprovalDetail(spNo string) (*oa.ApprovalDetail, error) {
	var resp respOAGetApprovalDetail
	err := a.c.PostJSON(UriPrefix+"/oa/getapprovaldetail", client.MustMarshal(&reqOAGetApprovalDetail{
		SpNo: spNo,
	}), &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Info, nil
}
