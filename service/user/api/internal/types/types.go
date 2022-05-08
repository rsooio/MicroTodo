// Code generated by goctl. DO NOT EDIT.
package types

type LoginReq struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type LoginReply struct {
	Id            int64  `json:"id"`
	Username      string `json:"username"`
	Gender        string `json:"gender"`
	AccessToken   string `json:"accessToken"`
	AccessExpire  int64  `json:"accessExpire"`
	RefreshToken  string `json:"refreshToken"`
	RefreshExpire int64  `json:"refreshExpire"`
}

type RegisterReq struct {
	Username string `form:"username"`
	Password string `form:"password"`
	Gender   string `form:"gender,default=none,options=male|female|none"`
}

type UpdateInfoReq struct {
	Username string `json:"username,optional"`
	Gender   string `json:"gender,optional"`
}

type UpdatePwdReq struct {
	OldPassword string `form:"old_password"`
	NewPassword string `form:"new_password"`
}

type TokenRefreshReply struct {
	AccessToken   string `json:"accessToken"`
	AccessExpire  int64  `json:"accessExpire"`
	RefreshToken  string `json:"refreshToken"`
	RefreshExpire int64  `json:"refreshExpire"`
}