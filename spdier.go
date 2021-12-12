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

type Result4 struct {
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

var client http.Client

func HttpGet(url string) (result string, err error) {
	reqSpider, err := http.NewRequest("GET", url, nil)
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
	var results Result4
	_ = json.Unmarshal(bodyText, &results)
	result = fmt.Sprintf("%v",results.Data)
	return
}


//爬取网页
func spiderPage(url string,i int) string {
	fmt.Println("正在爬取", url)
	_, err := HttpGet(url)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(i)
	return url + " 抓取成功"
}
/*
可以看到我定义了两个通道
一个是用来存入 url 的，另一个是用来存入爬取结果的
缓冲空间是 100在方法 doWork 中，
sendResult 会阻塞等待 results 通道的输出
匿名函数则是等待 page 通道的输出紧接着下面就是把 200 个 url 写入 page 通道
匿名函数得到 page 的输出就会执行 asyn_worker 函数，也就是爬取 html 的函数了(将其存入results 通道)
然后 sendResult 函数得到 results 通道的输出，将结果打印出来 可以看到 我在匿名函数中并发了 20 个 goroution
并且启用了同步等待组作为参数传入，理论上可以根据机器的性能来定义 并发数
*/

func asynWorker(page chan string, results chan string,wg *sync.WaitGroup,i int) {
	defer wg.Done() //defer wg.Done()必须放在go并发函数内
	for {
		v, ok := <-page //显示的调用close方法关闭通道。
		if !ok {
			fmt.Println("已经读取了所有的数据，", ok)
			break
		}
		results <- spiderPage(v,i)
	}
}

func doWork(start, end int,wg *sync.WaitGroup) {
	fmt.Printf("正在爬取第%d页到%d页\n", start, end)
	//因为很有可能爬虫还没有结束下面的循环就已经结束了，所以这里就需要且到通道
	page := make(chan string, 100)
	results := make(chan string, 100)
	go sendResult(results, start, end)
	go func() {
		for i := 0; i <= 10; i++ {
			wg.Add(1)
			go asynWorker(page, results, wg, i)
		}
	}()
	for i := start; i <= end; i++ {
		url := "http://localhost:3000/api/v1/products"
		page <- url
		println("加入" + url + "到page")
	}
	println("关闭通道")
	close(page)
	wg.Wait()
}


func sendResult(results chan string,start,end int) {
	// 发送抓取结果
	for {
		_, ok := <-results
		if !ok {
			fmt.Println("已经读取了所有的数据，", ok)
			break
		}
		//fmt.Println(v)
	}
}

func main() {
	var wg sync.WaitGroup
	start := time.Now()
	doWork(1, 11, &wg)
	//输出执行时间，单位为毫秒。
	elapsed := time.Since(start)
	fmt.Println("time",elapsed)
}