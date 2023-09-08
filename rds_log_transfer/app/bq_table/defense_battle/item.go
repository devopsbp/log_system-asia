package defense_battle

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"rds_log_transfer/bq"
	"rds_log_transfer/sb_log"
	shard_db "rds_log_transfer/shard-db"
)

var boolKeys = []string{}

type DefenseBattleItem struct {
	UserId                      *string               `json:"user_id,omitempty"`
	CharacterId                 *string               `json:"character_id,omitempty"`
	CharacterClass              *interface{}          `json:"character_class,omitempty"`
	CharacterClassLevel         *json.Number          `json:"character_class_level,omitempty"`
	UnalteredContributionScores *string               `json:"unaltered_contribution_scores,omitempty"`
	CharacterBattleScore        *json.Number          `json:"character_battle_score,omitempty"`
	EquippedWeapons             *string               `json:"equipped_weapons,omitempty"`
	EquippedBattleImagines      *string               `json:"equipped_battle_imagines,omitempty"`
	EquippedInnerImagines       *string               `json:"equipped_inner_imagines,omitempty"`
	EquippedTacticalSkills      *string               `json:"equipped_tactical_skills,omitempty"`
	EquippedTacticalAbilities   *string               `json:"equipped_tactical_abilities,omitempty"`
	EquippedClassSkills         *string               `json:"equipped_class_skills,omitempty"`
	CharacterDamage             *json.Number          `json:"character_damage,omitempty"`
	NumberOfDeaths              *json.Number          `json:"number_of_deaths,omitempty"`
	NumberOfEnemiesDefeated     *json.Number          `json:"number_of_enemies_defeated,omitempty"`
	RequestHeaderLocation       *interface{}          `json:"request_header_location,omitempty"`
	LogTime                     *string               `json:"log_time,omitempty"`
	Uuid                        *string               `json:"-"`
	Record                      *shard_db.SbLogRecord `json:"-"`
}

func (i DefenseBattleItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"user_id":                       i.Record.UserId,
		"character_id":                  i.Record.CharacterId,
		"character_class":               sb_log.EncodeJsonValue(i.CharacterClass),
		"character_class_level":         i.CharacterClassLevel,
		"unaltered_contribution_scores": i.UnalteredContributionScores,
		"character_battle_score":        i.CharacterBattleScore,
		"equipped_weapons":              i.EquippedWeapons,
		"equipped_battle_imagines":      i.EquippedBattleImagines,
		"equipped_inner_imagines":       i.EquippedInnerImagines,
		"equipped_tactical_skills":      i.EquippedTacticalSkills,
		"equipped_tactical_abilities":   i.EquippedTacticalAbilities,
		"equipped_class_skills":         i.EquippedClassSkills,
		"character_damage":              i.CharacterDamage,
		"number_of_deaths":              i.NumberOfDeaths,
		"number_of_enemies_defeated":    i.NumberOfEnemiesDefeated,
		"request_header_location":       sb_log.GetRequestHeaderLocationStr(i.Record),
		"log_time":                      i.LogTime,
		"uuid":                          i.Uuid,
	}, bigquery.NoDedupeID, nil
}

func InsertBqItems(clientWithContext *bq.ClientWithContext, records *[]*shard_db.SbLogRecord, uuidSeed int) ([]*shard_db.SbLogRecord, error) {
	var tableName string
	var items = make([]*DefenseBattleItem, 0)
	var failedRecords = make([]*shard_db.SbLogRecord, 0)
	for _, record := range *records {
		if tableName == "" {
			tableName = record.LogTag
		}
		var item DefenseBattleItem
		item.Record = record
		uuid := sb_log.GetUUID(uuidSeed, record.Id)
		item.Uuid = &uuid
		// escape JSON value
		jsonEscapedLog := sb_log.EscapeJsonValue(record.Log, &boolKeys)
		err := json.Unmarshal([]byte(jsonEscapedLog), &item)
		if err != nil {
			failedRecords = append(failedRecords, record)
		}
		items = append(items, &item)
	}

	err := bq.InsertRows(clientWithContext, tableName, &items)
	if err != nil {
		return failedRecords, err
	}

	return failedRecords, nil
}
