package connector

import (
	"cognix.ch/api/v2/core/model"
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/gocolly/colly/v2"
	"go.uber.org/zap"
	"jaytaylor.com/html2text"
	"net/url"
	"strings"
	"time"
)

type (
	Web struct {
		Base
		param   *WebParameters
		scraper *colly.Collector
		history map[string]string
	}
	WebParameters struct {
		URL string `url:"url"`
	}
)

func (c *Web) Config(connector *model.Connector) (Connector, error) {
	c.Base.Config(connector)
	c.param = &WebParameters{}
	c.history = make(map[string]string)
	if err := connector.ConnectorSpecificConfig.ToStruct(c.param); err != nil {
		return nil, err
	}
	c.scraper = colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)
	return c, nil
}

func (c *Web) Execute(ctx context.Context, param model.JSONMap) (*model.Connector, error) {
	zap.S().Debugf("Run web connector with param %s ...", c.param.URL)
	c.scraper.OnHTML("body", c.onBody)
	err := c.scraper.Visit(c.param.URL)
	if err != nil {
		zap.L().Error("Failed to scrape URL", zap.String("url", c.param.URL), zap.Error(err))
	}
	zap.S().Debugf("Complete web connector with param %s", c.param.URL)
	return c.model, err
}

func NewWeb(connector *model.Connector) (Connector, error) {
	web := Web{}
	return web.Config(connector)
}

func (c *Web) onBody(e *colly.HTMLElement) {
	child := e.ChildAttrs("a", "href")
	text, _ := html2text.FromString(e.ChildText("main"), html2text.Options{
		PrettyTables: true,
		PrettyTablesOptions: &html2text.PrettyTablesOptions{
			AutoFormatHeader: true,
			AutoWrapText:     true,
		},
		OmitLinks: true,
	})
	c.history[e.Request.URL.String()] = text
	c.processChildLinks(e.Request.URL, child)
	signature := fmt.Sprintf("%x", sha256.Sum256([]byte(text)))
	docID := e.Request.URL.String()
	doc, ok := c.model.DocsMap[docID]
	if !ok {
		doc = &model.Document{
			DocumentID:  docID,
			ConnectorID: c.model.ID,
			Link:        docID,
			CreatedDate: time.Now().UTC(),
			IsExists:    true,
			IsUpdated:   true,
		}
		c.model.DocsMap[docID] = doc
		c.model.Docs = append(c.model.Docs, doc)
	}
	doc.IsExists = true
	if doc.Signature == signature {
		return
	}
	doc.Signature = signature
	if doc.ID != 0 {
		doc.IsUpdated = true
		doc.UpdatedDate = pg.NullTime{time.Now().UTC()}
	}

	// todo send text for indexing
	//fmt.Println(text)

}

func (c *Web) processChildLinks(baseURL *url.URL, urls []string) {
	for _, u := range urls {
		if len(u) == 0 || u[0] == '#' || !strings.Contains(u, baseURL.Path) ||
			(strings.HasPrefix(u, "http") && !strings.Contains(u, baseURL.Host)) {
			continue
		}
		if strings.HasPrefix(u, baseURL.Path) {
			u = fmt.Sprintf("%s://%s%s", baseURL.Scheme, baseURL.Host, u)
		}
		if _, ok := c.history[u]; ok {
			continue
		}
		if err := c.scraper.Visit(u); err != nil {
			zap.S().Errorf("Failed to scrape URL: %s", u)
		}
	}
	return
}