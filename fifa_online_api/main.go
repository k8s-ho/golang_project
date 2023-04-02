package main

import (
	"encoding/json"
	"fmt"
	"github.com/thedevsaddam/gojsonq"
	"io/ioutil"
	"net/http"
	"strings"
)

type INFO struct {
	accessId string
	nickname string
	level    int
}

// 원하는 경기 데이터가 있다면 여기에 추가 
const(
	cns_nick = "nickname"
	cns_goal = "shoot.goalTotal"
	cns_result = "matchDetail.matchResult"
)

var user_nickname string

// access token file open하여 사용하는 함수
func get_token() string {
	file, err := ioutil.ReadFile("token")
	if err != nil {
		fmt.Println("Error reading token file:", err)
	}

	// file open 시, 줄바꿈이 들어가 payload에 문제발생 -> 제거 진행코드
	token := string(file)
	token = strings.Replace(token, "\n", "", -1)
	return token
}

// api request 및 response 받는 함수 
func get_request(url string) string { // string user id return
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

// nickname을 통해 user access id 가져오는 함수 
func get_user_info(nickname string) string {
	api_url := "https://api.nexon.co.kr/fifaonline4/v1.0/users/?nickname=" + nickname //nickname을 통한 user 정보 api
	user_info := get_request(api_url)                                                 //user_info 리턴받아서 바로 데이터받아오기

	access_id := gojsonq.New().FromString(user_info).Find("accessId")
	if access_id == nil {
		panic("[-] No such nickname exists. check your nickname :(")
	}
	return access_id.(string)
}

// match data 출력 함수 
func show_data(myidx int, runman interface{}, match_data string){
	var my_nickname, enemy_nickname, match_result, my_goal, enemy_goal interface{}
	my_nickname = gojsonq.New().FromString(match_data).Find(fmt.Sprintf("matchInfo.[%d].%s", myidx, cns_nick))	
	my_goal = gojsonq.New().FromString(match_data).Find(fmt.Sprintf("matchInfo.[%d].%s", myidx, cns_goal))
	match_result = gojsonq.New().FromString(match_data).Find(fmt.Sprintf("matchInfo.[%d].%s", myidx, cns_result))

	if myidx == 0 {
		if runman == nil{
			fmt.Println("[!] 상대방이 탈주하였습니다") // 탈주자 유무 확인
			enemy_nickname = "탈주자"
			enemy_goal = 0
		}else {
			enemy_nickname = gojsonq.New().FromString(match_data).Find(fmt.Sprintf("matchInfo.[%d].%s", myidx+1, cns_nick))
			enemy_goal = gojsonq.New().FromString(match_data).Find(fmt.Sprintf("matchInfo.[%d].%s", myidx+1, cns_goal))
		}	
	}else if myidx == 1 {
		enemy_nickname = gojsonq.New().FromString(match_data).Find(fmt.Sprintf("matchInfo.[%d].%s", myidx-1, cns_nick))
		enemy_goal = gojsonq.New().FromString(match_data).Find(fmt.Sprintf("matchInfo.[%d].%s", myidx-1, cns_goal))
	}

	fmt.Printf("[경기] %v(나) vs %v(상대)\n", my_nickname, enemy_nickname)
	fmt.Printf("[경기결과: '%v'] %v:%v\n", match_result, my_goal, enemy_goal)
	fmt.Println()
}

// match data를 찾아오는 함수
func get_match_data(match_id string, access_id string) {
	// match 정보 api
	api_url := "https://api.nexon.co.kr/fifaonline4/v1.0/matches/"

	// match_id -> array로 변환, 다수의 match_id를 array에 저장
	var match_arr []string
	err := json.Unmarshal([]byte(match_id), &match_arr)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// match data 가져옴
	match_cnt := len(match_arr) // 최근 경기수 최대 100건
	fmt.Printf("[ '%v'님의 최근 공식 경기 수: %v]\n", user_nickname, match_cnt)

	// 최근 경기 데이터 전부 가져옴
	for n := 0; n < match_cnt; n++{
		match_data := get_request(api_url + match_arr[n])
		match_date := gojsonq.New().FromString(match_data).Find("matchDate")
		fmt.Printf("[게임시간] %v\n", match_date)

		// 탈주자 여부 확인(matchInfo가 1개이면 탈주자 존재한다는 의미)
		chk := gojsonq.New().FromString(match_data).Find("matchInfo.[1].accessId")// 두번째 배열 존재 여부 확인으로 탈주자 검사 
		
		// 나의 matchinfo idx 확인 
		if gojsonq.New().FromString(match_data).Find("matchInfo.[0].accessId") == access_id {
			show_data(0,chk, match_data) // 나의 idx가 0번 일경우를 기준으로 match data 출력 
		}else {
			show_data(1,chk, match_data) // 나의 idx가 1번 일경우를 기준으로 match data 출력  
		}
	}

}

func get_match_id(id string) {
	// 리그 친선: 30  클래식: 40  공식: 50  감독: 52  공식친선: 60
	api_url := "https://api.nexon.co.kr/fifaonline4/v1.0/users/" + id + "/matches?matchtype=40" // get match id
	match_id := get_request(api_url)
	get_match_data(match_id, id)
}

func main() {
	fmt.Println("[     피파온라인4 경기이력 검색기     ]")
	fmt.Println("Data based on NEXON DEVELOPERS")
	fmt.Print("[+] 닉네임을 입력하세요 = ")
	fmt.Scan(&user_nickname)
	access_id := get_user_info(user_nickname)
	get_match_id(access_id)
}
