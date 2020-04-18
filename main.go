/*
 * @Description: In User Settings Edit
 * @Author: x
 * @Date: 2019-09-10 12:24:08
 * @LastEditTime: 2019-09-10 15:24:21
 * @LastEditors: Please set LastEditors
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"time"
	"weagentweb/gamecfg"
	"weagentweb/server"
	"weagentweb/servercfg"

	_ "weagentweb/handles/auth"
	_ "weagentweb/handles/file"
	_ "weagentweb/handles/money"

	log "github.com/sirupsen/logrus"
)

var (
	cfgFilepath    = ""
	redisServerURL = ""
	serverUUID     = ""
	genInPath      = ""
	genOutPath     = ""
)

var (
	// 显示版本
	showVer = false
	// BuildVersion 编译版本
	BuildVersion string
	//BuildTime 编译时间
	BuildTime string
	//CommitID 提供ID
	CommitID string
)

func init() {
	flag.StringVar(&cfgFilepath, "c", "servercfg/x.json", "specify the config file path name")
	flag.StringVar(&serverUUID, "u", "", "specify the server UUID")
	flag.StringVar(&redisServerURL, "r", "", "redis server address")
	flag.StringVar(&genInPath, "gi", "", "input path")
	flag.StringVar(&genOutPath, "go", "", "output path")
}

func main() {
	// only one thread
	runtime.GOMAXPROCS(1)

	flag.BoolVar(&showVer, "v", false, "show version")

	flag.Parse()

	if showVer {
		fmt.Println("Build Version:", BuildVersion)
		fmt.Println("Build Time:", BuildTime)
		fmt.Println("CommitID:", CommitID)
		os.Exit(0)
	}

	if genInPath != "" && genOutPath != "" {
		gamecfg.Gen(genInPath, genOutPath)
		os.Exit(0)
	}

	if redisServerURL != "" {
		servercfg.RedisServer = redisServerURL
	}

	if serverUUID != "" {
		serverInt, err := strconv.Atoi(serverUUID)
		if err != nil {
			log.Fatal("serverUUID must be integer")
		}
		servercfg.ServerID = serverInt
	}

	if cfgFilepath == "" {
		// 如果没有配置json文件，则必须提供uuid以及redis地址
		if serverUUID == "" || redisServerURL == "" {
			log.Fatal("must provide redis and uuid when json config file is omit")
		}
	}

	if cfgFilepath != "" {
		r := servercfg.ParseConfigFile(cfgFilepath)
		if r != true {
			log.Fatal("can't parse configure file:", cfgFilepath)
		}
	} else {
		log.Fatal("please specify a valid config file path")
	}

	// server startTime
	servercfg.StartTime = int(time.Now().Unix())

	// load game config
	gamecfg.LoadAll(servercfg.GameCfgsDir)

	log.Println("start server")

	// start cron job
	server.CronJob()

	// start http server
	server.CreateHTTPServer()
	log.Println("start server finish...")

	select {}
}

func dumpGoRoutinesInfo() {
	log.Println("current goroutine count:", runtime.NumGoroutine())
	// use DEBUG=2, to dump stack like golang dying due to an unrecovered panic.
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 2)
}
