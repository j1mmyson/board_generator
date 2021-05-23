## Done

- 글 쓰기
- 글 목록 불러오기

## ToDo

- [x] 페이징
- [x] 검색
- [ ] 글 수정/삭제
- [ ] 댓글

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