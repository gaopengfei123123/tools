package calculate

import (
	"context"
	"github.com/astaxie/beego/logs"
	"github.com/gaopengfei123123/tools/debug"
	"github.com/olivere/elastic/v7"
	"testing"
)

// 需要用到会话相关指标统计
const SceneConversation = "SceneConversation"

// ConversationMetricsConversationCnt 会话总数
const ConversationMetricsConversationCnt = "ConversationMetricsConversationCnt"

// 将获取函数注册到配置中
func init() {
	InitConfig()
}

// InitConfig 将指标相关配置加载到组件中
func InitConfig() {
	// 设置获取指标的方法
	SetMetricsAggFunc(GetAllAggFunc)
}

// 多场景下加载统计指标
func GetAllAggFunc(metricName string, sceneName ...string) AggFunc {
	scene := "default"
	if len(sceneName) > 0 {
		scene = sceneName[0]
	}

	switch scene {
	case SceneConversation:
		return GetConversationMetricsAgg(metricName, sceneName...)
	default:
	}
	return nil
}

func GetConversationMetricsAgg(metricName string, sceneName ...string) AggFunc {
	logs.Info("GetConversationMetricsAgg, metricName:%s, sceneName: %v", metricName, sceneName)

	switch metricName {
	case ConversationMetricsConversationCnt:
		return ConversationMetricsConversationCntFunc
	default:
	}
	return nil
}

// ConversationMetricsConversationCntFunc 会话数
func ConversationMetricsConversationCntFunc(params map[string]interface{}, currentTerm ...string) elastic.Aggregation {
	tmpQuery := elastic.NewBoolQuery()
	metrics := elastic.NewFilterAggregation().Filter(tmpQuery)
	aggCount := elastic.NewCardinalityAggregation().Field("conversation_id")
	metrics.SubAggregation(SignSingle, aggCount)
	return metrics
}

func TestGetConversationMetricsAgg(t *testing.T) {
	params := map[string]interface{}{
		"corp_wxid": "hujinglin",
	}

	metrics := []string{
		ConversationMetricsConversationCnt,
	}

	// 要聚合的指标层级
	termsList := []string{
		"corp_wxid",
	}

	client := getEsCline()
	ctx := context.TODO()

	indexList := []string{
		"wk_workwx_conversation_202308",
	}
	builder := new(EsQueryBuilder)

	result, _, err := builder.GetTermsMetrics(ctx, client, indexList, termsList, metrics, params, SceneConversation)

	debug.PrintJson("result", result, true)
	logs.Info("err: %v", err)

}
