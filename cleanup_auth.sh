#!/bin/bash

# 批量清理控制器中的重复鉴权逻辑

files=(
    "internal/controller/journal_controller.go"
    "internal/controller/stats_controller.go"
    "internal/controller/plan_controller.go"
)

for file in "${files[@]}"; do
    echo "Processing $file..."
    
    # 使用 perl 进行多行替换，处理鉴权逻辑
    perl -i -pe '
        # 标记替换区域开始
        if (/userIDStr := .*\.Get\("user_id"\)/) {
            $in_auth_block = 1;
            $_ = "	userID, err := middleware.GetUserIDFromContext(ctx)\n";
            next;
        }
        
        # 在替换区域中，跳过原有的鉴权代码直到找到 userID, err := uuid.Parse
        if ($in_auth_block) {
            if (/userID, err := uuid\.Parse/) {
                $_ = "	if err != nil {\n		return middleware.HandleUnauthorized(ctx, err)\n	}\n\n";
                $in_auth_block = 0;
            } else {
                # 跳过这些行
                next;
            }
        }
    ' "$file"
    
    echo "Completed $file"
done

echo "All files processed"
