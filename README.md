# aemail
这是一个极简的基于github.com/go-gomail/gomail的go-email

## 使用简介

```go
package main

import (
	"github.com/1LOS1/aemail"
	"fmt"
)

func main() {
	center := aemail.NewEmailCenter(&aemail.EmailConfig{
		ServerHost:   "smtp.163.com", //邮箱服务器
		ServerPort:   25,             //邮箱端口
		FromEmail:    "123@163.com",  // 发送邮箱
		FromPassword: "AAAAAAAAA",    //授权码
	})

	err := center.
		AddToers("123@qq.com", "123@hotmail.com", "12345@qq.com").
		SetMessage("los", "测试邮箱", "dasuki\nform http://lightos.cloud").
		Send()

	if err != nil {
		fmt.Println(err)
		return
	}
}

```

