# Stage 1: Build
FROM golang:1.22-alpine AS builder

# Thiết lập thư mục làm việc
WORKDIR /app

# Sao chép mã nguồn và cài đặt dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Sao chép toàn bộ mã nguồn
COPY . .

# Biên dịch ứng dụng
RUN go build -o myapp cmd/app/main.go

# Stage 2: Final
FROM scratch

# Sao chép file nhị phân từ stage 1
COPY --from=builder /app/myapp /myapp

# Chỉ định lệnh khởi chạy ứng dụng
CMD ["/myapp"]

# Cổng mà ứng dụng lắng nghe
EXPOSE 8080