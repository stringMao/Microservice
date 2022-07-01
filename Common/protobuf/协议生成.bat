@echo off
set curdir=%cd%

::指令学习https://zhuanlan.zhihu.com/p/83010418
protoc  --proto_path=   --go_out=../proto/   *.proto
:: protoc  --go_out=../proto/   *.proto


::xcopy   .\proto\*  ..\GatewayClient\proto\  /s/e/y

::xcopy   .\proto\*  ..\GatewayFuwuqi\proto\  /s/e/y

PAUSE

