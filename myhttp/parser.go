package myhttp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

func ParseHttpRequest(conn net.Conn) (Request, error) {
	reader := bufio.NewReader(conn)
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
			fmt.Println("len = ", len(RLine))

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
			fmt.Println("Headers ends")
			break
		}
	}
	var err error
	request.Body, err = ioutil.ReadAll(reader)
	if err != nil {
		return request, err
	}
	fmt.Printf("%#v\n", string(request.Body))

	// fmt.Printf("%v", data)
	fmt.Println("end")

	return request, nil
}
