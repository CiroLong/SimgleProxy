package myhttp

import (
	"bufio"
	"errors"
	"io"
	"log"
	"strings"
)

//好像没有bug了, body正确, url正确
func ParseHttpRequest(reader *bufio.Reader) (Request, error) {
	request := newRequest()
	var inhead bool = true
	var iter int = 1 //计数器
	for {            //头部
		data, _, err := reader.ReadLine()
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			} else {
				return request, err
			}
		}
		data_s := string(data)
		if iter == 1 {
			RLine := strings.Split(data_s, " ")
			if len(RLine) != 3 {
				return request, errors.New("not http")
			}

			request.Method = RLine[0]
			request.Url = RLine[1]
			request.Proto = RLine[2]
			iter++
			continue
		}

		if data_s == "" {
			inhead = false
		}
		if inhead {
			HeaderLine := strings.Split(data_s, ": ")
			if len(HeaderLine) != 2 {
				return request, errors.New("HeaderLine error")
			}
			HeaderDomain := HeaderLine[0]
			HeadValues := HeaderLine[1]
			HeadValue := strings.Split(HeadValues, ", ")
			request.Headers[HeaderDomain] = append(request.Headers[HeaderDomain], HeadValue...)

		} else {
			// fmt.Println("Headers ends")
			break
		}
	}
	var err error
	_, err = reader.Read(request.Body)
	if err != nil {
		return request, err
	}

	return request, nil
}

func ParseUrl(url string) {

}
