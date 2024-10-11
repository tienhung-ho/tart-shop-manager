package trainsaction

// Định nghĩa một loại khóa riêng để tránh xung đột với các khóa khác trong context
type contextKey string

// TransactionKey được sử dụng để lưu trữ transaction trong context
const TransactionKey contextKey = "transaction"
