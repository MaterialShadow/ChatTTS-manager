/*
* @desc:缓存相关
* @company:云南奇讯科技有限公司
* @Author: yixiaohu
* @Date:   2022/3/9 11:25
 */

package consts

const (
	// CachePrefix 应用缓存数据前缀
	CachePrefix = "APP:"

	CacheModelMem   = "memory"
	CacheModelRedis = "redis"
	CacheModelDist  = "dist"

	// CacheSysDict 字典缓存菜单KEY
	CacheSysDict = CachePrefix + "sysDict"

	// CacheSysDictTag 字典缓存标签
	CacheSysDictTag = CachePrefix + "sysDictTag"
	// CacheSysConfigTag 系统参数配置
	CacheSysConfigTag = CachePrefix + "sysConfigTag"
)
