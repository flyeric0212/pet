
# [错误码说明]

[data]
+ 200: OK
+ 500: Internal Server Error
+ 501: 参数异常
+ 502: 用户名重复
+ 503: 用户名不存在
+ 504: 密码错误

# 微信接口 api Doc #
---

# [用户电话注册 - `POST /api/users/phone_regist`]
+ **创建**(`liangbo`, `2017-10-27`)

+ Description

		电话号码注册，姓名/电话号码/微信用户id 不能为空

+ Request:

		{
			“name”: (required, string, 姓名),
			"gender": (required, int，性别，0: 无性别 1: 男 2: 女)
			"phone": (required, string, 电话号码)
			"email": (optional, string, 邮箱地址)
			"openid": (required, string, 微信用户标志id)
		}

+ Response Succ:

	     {
		 	"status": "OK",
		  	"data": {
                 "user_id": (int, 用户id)
                 "name": (string, 姓名),
                 "nickname": (string, 用户名),
                 "avatar": (string, 头像),
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


# Pet官网 api Doc #
---

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
                 "avatar": (string, 头像),
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
                 "avatar": (string, 头像),
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


# [Banner列表 - `GET /api/banner/get_banner_list`]
+ **创建**(`liangbo`, `2017-10-22`)

+ Description

		banner列表，分页

+ Request:

		{
			“type”: (optional, int, 类型，1:首页 2:展商 3: 合作媒体),
			"page_num": (optional, int，页码，默认1)
			"page_size": (optional, int, 分页大小，默认10)
		}

+ Response Succ:

	     {
		 	"status": "OK",
		  	"data": {
                "banner_list": [
                    {
                        "id": (int, 主键),
                        "pic": (string, 图片地址),
                        "ref_url": (string, 跳转地址),
                        "type": (int, banner类型)
                    },
                    ...
                ],
                "total_num": (int, 总数)
		  	}
		   	"desc": ""
	      }

+ Response Error:

		{
			"status": (required, string, 'Error', '返回状态 OK/Error'),
			"data": (required, string, '', '返回错误code'),
			"desc": (required, string, '', '返回描述，错误时描述')
		}


# [文章列表 - `GET /api/article/get_article_list`]
+ **创建**(`liangbo`, `2017-10-22`)

+ Description

		article列表，分页

+ Request:

		{
			“type”: (optional, int, 类型，1:会展动态),
			"page_num": (optional, int，页码，默认1)
			"page_size": (optional, int, 分页大小，默认10)
		}

+ Response Succ:

	     {
		 	"status": "OK",
		  	"data": {
                "article_list": [
                    {
                        "id": (int, 主键),
                        "titile": (string, 标题),
                        "content": (string, 内容),
                        "type": (int, 文章类型)
                    },
                    ...
                ],
                "total_num": (int, 总数)
		  	}
		   	"desc": ""
	      }

+ Response Error:

		{
			"status": (required, string, 'Error', '返回状态 OK/Error'),
			"data": (required, string, '', '返回错误code'),
			"desc": (required, string, '', '返回描述，错误时描述')
		}