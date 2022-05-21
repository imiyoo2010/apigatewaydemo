package dbstore

import (
	"context"
	"github.com/cihub/seelog"
	"gopkg.in/olivere/elastic.v5"
)

type ESOutput struct {
	client *elastic.Client
}

func NewESClient(url string) *ESOutput {

	es := new(ESOutput)

	es.init(url)

	return es
}

func (c *ESOutput) init(url string) {

	//host :="http://10.8.119.230:9200/"

	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))

	if err != nil {
		seelog.Error("elastic connect error")
	}

	c.client = client
}

func (c *ESOutput) CollectApiLog(data interface{}) {

	if c.client != nil {
		_, err := c.client.Index().
			Index("myapigateway").
			Type("log").
			BodyJson(data).
			Do(context.Background())
		if err != nil {
			seelog.Error("elastic collect stat data error")
		}
	}
}
