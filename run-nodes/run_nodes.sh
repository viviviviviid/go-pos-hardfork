#!/bin/bash

# logs 폴더가 존재하지 않으면 생성
mkdir -p logs

go run ../main.go -mode=rest -port=3000 > logs/log_3000.log 2>&1 &
for ((i=0; i<10; i++)); do
    port=$((4000 + $i))
    go run ../main.go -mode=rest -port=$port > logs/log_$port.log 2>&1 &
done
