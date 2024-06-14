/*
* @desc:缓存键
* @company:云南奇讯科技有限公司
* @Author: yixiaohu
* @Date:   2022/3/9 12:06
 */

package consts

import commonConsts "github.com/tiger1103/gfast/v3/internal/app/common/consts"

const (
	// CacheSysAuthMenu 缓存菜单key
	CacheSysAuthMenu = commonConsts.CachePrefix + "sysAuthMenu"
	// CacheSysDept 缓存部门key
	CacheSysDept = commonConsts.CachePrefix + "sysDept"

	// CacheSysRole 角色缓存key
	CacheSysRole = commonConsts.CachePrefix + "sysRole"
	// CacheSysWebSet 站点配置缓存key
	CacheSysWebSet = commonConsts.CachePrefix + "sysWebSet"
	// CacheSysCmsMenu cms缓存key
	CacheSysCmsMenu = commonConsts.CachePrefix + "sysCmsMenu"

	// CacheSysAuthTag 权限缓存TAG标签
	CacheSysAuthTag = commonConsts.CachePrefix + "sysAuthTag"
	// CacheSysModelTag 模型缓存标签
	CacheSysModelTag = commonConsts.CachePrefix + "sysModelTag"
	// CacheSysCmsTag cms缓存标签
	CacheSysCmsTag = commonConsts.CachePrefix + "sysCmsTag"
)
