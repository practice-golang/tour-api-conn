package main // import "tour-api-conn"

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
)

/*
OpenAPI : 관광지 정보 조회 API 구조체

	Uri : 요청 URL
	AuthKey : 인증키. https://www.data.go.kr
	Parameters : 요청 파라미터
*/
type OpenAPI struct {
	Uri        *url.URL
	AuthKey    string
	Parameters interface{}
}

/*
ParamsList : 관광지 목록 조회 파라미터
*/
type ParamsList struct {
	MapX      float64 `json:"mapX"`
	MapY      float64 `json:"mapY"`
	Radius    float64 `json:"radius"`
	ListYN    string  `json:"listYN"`
	Arrange   string  `json:"arrange"`
	NumOfRows int     `json:"numOfRows"`
	PageNo    int     `json:"pageNo"`

	MobileOS  string `json:"MobileOS"`
	MobileApp string `json:"MobileApp"`
	DataType  string `json:"_type"`
}

/*
ParamsInfo : 관광지 정보 조회 파라미터
*/
type ParamsInfo struct {
	ContentTypeID string `json:"contentTypeId"`
	ContentID     string `json:"contentId"`
	DefaultYN     string `json:"defaultYN"`
	FirstImageYN  string `json:"firstImageYN"`
	AreaCodeYN    string `json:"areacodeYN"`
	CatCodeYN     string `json:"catcodeYN"`
	AddrInfoYN    string `json:"addrinfoYN"`
	MapInfoYN     string `json:"mapinfoYN"`
	OverViewYN    string `json:"overviewYN"`

	MobileOS  string `json:"MobileOS"`
	MobileApp string `json:"MobileApp"`
	DataType  string `json:"_type"`
}

/*
initParams : 관광지 정보 조회 파라미터 초기화
*/
func (TourAPI *OpenAPI) initParams() {
	var params interface{}

	switch TourAPI.Parameters.(type) {
	case ParamsList:
		params = ParamsList{
			ListYN:    "Y",
			MobileOS:  "ETC",
			MobileApp: "TourAPI3.0_Guide",
			Arrange:   "A",
		}

	case ParamsInfo:
		params = ParamsInfo{
			MobileOS:     "ETC",
			MobileApp:    "TourAPI3.0_Guide",
			DefaultYN:    "Y",
			FirstImageYN: "Y",
			AreaCodeYN:   "Y",
			CatCodeYN:    "Y",
			AddrInfoYN:   "Y",
			MapInfoYN:    "Y",
			OverViewYN:   "Y",
		}
	}

	TourAPI.Parameters = params
}

/*
setParams : 요청 파라미터 설정
*/
func (TourAPI *OpenAPI) setParams() {
	query := TourAPI.Uri.Query()
	query.Add("ServiceKey", TourAPI.AuthKey)

	v := reflect.ValueOf(TourAPI.Parameters)
	for i := 0; i < v.NumField(); i++ {
		query.Add(v.Type().Field(i).Tag.Get("json"), fmt.Sprint(v.Field(i).Interface()))
	}

	TourAPI.Uri.RawQuery = query.Encode()
}

/*
requestTourData : 관광지 요청 및 데이터 취득
*/
func (TourAPI *OpenAPI) requestTourData() (string, error) {
	var err error

	// GET 호출
	resp, err := http.Get(TourAPI.Uri.String())
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// 결과 취득
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

/*
GetSightList : 관광지 정보 조회

	x : 위도. latitude
	y : 경도. longitude
	r : 반경. radius
	format : 응답 형식. json, xml
	authKey : 인증키. https://www.data.go.kr
*/
func GetSightList(x, y, r float64, format, authKey string) string {
	uri, err := url.Parse(`https://api.visitkorea.or.kr/openapi/service/rest/KorService/locationBasedList`)
	if err != nil {
		panic(err)
	}

	tourAPI := OpenAPI{Uri: uri}
	tourAPI.AuthKey = authKey

	tourAPI.Parameters = ParamsList{}
	tourAPI.initParams()

	params := tourAPI.Parameters.(ParamsList)
	params.MapX = x
	params.MapY = y
	params.Radius = r
	params.DataType = format
	params.NumOfRows = 10
	params.PageNo = 1

	tourAPI.Parameters = params
	tourAPI.setParams()

	data, err := tourAPI.requestTourData()
	if err != nil {
		panic(err)
	}

	return data
}

/*
GetSightInfo : 관광지 정보 조회

	id : 콘텐츠 ID
	typeid : 콘텐츠 타입 ID
	format : 응답 형식. json, xml
	authKey : 인증키. https://www.data.go.kr
*/
func GetSightInfo(id, typeid, format, authKey string) string {
	uri, err := url.Parse(`https://api.visitkorea.or.kr/openapi/service/rest/KorService/detailCommon`)
	if err != nil {
		panic(err)
	}

	tourAPI := OpenAPI{Uri: uri}
	tourAPI.AuthKey = authKey

	tourAPI.Parameters = ParamsInfo{}
	tourAPI.initParams()

	params := tourAPI.Parameters.(ParamsInfo)
	params.ContentTypeID = id
	params.ContentID = typeid
	params.DataType = format

	tourAPI.Parameters = params
	tourAPI.setParams()

	data, err := tourAPI.requestTourData()
	if err != nil {
		panic(err)
	}

	return data
}

func main() {
	// 인증키 - https://www.data.go.kr 에서 제공받은 인증키. 일반 인증키(Decoding) 사용
	authKey := ""
	format := "json"

	if len(os.Args) > 1 {
		authKey = os.Args[1]
		if len(os.Args) == 3 {
			format = os.Args[2]
		}
	} else {
		fmt.Println("Usage : go run main.go [key] [format: json, xml. default: json]")
		os.Exit(0)
	}

	fmt.Println(GetSightList(126.977969, 37.566535, 2000, format, authKey))
	fmt.Println(GetSightInfo("14", "129898", format, authKey))
}
