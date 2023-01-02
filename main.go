package main

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

var IpHub map[string][]int64

func webpage(w http.ResponseWriter, r *http.Request) {
	// 取得IP
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	// 取得當下timestamp
	now := time.Now().Unix()
	// 判斷IP是否已被紀錄
	if _, ok := IpHub[ip]; ok {
		// 清除60秒前紀錄
		find := false
		for i, v := range IpHub[ip] {
			if v >= now-60 {
				IpHub[ip] = IpHub[ip][i:]
				find = true
				break
			}
		}
		if !find {
			IpHub[ip] = IpHub[ip][:0]
		}
		// 檢查當前長度, 超過顯示Error
		if len(IpHub[ip]) >= 60 {
			fmt.Fprintf(w, "Error")
			return
		} else { // 長度未超過, append
			IpHub[ip] = append(IpHub[ip], now)
			fmt.Fprintf(w, strconv.Itoa(len(IpHub[ip])))
			return
		}
		// 不存在的IP新增
	} else {
		IpHub[ip] = []int64{now}
		fmt.Fprintf(w, "1")
		return
	}

}

func main() {
	IpHub = make(map[string][]int64)
	http.HandleFunc("/", webpage)
	http.ListenAndServe(":8888", nil)
}
