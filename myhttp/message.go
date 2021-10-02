package myhttp

import (
	"bufio"
	"fmt"
	"io"
)

func SendHttpRequest(w io.Writer, rq Request) error {
	buf := bufio.NewWriter(w)

	_, err := buf.WriteString(rq.Method + " " + rq.Url + " " + rq.Proto + "\r\n") // 请求行

	if err != nil {
		return err
	}

	for k, v := range rq.Headers {
		_, err = buf.WriteString(k + ": ")
		if err != nil {
			return err
		}

		for i, val := range v {
			_, err = buf.WriteString(val)
			if err != nil {
				return err
			}
			if i != len(v)-1 {
				_, err = buf.WriteString(", ")
				if err != nil {
					return err
				}
			}
		}

		_, err = buf.WriteString("\r\n")
		if err != nil {
			return err
		}

	} //头部行

	_, err = buf.WriteString("\r\n")

	if err != nil {
		return err
	}

	_, err = buf.Write(rq.Body)
	if err != nil {
		return err
	}

	err = buf.Flush()
	if err != nil {
		return err
	}
	buf.Write([]byte(io.EOF.Error()))
	return nil
}

func SendHttpResponse(w io.Writer, rep Response) error {
	buf := bufio.NewWriter(w)

	_, err := buf.WriteString(rep.Proto + " " + fmt.Sprint(rep.StatusCode) + " " + rep.Phrase + "\r\n") //响应行
	if err != nil {
		return err
	}

	for k, v := range rep.Headers {
		_, err = buf.WriteString(k + ": ")
		if err != nil {
			return err
		}

		for i := 0; i < len(v); i++ {
			_, err = buf.WriteString(v[0])
			if err != nil {
				return err
			}

			if i != len(v)-1 {
				_, err = buf.WriteString(", ")
				if err != nil {
					return err
				}
			}
		}
		_, err = buf.WriteString("\r\n")
		if err != nil {
			return err
		}
	} //头部行

	_, err = buf.WriteString("\r\n")

	if err != nil {
		return err
	}

	_, err = buf.Write(rep.Body)
	if err != nil {
		return err
	}

	err = buf.Flush()
	if err != nil {
		return err
	}
	return nil
}
