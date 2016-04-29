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
	JOB51.Register()
}

var JOB51 = &Spider{
	Name:        "JOB51",
	Description: "智联招聘职务  [http://51job.com//]",
	// Pausetime: 300,
	// Keyin:   KEYIN,
	// Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.Aid(map[string]interface{}{"loop": [2]int{1, 2}, "Rule": "请求列表"}, "请求列表")
		},

		Trunk: map[string]*Rule{

			"请求列表": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
				    for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {

					ctx.AddQueue(&request.Request{
						Url:  "http://search.51job.com/jobsearch/search_result.php?fromJs=1&jobarea=000000%2C00&district=000000&funtype=0000&industrytype=00&issuedate=9&providesalary=99&keyword=%E8%BD%AF%E4%BB%B6%E5%B7%A5%E7%A8%8B%E5%B8%88%28java%29&keywordtype=0&lang=c&stype=2&postchannel=0000&workyear=99&cotype=99&degreefrom=99&jobterm=99&companysize=99&lonlat=0%2C0&radius=-1&ord_field=0&list_type=0&fromType=14&dibiaoid=0&confirmdate=9&curr_page=" + strconv.Itoa(loop[0]),
						Rule: "请求列表",
					})
				    }
				    return nil
				},

				
			},

			"获取列表": {
				ParseFunc: func(ctx *Context) {
					logs.Log.Informational("获取列表log")

					ctx.GetDom().
						Find(".t1").
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

					thjob := query.Find(".tHjob").First()
					tCompany_main_jtag := query.Find(".tCompany_main").First().Find(".jtag").First()

					salary := thjob.Find("strong").First().Text()
					work_position := thjob.Find(".lname").First().Text()
					publish_date := tCompany_main_jtag.Find(".sp4").Eq(3).Text()
					job_type := ""
					job_years := tCompany_main_jtag.Find(".sp4").Eq(0).Text()
					education := tCompany_main_jtag.Find(".sp4").Eq(1).Text()
					number := tCompany_main_jtag.Find(".sp4").Eq(2).Text()
					job_category := thjob.Find("h1").First().Text()

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
