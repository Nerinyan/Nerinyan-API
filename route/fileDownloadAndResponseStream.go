package route

import (
	"Nerinyan-API/bodyStruct"
	"Nerinyan-API/config"
	"Nerinyan-API/db"
	"Nerinyan-API/fileHandler"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/pterm/pterm"
	"io"
	"net/http"
	"strconv"
	"time"
)

func DownloadBeatmapSet(c *fiber.Ctx) (err error) {

	//1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
	noVideo, _ := strconv.ParseBool(c.Query("noVideo"))
	noVideo2, _ := strconv.ParseBool(c.Query("nv"))
	noVideo = noVideo || noVideo2

	mid, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(bodyStruct.ErrorStruct{
			Code:      "DownloadBeatmapSet-001",
			Path:      c.Path(),
			RequestId: c.GetReqHeaders()["X-Request-ID"],
			Error:     err.Error(),
			Message:   "cannot parse beatmap set id",
		})
	}

	//go src.ManualUpdateBeatmapSet(mid)

	row := db.Maria.QueryRow(`SELECT beatmapset_id,artist,title,last_updated,video FROM osu.beatmapset WHERE beatmapset_id = ?`, mid)

	var a struct {
		Id          string
		Artist      string
		Title       string
		LastUpdated string
		Video       bool
	}

	if err = row.Scan(&a.Id, &a.Artist, &a.Title, &a.LastUpdated, &a.Video); err != nil {
		if err == sql.ErrNoRows {
			c.Status(http.StatusNotFound)
			return c.JSON(bodyStruct.ErrorStruct{
				Code:      "DownloadBeatmapSet-002",
				Path:      c.Path(),
				RequestId: c.GetReqHeaders()["X-Request-ID"],
				Error:     fiber.ErrNotFound.Error(),
				Message:   "",
			})
		}
		c.Status(http.StatusInternalServerError)
		return c.JSON(bodyStruct.ErrorStruct{
			Code:      "DownloadBeatmapSet-003",
			Path:      c.Path(),
			RequestId: c.GetReqHeaders()["X-Request-ID"],
			Error:     err.Error(),
			Message:   "RDBMS Request Error.",
		})
	}

	lu, err := time.Parse("2006-01-02 15:04:05", a.LastUpdated)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Status(http.StatusNotFound)
			return c.JSON(bodyStruct.ErrorStruct{
				Code:      "DownloadBeatmapSet-004",
				Path:      c.Path(),
				RequestId: c.GetReqHeaders()["X-Request-ID"],
				Error:     fiber.ErrNotFound.Error(),
				Message:   "",
			})
		}
		c.Status(http.StatusInternalServerError)
		return c.JSON(bodyStruct.ErrorStruct{
			Code:      "DownloadBeatmapSet-005",
			Path:      c.Path(),
			RequestId: c.GetReqHeaders()["X-Request-ID"],
			Error:     err.Error(),
			Message:   "RDBMS Request Error.",
		})
	}

	url := fmt.Sprintf("https://osu.ppy.sh/api/v2/beatmapsets/%d/download", mid)
	if a.Video && noVideo {
		mid *= -1
		a.Title += " [no video]"
		url += "?noVideo=1"
	}
	conf := config.Config
	serverFileName := fmt.Sprintf("%s/%d.osz", conf.TargetDir, mid)

	if fileHandler.FileList[mid].Unix() >= lu.Unix() { // 맵이 최신인경우
		c.Response().Header.Set("Content-Type", "application/x-osu-beatmap-archive")
		c.Attachment(serverFileName, fmt.Sprintf("%s %s - %s.osz", a.Id, a.Artist, a.Title))
		return
	}

	//==========================================
	//=        비트맵 파일이 서버에 없는경우        =
	//==========================================

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(bodyStruct.ErrorStruct{
			Code:      "DownloadBeatmapSet-006",
			Path:      c.Path(),
			RequestId: c.GetReqHeaders()["X-Request-ID"],
			Error:     err.Error(),
			Message:   "http request build Error.",
		})
	}
	req.Header.Add("Authorization", conf.Osu.Token.TokenType+" "+conf.Osu.Token.AccessToken)

	res, err := client.Do(req)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(bodyStruct.ErrorStruct{
			Code:      "DownloadBeatmapSet-007",
			Path:      c.Path(),
			RequestId: c.GetReqHeaders()["X-Request-ID"],
			Error:     err.Error(),
			Message:   "http request Error.",
		})
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		c.Status(http.StatusNotFound)
		return c.JSON(bodyStruct.ErrorStruct{
			Code:      "DownloadBeatmapSet-008",
			Path:      c.Path(),
			RequestId: c.GetReqHeaders()["X-Request-ID"],
			Error:     res.Status,
			Message:   "Please try again in a few seconds. OR map is not alive. check beatmapset id.",
		})
	}
	pterm.Info.Println("beatmapSet Downloading at", serverFileName)

	cLen, _ := strconv.Atoi(res.Header.Get("Content-Length"))
	c.Response().Header.Set("Content-Length", res.Header.Get("Content-Length"))
	c.Response().Header.Set("Content-Disposition", res.Header.Get("Content-Disposition"))
	c.Response().Header.Set("Content-Type", "application/x-osu-beatmap-archive")
	c.Response().Header.Set("Connection", "Keep-Alive")
	c.Response().Header.Set("Keep-Alive", "timeout=20")
	c.Response().Header.Set("Vary", "Origin")
	c.Response().Header.Set("Strict-Transport-Security", "max-age=15768000; includeSubdomains; preload")
	var buf bytes.Buffer

	for i := 0; i < cLen; { // 읽을 데이터 사이즈 체크
		var b = make([]byte, 64000) // 바이트 배열
		n, err := res.Body.Read(b)  // 반쵸 스트림에서 64k 읽어서 바이트 배열 b 에 넣음

		i += n // 현재까지 읽은 바이트
		if n > 0 {
			buf.Write(b[:n]) // 서버에 저장할 파일 버퍼에 쓴다
			//c.Response().SetBodyRaw(b[:n])
			//c.Response().AppendBody(b[:n])
			//fmt.Println(c.Send(b[:n]))
			if _, err := c.Response().BodyWriter().Write(b[:n]); err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(bodyStruct.ErrorStruct{
					Code:      "DownloadBeatmapSet-009",
					Path:      c.Path(),
					RequestId: c.GetReqHeaders()["X-Request-ID"],
					Error:     err.Error(),
					Message:   "fail to write body stream.",
				})
			}

		}

		if err == io.EOF {
			break
		} else if err != nil { //에러처리
			fmt.Println(err.Error())
			return err
		}
	}
	c.Response().SetConnectionClose()
	if cLen == buf.Len() {
		//go saveLocal(&buf, serverFileName, mid)
		return
	}
	errMsg := fmt.Sprintf("filesize not match: bancho response bytes : %d | downloaded bytes : %d", cLen, buf.Len())
	pterm.Error.Printfln(errMsg)
	return errors.New(errMsg)

}

//func saveLocal(data *bytes.Buffer, path string, id int)  {
//	tmp := path + ".down"
//	file, err := os.Create(tmp)
//	if err != nil {
//		return
//	}
//	if file == nil {
//		return errors.New("")
//	}
//	_, err = file.Write(data.Bytes())
//	if err != nil {
//		return
//	}
//	file.Close()
//
//	if _, err = os.Stat(path); !os.IsNotExist(err) {
//		err = os.Remove(path)
//		if err != nil {
//			return
//		}
//	}
//	err = os.Rename(tmp, path)
//	if err != nil {
//		return
//	}
//
//	fileHandler.FileList[id] = time.Now()
//	pterm.Info.Println("beatmapSet Downloading Finished", path)
//	return
//}
