# go-tencent-api
Client packages for tencent API with corporation sites, like ExMail, ExWechat, ...

应用于企业版微信（不是微信企业号）和QQ企业邮箱相关的操作库

## ExMail: 用于企业邮箱



### Environment 环境变量

Put Auth string as variable `EXMAIL_API_AUTHS` in Environment

### Features 已完成特性

- `exmail.GetUser(alias string) (*User, error)`
- `exmail.CountNewMail(alias string) (int, error)`


## wxwork: 用于企业版微信

https://work.weixin.qq.com/api/doc

### Environment 环境变量

- `EXWECHAT_CORP_ID`: corpId，可在管理后台“我的企业”-“企业信息”下查看
- `EXWECHAT_CORP_SECRET`: 通讯录接口的密钥，在“管理工具”-“通讯录同步助手”可找到

### Features 已完成特性

- `API.GetUser(userId string) (*User, error)`
- `API.DelUser(userId string) error`
- `API.AddUser(user *User) error`

