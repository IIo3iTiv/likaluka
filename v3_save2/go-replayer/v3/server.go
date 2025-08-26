package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

type LoopData struct {
	Id        uint64 `json:"Id"`
	StartFile int    `json:"Start"`
	EndFile   int    `json:"End"`
	Loop      bool   `json:"Loop"`
	Guid      int    `json:"Guid"`
	ch        chan bool
}

type LoopDataArr struct {
	Arr []LoopData `json:"Arr"`
}

var loopDataIdMap map[uint64]int = map[uint64]int{}

var x LoopDataArr
var id uint64 = 0

func InitServer() {
	r := gin.Default()
	r.GET("/files/transfer/:guid/:start/:end/:loop?", LoadDataset)
	r.GET("/loop/close/:id", CloseLoop)
	r.GET("/loop/list", ListLoop)
	r.Run(":8080")
}

func ListLoop(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	json, err := jsoniter.Marshal(&x)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "marshal json",
		})
		log.Printf("Error. json.Marshal: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, string(json))
}

func CloseLoop(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id is not valid",
		})
	}
	id, ok := loopDataIdMap[uint64(id)]
	if ok {
		x.Arr[id].ch <- true
		c.JSON(http.StatusOK, gin.H{
			"Success": "Loop is closed",
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "id does not exist",
	})
}

func LoadDataset(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	guidStr := c.Param("guid")
	startStr := c.Param("start")
	endStr := c.Param("end")
	loopStr := c.Param("loop")

	guid, err := strconv.Atoi(guidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "guid is not valid",
		})
		return
	}

	start, err := strconv.Atoi(startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "start is not number",
		})
		return
	}

	end, err := strconv.Atoi(endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "end is not number",
		})
		return
	}

	if start < 1 || start > 38 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("start is not valid 1-38. Get:%d", start),
		})
		return
	}

	if end < 1 || end > 38 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("start is not valid 1-38. Get:%d", start),
		})
		return
	}

	loop := false
	if loopStr != "" {
		loop = strings.TrimSpace(strings.ToLower(loopStr)) == "true"
	}

	id++
	if loop {
		ch := make(chan bool)
		x.Arr = append(x.Arr, LoopData{
			Id:        id,
			StartFile: start,
			EndFile:   end,
			Loop:      loop,
			Guid:      guid,
			ch:        ch,
		})
		loopDataIdMap[id] = len(x.Arr)
		go ReadCSVLoop(start, end, guid, loop, &ch)
		return
	}
	go ReadCSVLoop(start, end, guid, loop, nil)
}
