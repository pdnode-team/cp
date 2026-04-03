// util/idgen.go
package util

import (
	"log"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

// 在包被导入时自动初始化（也可以在 main 里面手动调用初始化）
func init() {
	// 这里的 1 是机器的节点 ID。如果是多台服务器，每台机器应分配不同的节点 ID (0-1023)
	// 本地单机测试直接写 1 即可
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		log.Fatalf("雪花节点初始化失败: %v", err)
	}
}

// GenID 提供给 Ent 调用的默认值生成函数
func GenID() int64 {
	return node.Generate().Int64()
}
