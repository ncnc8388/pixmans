package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"Golang/api"
	"Golang/api/yqk"
	"Golang/liveurls"
	"Golang/utils"
)

// liveHandler 处理 /live/ 前缀的旧格式路由，本地服务使用。
func liveHandler(w http.ResponseWriter, r *http.Request) {
	adurl := "https://cdn.jsdelivr.net/gh/feiyangdigital/testvideo/sdr1080pvideo/index.m3u8"
	path := r.URL.Path
	params := strings.Split(path, "/")

	enableTV := os.Getenv("TV") != "false"

	if len(params) >= 4 {
		platform := params[2]
		rid := params[3]
		ts := utils.DefaultQuery(r, "ts", "")
		switch platform {
		case "itv":
			if enableTV {
				itvobj := &liveurls.Itv{}
				cdn := utils.DefaultQuery(r, "cdn", "")
				if ts == "" {
					itvobj.HandleMainRequest(w, r, cdn, rid)
				} else {
					itvobj.HandleTsRequest(w, r)
				}
			} else {
				http.Error(w, "公共服务不提供TV直播", http.StatusForbidden)
			}
		case "ysptp":
			if enableTV {
				ysptpobj := &liveurls.Ysptp{}
				if ts == "" {
					ysptpobj.HandleMainRequest(w, r, rid)
				} else {
					ysptpobj.HandleTsRequest(w, ts, utils.DefaultQuery(r, "wsTime", ""))
				}
			} else {
				http.Error(w, "公共服务不提供TV直播", http.StatusForbidden)
			}
		case "douyin":
			douyinobj := &liveurls.Douyin{}
			douyinobj.Rid = rid
			douyinobj.Stream = utils.DefaultQuery(r, "stream", "flv")
			http.Redirect(w, r, utils.Duanyan(adurl, douyinobj.GetDouYinUrl()), http.StatusMovedPermanently)
		case "douyu":
			douyuobj := &liveurls.Douyu{}
			douyuobj.Rid = rid
			douyuobj.Stream_type = utils.DefaultQuery(r, "stream", "flv")
			http.Redirect(w, r, utils.Duanyan(adurl, douyuobj.GetRealUrl()), http.StatusMovedPermanently)
		case "huya":
			huyaobj := &liveurls.Huya{}
			huyaobj.Rid = rid
			huyaobj.Cdn = utils.DefaultQuery(r, "cdn", "hwcdn")
			huyaobj.Media = utils.DefaultQuery(r, "media", "flv")
			huyaobj.Type = utils.DefaultQuery(r, "cdntype", "nodisplay")
			if huyaobj.Type == "display" {
				fmt.Fprintf(w, huyaobj.GetLiveUrl().(string))
			} else {
				http.Redirect(w, r, utils.Duanyan(adurl, huyaobj.GetLiveUrl()), http.StatusMovedPermanently)
			}
		case "bilibili":
			biliobj := &liveurls.BiliBili{}
			biliobj.Rid = rid
			biliobj.Platform = utils.DefaultQuery(r, "platform", "web")
			biliobj.Quality = utils.DefaultQuery(r, "quality", "10000")
			biliobj.Line = utils.DefaultQuery(r, "line", "first")
			http.Redirect(w, r, utils.Duanyan(adurl, biliobj.GetPlayUrl()), http.StatusMovedPermanently)
		case "youtube":
			ytbObj := &liveurls.Youtube{}
			ytbObj.Rid = rid
			ytbObj.Quality = utils.DefaultQuery(r, "quality", "1080")
			http.Redirect(w, r, utils.Duanyan(adurl, ytbObj.GetLiveUrl()), http.StatusMovedPermanently)
		case "yy":
			yyObj := &liveurls.Yy{}
			yyObj.Rid = rid
			yyObj.Quality = utils.DefaultQuery(r, "quality", "4")
			http.Redirect(w, r, utils.Duanyan(adurl, yyObj.GetLiveUrl()), http.StatusMovedPermanently)
		default:
			fmt.Fprintf(w, "Unknown platform=%s, room=%s", platform, rid)
		}
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/yqk/", func(w http.ResponseWriter, r *http.Request) {
		yqk.Handler(w, r)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/live/") {
			liveHandler(w, r)
			return
		}
		api.Handler(w, r)
	})

	addr := ":" + port
	fmt.Printf("Server listening on http://0.0.0.0%s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}