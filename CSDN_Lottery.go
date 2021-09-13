package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

//返回的结构体
type Result struct {
	Code int64 `json:"code"`
	Data struct {
		Count      int64 `json:"count"`
		FloorCount int64 `json:"floorCount"`
		List       []struct {
			Info struct {
				ArticleID             int64  `json:"articleId"`
				Avatar                string `json:"avatar"`
				CommentFromTypeResult struct {
					Index int64  `json:"index"`
					Key   string `json:"key"`
					Title string `json:"title"`
				} `json:"commentFromTypeResult"`
				CommentID       int64         `json:"commentId"`
				CompanyBlog     interface{}   `json:"companyBlog"`
				CompanyBlogIcon interface{}   `json:"companyBlogIcon"`
				Content         string        `json:"content"`
				DateFormat      string        `json:"dateFormat"`
				Digg            int64         `json:"digg"`
				DiggArr         []interface{} `json:"diggArr"`
				Flag            interface{}   `json:"flag"`
				FlagIcon        interface{}   `json:"flagIcon"`
				LevelIcon       interface{}   `json:"levelIcon"`
				NickName        string        `json:"nickName"`
				ParentID        int64         `json:"parentId"`
				ParentNickName  interface{}   `json:"parentNickName"`
				ParentUserName  interface{}   `json:"parentUserName"`
				PostTime        string        `json:"postTime"`
				UserName        string        `json:"userName"`
				Vip             interface{}   `json:"vip"`
				VipIcon         interface{}   `json:"vipIcon"`
				Years           interface{}   `json:"years"`
			} `json:"info"`
			PointCommentID interface{}   `json:"pointCommentId"`
			Sub            []interface{} `json:"sub"`
		} `json:"list"`
		PageCount int64 `json:"pageCount"`
	} `json:"data"`
	Message string `json:"message"`
}

//结构体
type LuckyPerson struct {  // 存储这个选手的信息
	UserName string
	NickName string
}

//主函数
func main() {
	total := getTotalNum()  // 获取全部的评论数目
	var totalPerson []LuckyPerson
	for page := 1; page < total/10+1; page++ {
		for _, person := range Spider(page) {  // 根据目录进行爬取
			totalPerson = append(totalPerson, person)
		}
	}
	totalPersonSet := removeRepByMap(totalPerson)    // 去重
	fmt.Println("以下名单进行抽奖：")
	for _, item := range totalPersonSet {
		fmt.Println(item.NickName)
	}
	fmt.Printf("本次抽奖共有 %d 人参与! \n准备抽奖~\n",len(totalPersonSet))
	for i := 3; i >= 1; i-- {
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}
	fmt.Println("恭喜这位选手中奖~")
	luckDog := lottery(totalPersonSet)
	fmt.Printf("中奖是 %s \n", luckDog.NickName)
}

func getTotalNum() int {  // 获取全部的一级评论数目
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://blog.csdn.net/phoenix/web/v1/comment/list/120067923?page=1&size=10", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("content-length", "0")
	req.Header.Set("accept", "*/*")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var result Result
	_ = json.Unmarshal(bodyText, &result) //byte to json
	total := result.Data.Count
	return int(total)
}

func Spider(page int) []LuckyPerson {  // 传入页数，一页一页爬取
	var tmp []LuckyPerson
	p :=strconv.Itoa(page)
	client := &http.Client{}
	reqSpider, err := http.NewRequest("POST", "https://blog.csdn.net/phoenix/web/v1/comment/list/120067923?page="+p+"&size=10", nil)
	if err != nil {
		log.Fatal(err)
	}
	reqSpider.Header.Set("content-length", "0")
	reqSpider.Header.Set("accept", "*/*")
	reqSpider.Header.Set("x-requested-with", "XMLHttpRequest")
	respSpider, err := client.Do(reqSpider)
	if err != nil {
		log.Fatal(err)
	}
	defer respSpider.Body.Close()
	bodyText, err := ioutil.ReadAll(respSpider.Body)
	if err != nil {
		log.Fatal(err)
	}
	var result Result
	_ = json.Unmarshal(bodyText, &result)               //byte to json
	num := len(result.Data.List)
	commentList := result.Data.List
	for i:=0 ;i<num; i++ {
		var luckPerson LuckyPerson
		luckPerson.UserName = commentList[i].Info.UserName
		luckPerson.NickName = commentList[i].Info.NickName
		tmp = append(tmp, luckPerson)   // 存放每一个评论的人
	}
	return tmp
}

func removeRepByMap(slc []LuckyPerson) []LuckyPerson {  //去除重复的元素
	var result []LuckyPerson           //存放返回的不重复切片
	tempMap := map[LuckyPerson]byte{}  // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0 					//当e存在于tempMap中时，再次添加是添加不进去的，因为key不允许重复
		if len(tempMap) != l { 			// 加入map后，map长度变化，则元素不重复
			result = append(result, e)  //当元素不重复时，将元素添加到切片result中
		}
	}
	return result
}

func lottery(totalPerson []LuckyPerson) LuckyPerson {  // 抽取中奖选手
	rand.Seed(time.Now().UnixNano())    // 使用随机种子
	index:=rand.Intn(len(totalPerson))  // 生成0到这个列表的长度的一个数字
	return totalPerson[index]  			// 返回中奖选手
}
