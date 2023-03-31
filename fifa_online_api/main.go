package main
 
import (
    "net/http"
    "fmt"
    "io/ioutil"
    "strings"
    "github.com/thedevsaddam/gojsonq"
    "encoding/json"
)

type INFO struct{
    accessId string
    nickname string
    level int
}

func get_token()(string){
    file, err := ioutil.ReadFile("token")
    if err != nil {
    	fmt.Println("Error reading token file:", err)
    }

    // file open 시, 줄바꿈이 들어가 payload에 문제발생 -> 제거 진행코드
    token := string(file)
    token = strings.Replace(token, "\n", "", -1)
    return token
}

func get_request(url string)(string){ // string user id return
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
	panic(err)
    }

    token := get_token()

    // Header에 Authorization token 추가
    req.Header.Set("Authorization", token)

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
	panic(err)
    }
    defer resp.Body.Close()
    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
    result := string(data)
    return result
}


func get_user_info(nickname string)(string){
    api_url := "https://api.nexon.co.kr/fifaonline4/v1.0/users/?nickname=" + nickname  //nickname을 통한 user 정보 api
    user_info := get_request(api_url)  //user_info 리턴받아서 바로 데이터받아오기

    access_id := gojsonq.New().FromString(user_info).Find("accessId")
    if access_id == nil {
	panic("[-] No such nickname exists. check your nickname :(")
    }
    return access_id.(string)
}

// match data를 찾아오는 함수
func get_match_data(match_id string){
    // match 정보 api
    api_url := "https://api.nexon.co.kr/fifaonline4/v1.0/matches/"
	
    // match_id -> array로 변환
    var match_arr []string
    err := json.Unmarshal([]byte(match_id), &match_arr)
    if err != nil {
        fmt.Println("Error:", err)
    }

    //fmt.Println(match_arr[1])
    // match data 가져옴
    fmt.Println(get_request(api_url + match_arr[1]))
}

func get_match_id(id string){
    // 리그 친선: 30  클래식: 40  공식: 50  감독: 52  공식친선: 60
    api_url := "https://api.nexon.co.kr/fifaonline4/v1.0/users/" + id + "/matches?matchtype=50"  // get match id 
    match_id := get_request(api_url)
    get_match_data(match_id)
}


func main() {
    var user_nickname string
    fmt.Println("[     피파온라인4 경기이력 검색기     ]")
    fmt.Print("[+] 닉네임을 입력하세요 = ")
    fmt.Scan(&user_nickname)
    access_id := get_user_info(user_nickname)
    get_match_id(access_id)
}
