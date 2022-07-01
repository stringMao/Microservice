package codes



// 用户登入网关服结果
const (
	Code_LoginGateSvrAuthFail = 1  //网关服登入失败-身份验证未通过
	Code_LoginGateSvrReplaceFail=2  //网关服登入失败-顶号失败
)


// 服务器注册网关服失败
const (
	Code_SvrRegister_AuthFail = 1  //身份认证失败
	Code_SvrRegister_ExistedFail = 2  //重复注册
)

const (
	Code_JoinSvr_SvrNoFind =1  //服务器未找到
	Code_JoinSvr_AgainJoin =2  //重复加入

)