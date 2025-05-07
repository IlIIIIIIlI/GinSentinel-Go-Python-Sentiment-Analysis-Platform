#!/bin/bash

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # 无颜色

# 测试文本
TEST_TEXT="这是一个用于性能测试的样本文本"

# 测试单独顺序请求
test_sequential() {
    local count=$1
    echo -e "${BLUE}测试 $count 个顺序请求...${NC}"

    start_time=$(date +%s.%N)

    for (( i=1; i<=$count; i++ )); do
        curl -s -X POST http://localhost:9001/api/v1/sentiment/analyze \
          -H "Content-Type: application/json" \
          -d "{
            \"text\": \"$TEST_TEXT\",
            \"language\": \"zh\",
            \"store_result\": false
          }" > /dev/null

        # 每100个请求打印一次进度
        if (( i % 100 == 0 )); then
            echo -e "  已完成: ${i}/${count}"
        fi
    done

    end_time=$(date +%s.%N)
    duration=$(echo "$end_time - $start_time" | bc)

    echo -e "${GREEN}所有顺序请求完成!${NC}"
    echo -e "总时间: ${YELLOW}${duration}秒${NC}"
    echo -e "平均每个请求: ${YELLOW}$(echo "scale=4; $duration/$count" | bc)秒${NC}"
    echo -e "每秒请求数: ${YELLOW}$(echo "scale=2; $count/$duration" | bc)${NC}"
    echo "--------------------------------------------------------------"

    # 保存结果
    echo "sequential,$count,$duration,$(echo "scale=4; $duration/$count" | bc),$(echo "scale=2; $count/$duration" | bc)" >> performance_comparison.csv
}

# 测试并发请求
test_concurrent() {
    local count=$1
    local max_parallel=50  # 每批最大并发数

    echo -e "${BLUE}测试 $count 个并发请求...${NC}"
    start_time=$(date +%s.%N)

    # 创建临时文件存储结果
    tmp_file=$(mktemp)

    # 按批次发送请求（避免打开过多文件描述符）
    for (( i=0; i<$count; i+=$max_parallel )); do
        batch_size=$max_parallel
        if (( i + batch_size > count )); then
            batch_size=$((count - i))
        fi

        echo -e "  发送批次: $((i/max_parallel + 1)), 大小: $batch_size"

        # 并发发送这一批次的请求
        for (( j=0; j<$batch_size; j++ )); do
            (curl -s -X POST http://localhost:9001/api/v1/sentiment/analyze \
              -H "Content-Type: application/json" \
              -d "{
                \"text\": \"$TEST_TEXT\",
                \"language\": \"zh\",
                \"store_result\": false
              }" > /dev/null; echo "$?" >> $tmp_file) &
        done

        # 等待这一批次完成
        wait
    done

    end_time=$(date +%s.%N)
    duration=$(echo "$end_time - $start_time" | bc)

    # 计算成功请求数
    success_count=$(grep -c "^0$" $tmp_file)

    echo -e "${GREEN}所有并发请求完成!${NC}"
    echo -e "成功请求: ${GREEN}$success_count${NC}/${YELLOW}$count${NC}"
    echo -e "总时间: ${YELLOW}${duration}秒${NC}"
    echo -e "平均每个请求: ${YELLOW}$(echo "scale=4; $duration/$count" | bc)秒${NC}"
    echo -e "每秒请求数: ${YELLOW}$(echo "scale=2; $count/$duration" | bc)${NC}"
    echo "--------------------------------------------------------------"

    # 保存结果
    echo "concurrent,$count,$success_count,$duration,$(echo "scale=4; $duration/$count" | bc),$(echo "scale=2; $count/$duration" | bc)" >> performance_comparison.csv

    # 清理临时文件
    rm $tmp_file
}

# 主函数
main() {
    echo -e "${GREEN}开始1000个请求性能对比测试...${NC}"

    # 创建或清空结果文件
    echo "type,count,success_count,total_time,avg_time_per_request,requests_per_second" > performance_comparison.csv

    # 运行对比测试
    test_sequential 1000
    test_concurrent 1000

    # 输出对比总结
    echo -e "${YELLOW}=== 性能对比总结 ===${NC}"
    seq_time=$(awk -F, '/sequential/ {print $4}' performance_comparison.csv)
    seq_rps=$(awk -F, '/sequential/ {print $5}' performance_comparison.csv)
    con_time=$(awk -F, '/concurrent/ {print $5}' performance_comparison.csv)
    con_rps=$(awk -F, '/concurrent/ {print $6}' performance_comparison.csv)

    echo -e "顺序请求: 平均每请求 ${YELLOW}${seq_time}秒${NC}, 每秒 ${YELLOW}${seq_rps}${NC} 请求"
    echo -e "并发请求: 平均每请求 ${YELLOW}${con_time}秒${NC}, 每秒 ${YELLOW}${con_rps}${NC} 请求"

    # 计算性能提升
    speedup=$(echo "scale=2; $seq_time/$con_time" | bc)
    echo -e "并发处理性能提升: ${GREEN}${speedup}倍${NC}"
}

# 执行主函数
main