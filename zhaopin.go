package spider_lib

// 基础包
import (
	"github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	"github.com/henrylee2cn/pholcus/logs"               //信息输出
	. "github.com/henrylee2cn/pholcus/app/spider" //必需
	// . "github.com/henrylee2cn/pholcus/app/spider/common"          //选用

	// net包
	// "net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	// "encoding/json"

	// 字符串处理包
	// "regexp"
	"strconv"
	// 其他包
	// "fmt"
	// "math"
	// "time"
)

func init() {
	ZHILIAN.Register()
}

var ZHILIAN = &Spider{
	Name:        "zhaopin",
	Description: "智联招聘职务  [http://sou.zhaopin.com/]",
	// Pausetime: 300,
	// Keyin:   KEYIN,
	// Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.Aid(map[string]interface{}{"loop": [2]int{0, 1}, "Rule": "请求列表"}, "请求列表")
		},

		Trunk: map[string]*Rule{

			"请求列表": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
				    for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {

					ctx.AddQueue(&request.Request{
						Url:  "http://sou.zhaopin.com/jobs/searchresult.ashx?jl=%E5%8C%97%E4%BA%AC&kw=java%E9%AB%98%E7%BA%A7%E5%B7%A5%E7%A8%8B%E5%B8%88&sm=0&p=" + strconv.Itoa(loop[0]),
						Rule: "请求列表",
					})
				    }
				    return nil
				},

				ParseFunc: func(ctx *Context) {
					var curr int

					logs.Log.Informational("页码：" ,curr)
					logs.Log.Informational("页码：" ,strconv.Itoa(curr+1))

					ctx.AddQueue(&request.Request{
						Url:  "http://sou.zhaopin.com/jobs/searchresult.ashx?jl=%E5%8C%97%E4%BA%AC&kw=java%E9%AB%98%E7%BA%A7%E5%B7%A5%E7%A8%8B%E5%B8%88&sm=0&p=" + strconv.Itoa(curr+1),
						Rule: "请求列表",
						Temp: map[string]interface{}{"p": curr + 1},
					})

					// 用指定规则解析响应流
					ctx.Parse("获取列表")
				},
			},

			"获取列表": {
				ParseFunc: func(ctx *Context) {
					logs.Log.Informational("获取列表log")

					logs.Log.Informational("获取列表GetDom", ctx.GetDom())

					ctx.GetDom().
						Find(".zwmc").
						Each(func(i int, s *goquery.Selection) {
							url, _ := s.Find("a").Attr("href")

							logs.Log.Informational("url:", url)

							ctx.AddQueue(&request.Request{
								Url:      url,
								Rule:     "output",
								Priority: 1,
							})
						})
				},
			},

			"output": {
				//注意：有无字段语义和是否输出数据必须保持一致
				ItemFields: []string{
					"salary",
					"work_position",
					"publish_date",
					"job_type",
					"job_years",
					"education",
					"number",
					"job_category",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

					domresult := query.Find(".terminalpage-left").First().Find("li")

					salary := domresult.First().Find("strong").First().Text()
					work_position := domresult.Eq(1).Find("strong").First().Text()
					publish_date := domresult.Eq(2).Find("strong").First().Text()
					job_type := domresult.Eq(3).Find("strong").First().Text()
					job_years := domresult.Eq(4).Find("strong").First().Text()
					education := domresult.Eq(5).Find("strong").First().Text()
					number := domresult.Eq(6).Find("strong").First().Text()
					job_category := domresult.Eq(7).Find("strong").First().Text()

					// 结果存入Response中转
					ctx.Output(map[int]interface{}{
						0: salary,
						1: work_position,
						2: publish_date,
						3: job_type,
						4: job_years,
						5: education,
						6: number,
						7: job_category,
					})
				},
			},
		},
	},
}
