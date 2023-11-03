#!/bin/bash

# -mode=rest로 실행된 프로세스들의 PID를 가져온다
pids=$(ps aux | grep '\-mode=auto' | grep -v grep | awk '{print $2}')

# 각 PID를 종료한다
for pid in $pids; do
    kill -9 $pid
done

echo "All nodes with -mode=auto have been terminated."

# -mode=rest로 실행된 프로세스들의 PID를 가져온다
pids=$(ps aux | grep '\-mode=rest' | grep -v grep | awk '{print $2}')

# 각 PID를 종료한다
for pid in $pids; do
    kill -9 $pid
done

echo "All nodes with -mode=rest have been terminated."
