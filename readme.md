## 概述

该模块用于获取系统信息

## 调用方式

通过go get获取

```
go get github.com/xinyulu3344/os_info
```

## 导包

```go
import "github.com/xinyulu3344/os_info"
```

## 方法说明

```
获取内核版本
func (o *OsInfo) GetKernelVersion() string

获取Linux发行版，类似于Python中platform模块的linux_distribution()方法
func (o *OsInfo) GetLinuxDistribution() []string
```

## 调用示例

```
import "github.com/xinyulu3344/os_info"

osInfo := os_info.NewOsInfo()
fmt.Println(osInfo.GetLinuxDistribution())
```