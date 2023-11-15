package id

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
)

func GenerateId(nodeId int64) string {
	node, err := snowflake.NewNode(nodeId)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	id := node.Generate()
	return id.String()
}
