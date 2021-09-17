### 
  cpp  hpp   common.h   common.h  chain33   chain33 。   

    。

### 

#### Emscripten 
```bash
git clone https://github.com/juj/emsdk.git
cd emsdk
git checkout 6adb624e04b0c6a0f4c5c06d3685f4ca2be7691d # 
./emsdk install latest
./emsdk activate latest

#  

# on Linux or Mac OS X
source ./emsdk_env.sh

# on Windows
emsdk_env.bat
```

####  Emscripten 

```bash
em++ -o dice.wasm dice.cpp -s WASM=1 -O3 -s EXPORTED_FUNCTIONS="[_startgame, _deposit, _play, _draw, _stopgame]" -s ERROR_ON_UNDEFINED_SYMBOLS=0
```

- em++  wasm   c++    emcc  
- -o   .wasm   
- dice.cpp    
- -s WASM=1  wasm    
- -O3   wasm   1～3   
- -s EXPORTED_FUNCTIONS      "_"
- -s ERROR_ON_UNDEFINED_SYMBOLS=0   common.h  c/c++   chain33  go   

 ：https://developer.mozilla.org/en-US/docs/WebAssembly

####  wabt（the WebAssembly Binary Toolkit）

```bash
git clone --recursive https://github.com/WebAssembly/wabt
cd wabt
mkdir build
cd build
cmake ..
cmake --build .
```

```bash
#  wasm abi，ab impor  expor 。
wabt/bin/wasm2wat dice.wasm
```

### 
```bash
#  
./chain33-cli send wasm create -n  -p was  -k 
```

### 
```bash
# 
./chain33-cli wasm check -n 
```

### 
```bash
#  
./chain33-cli send wasm update -n  -p was  -k 
```

### 
```bash
  
./chain33-cli send wasm call -n  -m  -p  -v  -k   
```

### 
```bash
# statedb
./chain33-cli wasm query state -n  -k key  

# localdb
./chain33-cli wasm query local -n  -k key  
```

### 
```bash
   wasm 
./chain33-cli send coins send_exec -e wasm -a  -k 

 
./chain33-cli send coins withdraw -e wasm -a  -k 
```
### RP 
```bash
 
curl http://localhost:8801 -ksd '{"method":"wasm.CallContract", "params":[{"contract":"dice","method":"play","parameters":[1000000000, 10]}]}'

 
curl http://localhost:8801 -ksd '{"method":"Chain33.SignRawTx","params":[{"privkey":"0x4257d8692ef7fe13c68b65d6a52f03933db2fa5ce8faf210b5b8b80c721ced01","txhex":"0x0a047761736d1218180212140a04646963651204706c61791a068094ebdc030a20a08d0630ec91c19ede9ef4d1693a22314b3732554137393845775a66427855546b4265686864766b656f3277377446344c","expire":"300s"}]}'

 
curl http://localhost:8801 -ksd '{"method":"Chain33.SendTransaction","params":[{"data":"0a047761736d1218180212140a04646963651204706c61791a068094ebdc030a1a6d080112210320bbac09528e19c55b0f89cb37ab265e7e856b1a8c388780322dbbfd194b52ba1a46304402201dc04e89da9220e42b2768a23cd2e6a7c452b2bfd30e0799f5c6f1b035151d1402201160929f74feb26be4205cf4432bdf377eb775f189db2883556cedc31c4fb01920a08d0628b2cb90fb0530ec91c19ede9ef4d1693a22314b3732554137393845775a66427855546b4265686864766b656f3277377446344c"}]}'
```