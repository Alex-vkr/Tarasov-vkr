package consruct_namespaces

import (
	"diplom-vuls-server/internal/model"
	"errors"
	"github.com/restream/reindexer/v3"
	"github.com/rs/zerolog/log"
)

func ConstructNamespaces(db *reindexer.Reindexer) {
	if err := db.OpenNamespace("software_program", reindexer.DefaultNamespaceOptions(), model.SoftwareProgramItem{}); err != nil {
		log.Error().Err(errors.New("open namespace software_program: " + err.Error())).Msg("")
	}
	if err := db.OpenNamespace("version", reindexer.DefaultNamespaceOptions(), model.VersionItem{}); err != nil {
		log.Error().Err(errors.New("open namespace version: " + err.Error())).Msg("")
	}
	if err := db.OpenNamespace("platform", reindexer.DefaultNamespaceOptions(), model.PlatformItem{}); err != nil {
		log.Error().Err(errors.New("open namespace platform: " + err.Error())).Msg("")
	}
	if err := db.OpenNamespace("vulnerability", reindexer.DefaultNamespaceOptions(), model.VulnerabilityItem{}); err != nil {
		log.Error().Err(errors.New("open namespace vulnerability: " + err.Error())).Msg("")
	}
}
