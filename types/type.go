// 애플리케이션에서 공통으로 사용하는 타입과 그 타입을 생성하는 함수들을 정의

package types

import "strings"

const (
	SocketBufferSize  = 1024
	MessageBufferSize = 256
)

// JSON 응답의 공통 헤더 데이터를 포함하는 구조체
type header struct {
	Result int    `json:"result"`
	Data   string `json:"data"`
}

// header 구조체 인스턴스를 생성
func newHeader(result int, data ...string) *header {
	return &header{
		Result: result,
		Data:   strings.Join(data, ","),
	}
}

// header와 실제 응답 결과를 포함하는 구조체
type response struct {
	*header
	Result interface{} `json:"result"`
}

// response 구조체 인스턴스를 생성
func NewRes(result int, res interface{}, data ...string) *response {
	return &response{
		header: newHeader(result, data...),
		Result: res,
	}
}
