/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/11/5 15:20 
 */
package main

import (
    "fmt"
    "pet/utils"
)

func main() {
    s := []string{"18505921256", "13489594009", "12759029321", "18501790879", "", "123", "abc"}
    for _, v := range s {
        fmt.Println(utils.PhoneValid(v))
    }
}
