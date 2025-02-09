package api_db

import (
	"diplom-vuls-server/internal/model"
	custom_errors "diplom-vuls-server/pkg/custom-errors"
	"github.com/restream/reindexer/v3"
	"github.com/valyala/fasthttp"
)

func CheckExistingVersionDB(item *model.VersionItem, db *reindexer.Reindexer) bool {
	_, found := db.Query("version").WhereString("name", reindexer.EQ, item.Name).
		WhereInt64("software_program_id", reindexer.EQ, item.SoftwareProgramID).Get()

	return found
}

func CheckVersionDB(version string, programID int64, db *reindexer.Reindexer) (int64, *custom_errors.ErrHttp) {
	rec, found := db.Query("version").WhereInt64("software_program_id", reindexer.EQ, programID).
		Match("search", version).Get()
	if !found {
		inserted, err := db.Insert("version", &model.VersionItem{
			ID:                1,
			Name:              version,
			SoftwareProgramID: programID,
		}, "id=serial()")
		if err != nil {
			return -1, custom_errors.New(fasthttp.StatusInternalServerError, "insert version: "+err.Error())
		}

		if inserted == 0 {
			return -1, custom_errors.New(fasthttp.StatusInternalServerError, "insert version: something went wrong")
		}

		rec, found = db.Query("version").WhereInt64("software_program_id", reindexer.EQ, programID).
			WhereString("name", reindexer.EQ, version).Get()
		if !found {
			return -1, custom_errors.New(fasthttp.StatusNotFound, "software_program with this name doesn't exist")
		}
	}

	return rec.(*model.VersionItem).ID, nil
}

func CreateVersionDB(item *model.VersionItem, db *reindexer.Reindexer) *custom_errors.ErrHttp {
	inserted, err := db.Insert("version", item, "id=serial()")
	if err != nil {
		return custom_errors.New(fasthttp.StatusInternalServerError, "insert version: "+err.Error())
	}

	if inserted == 0 {
		return custom_errors.New(fasthttp.StatusInternalServerError, "insert program: something went wrong")
	}

	return nil
}
