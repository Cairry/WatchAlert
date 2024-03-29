package services

import (
	"watchAlert/alert/queue"
	"watchAlert/controllers/repo"
	"watchAlert/globals"
	models "watchAlert/models"
	"watchAlert/utils/cmd"
)

type RuleService struct {
	rule chan *models.AlertRule
}

type InterRuleService interface {
	Create(rule models.AlertRule) error
	Update(rule models.AlertRule) error
	Delete(id string) error
	List(ruleGroupId string) ([]models.AlertRule, error)
	Search(ruleId string) models.AlertRule
}

func NewInterRuleService() InterRuleService {
	return &RuleService{
		rule: queue.AlertRuleChannel,
	}
}

func (rs *RuleService) Create(rule models.AlertRule) error {

	rule.RuleId = "a-" + cmd.RandId()

	newRule := rule.ParserRuleToGorm()

	err := repo.DBCli.Create(&models.AlertRule{}, &newRule)
	if err != nil {
		return err
	}

	rs.rule <- newRule

	return nil

}

func (rs *RuleService) Update(rule models.AlertRule) error {

	newRule := rule.ParserRuleToGorm()

	/*
		重启协程
		判断当前状态是否是false 并且 历史状态是否为true
	*/
	alertInfo := models.AlertRule{}
	globals.DBCli.Model(&models.AlertRule{}).Where("rule_id = ?", rule.RuleId).Find(&alertInfo)

	if alertInfo.Enabled == "true" && newRule.EnabledBool == false {
		if cancel, exists := queue.WatchCtxMap[newRule.RuleId]; exists {
			cancel()
		}
	}
	if alertInfo.Enabled == "true" && newRule.EnabledBool == true {
		if cancel, exists := queue.WatchCtxMap[newRule.RuleId]; exists {
			cancel()
		}
	}

	// 删除缓存
	iter := globals.RedisCli.Scan(0, models.FiringAlertCachePrefix+rule.RuleId+"*", 0).Iterator()
	keys := make([]string, 0)
	for iter.Next() {
		key := iter.Val()
		keys = append(keys, key)
	}
	globals.RedisCli.Del(keys...)

	// 启动协程
	if newRule.EnabledBool {
		rs.rule <- newRule
		globals.Logger.Sugar().Infof("重启 RuleId 为 %s 的 Watch 进程", newRule.RuleId)
	}

	// 更新数据
	data := repo.Updates{
		Table:   models.AlertRule{},
		Where:   []string{"rule_id = ?", newRule.RuleId},
		Updates: &newRule,
	}
	err := repo.DBCli.Updates(data)
	if err != nil {
		return err
	}

	return nil

}

func (rs *RuleService) Delete(id string) error {

	var alertRule models.AlertRule
	data := repo.Delete{
		Table: alertRule,
		Where: []string{"rule_id = ?", id},
	}

	err := repo.DBCli.Delete(data)
	if err != nil {
		return err
	}

	// 退出该规则的协程
	if alertRule.EnabledBool {
		globals.Logger.Sugar().Infof("停止 RuleId 为 %s 的Watch 进程", id)
		if cancel, exists := queue.WatchCtxMap[id]; exists {
			cancel()
		}
		//rs.quit <- &id
	}

	iter := globals.RedisCli.Scan(0, models.FiringAlertCachePrefix+id+"*", 0).Iterator()
	keys := make([]string, 0)
	for iter.Next() {
		key := iter.Val()
		keys = append(keys, key)
	}

	globals.RedisCli.Del(keys...)
	globals.Logger.Sugar().Infof("删除队列数据 ->%s", keys)

	return nil

}

func (rs *RuleService) List(ruleGroupId string) ([]models.AlertRule, error) {

	var alertRuleList []models.AlertRule

	globals.DBCli.Model(&models.AlertRule{}).Where("rule_group_id = ?", ruleGroupId).Find(&alertRuleList)

	for k, v := range alertRuleList {
		newRule := v.ParserRuleToJson()
		alertRuleList[k] = *newRule
	}

	return alertRuleList, nil

}

func (rs *RuleService) Search(ruleId string) models.AlertRule {

	var alertRule models.AlertRule
	globals.DBCli.Where("rule_id", ruleId).Find(&alertRule)

	if alertRule.RuleName == "" {
		return models.AlertRule{}
	}

	newRule := alertRule.ParserRuleToJson()

	return *newRule

}
