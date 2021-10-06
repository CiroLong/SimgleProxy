package logger

import (
	"SIMGLEPROXY/myhttp"
	"fmt"
	"os"
	"time"
)

type Logger struct {
	Fp *os.File
}

func (l *Logger) PrintAccess(res *myhttp.Request, rep *myhttp.Response) {
	//它至少需要包括时间、用户http请求起始行信息、请求返
	//回的状态码、User-Agent。
	l.Fp.WriteString(time.Now().String())
	l.Fp.WriteString("\n " + res.Method + " " + res.Url + " ")
	l.Fp.WriteString("\n code:" + fmt.Sprint(rep.StatusCode))
	user_agent := res.Headers["User-Agent"]
	l.Fp.WriteString("\n User-Agent:" + user_agent[0])
	l.Fp.WriteString("\n")
}

func (l *Logger) PrintError(err error) {
	//时间、错误级别、错误信息。
	l.Fp.WriteString(time.Now().String())
	l.Fp.WriteString("\nerror: " + err.Error() + "\n")
}
