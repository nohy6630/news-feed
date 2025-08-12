# --------------------
# 1. 빌더 스테이지: 애플리케이션 컴파일
# --------------------
# GoLang 공식 이미지 중 빌드에 필요한 환경이 갖춰진 이미지를 사용합니다.
FROM golang:1.24-alpine AS builder

# 작업 디렉터리를 설정합니다.
WORKDIR /app

# go.mod와 go.sum 파일을 복사하여 의존성을 캐싱합니다.
# 애플리케이션 코드가 변경되더라도 의존성 파일이 변경되지 않았다면 캐시를 활용하여 빌드 속도를 높입니다.
COPY go.mod go.sum ./
RUN go mod download

# 모든 애플리케이션 소스 코드를 복사합니다.
COPY . .

# 애플리케이션을 빌드합니다.
# CGO_ENABLED=0는 CGO를 비활성화하여 정적으로 링크된 바이너리를 생성합니다.
# -o main은 출력 파일 이름을 'main'으로 지정합니다.
RUN CGO_ENABLED=0 GOOS=linux go build -o news-feed bin/main.go

# --------------------
# 2. 실행 스테이지: 컴파일된 바이너리만 복사
# --------------------
# 가벼운 경량 OS인 alpine 이미지를 사용합니다.
FROM alpine:latest

# 빌더 스테이지에서 컴파일된 바이너리를 복사합니다.
# --from=builder는 "builder"라는 이름의 이전 스테이지에서 파일을 가져온다는 의미입니다.
COPY --from=builder /app/news-feed /news-feed

# 애플리케이션이 사용할 포트를 외부에 노출합니다.
EXPOSE 8081

# 컨테이너 시작 시 실행될 명령어를 지정합니다.
CMD ["/news-feed"]
