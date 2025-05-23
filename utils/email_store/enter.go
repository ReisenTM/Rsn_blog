package email_store

import "sync"

// EmailStoreInfo 邮箱验证码存放容器
type EmailStoreInfo struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

var emailVerifyStore = sync.Map{}

// Set 存入验证码，以用来验证
func Set(id, email, code string) {
	emailVerifyStore.Store(id, EmailStoreInfo{
		Email: email,
		Code:  code,
	})
}

// Verify 验证
func Verify(id, code string) (info EmailStoreInfo, ok bool) {
	value, ok := emailVerifyStore.Load(id)
	if !ok {
		//如果没找到
		return
	}
	info, ok = value.(EmailStoreInfo)
	if !ok {
		//如果取出来的数据不对
		return
	}
	if info.Code != code {
		//如果验证码不对
		emailVerifyStore.Delete(id)
		ok = false
		return
	}
	emailVerifyStore.Delete(id)
	return
}
