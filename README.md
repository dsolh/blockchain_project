# 실행방법 및 순서

## 1. Network 실행

whozard 경로에서 (다음 쉘스크립트 실행):

- ./startWhozardNetwork.sh
- ./deployWhozardCC.sh

## 2. CCP 생성

application/ccp 경로에서 (다음 쉘스크립트 실행):

- ./ccp-generate.sh

## 3. 인증서 가져오기 및 지갑 생성

application 경로에서 

- ./addToWallet.js

## 4. node.js 모듈 설치

application 경로에서

- npm install

## 5. invoke/query 실행 또는 서버 실행

application 경로에서

- node server.js : whozard 어플리케이션 웹서버 실행

## 6. 서버 종료

application 경로에서

- rm -rf wallet

whozard 경로에서 (다음 쉘스크립트 실행):

- ./downWhozardNetwork.sh