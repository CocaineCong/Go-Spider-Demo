package main

import (
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	USERNAME = "root"
	PASSWORD = "root"
	HOST = "127.0.0.1"
	PORT = "3306"
	DBNAME = "douban_movie"
)

var DB *sql.DB


type MovieData struct {
	Title    string `json:"title"`
	Director string `json:"director"`
	Picture  string `json:"picture"`
	Actor    string `json:"actor"`
	Year     string `json:"year"`
	Score    string `json:"score"`
	Quote    string `json:"quote"`
}


func main() {
	InitDB()
	//for i := 0; i < 10; i++ {
	//	Spider(strconv.Itoa(i*25))
	//}
	ch := make(chan bool)
	for i := 0; i < 10; i++ {
		go Spider(strconv.Itoa(i*25), ch)
	}
	for i := 0; i < 10; i++ {
		<-ch
	}
}

func Spider(page string,ch chan bool) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://movie.douban.com/top250?start="+page, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", "https://movie.douban.com/chart")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	docDetail, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	docDetail.Find("#content > div > div.article > ol > li > div").
		Each(func(i int, s *goquery.Selection) {
			var movieData MovieData
			//#content > div > div.article > ol > li:nth-child(16) > div > div.info > div.hd > a > span:nth-child(1)
			title := s.Find("div.info > div.hd > a > span:nth-child(1)").Text()
			img := s.Find("div.pic > a > img")
			imgTmp, ok := img.Attr("src")
			info := strings.Trim(s.Find("div.info > div.bd > p:nth-child(1)").Text(), " ")
			director, actor, year := InfoSpite(info)
			score := strings.Trim(s.Find("div.info > div.bd > div > span.rating_num").Text(), " ")
			score = strings.Trim(score, "\n")
			quote := strings.Trim(s.Find("div.info > div.bd > p.quote > span").Text(), " ")
			if ok {
				movieData.Title = title
				movieData.Director = director
				movieData.Picture = imgTmp
				movieData.Actor = actor
				movieData.Year = year
				movieData.Score = score
				movieData.Quote = quote
				fmt.Println(movieData)
				InsertSql(movieData)
			}
		})
	if ch != nil{
		ch <- true
	}
}

func InfoSpite(info string) (director ,actor , year string) {
	directorRe, _ := regexp.Compile(`导演:(.*)主演:`)
	if len(director) < 8 {
		director = string(directorRe.Find([]byte(info)))
	} else {
		director = string(directorRe.Find([]byte(info)))[8:]
	}
	director = strings.Trim(director, "主演:")
	actorRe, _ := regexp.Compile(`主演:(.*)`)
	if len(actor) < 8 {
		actor = string(actorRe.Find([]byte(info)))
	} else {
		actor = string(actorRe.Find([]byte(info)))[8:]
	}
	yearRe, _ := regexp.Compile(`(\d+)`)
	year = string(yearRe.Find([]byte(info)))
	return
}

func InitDB(){
	path := strings.Join([]string{USERNAME, ":", PASSWORD, "@tcp(",HOST, ":", PORT, ")/", DBNAME, "?charset=utf8"}, "")
	DB, _ = sql.Open("mysql", path)
	DB.SetConnMaxLifetime(10)
	DB.SetMaxIdleConns(5)
	if err := DB.Ping(); err != nil{
		fmt.Println("opon database fail")
		return
	}
	fmt.Println("connnect success")
}

func InsertSql(movieData MovieData) bool {
	tx, err := DB.Begin()
	if err != nil {
		fmt.Println("tx fail")
		return false
	}
	stmt, err := tx.Prepare("INSERT INTO movie_data (`Title`,`Director`,`Picture`,`Actor`,`Year`,`Score`,`Quote`) VALUES (?, ?, ?,?,?,?,?)")
	if err != nil {
		fmt.Println("Prepare fail", err)
		return false
	}
	_, err = stmt.Exec(movieData.Title, movieData.Director, movieData.Picture, movieData.Actor, movieData.Year, movieData.Score, movieData.Quote)
	if err != nil {
		fmt.Println("Exec fail", err)
		return false
	}
	_ = tx.Commit()
	return true
}