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

//?????????
type LuckyBlinkPerson struct {  // ???????????????????????????
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
	fmt.Println("???????????????????????????")
	for _, item := range luckyBPList {
		fmt.Println(item.NickName)
	}
	fmt.Printf("?????????????????? %d ?????????! \n????????????~\n",len(luckyBPList))
	for i := 3; i >= 1; i-- {
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}
	fmt.Println("????????????????????????~")
	luckDog := lotteryBlink(luckyBPList)
	fmt.Printf("Blink??????????????????%s ???????????????%s ???????????????%s ???????????????%s\n", luckDog.NickName,luckDog.UserName,luckDog.CreateTime,luckDog.Content)
}

func lotteryBlink(totalPerson []LuckyBlinkPerson) LuckyBlinkPerson {  // ??????????????????
	rand.Seed(time.Now().UnixNano())    // ??????????????????
	index:=rand.Intn(len(totalPerson))  // ??????0???????????????????????????????????????
	return totalPerson[index]  			// ??????????????????
}


func removeBlinkRepByMap(slc []LuckyBlinkPerson) []LuckyBlinkPerson {  //?????????????????????
	var result []LuckyBlinkPerson           //??????????????????????????????
	tempMap := map[LuckyBlinkPerson]byte{}  // ?????????????????????
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0 					//???e?????????tempMap???????????????????????????????????????????????????key???????????????
		if len(tempMap) != l { 			// ??????map??????map?????????????????????????????????
			result = append(result, e)  //????????????????????????????????????????????????result???
		}
	}
	return result
}