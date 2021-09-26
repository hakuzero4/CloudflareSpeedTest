#XIU2/CloudflareSpeedTest 派生项目

功能保持同步更新,新增设置dnspod境内相对最快的cdn节点

## 使用说明

需要在执行文件目录新增config.yaml文件

配置如下:

```yaml
dnspod:
  id: xxx
  token: xxx
  domain: 99999.xyz
  record: 
    av: 19999
    ghs: 22999
```

参数说明:

| 参数      | 说明             |
| --------- | ---------------- |
| id        | dnspod api id    |
| token     | dnspod api token |
| domain    | 需要设置的域名   |
| record | 需要修改的记录(key:主机记录)(value:记录id) |



dnspod token 申请地址: https://console.dnspod.cn/account/token



dnspod 记录ID获取方式

`./CloudflareSpeedTest -dlist`

因无法找到ip isp api顾不提供线路区分,默认为境内

### changelog

新增-c参数,指定config.ymal配置文件路径

