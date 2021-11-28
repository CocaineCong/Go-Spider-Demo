package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

//https://blink-open-api.csdn.net/v1/pc/blink/allComment?pageNum=3&pageSize=100&blinkId=1260435


type BlinkResult struct {
	Code int64 `json:"code"`
	Data []struct {
		Anonymous bool   `json:"anonymous"`
		Avatar    string `json:"avatar"`
		BizNo     string `json:"bizNo"`
		Child     []struct {
			Anonymous  bool   `json:"anonymous"`
			Avatar     string `json:"avatar"`
			BizNo      string `json:"bizNo"`
			ChildCount int64  `json:"childCount"`
			Content    string `json:"content"`
			CreateTime string `json:"createTime"`
			Employee   struct {
				Bala            string `json:"bala"`
				IsCertification bool   `json:"isCertification"`
				Org             string `json:"org"`
			} `json:"employee"`
			Ext        string `json:"ext"`
			FlowerName struct {
				FlowerName string `json:"flowerName"`
				Level      string `json:"level"`
			} `json:"flowerName"`
			ID              int64  `json:"id"`
			Level           int64  `json:"level"`
			LikeCount       int64  `json:"likeCount"`
			Nickname        string `json:"nickname"`
			ParentID        int64  `json:"parentId"`
			Platform        string `json:"platform"`
			ReplyFlowerName struct {
				FlowerName string `json:"flowerName"`
				Level      string `json:"level"`
			} `json:"replyFlowerName"`
			ReplyNickname string `json:"replyNickname"`
			ReplyUsername string `json:"replyUsername"`
			ResourceGroup string `json:"resourceGroup"`
			ResourceID    string `json:"resourceId"`
			ResourceOrder string `json:"resourceOrder"`
			ResourceUser  string `json:"resourceUser"`
			Score         int64  `json:"score"`
			Status        int64  `json:"status"`
			Student       struct {
				Bala            string `json:"bala"`
				IsCertification bool   `json:"isCertification"`
				Org             string `json:"org"`
			} `json:"student"`
			Title struct {
				ID       int64  `json:"id"`
				TitleURL string `json:"titleUrl"`
				Used     bool   `json:"used"`
				Username string `json:"username"`
			} `json:"title"`
			Top      bool   `json:"top"`
			UserLike bool   `json:"userLike"`
			Username string `json:"username"`
		} `json:"child"`
		ChildCount int64  `json:"childCount"`
		Content    string `json:"content"`
		CreateTime string `json:"createTime"`
		Employee   struct {
			Bala            string `json:"bala"`
			IsCertification bool   `json:"isCertification"`
			Org             string `json:"org"`
		} `json:"employee"`
		Ext        string `json:"ext"`
		FlowerName struct {
			FlowerName string `json:"flowerName"`
			Level      string `json:"level"`
		} `json:"flowerName"`
		FromType      string `json:"fromType"`
		ID            int64  `json:"id"`
		Level         int64  `json:"level"`
		LikeCount     int64  `json:"likeCount"`
		Nickname      string `json:"nickname"`
		ParentID      int64  `json:"parentId"`
		Platform      string `json:"platform"`
		ReplyUsername string `json:"replyUsername"`
		ResourceGroup string `json:"resourceGroup"`
		ResourceID    string `json:"resourceId"`
		ResourceOrder string `json:"resourceOrder"`
		ResourceUser  string `json:"resourceUser"`
		Score         int64  `json:"score"`
		Status        int64  `json:"status"`
		Student       struct {
			Bala            string `json:"bala"`
			IsCertification bool   `json:"isCertification"`
			Org             string `json:"org"`
		} `json:"student"`
		Title struct {
			ID       int64  `json:"id"`
			TitleURL string `json:"titleUrl"`
			Used     bool   `json:"used"`
			Username string `json:"username"`
		} `json:"title"`
		Top      bool   `json:"top"`
		UserLike bool   `json:"userLike"`
		Username string `json:"username"`
	} `json:"data"`
	Msg string `json:"msg"`
}

//结构体
type LuckyBlinkPerson struct {  // 存储这个选手的信息
	UserName string
	NickName string
	CreateTime string
	Content string
}



func main() {
	client := &http.Client{}
	reqSpider, err := http.NewRequest("GET", "https://blink-open-api.csdn.net/v1/pc/blink/allComment?pageNum=1&pageSize=600&blinkId=1260435", nil)
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
	bodyText, _ := ioutil.ReadAll(respSpider.Body)
	var result BlinkResult
	_ = json.Unmarshal(bodyText, &result)
	commentList := result.Data
	var luckyBPList []LuckyBlinkPerson
	for _,v := range commentList{
		var luckBlinkPerson LuckyBlinkPerson
		luckBlinkPerson.UserName=v.Username
		luckBlinkPerson.NickName=v.Nickname
		luckBlinkPerson.CreateTime=v.CreateTime
		luckBlinkPerson.Content=v.Content
		luckyBPList = append(luckyBPList, luckBlinkPerson)
	}
	luckyBPList = removeBlinkRepByMap(luckyBPList)
	fmt.Println("以下名单进行抽奖：")
	for _, item := range luckyBPList {
		fmt.Println(item.NickName)
	}
	fmt.Printf("本次抽奖共有 %d 人参与! \n准备抽奖~\n",len(luckyBPList))
	for i := 3; i >= 1; i-- {
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}
	fmt.Println("恭喜这位选手中奖~")
	luckDog := lotteryBlink(luckyBPList)
	fmt.Printf("Blink中奖昵称是：%s 用户名是：%s 评论时间：%s 评论内容：%s\n", luckDog.NickName,luckDog.UserName,luckDog.CreateTime,luckDog.Content)
}

func lotteryBlink(totalPerson []LuckyBlinkPerson) LuckyBlinkPerson {  // 抽取中奖选手
	rand.Seed(time.Now().UnixNano())    // 使用随机种子
	index:=rand.Intn(len(totalPerson))  // 生成0到这个列表的长度的一个数字
	return totalPerson[index]  			// 返回中奖选手
}


func removeBlinkRepByMap(slc []LuckyBlinkPerson) []LuckyBlinkPerson {  //去除重复的元素
	var result []LuckyBlinkPerson           //存放返回的不重复切片
	tempMap := map[LuckyBlinkPerson]byte{}  // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0 					//当e存在于tempMap中时，再次添加是添加不进去的，因为key不允许重复
		if len(tempMap) != l { 			// 加入map后，map长度变化，则元素不重复
			result = append(result, e)  //当元素不重复时，将元素添加到切片result中
		}
	}
	return result
}