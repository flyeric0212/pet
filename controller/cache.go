/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/11/5 17:28 
 */
package controller

import "pet/utils"

var g_cache *utils.Cache

func InitCachePool(redis_conf *utils.RedisConfig) (err error) {
    g_cache, err = utils.InitRedisPool(redis_conf)
    return err
}
