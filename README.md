# Hướng dẫn sử dụng HaiTranslator

HaiTranslator là ứng dụng dịch tài liệu bằng Gemini. Ứng dụng gửi nội dung tệp
đã chọn đến Gemini, hiển thị bản dịch và lưu bản dịch thành tệp mới.

![Ảnh chụp màn hình ứng dụng HaiTranslator](screenshot-2026-10-19_17-39-48.png)
## Các phần trên giao diện
1. Vị trí số 1: Đây là nơi xuất hiện nội dung bản dịch
2. Vị trí số 2: Đây là nơi chọn các thông số đầu vào
**API KEY** : là đoạn mã dùng để xác thực khi gửi tin nhắn đến gemini.
**Ngôn ngữ nguồn**: chọn ngôn ngữ của tài liệu hoặc chọn **Tự động phát hiện**.
**Ngôn ngữ đích**: chọn ngôn ngữ cần dịch sang.
 **Chọn định dạng file dịch**: chọn định dạng tệp đầu ra: `.txt`, `.docx`, `.pdf`, `.md`, `.tex` hoặc `.odt`.
**Mô hình AI**: chọn mô hình Gemini dùng để dịch.
**Nút DỊCH NGAY VIA AI**: gửi tài liệu cùng các lựa chọn lên Gemini. Trong khi chờ hản hồi, thanh tiến trình uất hiện và nút dịch tạm thời bị khóa. Khi xong,bản dịch hiện ở khung phía trên và ứng dụng thông báo nơi lưu tệp.
**Mở thư mục input**: mở thư mục chứa các tệp đầu vào đã chọn.
**Mở thư mục output**: mở thư mục lưu bản dịch. 
3. Vị trí sô 3Đây là bảng thống kê số token đã sử dụng, và chi phí phải trả cho mỗi phiên dịch.
### Windows
Để cho đơn giản, tôi đã biên dịch thành file .exe, người dùng chỉ click đúp chuột vào là ứng dụng chạy.
Ứng dụng được biên dịch theo phương pháp plug and play nên không cài đặt vào hệ thống nên không chiếm dung lượng của ổ cứng.

## Linux 
Người dùng linux giải nén tập tin
Làm theo hướng dẫn ghi trong file Makefile

## Macos
Người dùng macOS chỉ cần kéo file cài đặt vào thư mục ứng dụng là xong.

## Tạo Gemini API key

Quy trình tạo key trên trang Google AI Studio giống nhau trên Windows, macOS
và Linux. Truy cập [Google AI Studio - API keys](https://aistudio.google.com/app/apikey)
bằng tài khoản Google, tạo hoặc chọn một Google Cloud project rồi chọn **Create
API key**. Lưu key an toàn; không đăng key lên mã nguồn, ảnh chụp màn hình hoặc
kho lưu trữ công khai. Ứng dụng hiện cần cùng một key ở biến môi trường
`GEMINI_API_KEY` và ô **API key** trên giao diện.

### API key với quota miễn phí

1. Mở trang API keys trong Google AI Studio và đăng nhập.
2. Chọn một project có thể dùng Gemini API hoặc tạo project mới.
3. Chọn **Create API key**, sao chép key và cấu hình theo bước 3 của hệ điều
   hành tương ứng ở trên.
4. Kiểm tra [quota và giới hạn sử dụng](https://ai.google.dev/gemini-api/docs/rate-limits)
   của project/mô hình trước khi dịch.

Quota miễn phí phụ thuộc vào khu vực, project, mô hình và chính sách hiện hành
của Google; không phải mọi mô hình hoặc tài khoản đều có cùng hạn mức. Kiểm tra
trang giá và điều khoản hiện tại trước khi gửi tài liệu:
[Gemini API pricing](https://ai.google.dev/gemini-api/docs/pricing).

### API key với thanh toán

API key không có một loại riêng gọi là “key trả phí”: mức tính phí gắn với
project Google Cloud mà key thuộc về. Để sử dụng theo mức trả phí:

1. Trong Google AI Studio, tạo key thuộc project cần dùng hoặc chọn project đó.
2. Mở [Google Cloud Billing](https://console.cloud.google.com/billing), thiết
   lập tài khoản thanh toán và liên kết tài khoản thanh toán với đúng project.
3. Kiểm tra mô hình có hỗ trợ cho project, giá, quota và giới hạn chi tiêu tại
   [Gemini API pricing](https://ai.google.dev/gemini-api/docs/pricing) và
   [Google AI Studio](https://aistudio.google.com/).
4. Cấu hình key như hướng dẫn hệ điều hành ở trên. Theo dõi mức sử dụng và hóa
   đơn trong Google Cloud; giữ key riêng tư và không đưa vào source code.
