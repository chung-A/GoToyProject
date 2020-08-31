package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type extractedJob struct{
	id string
	title string
	location string
	salary string
	summary string
}

var errNoFound=errors.New("no found url")
var baseUrl="https://kr.indeed.com/%EC%B7%A8%EC%97%85?q=Spring&limit=50"

/*삽질 기록*/
//1. 패키지 명은 가장 끝단의 폴더명과 동일하게 해야함.

func main() {
	webScrapper()
}

func webScrapper(){
	c:=make(chan []extractedJob)
	var jobs []extractedJob
	totalPages:=getPages()

	for i:=0;i<totalPages;i++{
		go getPage(i,c)
	}

	for i:=0;i<totalPages;i++{
		jobs=append(jobs,<-c...)
	}

	writeJobs(jobs)
	fmt.Println("추출 완료, 추출한 공고 수: ", len(jobs))
}

func writeJobs(jobs []extractedJob){
	file,err:=os.Create("indeed 직업.csv")
	errCheck(err)

	w:=csv.NewWriter(file)
	defer w.Flush()//이시점에서 파일에 데이터를 입력하고 저장됨.

	headers:=[]string{"ID","Title","Location","Salary","Summary"}
	wErr:=w.Write(headers)
	errCheck(wErr)

	for _,job:=range jobs{
		go writeRows(job,w)
	}
}

func writeRows(job extractedJob,w *csv.Writer){
	var buffer bytes.Buffer
	buffer.WriteString("https://kr.indeed.com/%EC%B1%84%EC%9A%A9%EB%B3%B4%EA%B8%B0?jk=")
	buffer.WriteString(job.id)
	jobSlice :=[]string{buffer.String(),job.title,job.location,job.salary,job.summary}
	jwErr:=w.Write(jobSlice)
	errCheck(jwErr)
	buffer.Reset()
}

func getPage(pageNumber int,c2 chan []extractedJob){
	c:=make(chan extractedJob)

	var jobs []extractedJob
	pageURL:=baseUrl+"&start="+strconv.Itoa(pageNumber*50)
	fmt.Println("request ",pageURL)
	res,err:=http.Get(pageURL)
	errCheck(err)
	checkCode(res)
	defer res.Body.Close()

	doc:=getDocument(res)

	maxCnt:=0
	doc.Find(".jobsearch-SerpJobCard").Each(func(i int, selection *goquery.Selection) {
		go extractJob(selection,c)
		maxCnt=i
	})

	for i:=0;i<maxCnt;i++{
		jobs=append(jobs,<-c)
	}
	c2<-jobs
}

func extractJob(selection *goquery.Selection,c chan extractedJob){
	id,_:=selection.Attr("data-jk")
	id=cleanString(id)
	title:=cleanString(selection.Find(".title>a").Text())
	location:=cleanString(selection.Find(".sjcl").Text())
	salary:=cleanString(selection.Find(".salaryText").Text())
	summary:=cleanString(selection.Find(".summary").Text())
	c<-extractedJob{id: id, title: title, location: location, salary: salary, summary: summary}
	return
}

func cleanString(str string) string{
	return strings.Join(strings.Fields(strings.TrimSpace(str)),"")
}

func getPages() (pages int){
	res,err:=http.Get(baseUrl)
	errCheck(err)
	checkCode(res)
	defer res.Body.Close()

	doc:=getDocument(res)

	doc.Find(".pagination").Each(func(i int, selection *goquery.Selection) {
		pages=selection.Find("a").Length()
	})

	fmt.Println("전체 페이지: ",pages)
	return pages
}

func getDocument(res *http.Response) (doc *goquery.Document){
	doc,err:=goquery.NewDocumentFromReader(res.Body)
	errCheck(err)
	return
}

func errCheck(err error){
	if err!=nil{
		log.Fatal(err)
	}
}

func checkCode(res *http.Response){
	if res.StatusCode!=200 {
		log.Fatal("Request failed with status: ",res.StatusCode)
	}
}

//url checker-goroutine
func checkUrl(){
	urls:=[]string{
		"https://www.google.com/",
		"https://www.naver.com/",
		"https://www.amazon.com/",
		"https://www.airbnb.com/",
		"https://www.reddit.com/",
		"http://203.252.161.219:8081/",
		"https://www.soundcloud.com/",
		"https://www.facebook.com/",
	}

	c:=make(chan string)
	var result = map[string]string{}
	for _,url := range urls{
		fmt.Println("check url",url)
		go hitUrlWithChannel(url,c)

		//err:=hitUrl(url)
		//if err!=nil{
		//	result[url]="Fail"
		//}else{
		//	result[url]="Ok"
		//}
	}

	for i:=0;i< len(urls);i++{
		result[urls[i]]=<-c
	}

	fmt.Println("검사 완료! ****************************************")
	for key,value:=range result {
		fmt.Println("key: ",key,"/value: "+value)
	}
}

func hitUrl(url string) error{
	response,err:=http.Get(url)
	if err==nil&&response!=nil&&response.StatusCode<400 {
		//fmt.Println("check success!")
		return nil
	}

	if response!=nil{
		fmt.Println("Check Failed! ","/statusCode: ",response.StatusCode,"/err: ",err)
	}else{
		fmt.Println("there is no such host")
	}
	return errNoFound
}

func hitUrlWithChannel(url string,c chan string) {
	var buffer bytes.Buffer
	response, err := http.Get(url)
	if err==nil&&response!=nil&&response.StatusCode<400 {
		c<-"check success!"
	}

	if response!=nil{
		buffer.WriteString("Check Failed! /statusCode: ")
		buffer.WriteString(strconv.Itoa(response.StatusCode))
	}else{
		buffer.WriteString("there is no such host")
	}
	s:=buffer.String()
	c<-s
}

	//...으로 하면 동일한 타입의 여러 매개변수를 받을 수 있다.
//_하면 무시가능.
func arrayArgFunc(len int, txt ...string) int {
	fmt.Println("myFunc")
	fmt.Println(txt)
	return 0
}

//naked 기능.
func nakedFunc(name string) (length int, orignTxt string) {
	length = len(name)
	orignTxt = name
	return
}

//defer기능-함수가 끝나고 원래자리로 가기전에 defer내용을 실행시키고 감.
//defer 내부에 함수안에 매개변수로 함수가 있으면 미리 실행시키고 가는듯...?
//println(함수()) 이런식으로는 안쓰는게 좋을 듯 하다.
func deferFunc(testStr string) (rst string) {
	defer println(len(testStr))
	println(testStr + "/ deferFunc")
	return testStr
}

//loop
