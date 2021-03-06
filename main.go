package main

import (
	"Nerinyan-API/bancho"
	"Nerinyan-API/config"
	"Nerinyan-API/db"
	"Nerinyan-API/fileHandler"
	"Nerinyan-API/route"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"log"
)

var skipLog = map[string]bool{
	"/monitor": true,
}

func init() {
	config.LoadConfig()
	go fileHandler.StartIndex()
	db.ConnectMaria()
	ch := make(chan struct{})
	go bancho.LoadBancho(ch)
	_ = <-ch
}
func main() {
	f := fiber.New(fiber.Config{
		Prefork:     false,
		IdleTimeout: 60,
	})
	f.Use(requestid.New())
	f.Use(logger.New(logger.Config{
		Next: func(c *fiber.Ctx) bool {
			return skipLog[c.Path()]
		},
		TimeZone: "asia/Seoul",
	}))

	f.Get("/monitor", monitor.New())

	// 로드벨런서.========================================================================================================
	//f.Get("/d/:id", route.BeatmapDownloadServerLoadBalance)
	// docs ============================================================================================================
	f.Get("/", nil)

	// 서버상태 체크용 ====================================================================================================
	f.Get("/health", route.Health)
	f.Get("/robots.txt", route.Robots)
	// 맵 파일 다운로드 ===================================================================================================
	f.Get("/d/:id", route.DownloadBeatmapSet)

	// 비트맵 리스트 검색용 ================================================================================================
	f.Get("/search", route.Search)
	f.Get("/search/beatmap/:mi", route.SearchByBeatmapId)
	f.Get("/search/beatmapset/:si", route.SearchByBeatmapSetId)

	log.Fatal(f.Listen(":" + config.Config.Port))

}
