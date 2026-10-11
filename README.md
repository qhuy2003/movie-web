# CineStream — Website đặt vé xem phim

CineStream là website khám phá phim và đặt vé xem phim được xây dựng theo mô hình full-stack. Người dùng có thể xem thông tin phim, chọn rạp, suất chiếu, ghế ngồi và quản lý các vé đã đặt.

## Tính năng

- Duyệt và tìm kiếm phim; xem thể loại, năm phát hành, điểm đánh giá IMDb, thời lượng, nội dung, diễn viên và poster.
- Đăng ký, đăng nhập và quản lý danh sách phim yêu thích.
- Trang quản trị để quản lý phim, thể loại và người dùng.
- Quy trình đặt vé gồm chọn rạp, suất chiếu và ghế ngồi.
- Tạm giữ ghế trong khi chờ thanh toán và theo dõi trạng thái đơn vé.
- Tích hợp thanh toán qua ZaloPay Sandbox.
- Trang “Vé của tôi” phân loại vé đã mua, chờ thanh toán, đã hủy và hết hạn.
- Tự tạo dữ liệu mẫu cho rạp, phòng chiếu và ghế khi khởi động backend.

## Công nghệ sử dụng

- **Frontend:** Vue 3, Vue Router, Vite
- **Backend:** Go, Gin, GORM
- **Cơ sở dữ liệu:** PostgreSQL
- **Thanh toán:** ZaloPay Sandbox

## Yêu cầu

- Go theo phiên bản khai báo trong `go.mod`
- Node.js và npm
- PostgreSQL

## Chạy dự án trên máy cá nhân

### 1. Cấu hình PostgreSQL

Tạo cơ sở dữ liệu PostgreSQL (tên mặc định: `movie_db`). Backend đọc các biến môi trường sau; nếu không khai báo, ứng dụng sẽ dùng giá trị mặc định:

| Biến môi trường | Giá trị mặc định |
| --- | --- |
| `DB_HOST` | `127.0.0.1` |
| `DB_PORT` | `5433` |
| `DB_USER` | `postgres` |
| `DB_PASSWORD` | `postgres` |
| `DB_NAME` | `movie_db` |

Hãy cấu hình các giá trị phù hợp với PostgreSQL trên máy của bạn trước khi chạy backend. Khi khởi động, GORM sẽ tạo hoặc cập nhật các bảng cần thiết và thêm dữ liệu mẫu về rạp chiếu.

### 2. Cấu hình ZaloPay (không bắt buộc)

Sao chép `.env.example` thành `.env`, sau đó điền thông tin Sandbox của ZaloPay vào `ZALOPAY_KEY1` và `ZALOPAY_KEY2`. File mẫu đã có sẵn App ID và địa chỉ Sandbox.

Nếu chưa có thông tin merchant hợp lệ, bạn vẫn có thể khám phá phim và thử quy trình đặt vé, nhưng không thể hoàn tất thanh toán ZaloPay.

Không đưa file `.env`, API key, mật khẩu cơ sở dữ liệu hoặc thông tin bí mật vào GitHub.

### 3. Chạy backend Go

Tại thư mục gốc của dự án:

```bash
go run ./cmd/server
```

API chạy tại `http://localhost:8080`.

### 4. Chạy frontend Vue

Mở terminal khác và chạy:

```bash
cd frontend
npm install
npm run dev
```

Mở địa chỉ mà Vite hiển thị trong terminal, thường là `http://localhost:5173`.

## Thiết lập tài khoản quản trị

Trên cơ sở dữ liệu mới chưa có tài khoản, tài khoản đầu tiên đăng ký qua website sẽ được cấp quyền `admin`. Hãy đăng ký email quản trị trước khi tạo tài khoản người dùng khác.

Nếu cơ sở dữ liệu đã có người dùng, hãy đăng nhập bằng tài khoản admin hiện có và đổi quyền người dùng trong trang quản trị.

Vì lý do bảo mật, README không chứa mật khẩu admin. Hãy dùng mật khẩu riêng, không công khai mật khẩu trong README, mã nguồn, ảnh chụp màn hình hoặc lịch sử commit. Nếu mật khẩu thật đã được chia sẻ công khai, hãy đổi mật khẩu trước khi đưa dự án lên GitHub.

## Cấu trúc thư mục

```text
cmd/server/       Điểm khởi chạy Go server và khai báo API routes
internal/handler/ HTTP handlers, xác thực và xử lý đặt vé
internal/model/   Các model cơ sở dữ liệu và API
internal/db/      Kết nối PostgreSQL, migration và dữ liệu rạp mẫu
frontend/src/     Ứng dụng Vue, các trang, component và API services
```

## Ghi chú

- Dữ liệu phim có các trường thể loại, năm phát hành và điểm đánh giá; dự án có dữ liệu phim mẫu để dùng khi phát triển.
- ZaloPay Sandbox chỉ dành cho môi trường thử nghiệm.
- Frontend hiện gọi API tại `http://localhost:8080`.

## Giấy phép

Nếu muốn chia sẻ hoặc cho phép người khác sử dụng dự án, hãy bổ sung file giấy phép phù hợp.
