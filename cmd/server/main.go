package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	// "runtime/pprof"
	"github.com/naya0000/Advertisement_Manage.git/pkg/api"
	"github.com/naya0000/Advertisement_Manage.git/pkg/db"
)

func main() {
	// 启动 Profiling 服务
	// http.ListenAndServe("localhost:6060", nil)

	// // 创建 CPU Profiling 文件
	// f, err := os.Create("cpu_profile.prof")
	// if err != nil {
	// 	log.Fatal("could not create CPU profile: ", err)
	// }
	// defer f.Close()

	// // 启动 CPU Profiling
	// if err := pprof.StartCPUProfile(f); err != nil {
	// 	log.Fatal("could not start CPU profile: ", err)
	// }
	// defer pprof.StopCPUProfile()

	log.Print("server has started")
	//start the db
	db, err := db.StartDB()
	if err != nil {
		log.Printf("error starting the database %v", err)
		panic("error starting the database")
	}
	//get the router of the API by passing the db
	router := api.StartAPI(db)
	//get the port from the environment variable
	port := os.Getenv("PORT")
	//pass the router and start listening with the server
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		log.Printf("error from router %v\n", err)
		return
	}
}
