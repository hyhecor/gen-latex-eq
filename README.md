# gen-latex-eq

## 기능 

- converter

- Usage
    
      $ gen-latex-eq < [스트림으로 넘겨줄 파일 이름]

## 테스트

    ## eq 파일 만들기
    cat <<EOF > latex-eq
    2n.svg          = 2n
    factorial.svg   = n!= \prod_{k=1}^{n} = n \cdot (n-1) \cdot (n-2) \cdot \cdot \cdot \cdot \cdot 3 \cdot 2 \cdot 1
    EOF
    ## 테스트 gen-latex-eq을 실행하여 eq파일을 latex 수식 이미지로 변환
    ./gen-latex-eq < latex-eq


   
 
