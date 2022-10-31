package calculate

import "github.com/olivere/elastic/v7"

func getEsCline() *elastic.Client {
	client, _ := elastic.NewClient(elastic.SetURL("http://172.17.33.40:9200"), elastic.SetTraceLog(new(tracelog)))
	return client
}
