package api_db

import (
	"diplom-vuls-server/internal/model"
	custom_errors "diplom-vuls-server/pkg/custom-errors"
	"github.com/restream/reindexer/v3"
	"github.com/valyala/fasthttp"
)

func CheckExistingPlatformDB(item *model.PlatformItem, db *reindexer.Reindexer) bool {
	_, found := db.Query("platform").WhereString("name", reindexer.EQ, item.Name).
		WhereInt64("software_program_id", reindexer.EQ, item.SoftwareProgramID).Get()

	return found
}

func CheckPlatformDB(platform string, programID int64, db *reindexer.Reindexer) (int64, *custom_errors.ErrHttp) {
	rec, found := db.Query("platform").WhereInt64("software_program_id", reindexer.EQ, programID).
		Match("search", platform).Get()
	if !found {
		inserted, err := db.Insert("platform", &model.PlatformItem{
			ID:                1,
			Name:              platform,
			SoftwareProgramID: programID,
		}, "id=serial()")
		if err != nil {
			return -1, custom_errors.New(fasthttp.StatusInternalServerError, "insert platform: "+err.Error())
		}

		if inserted == 0 {
			return -1, custom_errors.New(fasthttp.StatusInternalServerError, "insert platform: something went wrong")
		}

		rec, found = db.Query("platform").WhereInt64("software_program_id", reindexer.EQ, programID).
			WhereString("name", reindexer.EQ, platform).Get()
		if !found {
			return -1, custom_errors.New(fasthttp.StatusNotFound, "platform with this name doesn't exist")
		}
	}

	return rec.(*model.PlatformItem).ID, nil
}

func CreatePlatformDB(item *model.PlatformItem, db *reindexer.Reindexer) *custom_errors.ErrHttp {
	inserted, err := db.Insert("platform", item, "id=serial()")
	if err != nil {
		return custom_errors.New(fasthttp.StatusInternalServerError, "insert platform: "+err.Error())
	}

	if inserted == 0 {
		return custom_errors.New(fasthttp.StatusInternalServerError, "insert platform: something went wrong")
	}

	return nil
}
