package myhttp

import (
	"bufio"
	"errors"
	"io"
	"log"
	"strconv"
	"strings"
)

//好像没有bug了, body正确, url正确
//好吧要根据Content-Length 读取指定字节的数据
//难过了
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
			request.UrlParsed, err = ParseUrl(request.Url)
			if err != nil {
				return request, err
			}
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

	Length, ok := request.Headers["Content-Length"]
	if !ok {
		return request, nil
	}
	length, _ := strconv.Atoi(Length[0])

	var err error
	body := make([]byte, length)
	_, err = reader.Read(body)
	if err != nil && err != io.EOF {
		return request, err
	}
	request.Body = body[:]

	return request, nil
}

func ParseUrl(url string) (Url, error) {
	U := new(Url)
	Sli := strings.Split(url, "?")
	U.Router = Sli[0]
	if len(Sli) == 1 {
		return *U, nil
	}
	querys := strings.Split(Sli[1], "&")
	for _, one := range querys {
		kv := strings.Split(one, "=")
		if len(kv) != 2 {
			return *U, errors.New("query error")
		}
		k := kv[0]
		v := kv[1]
		U.Query[k] = append(U.Query[k], v)
	}

	return *U, nil
}

func ParseHttpResponse(r io.Reader) (Response, error) {
	reader := bufio.NewReader(r)
	response := NewResponse()
	var inhead bool = true
	var iter int = 1 //计数器
	for {            //头部
		data, _, err := reader.ReadLine()
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			} else {
				return response, err
			}
		}
		data_s := string(data)
		if iter == 1 {
			RLine := strings.Split(data_s, " ")
			if len(RLine) != 3 {
				return response, errors.New("not http")
			}

			response.Proto = RLine[0]
			response.StatusCode, _ = strconv.Atoi(RLine[1])
			response.Phrase = RLine[2]
			iter++
			continue
		}

		if data_s == "" {
			inhead = false
		}
		if inhead {
			HeaderLine := strings.Split(data_s, ": ")
			if len(HeaderLine) != 2 {
				return response, errors.New("HeaderLine error")
			}
			HeaderDomain := HeaderLine[0]
			HeadValues := HeaderLine[1]
			HeadValue := strings.Split(HeadValues, ", ")
			response.Headers[HeaderDomain] = append(response.Headers[HeaderDomain], HeadValue...)

		} else {
			// fmt.Println("Headers ends")
			break
		}
	}

	Length, ok := response.Headers["Content-Length"]
	if !ok {
		return response, nil
	}
	length, _ := strconv.Atoi(Length[0])

	var err error
	body := make([]byte, length)
	_, err = reader.Read(body)
	if err != nil && err != io.EOF {
		return response, err
	}
	response.Body = body[:]

	return response, nil
}
