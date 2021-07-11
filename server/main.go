package server

import (
	"fmt"
	"log"
	"net/http"
)

func cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                                // 允许访问所有域，可以换成具体url，注意仅具体url才能带cookie信息
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")                    //header的类型
		w.Header().Add("Access-Control-Allow-Credentials", "true")                        //设置为true，允许ajax异步请求带cookie信息
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE") //允许请求方法
		w.Header().Set("content-type", "application/json;charset=UTF-8")                  //返回数据格式是json
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		f(w, r)
	}
}

func main() {
	//IDEMIXPLUS handlefunc
	var err error
	http.HandleFunc("/initIssuer", cors(initIssuer))

	fmt.Println("初始化handler成功")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("err:", err)
	}
}
