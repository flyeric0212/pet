/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/10/25 21:30 
 */
package main

import (
    "github.com/chanxuehong/wechat.v2/mp/menu"
    "fmt"
)

// 初始化微信菜单
// /var/www/go_workspace/bin/pet -a init_weixin_menu
func InitWinxinMenuList() {
    // 先清空菜单
    err := menu.Delete(wechatClient)
    if nil != err {
        fmt.Printf("delete menu err: %v \n", err)
    }
    // 初始化菜单
    menu_val := new(menu.Menu)

    button_list := make([]menu.Button, 0)

    button := menu.Button{}
    button.Type = "view"
    button.Name = "观众中心"
    button.URL = "http://mp.petfair.cc/api/vistor_center_auth"
    button_list = append(button_list, button)


    menu_val.Buttons = button_list
    err = menu.Create(wechatClient, menu_val)
    if nil != err {
        fmt.Printf("create menu err: %v \n", err)
    }
}