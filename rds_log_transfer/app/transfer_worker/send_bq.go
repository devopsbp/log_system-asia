package transfer_worker

import (
	"rds_log_transfer/bq"
	"rds_log_transfer/bq_table/achievement_clear"
	"rds_log_transfer/bq_table/aesthe_update_character"
	"rds_log_transfer/bq_table/auto_order_quest_reload"
	"rds_log_transfer/bq_table/batch_coin_delete_expired"
	"rds_log_transfer/bq_table/batch_currency_delete_expired"
	"rds_log_transfer/bq_table/buy_character_slot"
	"rds_log_transfer/bq_table/buy_season_pass"
	"rds_log_transfer/bq_table/character_add_engram"
	"rds_log_transfer/bq_table/character_ban"
	"rds_log_transfer/bq_table/character_buy_item"
	"rds_log_transfer/bq_table/character_class_add_exp"
	"rds_log_transfer/bq_table/character_create"
	"rds_log_transfer/bq_table/character_delete"
	"rds_log_transfer/bq_table/character_generic_get"
	"rds_log_transfer/bq_table/character_generic_remove"
	"rds_log_transfer/bq_table/character_item_update"
	"rds_log_transfer/bq_table/character_levelup"
	"rds_log_transfer/bq_table/character_money_history"
	"rds_log_transfer/bq_table/character_restore"
	"rds_log_transfer/bq_table/character_sell_item"
	"rds_log_transfer/bq_table/character_use_item"
	"rds_log_transfer/bq_table/client_health_check"
	"rds_log_transfer/bq_table/currency_history"
	"rds_log_transfer/bq_table/currency_purchase"
	"rds_log_transfer/bq_table/defense_battle"
	"rds_log_transfer/bq_table/finish_send_job"
	"rds_log_transfer/bq_table/gasha_reward_receipt"
	"rds_log_transfer/bq_table/get_season_pass_reward"
	"rds_log_transfer/bq_table/guild_accept_invitation"
	"rds_log_transfer/bq_table/guild_accept_member"
	"rds_log_transfer/bq_table/guild_change_name"
	"rds_log_transfer/bq_table/guild_create"
	"rds_log_transfer/bq_table/guild_dissolution"
	"rds_log_transfer/bq_table/guild_entry"
	"rds_log_transfer/bq_table/guild_exchange_authority"
	"rds_log_transfer/bq_table/guild_invitation"
	"rds_log_transfer/bq_table/guild_kick_member"
	"rds_log_transfer/bq_table/guild_withdraw"
	"rds_log_transfer/bq_table/imagine_craft"
	"rds_log_transfer/bq_table/item_craft"
	"rds_log_transfer/bq_table/item_equip"
	"rds_log_transfer/bq_table/liquid_memory_accumulate"
	"rds_log_transfer/bq_table/liquid_memory_use"
	"rds_log_transfer/bq_table/login"
	"rds_log_transfer/bq_table/logout"
	"rds_log_transfer/bq_table/map_conn"
	"rds_log_transfer/bq_table/player_using_item"
	"rds_log_transfer/bq_table/quest_cancel"
	"rds_log_transfer/bq_table/quest_complete"
	"rds_log_transfer/bq_table/quest_start"
	"rds_log_transfer/bq_table/quest_update"
	"rds_log_transfer/bq_table/remove_appearance_weapon_sticker"
	"rds_log_transfer/bq_table/save_appearance_weapon_sticker"
	"rds_log_transfer/bq_table/season_pass_change_points"
	"rds_log_transfer/bq_table/season_pass_complete_quest"
	"rds_log_transfer/bq_table/season_pass_rank_up"
	"rds_log_transfer/bq_table/second_password_expired"
	"rds_log_transfer/bq_table/second_password_registration"
	"rds_log_transfer/bq_table/second_password_reset"
	"rds_log_transfer/bq_table/send_mail_error"
	"rds_log_transfer/bq_table/skill_acquisition"
	"rds_log_transfer/bq_table/soul_rank_up"
	"rds_log_transfer/bq_table/stamp_acquisition"
	"rds_log_transfer/bq_table/start_send_job"
	"rds_log_transfer/bq_table/weapon_craft"
	"rds_log_transfer/bq_table/weapon_fusion"
	"rds_log_transfer/bq_table/weapon_recycle"
	"rds_log_transfer/bq_table/weapon_tuning"
	"rds_log_transfer/bq_table/weapon_unlock_slot"
	shard_db "rds_log_transfer/shard-db"
)

// workerのBQ投入処理

