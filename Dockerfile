# 빌드 스테이지
FROM golang:1.23-alpine AS builder

# 작업 디렉토리 설정
WORKDIR /app

# 모듈 초기화 및 의존성 설치
COPY go.mod go.sum ./
RUN go mod download

# 소스 코드 복사
COPY . .

# 애플리케이션 빌드
RUN go build -o /app/sendmind-hub

# 런타임 스테이지
FROM alpine:latest

# CA 인증서 설치
RUN apk --no-cache add ca-certificates

# 빌드된 바이너리를 런타임 스테이지로 복사
COPY --from=builder /app/sendmind-hub /sendmind-hub

# 컨테이너가 시작될 때 실행될 명령
CMD ["/sendmind-hub"]
