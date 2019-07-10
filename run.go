package main

//自动生成日志文件
import (
	"flag"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"
)

var uaList = []string{
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/22.0.1207.1 Safari/537.1",
	"Mozilla/5.0 (X11; CrOS i686 2268.111.0) AppleWebKit/536.11 (KHTML, like Gecko) Chrome/20.0.1132.57 Safari/536.11",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.6 (KHTML, like Gecko) Chrome/20.0.1092.0 Safari/536.6",
	"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/536.6 (KHTML, like Gecko) Chrome/20.0.1090.0 Safari/536.6",
	"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/19.77.34.5 Safari/537.1",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/536.5 (KHTML, like Gecko) Chrome/19.0.1084.9 Safari/536.5",
	"Mozilla/5.0 (Windows NT 6.0) AppleWebKit/536.5 (KHTML, like Gecko) Chrome/19.0.1084.36 Safari/536.5",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1063.0 Safari/536.3",
	"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1063.0 Safari/536.3",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Trident/4.0; SE 2.X MetaSr 1.0; SE 2.X MetaSr 1.0; .NET CLR 2.0.50727; SE 2.X MetaSr 1.0)",
	"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1062.0 Safari/536.3",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1062.0 Safari/536.3",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; 360SE)",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1061.1 Safari/536.3",
	"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1061.1 Safari/536.3",
	"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1061.0 Safari/536.3",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/535.24 (KHTML, like Gecko) Chrome/19.0.1055.1 Safari/535.24",
	"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/535.24 (KHTML, like Gecko) Chrome/19.0.1055.1 Safari/535.24",
}

func randInt(min, max int) int {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if min > max {
		return max
	}

	return r.Intn(max-min) + min
}

func makeLog(currenturl, referUrl, ua string) string {
	u := url.Values{}
	u.Set("time", "1")
	u.Set("url", currenturl)
	u.Set("referUrl", referUrl)
	u.Set("ua", ua)
	paramsStr := u.Encode()

	logTemplate := "127.0.0.1 - - [09/Jul/2019:19:40:02 +0800] \"GET /dig.gif/?{$paramsStr}HTTP/1.1\" 200 43 \"{$ua}\""
	log := strings.Replace(logTemplate, "{$paramsStr}", paramsStr, -1)
	log = strings.Replace(log, "{$ua}", ua, -1)
	return log
}

func main() {
	total := flag.Int("total", 100, "how many rows")
	filePath := flag.String("filePath", "/var/log/nginx/test.emptygif.com.log", "path")
	flag.Parse()
	// fmt.Println(*total, *filePath)

	//生成total 行日志内容
	logStr := ""
	for i := 0; i < *total; i++ {
		currenturl := "house.leju.com/search"
		referURL := "house.leju.com"
		ua := uaList[randInt(0, len(uaList)-1)]
		logStr = logStr + makeLog(currenturl, referURL, ua) + "\n"

	}
	fd, _ := os.OpenFile(*filePath, os.O_RDWR|os.O_APPEND, 0644)
	fd.Write([]byte(logStr))
	fd.Close()
	fmt.Println("done")
}
