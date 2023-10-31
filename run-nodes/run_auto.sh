#!/bin/bash

echo "Setting up the node"
sleep 1
go run ../main.go -mode=auto -port=3000 > logs/log_3000.log 2>&1 &

for ((i=0; i<10; i++)); do
    port=$((4000 + $i))
    go run ../main.go -mode=auto -port=$port > logs/log_$port.log 2>&1 &
done

