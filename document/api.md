
# [错误码说明]

[data]
+ 200: OK
+ 500: Internal Server Error
+ 501: 参数异常
+ 502: 用户名重复
+ 503: 用户名不存在
+ 504: 密码错误


# [用户名注册 - `POST /api/users/nickname_regist`]
+ **创建**(`liangbo`, `2017-10-20`)

+ Description

		用户名注册，包含错误码有 501，502

+ Request:

		{
			“nickname”: (required, string, 用户名),
			"password": (required, string，密码)
		}

+ Response Succ:

	     {
		 	"status": "OK",
		  	"data": {
                 "user_id": (int, 用户id)
                 "name": (string, 姓名),
                 "nickname": (string, 用户名),
                 "gender": (string, 性别，0: 无性别 1: 男 2: 女),
                 "phone": (string, 电话号码),
                 "email": (string, 邮件地址),
                 "openid": (string, 微信公共号用户唯一标志)
		  	}
		   	"desc": ""
	      }

+ Response Error:

		{
			"status": (required, string, 'Error', '返回状态 OK/Error'),
			"data": (required, string, '', '返回错误code'),
			"desc": (required, string, '', '返回描述，错误时描述')
		}


# [用户名登陆 - `POST /api/users/nickname_login`]
+ **创建**(`liangbo`, `2017-10-20`)

+ Description

		用户名登陆，包含错误码有 501，503，504

+ Request:

		{
			“nickname”: (required, string, 用户名),
			"password": (required, string，密码)
		}

+ Response Succ:

	     {
		 	"status": "OK",
		  	"data": {
                 "user_id": (int, 用户id)
                 "name": (string, 姓名),
                 "nickname": (string, 用户名),
                 "gender": (string, 性别，0: 无性别 1: 男 2: 女),
                 "phone": (string, 电话号码),
                 "email": (string, 邮件地址),
                 "openid": (string, 微信公共号用户唯一标志)
		  	}
		   	"desc": ""
	      }

+ Response Error:

		{
			"status": (required, string, 'Error', '返回状态 OK/Error'),
			"data": (required, string, '', '返回错误code'),
			"desc": (required, string, '', '返回描述，错误时描述')
		}