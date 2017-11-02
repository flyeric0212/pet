
# [错误码说明]

[data]
+ 200: OK
+ 500: Internal Server Error
+ 501: 参数异常
+ 502: 电话号码重复


# [ 微信接口 api Doc ] #
---

# [用户电话注册 - `POST /api/users/phone_regist`]
+ **创建**(`liangbo`, `2017-10-27`)

+ Description

		电话号码注册，姓名/电话号码不能为空，如果regist_type=1，openid不能为空

+ Request:

		{
			“name”: (required, string, 姓名),
			"gender": (required, int，性别，0: 无性别 1: 男 2: 女)
			"phone": (required, string, 电话号码)
			"email": (optional, string, 邮箱地址)
			"regist_type": (required, int, 注册方式，1：微信  2：官网)
			// 微信数据
			"openid": (required, string, 微信用户标志id)
			"avatar": (optional, string, 微信头像)
			"nickname": (optional, string, 微信昵称)
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

# [发送验证码 - `POST /api/users/send_verify_code`]
+ **创建**(`liangbo`, `2017-11-02`)


# [校验验证码 - `POST /api/users/check_verify_code`]
+ **创建**(`liangbo`, `2017-11-02`)


# [ Pet官网 api Doc ] #
---

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