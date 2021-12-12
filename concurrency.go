package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)


type Result struct {
	Data struct {
		Item []struct {
			BossAvatar    string `json:"boss_avatar"`
			BossID        int64  `json:"boss_id"`
			BossName      string `json:"boss_name"`
			CategoryID    int64  `json:"category_id"`
			CreatedAt     int64  `json:"created_at"`
			DiscountPrice int64  `json:"discount_price"`
			ImgPath       string `json:"img_path"`
			Info          string `json:"info"`
			Name          string `json:"name"`
			Num           int64  `json:"num"`
			OnSale        bool   `json:"on_sale"`
			Price         int64  `json:"price"`
			ProductID     int64  `json:"product_id"`
			Title         string `json:"title"`
			View          int64  `json:"view"`
		} `json:"item"`
		Total int64 `json:"total"`
	} `json:"data"`
	Error  string `json:"error"`
	Msg    string `json:"msg"`
	Status int64  `json:"status"`
}

var Client http.Client
var wg sync.WaitGroup

func main() {
	url := "http://localhost:3000/api/v1/products"
	NormalStart(url) // 单线程爬虫
	ChannelStart(url) // Channel多线程爬虫
	WaitGroupStart(url) // Wait 多线程爬虫
}

func NormalStart(url string) {
	start := time.Now()
	for i := 0; i < 10; i++ {
		Spider(url, nil, i)
	}
	elapsed := time.Since(start)
	fmt.Printf("NormalStart Time %s \n", elapsed)
}


func ChannelStart(url string) {
	ch := make(chan bool)
	start := time.Now()
	for i := 0; i < 10; i++ {
		go Spider(url, ch, i)
	}
	for i := 0; i < 10; i++ {
		<-ch
	}
	elapsed := time.Since(start)
	fmt.Printf("ChannelStart Time %s \n", elapsed)
}

func WaitGroupStart(url string) {
	start := time.Now()
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			Spider(url,nil,i)
		}(i)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("WaitGroupStart Time %s\n ", elapsed)
}

func Spider(url string,  ch chan bool, i int) {
	reqSpider, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	reqSpider.Header.Set("content-length", "0")
	reqSpider.Header.Set("accept", "*/*")
	reqSpider.Header.Set("x-requested-with", "XMLHttpRequest")
	respSpider, err := Client.Do(reqSpider)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, _ := ioutil.ReadAll(respSpider.Body)
	var result Result
	_ = json.Unmarshal(bodyText, &result)
	//fmt.Println(i,result.Data)
	if ch != nil{
		ch <- true
	}
}