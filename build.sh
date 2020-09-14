#!/bin/bash
## set version `git describe --tags $(git rev-list --tags --max-count=1)`
VERSION=`git describe --tags $(git rev-list --tags --max-count=1)`
## set build `git log -1 --pretty=format:%h`
BUILD=`git log -1 --pretty=format:%h`

## go build
go build -ldflags "-X main.version=${VERSION}@${BUILD}" 

## test help
./gen-latex-eq -h

## test function
## latex-eq 파일 만들기
cat <<EOF > latex-eq@${BUILD}
2n@${BUILD}.svg          = 2n
factorial@${BUILD}.svg   = n!= \prod_{k=1}^{n} = n \cdot (n-1) \cdot (n-2) \cdot \cdot \cdot \cdot \cdot 3 \cdot 2 \cdot 1
EOF
## 테스트 gen-latex-eq을 실행하여 latex-eq파일을 latex 수식 이미지로 변환
./gen-latex-eq < latex-eq@${BUILD}
## latex-eq 파일제거
rm *@${BUILD}*
