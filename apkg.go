package apkg

//github.com/u-on/apkg
import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
)

// GetDir 获取程序自身所在目录 不含\
func GetDir() string {
	dir, _ := os.Executable()
	exPath := filepath.Dir(dir)
	return exPath
}

// PauseExit Press Ctrl+c to exit
func PauseExit() {
	fmt.Printf("Press Ctrl+c to exit...")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	for {
		if sig.String() == "interrupt" {
			break
		}
	}
	return

}

// PathExists 判断所给路径文件/文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	//nonexistent来判断，是不是不存在的错误
	if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
		return false, nil
	}
	return false, err //如果有错误了，但是不是不存在的错误，所以把这个错误原封不动的返回
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {

		return false
	}
	return s.IsDir()

}

// IsFile 判断所给路径是否为文件
func IsFile(path string) bool {

	return !IsDir(path)

}

type Gtihub struct{}

func (Gtihub) GetReleases(gurl string) gjson.Result {

	req, err := http.NewRequest("GET", gurl, nil)
	if err != nil {
		fmt.Println("[ERROR]: " + err.Error())
		os.Exit(1)
	}
	//req.Header.Add("Basic", "Z2hwX1JLeWZ4Ym9ERmVFZ09JS1ZPcGxIdFJSU2s1ZnFEejJJRlpBYjo=")
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	result := gjson.Get(string(body), "assets.#.browser_download_url")
	return result
}

func (Gtihub) GetReleasesEx(gurl string, Ex string) string {

	req, err := http.NewRequest("GET", gurl, nil)
	if err != nil {
		fmt.Println("[ERROR]: " + err.Error())
		os.Exit(1)
	}
	//req.Header.Add("Basic", "Z2hwX1JLeWZ4Ym9ERmVFZ09JS1ZPcGxIdFJSU2s1ZnFEejJJRlpBYjo=")
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	result := gjson.Get(string(body), "assets.#.browser_download_url")

	for _, name := range result.Array() {
		result, err := regexp.MatchString(Ex, name.String())
		if err != nil {
			fmt.Println(err.Error())
		}
		if result {
			return name.String()
			//fmt.Printf("%s matches\n", tt)
		}
	}
	return ""
}

// Download 下载文件
func Download(url string, savepath string) {
	var d string
	if strings.Contains(savepath, "\\") {
		d = savepath[:strings.LastIndex(savepath, "\\")]
	} else if strings.Contains(savepath, "/") {
		d = savepath[:strings.LastIndex(savepath, "/")]
	}

	if d != "" {
		d2, _ := PathExists(d)
		if d2 == false {
			err := os.Mkdir(d, 0666)
			if err != nil {
				fmt.Println(err)
			}

			//fmt.Println("dir '" + d + "' is not exist!")
		}

	}
	// progressbar
	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	f, _ := os.OpenFile(savepath, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)

	//io.Copy(io.MultiWriter(f, bar), resp.Body)

	if _, err := io.Copy(io.MultiWriter(f, bar), resp.Body); err != nil {
		fmt.Println(err)
	}

}
