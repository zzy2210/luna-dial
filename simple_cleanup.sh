#!/bin/bash

# 简化的批量替换脚本 

echo "Processing journal_controller.go..."

# 1. 添加 middleware 导入
sed -i '/import (/,/)/{ /^[[:space:]]*"okr-web\/internal\/service"$/a\
	"okr-web/internal/middleware"
}' internal/controller/journal_controller.go

# 2. 替换鉴权逻辑
# 第一次替换模式
sed -i 's/userIDStr := ctx\.Get("user_id")/userID, err := middleware.GetUserIDFromContext(ctx)/g' internal/controller/journal_controller.go

# 替换空检查和错误处理
sed -i '/if userIDStr == nil {/{
N
N
N
N
c\
	if err != nil {\
		return middleware.HandleUnauthorized(ctx, err)\
	}
}' internal/controller/journal_controller.go

# 删除 uuid.Parse 相关代码
sed -i '/userID, err := uuid\.Parse(userIDStr\.(string))/,+4d' internal/controller/journal_controller.go

echo "Processing stats_controller.go..."

# 同样的处理步骤用于 stats_controller.go
sed -i '/import (/,/)/{ /^[[:space:]]*"okr-web\/internal\/service"$/a\
	"okr-web/internal/middleware"
}' internal/controller/stats_controller.go

sed -i 's/userIDStr := ctx\.Get("user_id")/userID, err := middleware.GetUserIDFromContext(ctx)/g' internal/controller/stats_controller.go

sed -i '/if userIDStr == nil {/{
N
N
N
N
c\
	if err != nil {\
		return middleware.HandleUnauthorized(ctx, err)\
	}
}' internal/controller/stats_controller.go

sed -i '/userID, err := uuid\.Parse(userIDStr\.(string))/,+4d' internal/controller/stats_controller.go

echo "Done"
