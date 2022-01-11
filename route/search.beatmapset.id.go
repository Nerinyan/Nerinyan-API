package route

import (
	"Nerinyan-API/bodyStruct"
	"Nerinyan-API/db"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strings"
)

func SearchByBeatmapSetId(c *fiber.Ctx) (err error) {
	row := db.Maria.QueryRow(`select * from osu.beatmapset where beatmapset_id = ?;`, c.Params("si", ""))
	var set bodyStruct.BeatmapSetsOUT

	var mapids []int

	err = row.Scan(
		// beatmapset_id, artist, artist_unicode, creator, favourite_count, hype_current,
		//hype_required, nsfw, play_count, source, status, title, title_unicode, user_id,
		//video, availability_download_disabled, availability_more_information, bpm, can_be_hyped,
		//discussion_enabled, discussion_locked, is_scoreable, last_updated, legacy_thread_url,
		//nominations_summary_current, nominations_summary_required, ranked, ranked_date, storyboard,
		//submitted_date, tags, has_favourited, description, genre_id, genre_name, language_id, language_name, ratings

		&set.Id, &set.Artist, &set.ArtistUnicode, &set.Creator, &set.FavouriteCount, &set.Hype.Current, &set.Hype.Required, &set.Nsfw, &set.PlayCount, &set.Source, &set.Status, &set.Title, &set.TitleUnicode, &set.UserId, &set.Video, &set.Availability.DownloadDisabled, &set.Availability.MoreInformation, &set.Bpm, &set.CanBeHyped, &set.DiscussionEnabled, &set.DiscussionLocked, &set.IsScoreable, &set.LastUpdated, &set.LegacyThreadUrl, &set.NominationsSummary.Current, &set.NominationsSummary.Required, &set.Ranked, &set.RankedDate, &set.Storyboard, &set.SubmittedDate, &set.Tags, &set.HasFavourited, &set.Description.Description, &set.Genre.Id, &set.Genre.Name, &set.Language.Id, &set.Language.Name, &set.RatingsString)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Status(http.StatusNotFound)
			return c.JSON(bodyStruct.ErrorStruct{
				Code:      "SearchByBeatmapSetId-002",
				Path:      c.Path(),
				RequestId: c.GetReqHeaders()["X-Request-ID"],
				Error:     fiber.ErrNotFound.Error(),
				Message:   "",
			})
		}
		c.Status(http.StatusInternalServerError)
		return c.JSON(bodyStruct.ErrorStruct{
			Code:      "SearchByBeatmapSetId-003",
			Path:      c.Path(),
			RequestId: c.GetReqHeaders()["X-Request-ID"],
			Error:     err.Error(),
			Message:   "RDBMS Request Error.",
		})
	}
	mapids = append(mapids, *set.Id)

	if *set.Id == 0 {
		c.Status(http.StatusNotFound)
		return c.JSON(bodyStruct.ErrorStruct{
			Code:      "SearchByBeatmapSetId-004",
			Path:      c.Path(),
			RequestId: c.GetReqHeaders()["X-Request-ID"],
			Error:     fiber.ErrNotFound.Error(),
			Message:   "",
		})
	}

	rows, err := db.Maria.Query(fmt.Sprintf(`select * from osu.beatmap where beatmapset_id in( %s ) order by difficulty_rating asc;`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(mapids)), ", "), "[]")))

	if err != nil {
		if err == sql.ErrNoRows {
			c.Status(http.StatusNotFound)
			return c.JSON(bodyStruct.ErrorStruct{
				Code:      "SearchByBeatmapSetId-005",
				Path:      c.Path(),
				RequestId: c.GetReqHeaders()["X-Request-ID"],
				Error:     fiber.ErrNotFound.Error(),
				Message:   "",
			})
		}
		c.Status(http.StatusInternalServerError)
		return c.JSON(bodyStruct.ErrorStruct{
			Code:      "SearchByBeatmapSetId-006",
			Path:      c.Path(),
			RequestId: c.GetReqHeaders()["X-Request-ID"],
			Error:     err.Error(),
			Message:   "RDBMS Request Error.",
		})
	}
	defer rows.Close()

	for rows.Next() {
		var Map bodyStruct.BeatmapOUT
		err = rows.Scan(
			//beatmap_id, beatmapset_id, mode, mode_int, status, ranked, total_length, max_combo, difficulty_rating,
			//version, accuracy, ar, cs, drain, bpm, convert, count_circles, count_sliders, count_spinners, deleted_at,
			//hit_length, is_scoreable, last_updated, passcount, playcount, checksum, user_id
			&Map.Id, &Map.BeatmapsetId, &Map.Mode, &Map.ModeInt, &Map.Status, &Map.Ranked, &Map.TotalLength, &Map.MaxCombo, &Map.DifficultyRating, &Map.Version, &Map.Accuracy, &Map.Ar, &Map.Cs, &Map.Drain, &Map.Bpm, &Map.Convert, &Map.CountCircles, &Map.CountSliders, &Map.CountSpinners, &Map.DeletedAt, &Map.HitLength, &Map.IsScoreable, &Map.LastUpdated, &Map.Passcount, &Map.Playcount, &Map.Checksum, &Map.UserId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(http.StatusNotFound)
				return c.JSON(bodyStruct.ErrorStruct{
					Code:      "SearchByBeatmapSetId-007",
					Path:      c.Path(),
					RequestId: c.GetReqHeaders()["X-Request-ID"],
					Error:     fiber.ErrNotFound.Error(),
					Message:   "",
				})
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(bodyStruct.ErrorStruct{
				Code:      "SearchByBeatmapSetId-008",
				Path:      c.Path(),
				RequestId: c.GetReqHeaders()["X-Request-ID"],
				Error:     err.Error(),
				Message:   "RDBMS Request Error.",
			})
		}
		set.Beatmaps = append(set.Beatmaps, Map)

	}

	c.Status(http.StatusOK)
	return c.JSON(set)
}
