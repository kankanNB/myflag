package maigi

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

var r = flag.Bool("r", false, "加入参数会自动解压在/root/web/下的zip文件到项目目录")
var b = flag.Bool("b", false, "在指定-t参数的时候可以用-b来删除文件内对应的文件夹")
var s = flag.Bool("s", false, "设置为ssl模板")
var a = flag.String("a", "", "添加域名对应文件夹")
var e = flag.String("e", "", "删除域名对应文件夹")
var d = flag.String("dir", "/www", "设置文件目录 默认为/www")
var t = flag.String("t", "", "批量参数后面带文件")
var list = flag.Bool("list", false, "查看域名状态")

func main() {
	flag.Parse()
	switch {
	case *a != "" && *e != "": //同时使用a和e参数就会返回
		return
	case *list != false: //使用list参数就会测试项目状态
		_list(_dir(*d))
	case *t != "" && *b == false: //使用t参数而且未使用b参数就会批量生产项目文件夹g
		for _, context := range _batch(*t) {
			_mkdir(context, _dir(*d), *s)
		}
	case *t != "" && *b != false: //使用t参数而且使用b参数就会批量删除项目文件夹
		for _, context := range _batch(*t) {
			_delete(context, _dir(*d))
		}
	case *a != "" && *e == "":
		_mkdir(*a, _dir(*d), *s)
	case *e != "" && *a == "":
		_delete(*e, _dir(*d))

	}

}

// 快捷解压函数调用
func _unzip(do, dir string, r bool) {

}

// 创建快捷zip
func _mkdir_root(do string) {
	var dir = "/root/web/" + do
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("文件夹" + dir + "已经存在")
}

// 删除快捷zip
func _delete_root(do string) {
	var dir = "/root/web/" + do
	err := os.RemoveAll(dir)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("文件夹" + dir + "已经删除")
}

// 检测文件夹有无带/ 没有就会加上
func _dir(dir string) string {
	// 检查字符串末尾是否有斜杠
	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
		return dir
	}
	return dir
}

// 判断文件夹在不在不存在返回flase
func _directoryExists(dirname string) bool {
	fileInfo, err := os.Stat(dirname)
	if err != nil {
		if os.IsNotExist(err) {
			return false // 文件夹不存在
		}
	}
	return fileInfo.IsDir()
}

// 检查项目文件夹是否存在
func _mkdir(do, dir string, ssl bool) {
	//dirname := "example_directory" // 要检查的文件夹路径
	var DirConf = dir + "http/" + do
	if _directoryExists(DirConf) {
		fmt.Printf("文件夹 %s 已经存在\n", DirConf)
	} else {
		_mkdiring(do, dir, ssl)
	}
}

// 创建项目文件夹
func _mkdiring(do, dir string, ssl bool) {
	//dir=./www
	var DirConf = dir + "http/" + do
	var DirData = DirConf + "/data/"
	var List = dir + "list/"
	var FILE = dir + "list/" + "list.log"

	err1 := os.MkdirAll(DirData, os.ModePerm)
	if err1 != nil {
		fmt.Println(err1.Error())
	}
	fmt.Println("data路径：" + DirData)
	err2 := os.MkdirAll(DirConf, os.ModePerm)
	if err2 != nil {
		fmt.Println(err2.Error())
	}
	err3 := os.MkdirAll(List, os.ModePerm)
	if err3 != nil {
		fmt.Println(err3.Error())
	}
	filename := fmt.Sprintf("%s/%s.conf", DirConf, do)
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("无法创建文件:", err)
		return
	}
	defer file.Close()

	//text := _ssl(do, dir, ssl)

	text := func() string {
		if ssl != false {
			return _touch_ssl(do, dir)
		}
		return _touch_web(do, dir)
	}()

	_, err = file.WriteString(text)
	if err != nil {
		fmt.Println("无法写入文件:", err)
		return
	}

	fmt.Println("文件创建并写入成功:", filename)

	_createlist(FILE, []byte(do))

	_mkdir_root(do)

}

// 删除项目文件夹
func _delete(del, dir string) {
	var DirData = dir + "http/" + del
	var FILE = dir + "list/list.log"
	var NEWFILE = dir + "list/newlist.log"
	//err := os.Remove(DirConf + del + ".conf")
	//if err != nil {
	//	fmt.Println("无法删除文件:", err)
	//} else {
	//	fmt.Println("文件删除成功:", DirConf+del+".conf")
	//}

	// 删除文件夹及其内容
	err := os.RemoveAll(DirData)
	if err != nil {
		fmt.Println("无法删除文件夹:", err)
	} else {
		fmt.Println("文件夹删除成功:", DirData)
	}
	_dellist(FILE, NEWFILE, del)

	_delete_root(del)
}

// 扫描文件逐行返回文件内容
func _batch(do string) []string {
	file, err := os.Open(do)
	if err != nil {
		fmt.Println("打开文件失败：", err)
	}
	defer file.Close()

	// 创建 scanner 对象，用于逐行扫描文件
	scanner := bufio.NewScanner(file)

	var content []string

	// 逐行读取文件并输出字符串
	for scanner.Scan() {

		// 检查是否有错误发生
		if err = scanner.Err(); err != nil {
			fmt.Println("扫描文件出错：", err)

		}
		content = append(content, scanner.Text())
	}
	return content
}

