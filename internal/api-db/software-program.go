package api_db

import (
	"diplom-vuls-server/internal/model"
	custom_errors "diplom-vuls-server/pkg/custom-errors"
	"github.com/restream/reindexer/v3"
	"github.com/valyala/fasthttp"
	"strings"
)

func CheckExistingProgramDB(item *model.SoftwareProgramItem, db *reindexer.Reindexer) bool {
	_, found := db.Query("software_program").WhereString("name", reindexer.EQ, item.Name).Get()

	return found
}

func CheckProgramDB(program string, db *reindexer.Reindexer) (int64, *custom_errors.ErrHttp) {
	programArr := strings.Split(program, " ")
	for i := len(programArr) - 1; i > -1; i-- {
		rec, found := db.Query("software_program").Match("search", strings.Join(programArr[:i], " ")).Get()
		if !found {
			continue
		}

		return rec.(*model.SoftwareProgramItem).ID, nil
	}

	inserted, err := db.Insert("software_program", &model.SoftwareProgramItem{
		ID:   1,
		Name: program,
	}, "id=serial()")
	if err != nil {
		return -1, custom_errors.New(fasthttp.StatusInternalServerError, "insert software_program: "+err.Error())
	}

	if inserted == 0 {
		return -1, custom_errors.New(fasthttp.StatusInternalServerError, "insert software_program: something went wrong")
	}

	rec, found := db.Query("software_program").WhereString("name", reindexer.EQ, program).Get()
	if !found {
		return -1, custom_errors.New(fasthttp.StatusNotFound, "software_program with this name doesn't exist")
	}

	return rec.(*model.SoftwareProgramItem).ID, nil
}

func CreateSoftwareProgramDB(item *model.SoftwareProgramItem, db *reindexer.Reindexer) *custom_errors.ErrHttp {
	inserted, err := db.Insert("software_program", item, "id=serial()")
	if err != nil {
		return custom_errors.New(fasthttp.StatusInternalServerError, "insert software_program: "+err.Error())
	}

	if inserted == 0 {
		return custom_errors.New(fasthttp.StatusInternalServerError, "insert software_program: something went wrong")
	}

	return nil
}
