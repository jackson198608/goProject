package main

import (
	// "bytes"
	// "encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	// "io"
	// "io/ioutil"
	// "math/rand"
	// "net/http"
	// "os"
	"reflect"
	// "regexp"
	// "strings"
	// // "time"
	// "github.com/jackson198608/goProject/common/http/abuyunHttpClient"
	// "strconv"
)

// func getImg(url string) (n int64, err error) {
// 	path := strings.Split(url, "/")
// 	var name string
// 	if len(path) > 1 {
// 		name = path[len(path)-1]
// 	}
// 	fmt.Println(name)
// 	out, err := os.Create(name)
// 	defer out.Close()
// 	resp, err := http.Get(url)
// 	defer resp.Body.Close()
// 	pix, err := ioutil.ReadAll(resp.Body)
// 	n, err = io.Copy(out, bytes.NewReader(pix))
// 	return

// }

// func GenerateRangeNum(min, max int) int {
// 	// rand.Seed(time.Now().Unix())
// 	randNum := rand.Intn(max-min) + rand.Intn(min)
// 	return randNum
// }

const proxyServer = ""
const proxyUser = ""
const proxyPasswd = ""

var driverName string = "mysql"

var dbAuth string = "dog123:dog123"

var dbDsn string = "192.168.86.193:3307"

var dbName string = "process"

type ExecuteRecord struct {
	Id int
}

