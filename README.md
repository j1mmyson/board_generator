## Link

<http://ec2-3-17-39-222.us-east-2.compute.amazonaws.com/>

![Index Page](https://github.com/j1mmyson/board_generator/blob/main/img/indexPage.PNG?raw=true)

## Done

- 글 쓰기
- 글 목록 불러오기

## ToDo

- [x] 페이징
- [x] 검색
- [x] 글 수정
- [x] 글 삭제
- [ ] 댓글
- [x] 페이지 렌더링 시 현재 페이지를 굵고 클릭 불가능하게 변경
- [ ] `/board`에서 중복되는 코드 부분 함수로 묶기, 여러 파일로 나누기(main.go를 간단하게)

## How to run?

1. Clone this Repo

2. Create `account.go`  

   ```go
   package main
   
   const (
   	host     = <Host of Your DB>
   	database = <Name of Your DB>
   	user     = <User Name>
   	password = <Password>
   )
   
   ```

3. run `go build -o server`

4. run binary file by `./server`

5. Connect to `localhost:8080`