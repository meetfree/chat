package util

import "github.com/bwmarrin/snowflake"

func GenId() int64 {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	return node.Generate().Int64()
}
