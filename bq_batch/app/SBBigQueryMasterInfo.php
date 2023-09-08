<?php

namespace App;

/**
 * adminツールのログビューアーで使用するSB_LOGのタグ定義
 * resources/gen_src/gamelog/skyblue_game_log.xlsxから自動生成
 */
class SBBigQueryMasterInfo
{
    /**
     * ログEXCELに定義されているデータ情報（マスタ名）を取得する
     * @param   string $name データ名
     * @return  array|null
     */
    public static function getDataInfo(string $name): ?array {
        for ($index = 0; $index < count(SBBigQueryMasterInfo::MASTER_LIST); $index++) {
            if (SBBigQueryMasterInfo::MASTER_LIST[$index]['name'] === $name) {
				return SBBigQueryMasterInfo::MASTER_LIST[$index];
            }
        }

        return(null);
    }

    // マスタ参照が必要なキーのリスト
    const MASTER_LIST = [
        // アチーブメントID
        [ 'name' => 'achievement_id', 'master_name' => 'master_achievements', 'key' => 'id', 'text_key' => 'name', 'text_table_name' => 'achievements_text' ],
        // 装備したアーツ（スキル）ID
        [ 'name' => 'arts_id', 'master_name' => 'master_skill_datas', 'key' => 'skill_id', 'text_key' => 'skill_name', 'text_table_name' => 'master_skill_data_text' ],
        // 利用したチケット
        [ 'name' => 'change_token_id', 'master_name' => 'master_token', 'key' => 'id', 'text_key' => 'name', 'text_table_name' => 'master_token_text' ],
        // 装備ID：アクセ1
        [ 'name' => 'cos_acc1_id', 'master_name' => 'master_costume', 'key' => 'id', 'text_key' => 'name', 'text_table_name' => 'costume_text' ],
        // 装備ID：アクセ2
        [ 'name' => 'cos_acc2_id', 'master_name' => 'master_costume', 'key' => 'id', 'text_key' => 'name', 'text_table_name' => 'costume_text' ],
        // 装備ID：アクセ3
        [ 'name' => 'cos_acc3_id', 'master_name' => 'master_costume', 'key' => 'id', 'text_key' => 'name', 'text_table_name' => 'costume_text' ],
        // 装備ID：衣装・胴
        [ 'name' => 'cos_body_id', 'master_name' => 'master_costume', 'key' => 'id', 'text_key' => 'name', 'text_table_name' => 'costume_text' ],
        // 装備ID：衣装・手
        [ 'name' => 'cos_hand_id', 'master_name' => 'master_costume', 'key' => 'id', 'text_key' => 'name', 'text_table_name' => 'costume_text' ],
        // 装備ID：衣装・頭
        [ 'name' => 'cos_head_id', 'master_name' => 'master_costume', 'key' => 'id', 'text_key' => 'name', 'text_table_name' => 'costume_text' ],
        // 装備ID：衣装・足
        [ 'name' => 'cos_leg_id', 'master_name' => 'master_costume', 'key' => 'id', 'text_key' => 'name', 'text_table_name' => 'costume_text' ],
        // 装備ID：衣装・靴
        [ 'name' => 'cos_shoes_id', 'master_name' => 'master_costume', 'key' => 'id', 'text_key' => 'name', 'text_table_name' => 'costume_text' ],
        // コースID
        [ 'name' => 'course_id', 'master_name' => 'master_aesthe_course', 'key' => 'id', 'text_key' => 'name', 'text_table_name' => 'master_aesthe_course_text' ],
        // ガシャID
        [ 'name' => 'gasha_id', 'master_name' => 'master_gasha', 'key' => 'id', 'text_key' => 'title', 'text_table_name' => 'master_gasha_text' ],
        // イマジンID
        [ 'name' => 'imagine_id', 'master_name' => 'master_imagine', 'key' => 'id', 'text_key' => 'imagine_name', 'text_table_name' => 'master_imagine_text' ],
        // アイテムID
        [ 'name' => 'item_id', 'master_name' => 'master_items', 'key' => 'id', 'text_key' => 'name', 'text_table_name' => 'item_text' ],
        // 変更後クラスの武器アイテムID
        [ 'name' => 'now_weapon_item_id', 'master_name' => 'master_weapons', 'key' => 'id', 'text_key' => 'name', 'text_table_name' => 'weapon_text' ],
        // 能力ID
        [ 'name' => 'perk_id', 'master_name' => 'master_perks', 'key' => 'id', 'text_key' => 'ability_name1', 'text_table_name' => 'perk_text' ],
        // シーズンID
        [ 'name' => 'season_id', 'master_name' => 'master_season', 'key' => 'id', 'text_key' => 'season_name', 'text_table_name' => 'master_season_text' ],
        // シーズンパスクエストID
        [ 'name' => 'season_pass_quest_id', 'master_name' => 'master_quest_auto_order_season', 'key' => 'id', 'text_key' => 'quest_name', 'text_table_name' => 'master_quest_auto_order_season_text' ],
        // スキルID
        [ 'name' => 'skill_id', 'master_name' => 'master_skill_datas', 'key' => 'skill_id', 'text_key' => 'skill_name', 'text_table_name' => 'master_skill_data_text' ],
        // シールを外すときの使ったアイテムＩＤ
        [ 'name' => 'sticker_remover_item_id', 'master_name' => 'master_items', 'key' => 'id', 'text_key' => 'name', 'text_table_name' => 'item_text' ],
        // 対象見た目武器のアイテムＩＤ
        [ 'name' => 'weapon_appearance_weapon_id', 'master_name' => 'master_weapons', 'key' => 'id', 'text_key' => 'name', 'text_table_name' => 'weapon_text' ],
        // 使ったシールアイテムのアイテムＩＤ
        [ 'name' => 'weapon_appearance_weapon_sticker_id', 'master_name' => 'master_items', 'key' => 'id', 'text_key' => 'name', 'text_table_name' => 'item_text' ],
        // 装備ID：武器
        [ 'name' => 'weapon_id', 'master_name' => 'master_weapons', 'key' => 'id', 'text_key' => 'name', 'text_table_name' => 'weapon_text' ],
    ];
}
