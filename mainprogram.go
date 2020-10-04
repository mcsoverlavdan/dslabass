


// https://www.indiatoday.in/rss/1206584

// https://timesofindia.indiatimes.com/rssfeeds/7503091.cms

package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"runtime"
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

//times of india





type valueWeNeed struct{
	Title       string
	//Link        string
	Link template.URL
	//Description testtype
	Description string
	PubDate     string
}
//type (r Rss) String()
// xml to struct https://www.onlinetool.io/xmltogo/
var wg sync.WaitGroup
func main() {

	http.HandleFunc("/", indexPage)
	http.ListenAndServe(":8000", nil)
}

//func extractData(queue chan valueWeNeed,title,link,description,pubdate string) {
//	defer wg.Done()
//	e := strings.Index(description, "</a>")
//	description=description[e+4:]
//	p := valueWeNeed{Title: title,
//		Link        :link,
//		Description :description,
//		PubDate     :pubdate,
//	}
//	fmt.Println(p)
//	//runtime.Gosched()
//	queue<-p
//
//}
//func dataToWebsite(s Rss){
//	//finalData := make(map[int]valueWeNeed)
//	var finalData []valueWeNeed
//	queue := make(chan valueWeNeed)
//
//	for i:=0;i<50;i++ {
//		value:=s.Channel.Item[i]
//		//fmt.Println("entering....")
//		wg.Add(1)
//		go extractData(queue,value.Title,value.Link,value.Description,value.PubDate)
//		fmt.Println(runtime.NumGoroutine())
//
//	}
//	wg.Wait()
//	fmt.Println("heloooo")
//
//
//
//	for elem := range queue {
//		finalData=append(finalData,elem)
//
//	}
//	close(queue)
//
//	fmt.Println(finalData)
//
//
//
//	//fmt.Println(finalData[0].Link)
//	//fmt.Println(finalData[0].Description)
//	//str:=finalData[0].Description
//	//res1 := strings.TrimLeft(str, "</a>")
//	//e := strings.Index(str, "</a>")
//	//fmt.Println(res1)
//	//fmt.Println(str[e+4:])
//	//https://www.indiatoday.in/sports/ipl-2020/story/ipl-2020-cricket-simon-katich-rcb-fielding-devdutt-padikkal-banaglore-win-vs-rr-1728156-2020-10-04?utm_source=rss
////	<a href="https://www.indiatoday.in/sports/ipl-2020/story/ipl-2020-cricket-simon-katich-rcb-fielding-devdutt-padikkal-banaglore-win-vs-rr-1728156-2020-10-04"> <img align="left" hspace="2" height="180" width="305" alt="" border="0" src="https://akm-img-a-in.tosshub.com/indiatoday/images/story/202010/PTI03-10-2020_000230A-647x363.jpeg?5VykHvBFMeV2y.HV1i4Vg87ADnceociL"> </a> Asked for ruthless performance, RCB delivered: Katich after win against RR
//
//	//for _,v:=range finalData{
//	//	fmt.Println(v.Link)
//	//}
////description value
//	//<a href="https://www.indiatoday.in/india/story/air-india-one-boeing-prime-minister-narendra-modi-new-plane-inside-pics-president-kovind-vp-naidu-1727768-2020-10-02"> <img align="left" hspace="2" height="180" width="305" alt="" border="0" src="https://akm-img-a-in.tosshub.com/indiatoday/images/story/202010/Air_India_One_1-647x363.jpeg?70v1Duh8FBcSXjwmj4cX6GVqiKTl3uPk"> </a>
//	//EXCLUSIVE: Inside Boeing's Air India One bought for VVIPs, including PM, President
//
//
//}

func extractData(queue chan<- valueWeNeed,title,link,description,pubdate string) {
	defer func(){
		fmt.Println("Done!")
	}()
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
		fmt.Println("i value ",i,runtime.NumGoroutine())

	}
	wg.Wait()
	close(queue)
	fmt.Println("heloooo")



	for elem := range queue {
		finalData=append(finalData,elem)

	}


	return finalData



	//fmt.Println(finalData[0].Link)
	//fmt.Println(finalData[0].Description)
	//str:=finalData[0].Description
	//res1 := strings.TrimLeft(str, "</a>")
	//e := strings.Index(str, "</a>")
	//fmt.Println(res1)
	//fmt.Println(str[e+4:])
	//https://www.indiatoday.in/sports/ipl-2020/story/ipl-2020-cricket-simon-katich-rcb-fielding-devdutt-padikkal-banaglore-win-vs-rr-1728156-2020-10-04?utm_source=rss
	//	<a href="https://www.indiatoday.in/sports/ipl-2020/story/ipl-2020-cricket-simon-katich-rcb-fielding-devdutt-padikkal-banaglore-win-vs-rr-1728156-2020-10-04"> <img align="left" hspace="2" height="180" width="305" alt="" border="0" src="https://akm-img-a-in.tosshub.com/indiatoday/images/story/202010/PTI03-10-2020_000230A-647x363.jpeg?5VykHvBFMeV2y.HV1i4Vg87ADnceociL"> </a> Asked for ruthless performance, RCB delivered: Katich after win against RR

	//for _,v:=range finalData{
	//	fmt.Println(v.Link)
	//}
	//description value
	//<a href="https://www.indiatoday.in/india/story/air-india-one-boeing-prime-minister-narendra-modi-new-plane-inside-pics-president-kovind-vp-naidu-1727768-2020-10-02"> <img align="left" hspace="2" height="180" width="305" alt="" border="0" src="https://akm-img-a-in.tosshub.com/indiatoday/images/story/202010/Air_India_One_1-647x363.jpeg?70v1Duh8FBcSXjwmj4cX6GVqiKTl3uPk"> </a>
	//EXCLUSIVE: Inside Boeing's Air India One bought for VVIPs, including PM, President


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
	finalFinalData:=dataToWebsite(s)
	fmt.Println(finalFinalData[0].Link)

	//fmt.Printf("%T/n",finalFinalData)

	//t.Execute(w, p)
	//fmt.Fprintf(w, "<h1>testingggg</h1>")

	//p := valueWeNeed{Title:"this is the Title",
	//	Link        :"this is the Link",
	//	Description :"this is the Description",
	//	PubDate     :"this is the PubDate",
	//}

//https://stackoverflow.com/questions/54156119/range-over-string-slice-in-golang-template
//	{{range $element := .DataFields}} {{$element}} {{end}}
	page:=struct{
		News []valueWeNeed
	}{
		News:finalFinalData,
	}
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, page)



}