// 模板调用
func _touch_web(do, dir string) string {
	var DirData = dir + "http/" + do + "/data"
	touch := fmt.Sprintf(`server{
        listen 80;
        server_name %s;
  		client_max_body_size 100m;
        location /{
        	root %s;	#
        	index  index.html index.htm;
        }
		#location / {
			#if ($host = "%s"){	
        	#proxy_pass   http://127.0.0.1;
			#return 301   http://127.0.0.1;
			#}
		#}
}
`, do, DirData, do)
	return touch
}

// do=server dir=path
func _touch_ssl(do, dir string) string {
	var DirData = dir + "http/" + do + "/data/"
	var DirCert = dir + "cert"
	err1 := os.MkdirAll(DirCert, os.ModePerm)
	if err1 != nil {
		fmt.Println(err1.Error())
	}
	fmt.Println("cert路径：" + DirCert)
	touch := fmt.Sprintf(`server {
    listen 80;
    server_name %s;
    client_max_body_size 100m;
    return 301 https://$host$request_uri;
	#location / {
        	#proxy_pass   http://127.0.0.1;
			#return 301   http://127.0.0.1;
		#}
}

server {
	listen 443 ssl;
	server_name %s; 
	# SSL 证书和私钥的路径
	ssl_certificate     %s/chain1.pem;
	ssl_certificate_key %s/privkey1.pem;
	# 其他ssl配置
	ssl_session_timeout  5m;
	ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE:ECDH:AES:HIGH:!NULL:!aNULL:!MD5:!ADH:!RC4;
	ssl_protocols TLSv1 TLSv1.1 TLSv1.2;
	ssl_prefer_server_ciphers on;
	add_header X-Frame-Options DENY;
	add_header X-Content-Type-Options nosniff;
	add_header X-XSS-Protection "1; mode=block";
	add_header Strict-Transport-Security "max-age=16070400; includeSubdomains; preload";
	client_max_body_size 100m;
	proxy_set_header X-Real-IP $remote_addr;
	proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
	location / {
			root %s;
		   index  index.html index.htm;
	}
	#location / {
		#if ($host = "%s"){
		#proxy_pass   http://127.0.0.1;
		#return 301   http://127.0.0.1;
		#} 
	#}
}`, do, do, DirCert, DirCert, DirData, do)
	return touch
}

// list命令 用于测试页面
func _createlist(filename string, content []byte) {
	//file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	contentWithNewline := append(content, '\n')
	_, err1 := file.Write(contentWithNewline)
	if err1 != nil {
		log.Fatal(err1)
	}

}
func _list(dir string) {
	file, err := os.Open(dir + "list/list.log")
	if err != nil {
		fmt.Printf("无法打开文件：%s\n", err)
		return
	}
	defer file.Close()

	var domains []string

	// 从文件中读取域名
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		if domain != "" {
			domains = append(domains, domain)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("读取文件时出错：%s\n", err)
		return
	}

	var wg sync.WaitGroup

	for _, domain := range domains {
		wg.Add(1)
		go fetchURL(domain, &wg)
	}

	wg.Wait()
}
func _dellist(filename, newfile string, context string) {

	inputFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return
	}
	defer inputFile.Close()
	// 2. 创建一个新的临时文件以保存结果。
	outputFile, err := os.Create(newfile)
	if err != nil {
		fmt.Println("无法创建临时文件:", err)
		return
	}
	defer outputFile.Close()

	// 3. 逐行读取文件内容，跳过包含特定字符串的行。
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, "#") && !strings.Contains(line, context) && len(line) != 0 && line != "\r\n" {
			// 4. 将未包含特定字符串的行写回到临时文件中。
			_, err := outputFile.WriteString(line + "\n")
			if err != nil {
				fmt.Println("写入文件时出错:", err)
				return
			}
		}
	}
	// 处理扫描时的错误。
	if err := scanner.Err(); err != nil {
		fmt.Println("扫描文件时出错:", err)
		return
	}
	// 5. 关闭原始文件和新文件，并将新文件重命名为原始文件。
	inputFile.Close()
	outputFile.Close()
	err = os.Rename(newfile, filename)
	if err != nil {
		fmt.Println("无法重命名文件:", err)
		return
	}

	fmt.Println("已更新状态")

}
func fetchURL(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	// 添加HTTP协议头
	httpURL := "http://" + url
	responseHTTP, err := http.Get(httpURL)
	if err != nil {
		fmt.Printf("无法发送HTTP GET请求：%s\n", err)
		return
	}
	defer responseHTTP.Body.Close()

	// 添加HTTPS协议头
	httpsURL := "https://" + url
	responseHTTPS, err := http.Get(httpsURL)
	if err != nil {
		fmt.Printf("无法发送HTTPS GET请求：%s\n", err)
		return
	}
	defer responseHTTPS.Body.Close()

	fmt.Printf(`URL: %s    HTTP_CODE:%d    HTTPS_CODE:%d
`, url, responseHTTP.StatusCode, responseHTTPS.StatusCode)
}
