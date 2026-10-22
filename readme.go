package main

var readme string = `[Title] Hướng dẫn sử dụng phần mềm dịch bằng AI
1. Giới thiệu về phần mềm
- Đây là ứng dụng tân dụng sức mạnh của AI để dịch tài liệu, sách giữa các ngôn ngữ khác nhau như Anh <=> Việt, Việt <=> Đức, Việt <=> Trung,...
- Ngôn ngữ hỗ trợ bao gồm tất cả các ngôn ngữ trên thế giới (Nếu thiếu ngôn ngữ mà bạn cần hãy liên hệ để tôi bổ sung).
- Phần mềm hỗ trợ tất cả mô hình, hiện tại chỉ mới hỗ trợ Gemini là mô hình AI do Google pahts triển, và Gemma-4 một biến thể của Gemini.

2. Hướng dẫn sử dụng
2.1. Kết quả dịch sẽ xuất hiện ở đây.
Phía bên dưới:
2.2. Bên trái là bảng điều khiển, người dùng cần nhập các thông tin sau:
 	- Api key: đây là mã bí mật dùng để giao tiếp với mô hình AI, 
	- Ngôn ngữ nguồn: là ngôn ngữ của tài liệu, chế độ mặc định là "tự động phát hiện"
	- Ngôn ngữ dịch: là ngôn ngữ mà bạn muốn tài liệu được dịch sang. Lưu ý ngôn ngữ nguồn phải khác ngôn ngữ đích, nếu không phần mềm sẽ báo lỗi.
	- Định dạng tài liệu dịch: phần mềm có thể tạo ra các file như word, pdf, text, markdown,..
	- Chọn mô hình AI: người dùng nên tham khảo đơn giá kế bên để chọn mô hình phù hợp, mô hình các mới và càng thông minh thì mức giá càng cao.
# Các nút bên dưới
- Nút "Dịch byAI": là nút để ứng dụng tiến hành dịch tài liệu sau khi người dùng đã chọn các thông số ở trên.
- Nút Mở thư mục tài liệu gốc: dùng để mở thư mục lưu tài liệu gốc do người dùng tải lên.
- Nút Mở thư mục tài liệu dịch: dùng để mở thư mực chứa các tài liệu đã dịch thành công.
- Nút Lịch sử dịch: mở lên bảng theo dõi các tài liệu đã được dịch.
- Nút Chi phí đã sử dụng: mở lên bảng thống kê số token đã sử dụng và số tiền mà người dùng phải trả cho mỗi phiên dịch, được tổng hợp theo mô hình AI.
=> Thư mục và thông tin trên được lưu tại thư mục ~\user\haitranslator.
1.3. Bên phải là bảng đơn giá của từng mô hình, được tính trên 1 triệu token.
Để người dùng mới dễ hình dung, với 1 triệu token thì bạn có thể dịch toàn bộ một cuốn sách 300 - 400 trang. Và chi phí phải trả vào tầm 3.75usd x 26,000 ~ 45,500 vnđ. 


3. Hướng dẫn cài đặt và chạy phần mềm
3.1. Window và MacOS
Đối với người đùng sử dụng hệ điều hành window và macOs thì chỉ cần double click chuột là chạy.
- Đối với Window trong lần đầu chạy thì hệ thống bảo mật của Window yêu cầu xác nhận.
- Đối với MacOs: tôi chưa thử vì không có thiết bị chạy MacOs.
3.2. Linux
Đối với người dùng sử dụng bản phân phối của Linux tiến hành như sau:
- Giải nén file
- Mở Terminal và di chuyển vào thư mục vừa giải nén ở trên.
- Dùng lệnh make user-install: nếu chỉ muốn cài cho người dùng hiện tại
- Dùng lệnh sudo make install để cài đặt cho toàn hệ thống.
* Đọc thêm file MakeFile để tham khảo.


4. Api key ![VIDEO][https://youtu.be/vbzTuMhg6lA]
Để có thể dịch được tài liệu thì yêu cầu bắt buộc là phải có Apikey.
4.1. Gemini
Người dùng có thể sử dụng phiên bản miễn phí để test phần mềm.
- Để lấy apikey miễn phí thì người dùng truy cập trang Google AI studio ( https://aistudio.google.com ), đăng nhập bằng tài khoản Gmail.
Sau khi đăng nhập thì chọn biểu tưởng chìa khóa (góc bên dưới bên trái của màn hình), kế đến chọn Tạo khóa mới.
Sau khi test thành công và chuyển qua sử dụng chính thức thì người dùng cần sử dụng apikey có trả phí (lưu ý là bạn đnag trả tiền cho Google)
Hãy hỏi bất kỳ mô hình AI nào cách đăng ký tài khoản trả phí của google cloud.
4.2 OpenAI

5. Bản quyền
Đây là phần mềm do tôi (Ngô Hoài Nam) viết bằng golang và biên dịch.
Tôi đồng ý cho tất cả sử dụng MIỄN PHÍ với mục đích phi thương mại, nghĩa là bạn có thể sử dụng phần mềm này để tạo ra các tài liệu, sách biên dịch mà không phải trả phí bản quyền.
Đồng thời tài liệu dịch của bạn do bạn tự ý quyết định, bạn có thể bán hoặc cho thuê các tài liệu đó.


6. Tài trợ và Tip
- Việc phát triển phần mềm cần một số chi phí như chi trả điện, mua các thiết bị cần thiết để thử nghiệm,... vì vậy tôi cũng cần các nguồn tài trợ.
Đổi lại thương hiệu, sản phẩm của nhà tài trợ sẽ xuất hiện ở đây.

- Tip: bạn có thể mời tôi ly cà phê nếu phần mềm này hữu ích với bạn.
	Số tài khoản: 0121002153944
	Ngần hành VietComBank chi nhánh Đồng Nai
	Chủ tài khoản: Ngô Hoài Nam

`
