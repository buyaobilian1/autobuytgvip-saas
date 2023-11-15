package id

import (
	"github.com/bwmarrin/snowflake"
	"log"
)

func GenerateId(nodeId int64) string {
	node, err := snowflake.NewNode(nodeId)
	if err != nil {
		log.Println(err)
		return ""
	}
	id := node.Generate()
	return id.String()
}