func main() {

	dataSourceName := dbAuth + "@tcp(" + dbDsn + ")/" + dbName + "?charset=utf8mb4"
	engine, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		fmt.Println(err)
	}
	var exec = new(ExecuteRecord)
	fmt.Println(reflect.TypeOf(exec))
	counts, err := engine.Count(exec)
	fmt.Println(counts)
	fmt.Println(err)
	/**var abuyun *AbuyunProxy = abuyunHttpClient.NewAbuyunProxy(proxyServer,
		proxyUser,
		proxyPasswd)

	// t.Log("begin the test")

	if abuyun == nil {
		// t.Error("create abuyun error")
		return
	}

	targetUrl := "http://m.goumin.com/"

	var h http.Header = make(http.Header)
	h.Set("a", "1")
	statusCode, responseHeader, _, err := abuyun.SendRequest(targetUrl, h, true)
	fmt.Println(statusCode)
	fmt.Println(responseHeader)
	//fmt.Println(body)
	fmt.Println(err)
	s := "[size=2][color=#999999]go_u_m_admin 发表于 17/8/29 14:14[/color] [url=forum.php?mod=redirect&goto=findpost&pid=65821373&ptid=4437190][img]static/image/common/back.gif[/img][/url][/size]\n\rsdfsdfsdf[/quote][size=4][color=purple]    [img=255,90]http://i7.photobucket.com/albums/y253/rozzi/Hello.gif[/img][/color][/size][size=4][color=#800080][/color][/size][url=http://bbs.goumin.com/thread-450885-1-1.html][size=4][color=red][b]★★★★ 狗民网边境牧羊犬俱乐部QQ群以及开心账号 ★★★★[/b][/color][/size][/url][size=4][color=red][/color][/size]\n[size=4][color=purple]大家都来为狗狗建立小档案吧 也方便大家互相认识 \n格式如下[/color]:[/size]\n狗狗  [size=2] 姓名:\n性别:\n出生日期:\n所在地:\n爱好:\n联系方式(自愿):[/size]\n[color=green][size=5]最后别忘了狗狗的靓照哦[/size][/color]\n[[i] 本帖最后由 sunxiaochou 于 2011-5-9 09:10 编辑 [/i]]"
	m := strings.Index(s, "[quote]")
	n := strings.Index(s, "[/quote]")
	// if m >= 0 && n > 0 {

	fmt.Println(reflect.TypeOf(m), n)
	cont := s[m:n]
	fmt.Println(cont)
	// }
	// dir := strconv.Itoa(tid % 1000)
	tid := 16435695
	dir := ""
	if tid < 5000000 {
		dir = "4/" + strconv.Itoa(tid%1000)

	} else {
		dir = strconv.Itoa(tid/1000000) + "/" + strconv.Itoa(tid%500)
	}
	filename := dir + "/thread-" + strconv.Itoa(tid) + "-1-1.html"
	fmt.Println(tid % 500)
	fmt.Println(filename)
	// for i := 0; i < 10; i++ {
	// 	fmt.Println(GenerateRangeNum(500, 1000))
	// }

	// content := "[color=#EE1B2E][b][url=http://bbs.goumin.com/thread-4423738-1-1.html][color=#EE1B2E][b][color=#EE1B2E][b][color=#EE1B2E][b][color=#EE1B2E][b][color=#EE1B2E][b][color=#EE1B2E][b][size=3][color=Red][b]要求：只要是自家宝贝即可\n，图片1-3张即可每位有效参与ID奖励5分 5米[/b][/color][/size][/b][/color][/b][/color][/b][/color][/b][/color][/b][/color][/b][/color][/url][/b][/color]"
	content := "[color=#EE1B2E][b][url=http://bbs.goumin.com/thread-4423738-1-1.html][color=#EE1B2E][b][color=#EE1B2E][b][color=#EE1B2E][b][color=#EE1B2E][b][color=#EE1B2E][b][color=#EE1B2E][b][size=3][color=Red][b]要求：只要是自家宝贝即可，图片1-3张即可每位有效参与ID奖励5分 5米[/b][/color][/size][/b][/color][/b][/color][/b][/color][/b][/color][/b][/color][/b][/color][/url][/b][/color][url]www.wopet.cn[/url][media=wmv,400,300,0]http://player.youku.com/player.php/sid/XMTk1MDcwMzI=/v.swf[/media][quote][size=2][color=#999999]go_u_m_admin 发表于 17/8/29 14:14[/color] [url=forum.php?mod=redirect&goto=findpost&pid=65821373&ptid=4437190][img]static/image/common/back.gif[/img][/url][/size]\n\rsdfsdfsdf[/quote][size=4][color=purple]    [img=255,90]http://i7.photobucket.com/albums/y253/rozzi/Hello.gif[/img][/color][/size][size=4][color=#800080][/color][/size][url=http://bbs.goumin.com/thread-450885-1-1.html][size=4][color=red][b]★★★★ 狗民网边境牧羊犬俱乐部QQ群以及开心账号 ★★★★[/b][/color][/size][/url][size=4][color=red][/color][/size]\n[size=4][color=purple]大家都来为狗狗建立小档案吧 也方便大家互相认识 \n格式如下[/color]:[/size]\n狗狗  [size=2] 姓名:\n性别:\n出生日期:\n所在地:\n爱好:\n联系方式(自愿):[/size]\n[color=green][size=5]最后别忘了狗狗的靓照哦[/size][/color]\n[[i] 本帖最后由 sunxiaochou 于 2011-5-9 09:10 编辑 [/i]]"

	// bytes := []byte(content)
	// for {
	// m := strings.Index(content, "[quote]")
	// n := strings.Index(content, "[/quote]")
	// fmt.Println(m, n)
	// if m >= 0 && n > 0 {
	// 	substr := content[m:n]
	// 	content = strings.Replace(content, substr+"[/quote]", "", -1)
	// }

	// }
	// fmt.Println(content)
	// re, _ := regexp.Compile("\\[img\\](.*?)\\[\\/img\\]")
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	content = re.ReplaceAllStringFunc(content, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	content = re.ReplaceAllString(content, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	content = re.ReplaceAllString(content, "")
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	content = re.ReplaceAllString(content, "\n")

	re, _ = regexp.Compile("\\[img(.*?)\\](.*?)\\[/img\\]")
	content = re.ReplaceAllString(content, "<img alt=\"333\" src=\"$2\">")
	//去除所有尖括号内的HTML代码，并换成换行符

	// re, _ = regexp.Compile("\\[(quote|flash|audio|wma|wmv|rm|media)\\](.*?)\\[/(quote|flash|audio|wma|wmv|rm|media)\\]")
	// content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[hr\\]")
	content = re.ReplaceAllString(content, "<hr>")

	re, _ = regexp.Compile("\\[size\\=(.*?)\\](.*?)")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[font\\=(.*?)\\](.*?)")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[color\\=(.*?)\\](.*?)")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[backcolor\\=(.*?)\\](.*?)")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[/size\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[/font\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[/color\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[/backcolor\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[list=(\\d+)\\]")
	content = re.ReplaceAllString(content, "<ul class='numeric'>")

	re, _ = regexp.Compile("\\[list\\]")
	content = re.ReplaceAllString(content, "<ul class='dot'>")

	re, _ = regexp.Compile("\\[/list\\]")
	content = re.ReplaceAllString(content, "</ul>")

	re, _ = regexp.Compile("\\[url=.*?goto=findpost&pid=\\d+&ptid=\\d+\\](.*?)\\[/url\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[url=(http://bbs.goumin.com/)?forum-\\d+-\\d+\\.html\\](.*?)\\[/url\\]")
	content = re.ReplaceAllString(content, "$2")

	re, _ = regexp.Compile("\\[url=home.php\\?mod=space&uid=\\d+\\](.*?)\\[/url\\]")
	content = re.ReplaceAllString(content, "$1")

	re, _ = regexp.Compile("\\[url=(http.*?)\\](.*?)\\[/url\\]")
	content = re.ReplaceAllString(content, "<a href='$1'>$2</a>")

	re, _ = regexp.Compile("\\[url\\](.*?)\\[/url\\]")
	content = re.ReplaceAllString(content, "<a href='$1'>$1</a>")

	re, _ = regexp.Compile("\\[url=(mailto:.*?)\\](.*?)\\[/url\\]")
	content = re.ReplaceAllString(content, "<a href='$1'>$2</a>")

	// re, _ = regexp.Compile("(http:\\/\\/bbs\\.goumin\\.com\\/thread-\\d+-\\d+-\\d+.html)")
	// content = re.ReplaceAllString(content, "<a href='$1'>$2</a>")

	//代码 code
	re, _ = regexp.Compile("\\[code\\](.*?)\\[/code\\]")
	content = re.ReplaceAllString(content, "<div class='code'><pre>$1</pre></div>")
	// 音视频 audio
	re, _ = regexp.Compile("\\[audio\\](.*?)\\[/audio\\]")
	content = re.ReplaceAllString(content, "<audio src='$1' controls></audio>")

	// 视频 media
	re, _ = regexp.Compile("\\[media\\](.*?\\.mp4)\\[/media\\]")
	content = re.ReplaceAllString(content, "<video src='$1' controls></video>")

	re, _ = regexp.Compile("\\[media=(.*?)\\]http://v.youku.com/v_show/id_(.*?).html.*?\\[/media\\]")
	content = re.ReplaceAllString(content, "<a class='post_content_link' href='http://v.youku.com/v_show/id_$1.html'>***优酷视频点击播放***</a>")

	re, _ = regexp.Compile("\\[media.*?\\]http:\\/\\/player\\.youku\\.com\\/player\\.php\\/sid\\/(.*?)\\/v\\.swf\\[/media\\]")
	content = re.ReplaceAllString(content, "<a class='post_content_link' href='http://v.youku.com/v_show/id_$1.html'>***优酷视频点击播放***</a>")

	re, _ = regexp.Compile("\\[media.*?\\]http:\\/\\/www\\.tudou\\.com\\/programs\\/view\\/(.*?)\\/\\[/media\\]")
	content = re.ReplaceAllString(content, "<a class='post_content_link' href='http://www.tudou.com/programs/view/$1/'>***土豆视频点击播放***</a>")

	re, _ = regexp.Compile("\\[media\\](.*?)\\[/media\\]")
	content = re.ReplaceAllString(content, "")

	// 表格 table td tr
	re, _ = regexp.Compile("\\[table.*?\\]")
	content = re.ReplaceAllString(content, "<table class='post_table'>")

	re, _ = regexp.Compile("\\[\\/table\\]")
	content = re.ReplaceAllString(content, "</table>")

	re, _ = regexp.Compile("\\[tr(=.*?)?\\]")
	content = re.ReplaceAllString(content, "<tr>")

	re, _ = regexp.Compile("\\[\\/tr\\]")
	content = re.ReplaceAllString(content, "</tr>")

	re, _ = regexp.Compile("\\[td(=.*?)?\\]")
	content = re.ReplaceAllString(content, "<td>")

	re, _ = regexp.Compile("\\[\\/td\\]")
	content = re.ReplaceAllString(content, "</td>")

	// 左右对齐align，浮动 float
	re, _ = regexp.Compile("\\[align=(.*?)\\](.*?)\\[/align\\]")
	content = re.ReplaceAllString(content, "<div align='$1'><p>$2</p></div>")

	re, _ = regexp.Compile("\\[align.*?\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[/align\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[float=(.*?)\\](.*?)\\[/float\\]")
	content = re.ReplaceAllString(content, "<div style='float:$1'>$2</div>")

	re, _ = regexp.Compile("\\[float.*?\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[quote\\][\n\r]*(.+?)[\n\r]*\\[\\/quote\\]")
	content = re.ReplaceAllString(content, "<div class='quote'><blockquote>“$1”</blockquote></div>")

	re, _ = regexp.Compile("\\[free\\](.*?)\\[/free\\]")
	content = re.ReplaceAllString(content, "<div class='quote'><blockquote>“$1”</blockquote></div>")

	re, _ = regexp.Compile("\\[hide\\](.*?)\\[/hide\\]")
	content = re.ReplaceAllString(content, "<div class='quote'><blockquote>“$1”</blockquote></div>")

	// i
	re, _ = regexp.Compile("\\[i.*?\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[/i\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[u.*?\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[/u\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[b.*?\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[/b\\]")
	content = re.ReplaceAllString(content, "")

	if strings.Contains(content, "[attach]") == true && strings.Contains(content, "[/attach]") == true {
		// count := strings.Count(content, "[attach]")
		// re = regexp.MustCompile("\\[attach\\](\\d+)\\[\\/attach\\]")
		// str := re.FindAllStringSubmatch(content, -1)
		// // fmt.Println(str)
		// for i := 0; i < count; i++ {
		// 	fmt.Println(str[i][1])
		// 	fmt.Println(i)
		// }
	}
	fmt.Println(content)
	// data := make(map[string]interface{}, 1)
	// data["name"] = "xiaochuan"
	// data["age"] = 23
	// //序列化
	// json_obj, err := json.Marshal(data)
	// if err != nil {
	// 	fmt.Println("err :", err)
	// }
	// fmt.Println(reflect.TypeOf(data))

	/*json_obj := []byte(`a:6:{s:8:"required";b:0;s:8:"listable";b:1;s:6:"prefix";s:1:"1";s:5:"types";a:5:{i:1850;s:12:"经验分享";i:1851;s:12:"咨询求助";i:1852;s:12:"养宠故事";i:1853;s:12:"版内活动";i:1854;s:12:"版主公告";}s:5:"icons";a:5:{i:1850;s:0:"";i:1851;s:0:"";i:1852;s:0:"";i:1853;s:0:"";i:1854;s:0:"";}s:10:"moderators";a:5:{i:1850;N;i:1851;N;i:1852;N;i:1853;N;i:1854;N;}}`)
	marshal_data := make(map[string]interface{}, 1)
	//反序列化
	json_err := json.Unmarshal(json_obj, &marshal_data)
	if json_err != nil {
		fmt.Println(json_err)
	}
	fmt.Println(marshal_data["name"])

	start := 0
	end := 0
	content := "sfsdflsdkfdjsk<img src=''>sdfsdlfsdf<img ksdjfhsdf >sdfsdlfsdf<img ksdjfhsdf >"
	lenstr := len(content)
	contentBytes := []byte(content)
	// a = strings.Index(string(str), "<img")
	for {
		// aa := strings.IndexRune(string(contentBytes), "<img")
		a := strings.Index(string(contentBytes), "<img")
		if a < 0 {
			contentBytes[start] = '<'
			contentBytes[end] = '>'
			break
		}
		contentBytes[a] = '['
		for i := a; i < lenstr; i++ {
			if contentBytes[i] == '>' {
				contentBytes[i] = ']'
				start = a
				end = i
				break
			}
		}
	}
	b := strings.Index(string(contentBytes), "[img")
	if b > 0 {
		contentBytes[b] = '<'
		for j := b; j < lenstr; j++ {
			if contentBytes[j] == ']' {
				contentBytes[j] = '>'
				break
			}
		}
	}
	content = string(contentBytes)
	fmt.Println(content)
	// getImg("http://read.html5.qq.com/image?src=forum&q=5&r=0&imgflag=7&imageUrl=http://mmbiz.qpic.cn/mmbiz_jpg/hQaRzgJkxVVFVQWcNFytd2txqt2ww4MtS6UYtQhc9VPaZAY4qZMoxI7VxG00dc1HWIrxq2DdehCLgyrJKOGia8w/0.jpeg")
	*/

}

//该片段来自于http://outofmemory.cn
