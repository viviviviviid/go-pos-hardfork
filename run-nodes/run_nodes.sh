#!/bin/bash

echo "Setting up the node. It will take approximately 15 seconds."

# 로딩바의 길이를 정의
bar_length=30
loading_time=15 # 총 로딩 시간
interval=$((loading_time * 2 / bar_length))

# 로딩바 시작
(
for ((i=0; i<=$bar_length; i++)); do
    # 진행 상황에 따라 = 문자 출력
    printf "["
    for ((j=0; j<$i; j++)); do printf "="; done

    # 나머지 부분은 공백으로 출력
    for ((j=i; j<$bar_length; j++)); do printf " "; done
    printf "]"

    # 커서를 처음으로 돌려서 로딩바를 업데이트
    echo -ne '\r'

    # interval 동안 대기
    sleep 0.5
done
) & # 백그라운드에서 로딩바 프로세스를 실행

mkdir -p logs

sleep 1
go run ../main.go -mode=auto -port=3000 > logs/log_3000.log 2>&1 &
sleep 1

for ((i=0; i<10; i++)); do
    port=$((4000 + $i))
    go run ../main.go -mode=rest -port=$port > logs/log_$port.log 2>&1 &
    sleep 1
    if [ "$port" -ne 4000 ]; then
        curl -X POST http://localhost:4000/peer -H "Content-Type: application/json" -d "{\"address\": \"127.0.0.1\", \"port\": \"$port\"}"
    fi
done

sleep 1
curl -X POST http://localhost:4000/peer -H "Content-Type: application/json" -d '{"address": "127.0.0.1", "port": "3000"}'
sleep 1

curl http://localhost:3000/wallet
for ((i=0; i<10; i++)); do
    port=$((4000 + $i))
    curl http://localhost:$port/wallet
done

sleep 3
echo
echo
echo
echo "Now, 10 Nodes are Running! Good Luck!"