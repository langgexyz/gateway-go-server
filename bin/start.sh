# 如果有旧进程，先停止
if [ -f ./pid ]; then
    kill -TERM `cat ./pid`
fi

# 使用nohup启动Linux版本
nohup ./go-build-stream-gateway-go-server-linux -c ./config.debug.json > ./error.log 2>&1 & echo $! > pid
