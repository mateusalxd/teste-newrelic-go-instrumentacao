package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	nrgin "github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func endpoint1(g *gin.Context) {
	fmt.Println("endpoint1")
	endpoint11(g)
}

func endpoint11(g *gin.Context) {
	txn := nrgin.Transaction(g)
	defer txn.StartSegment("endpoint11").End()

	time.Sleep(time.Second * 5)
	endpoint2(g)
}

func endpoint2(g *gin.Context) {
	txn := nrgin.Transaction(g)
	defer txn.StartSegment("endpoint2").End()
	time.Sleep(time.Second * 1)

	endpoint21(g)
	endpoint22(g)
	endpoint23(g)
	endpoint24(g)
}

func endpoint21(g *gin.Context) {
	txn := nrgin.Transaction(g)
	ds := &newrelic.DatastoreSegment{
		StartTime:  txn.StartSegmentNow(),
		Product:    newrelic.DatastoreOracle,
		Collection: "SELECT COUNT(1) FROM tabela1",
		Operation:  "SELECT",
	}
	defer ds.End()

	time.Sleep(time.Second * 3)
}

func endpoint22(g *gin.Context) {
	txn := nrgin.Transaction(g)
	ds := &newrelic.DatastoreSegment{
		StartTime:  txn.StartSegmentNow(),
		Product:    newrelic.DatastoreOracle,
		Collection: "UPDATE tabela2 SET campo1 = 2",
		Operation:  "UPDATE",
	}
	defer ds.End()

	time.Sleep(time.Second * 10)
}

func endpoint23(g *gin.Context) {
	txn := nrgin.Transaction(g)
	ds := &newrelic.DatastoreSegment{
		StartTime:  txn.StartSegmentNow(),
		Product:    newrelic.DatastoreOracle,
		Collection: "BEGIN; PKG(123) END;",
		Operation:  "ANONYMOUS BLOCK",
	}
	defer ds.End()

	time.Sleep(time.Second * 2)
}

func endpoint24(g *gin.Context) {
	txn := nrgin.Transaction(g)
	s := newrelic.ExternalSegment{
		StartTime: txn.StartSegmentNow(),
		URL:       "http://www.example.com",
	}
	defer s.End()

	time.Sleep(time.Second * 4)
}

func main() {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("teste-primeira-api"),
		newrelic.ConfigLicense(os.Getenv("NEWRELIC_LICENSE")),
		newrelic.ConfigDistributedTracerEnabled(true),
	)

	if err != nil {
		fmt.Println(err.Error())
	}

	r := gin.Default()
	r.Use(nrgin.Middleware(app))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/e1", endpoint1)

	r.Run()
}
