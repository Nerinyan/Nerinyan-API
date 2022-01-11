package route

import (
	"Nerinyan-API/bodyStruct"
	"Nerinyan-API/db"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func SearchByBeatmapId(c *fiber.Ctx) (err error) {

	row := db.Maria.QueryRow(`select * from osu.beatmap where beatmap_id = ?;`, c.Params("mi", ""))
	var Map bodyStruct.BeatmapOUT
	err = row.Scan(
		//beatmap_id, beatmapset_id, mode, mode_int, status, ranked, total_length, max_combo, difficulty_rating,
		//version, accuracy, ar, cs, drain, bpm, convert, count_circles, count_sliders, count_spinners, deleted_at,
		//hit_length, is_scoreable, last_updated, passcount, playcount, checksum, user_id
		&Map.Id, &Map.BeatmapsetId, &Map.Mode, &Map.ModeInt, &Map.Status, &Map.Ranked, &Map.TotalLength, &Map.MaxCombo, &Map.DifficultyRating,
		&Map.Version, &Map.Accuracy, &Map.Ar, &Map.Cs, &Map.Drain, &Map.Bpm, &Map.Convert, &Map.CountCircles, &Map.CountSliders, &Map.CountSpinners, &Map.DeletedAt,
		&Map.HitLength, &Map.IsScoreable, &Map.LastUpdated, &Map.Passcount, &Map.Playcount, &Map.Checksum, &Map.UserId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Status(http.StatusNotFound)
			return c.JSON(bodyStruct.ErrorStruct{
				Code:      "SearchByBeatmapId-002",
				Path:      c.Path(),
				RequestId: c.GetReqHeaders()["X-Request-ID"],
				Error:     fiber.ErrNotFound.Error(),
				Message:   "",
			})
		}
		c.Status(http.StatusInternalServerError)
		return c.JSON(bodyStruct.ErrorStruct{
			Code:      "SearchByBeatmapId-003",
			Path:      c.Path(),
			RequestId: c.GetReqHeaders()["X-Request-ID"],
			Error:     err.Error(),
			Message:   "RDBMS Request Error.",
		})
	}
	c.Status(http.StatusOK)
	return c.JSON(Map)
}
