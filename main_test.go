/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/11/5 22:49 
 */
package main

import (
    "testing"
    "fmt"
    "pet/utils"
)

func TestPhoneValid(t *testing.T) {
    s := []string{"18505921256", "13489594009", "12759029321", "17358547087"}
    for _, v := range s {
        fmt.Println(utils.PhoneValid(v))
    }
}