func sendBq(tag string, records []*shard_db.SbLogRecord, workerNo int, clientWithContext *bq.ClientWithContext) ([]*shard_db.SbLogRecord, error) {
	var err error
	var parseFailedRecords []*shard_db.SbLogRecord

	switch tag { // caseはgeneratorで出力したものをコピー
	case "SB_LOG_item_craft":
		parseFailedRecords, err = item_craft.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_liquid_memory_use":
		parseFailedRecords, err = liquid_memory_use.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_liquid_memory_accumulate":
		parseFailedRecords, err = liquid_memory_accumulate.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_soul_rank_up":
		parseFailedRecords, err = soul_rank_up.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_skill_acquisition":
		parseFailedRecords, err = skill_acquisition.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_weapon_craft":
		parseFailedRecords, err = weapon_craft.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_weapon_unlock_slot":
		parseFailedRecords, err = weapon_unlock_slot.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_imagine_craft":
		parseFailedRecords, err = imagine_craft.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_weapon_fusion":
		parseFailedRecords, err = weapon_fusion.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_login":
		parseFailedRecords, err = login.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_logout":
		parseFailedRecords, err = logout.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_client_health_check":
		parseFailedRecords, err = client_health_check.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_character_create":
		parseFailedRecords, err = character_create.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_character_delete":
		parseFailedRecords, err = character_delete.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_character_restore":
		parseFailedRecords, err = character_restore.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_send_mail_error":
		parseFailedRecords, err = send_mail_error.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_start_send_job":
		parseFailedRecords, err = start_send_job.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_finish_send_job":
		parseFailedRecords, err = finish_send_job.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_map_conn":
		parseFailedRecords, err = map_conn.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_character_ban":
		parseFailedRecords, err = character_ban.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_character_class_add_exp":
		parseFailedRecords, err = character_class_add_exp.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_character_item_update":
		parseFailedRecords, err = character_item_update.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_character_buy_item":
		parseFailedRecords, err = character_buy_item.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_character_sell_item":
		parseFailedRecords, err = character_sell_item.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_character_use_item":
		parseFailedRecords, err = character_use_item.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_character_add_engram":
		parseFailedRecords, err = character_add_engram.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_character_levelup":
		parseFailedRecords, err = character_levelup.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_player_using_item":
		parseFailedRecords, err = player_using_item.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_quest_start":
		parseFailedRecords, err = quest_start.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_quest_update":
		parseFailedRecords, err = quest_update.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_quest_cancel":
		parseFailedRecords, err = quest_cancel.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_quest_complete":
		parseFailedRecords, err = quest_complete.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_item_equip":
		parseFailedRecords, err = item_equip.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_save_appearance_weapon_sticker":
		parseFailedRecords, err = save_appearance_weapon_sticker.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_remove_appearance_weapon_sticker":
		parseFailedRecords, err = remove_appearance_weapon_sticker.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_currency_history":
		parseFailedRecords, err = currency_history.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_currency_purchase":
		parseFailedRecords, err = currency_purchase.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_character_money_history":
		parseFailedRecords, err = character_money_history.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_batch_currency_delete_expired":
		parseFailedRecords, err = batch_currency_delete_expired.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_batch_coin_delete_expired":
		parseFailedRecords, err = batch_coin_delete_expired.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_guild_create":
		parseFailedRecords, err = guild_create.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_guild_change_name":
		parseFailedRecords, err = guild_change_name.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_guild_entry":
		parseFailedRecords, err = guild_entry.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_guild_accept_member":
		parseFailedRecords, err = guild_accept_member.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_guild_invitation":
		parseFailedRecords, err = guild_invitation.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_guild_accept_invitation":
		parseFailedRecords, err = guild_accept_invitation.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_guild_withdraw":
		parseFailedRecords, err = guild_withdraw.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_guild_kick_member":
		parseFailedRecords, err = guild_kick_member.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_guild_dissolution":
		parseFailedRecords, err = guild_dissolution.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_guild_exchange_authority":
		parseFailedRecords, err = guild_exchange_authority.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_achievement_clear":
		parseFailedRecords, err = achievement_clear.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_second_password_registration":
		parseFailedRecords, err = second_password_registration.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_second_password_expired":
		parseFailedRecords, err = second_password_expired.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_second_password_reset":
		parseFailedRecords, err = second_password_reset.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_buy_season_pass":
		parseFailedRecords, err = buy_season_pass.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_season_pass_complete_quest":
		parseFailedRecords, err = season_pass_complete_quest.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_season_pass_change_points":
		parseFailedRecords, err = season_pass_change_points.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_auto_order_quest_reload":
		parseFailedRecords, err = auto_order_quest_reload.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_get_season_pass_reward":
		parseFailedRecords, err = get_season_pass_reward.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_season_pass_rank_up":
		parseFailedRecords, err = season_pass_rank_up.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_weapon_recycle":
		parseFailedRecords, err = weapon_recycle.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_weapon_tuning":
		parseFailedRecords, err = weapon_tuning.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_gasha_reward_receipt":
		parseFailedRecords, err = gasha_reward_receipt.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_stamp_acquisition":
		parseFailedRecords, err = stamp_acquisition.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_aesthe_update_character":
		parseFailedRecords, err = aesthe_update_character.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_buy_character_slot":
		parseFailedRecords, err = buy_character_slot.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_defense_battle":
		parseFailedRecords, err = defense_battle.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_character_generic_get":
		parseFailedRecords, err = character_generic_get.InsertBqItems(clientWithContext, &records, workerNo)
	case "SB_LOG_character_generic_remove":
		parseFailedRecords, err = character_generic_remove.InsertBqItems(clientWithContext, &records, workerNo)
	} // end switch

	return parseFailedRecords, err
}
