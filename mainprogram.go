


// https://www.indiatoday.in/rss/1206584

// https://timesofindia.indiatimes.com/rssfeeds/7503091.cms

package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Media   string   `xml:"media,attr"`
	Atom    string   `xml:"atom,attr"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Text          string `xml:",chardata"`
		Title         string `xml:"title"`
		Description   string `xml:"description"`
		Link          string `xml:"link"`
		LastBuildDate string `xml:"lastBuildDate"`
		Generator     string `xml:"generator"`
		Image         struct {
			Text        string `xml:",chardata"`
			URL         string `xml:"url"`
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			Description string `xml:"description"`
		} `xml:"image"`
		Item []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			//Description testtype `xml:"description"`
			Description string `xml:"description"`
			PubDate     string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}
//type testtype struct{
//	Description string `xml:"description"`
//}
//func (l testtype) String () string {
//	return fmt.Sprintf(l.Description)
//}







type valueWeNeed struct{
	Title       string
	//Link        string
	Link template.URL
	//Description testtype
	Description string
	PubDate     string
}

var finalFinalData []valueWeNeed

var wg sync.WaitGroup

var politicList []secondPageData
var sushantList []secondPageData

func main() {

	http.HandleFunc("/", indexPage)
	http.HandleFunc("/tabpage/",tabPage)
	http.ListenAndServe(":8000", nil)
}



func extractData(queue chan<- valueWeNeed,title,link,description,pubdate string) {
	//defer func(){
	//	fmt.Println("Done!")
	//}()
	defer wg.Done()
	e := strings.Index(description, "</a>")
	description=description[e+4:]
	p := valueWeNeed{Title: title,
		Link        :template.URL(link),
		Description :description,
		PubDate     :pubdate,
	}
	//fmt.Println(p)
	//runtime.Gosched()
	queue<-p

}
func dataToWebsite(s Rss) []valueWeNeed{
	//finalData := make(map[int]valueWeNeed)
	var finalData []valueWeNeed
	queue := make(chan valueWeNeed,50)

	for i:=0;i<50;i++ {
		value:=s.Channel.Item[i]
		//fmt.Println("entering....")
		wg.Add(1)
		go extractData(queue,value.Title,value.Link,value.Description,value.PubDate)
		//fmt.Println("i value ",i,runtime.NumGoroutine())

	}
	wg.Wait()
	close(queue)
	fmt.Println("heloooo")



	for elem := range queue {
		finalData=append(finalData,elem)

	}


	return finalData


}


func indexPage(w http.ResponseWriter, r *http.Request) {
	resp, _ := http.Get("https://www.indiatoday.in/rss/1206584")
	bytes, _ := ioutil.ReadAll(resp.Body)

	var s Rss
	xml.Unmarshal(bytes, &s)
	resp.Body.Close()
	//https://stackoverflow.com/questions/38673673/access-http-response-as-string-in-go
	//xml.Unmarshal([]byte(resp.Body), &s)  --error

	//fmt.Println(s )
	finalFinalData=dataToWebsite(s)

	fmt.Println(finalFinalData[0].Link)

	page:=struct{
		News []valueWeNeed
	}{
		News:finalFinalData,
	}
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, page)



}
type secondPageData struct{
	Title string
	Link template.URL
}



func findindTabs(politicchannel ,shushantchannel chan secondPageData,quit chan int,newdata []valueWeNeed){
	politicWords:=[]string{
		"congress","BJP","sonia","rahul","PM","presdent",
	}
	sushant:=[]string{
		"sushant","rhea","suicide","CBI","singh","Singh","Sushant",
	}
	//var secondFinal []secondPageData
	for _,value:=range newdata{
		fmt.Println("------------------the babababbsaisijd;oadsoj----------------------",value)
		wg.Add(1)
		for _, word := range politicWords {
			if strings.Contains(value.Title,word){
				p:=secondPageData{
					Title:value.Title,
					Link:value.Link,
				}
				fmt.Println(p)
				politicchannel<-p

			}
		}
		for _, word := range sushant {
			if strings.Contains(value.Title,word){
				p:=secondPageData{
					Title:value.Title,
					Link:value.Link,
				}
				shushantchannel<-p
				fmt.Println(p)

			}
		}
		wg.Done()
	}
	close(shushantchannel)
	close(politicchannel)
	quit<-1

	return

}
func received(politicchannel,shushantchannel chan secondPageData,quit chan int){
	for{
		select{
			case v:=<-politicchannel:
				politicList=append(politicList,v)
				fmt.Println("hello case 1 ")
			case v:=<-shushantchannel:
				sushantList=append(sushantList,v)
				fmt.Println("hello case 2 ")
			case v:=<-quit:
				fmt.Println("quitingg.....",v)
				return
		}
	}


}


func tabPage(w http.ResponseWriter, r *http.Request){
	fmt.Println("enetering....")
	politicchannel := make(chan secondPageData,50)
	shushantchannel :=make(chan secondPageData,50)
	quit :=make(chan int)

	go findindTabs(politicchannel,shushantchannel,quit,finalFinalData)
	received(politicchannel,shushantchannel,quit)
	wg.Wait()
	fmt.Println(politicList)
	fmt.Println(sushantList)
	page:=struct{
		Politics []secondPageData
		Sushant []secondPageData
	}{
		Politics:politicList,
		Sushant:sushantList,
	}
	t, _ := template.ParseFiles("tabpage.html")
	t.Execute(w, page)

}

