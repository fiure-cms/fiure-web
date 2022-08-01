package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fiure-cms/fiure-web/internal/fcore"
	"github.com/fiure-cms/fiure-web/loggers"
	"github.com/fiure-cms/fiure-web/middlewares"
	"github.com/fiure-cms/fiure-web/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
	"github.com/peterbourgon/ff/v3"
)

var ServiceBuild string
var ServiceCommitId string

var ServiceName string = "buwebdefault"
var ServiceVersion string = "1.0.0"

var (
	fs = flag.NewFlagSet(ServiceName, flag.ExitOnError)

	domainName                                                           = fs.String("domain", "localhost:3001", "domain name")
	httpAddr                                                             = fs.String("http-addr", "localhost:3001", "http server")
	httpPrefork                                                          = fs.Bool("http-prefork", false, "http server prefork option")
	cdnHost                                                              = fs.String("cdn-host", "localhost:3005", "cdn hostname")
	redisClients, redisAddrs                                             fcore.ArrayFlagString
	redisPoolSizes, redisMinIdleConns, redisMaxRetries                   fcore.ArrayFlagInt
	sonicsearchClients, sonicsearchAddrs                                 fcore.ArrayFlagString
	sonicsearchChannelMode, sonicsearchAuthPassword                      fcore.ArrayFlagString
	sonicsearchPoolSizes, sonicsearchMinIdleConns, sonicsearchMaxRetries fcore.ArrayFlagInt
	_                                                                    = fs.String("env-file", ".env", "env file")
)

func main() {

	// Zap Logger Init
	loggers.SetupSugarLogger(ServiceName, ServiceVersion)

	// Start Logger: Haydaaa
	loggers.Sugar.Info("service started")

	// Flag Parse with Env
	fs.Var(&redisClients, "redis-clients", "redis client names")                               // Multiple
	fs.Var(&redisAddrs, "redis-addrs", "redis addrs for sync data to storage service")         // Multiple
	fs.Var(&redisPoolSizes, "redis-pool-sizes", "redis pool size config")                      // Multiple
	fs.Var(&redisMinIdleConns, "redis-min-idle-conns", "redis minimum idle conn count config") // Multiple
	fs.Var(&redisMaxRetries, "redis-max-retries", "redis max retries to connect")              // Multiple

	fs.Var(&sonicsearchClients, "sonicsearch-clients", "sonic client names")                               // Multiple
	fs.Var(&sonicsearchAddrs, "sonicsearch-addrs", "sonic addrs for sync data to storage service")         // Multiple
	fs.Var(&sonicsearchChannelMode, "sonicsearch-channelmode", "sonic addrs for working mode")             // Multiple
	fs.Var(&sonicsearchAuthPassword, "sonicsearch-authpassword", "sonic addrs for auth password")          // Multiple
	fs.Var(&sonicsearchPoolSizes, "sonicsearch-pool-sizes", "sonic pool size config")                      // Multiple
	fs.Var(&sonicsearchMinIdleConns, "sonicsearch-min-idle-conns", "sonic minimum idle conn count config") // Multiple
	fs.Var(&sonicsearchMaxRetries, "sonicsearch-max-retries", "sonic max retries to connect")              // Multiple

	err := ff.Parse(
		fs, os.Args[1:],
		ff.WithConfigFileFlag("env-file"),
		ff.WithConfigFileParser(ff.PlainParser),
		ff.WithEnvVarPrefix(strings.ToUpper(ServiceName)),
	)

	if err != nil {
		loggers.Sugar.With("error", err).Fatal("configration error")
	}

	// RedisClientRing Conn
	services.SetupRedisClientRing(redisClients, redisAddrs, redisPoolSizes, redisMinIdleConns, redisMaxRetries)

	// SonicSearch Client
	services.SetupSonicSearchClient(ServiceName, redisClients, sonicsearchAddrs, sonicsearchChannelMode, sonicsearchAuthPassword, sonicsearchPoolSizes, sonicsearchMinIdleConns, sonicsearchMaxRetries)

	// Managers Init
	services.SetupManagers()

	// Fiber Init
	engine := html.New("./templates", ".tmpl")
	engine.AddFunc("cdnUrl", func(filePath string) string {
		if strings.Contains(filePath, "http") {
			return filePath
		}

		return fmt.Sprintf("%s%s", *cdnHost, filePath)
	}).AddFunc("serviceVersion", func() string {
		return ServiceVersion
	}).AddFunc("add", func(a int, b int) int {
		return a + b
	}).AddFunc("canonicalUrl", func(ctype string, cpath string) string {
		var curl []string

		curl = append(curl, *domainName)

		if ctype == "page" {
			curl = append(curl, "p", cpath)
		} else if ctype == "blog" {
			curl = append(curl, "blog", cpath)
		} else if ctype == "search" {
			curl = append(curl, "s", cpath)
		} else if ctype == "mylist" {
			curl = append(curl, cpath)
		} else if ctype == "term" {
			curl = append(curl, cpath)
		} else if ctype == "single" {
			curl = append(curl, cpath)
		}

		return strings.Join(curl, "/")
	}).AddFunc("truncateText", func(content string, limitter int) string {
		return fcore.GetTruncateText(content, limitter)
	}).AddFunc("convertUpdateTime", func(updated int64) string {
		return time.Unix(updated, 0).Format("02-01-2006")
	}).AddFunc("option", func(key string) string {
		option := map[string]string{
			"twitter":       "bakbibu",
			"facebook":      "profile.php?id=",
			"instagram":     "bakbibu",
			"footer_slogan": "Bakbibu is the new search engine of quick finding millions of best 3D printed models.",
		}
		return option[key]
	}).AddFunc("unescape", func(value string) template.HTML {
		return template.HTML(value)
	})

	app := fiber.New(fiber.Config{
		Prefork:       *httpPrefork,
		CaseSensitive: true,
		AppName:       ServiceName + " v" + ServiceVersion,
		Views:         engine,
		ViewsLayout:   "layouts/base",
		ErrorHandler:  middlewares.ErrorHandler,
	})

	// Assets Files
	app.Static("/assets", "./assets")

	// Middlewares
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1: Best
	}))
	app.Use(middlewares.HeaderConf(ServiceVersion))
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format: "[${ip}]:${port} ${status} - ${latency} ${method} ${path} ${queryParams}\n",
	}))
	app.Use(recover.New()) // Sometime need to recover all

	// Routes
	app.Use(favicon.New())
	setupRoutes(app)

	// 404 Middleware
	app.Use(middlewares.NotFound())

	go func() {
		err = app.Listen(*httpAddr)
		if err != nil {
			loggers.Sugar.With("error", err).Fatal("http server stopped")
		}
	}()

	// Listen server quit or something happened and notify channel
	close := make(chan os.Signal, 1)
	signal.Notify(close, syscall.SIGINT, syscall.SIGTERM)

	<-close

	// Bye bye
	loggers.Sugar.Info("im shutting down. see you later")
}
