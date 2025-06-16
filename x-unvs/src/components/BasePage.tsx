// src/components/BasePage.tsx

import React from 'react';
import Caller from '../utils/Caller'; // Import class Caller

// --- THAY ĐỔI Ở ĐÂY ---
// Định nghĩa props cơ bản cho BasePage, BẮT BUỘC phải có featureId
interface BasePageProps {
  featureId: string; // Đây là Feature ID bắt buộc cho mỗi trang
  // Các props chung khác (nếu có)
  // Ví dụ: tenantName?: string;
}

// Định nghĩa state rỗng hoặc chung cho BasePage nếu có
interface BasePageState {
  // Có thể thêm các state chung cho tất cả các trang ở đây nếu cần
}

// BasePage là một Class Component kế thừa từ React.Component
// P là kiểu của props của component con, S là kiểu của state của component con
class BasePage<P = {}, S = {}> extends React.Component<P & BasePageProps, S & BasePageState> {
  // Khai báo thuộc tính 'caller' với kiểu là Caller
  protected caller: typeof Caller = Caller;

  // Thuộc tính để lưu trữ Feature ID của trang
  protected featureId: string; // Thêm thuộc tính featureId

  // --- THAY ĐỔI Ở ĐÂY ---
  // Constructor BẮT BUỘC phải nhận featureId
  constructor(props: P & BasePageProps) {
    super(props);
    this.featureId = props.featureId; // Lấy featureId từ props và gán vào thuộc tính
    
    // Ghi log để kiểm tra featureId đã được truyền đúng chưa
    console.log(`BasePage: Initialized with Feature ID: ${this.featureId}`);

    // Bạn có thể thêm logic kiểm tra quyền ở đây hoặc ở ProtectedRoute
    // Ví dụ: if (!this.checkPermissions(this.featureId)) { /* redirect */ }

    // Nếu bạn có constructor ở đây, hãy đảm bảo gọi super(props);
    // và khởi tạo state nếu cần
    // this.state = {
    //   ...(this.state as any) // Để các class con có thể thêm state của riêng chúng
    // } as S & BasePageState;
  }

  // --- TÙY CHỌN: Phương thức kiểm tra quyền trong BasePage ---
  // Đây là một ví dụ, logic thực tế có thể phức tạp hơn (ví dụ: gọi API kiểm tra quyền)
  protected checkPermissions(featureId: string): boolean {
    // Trong thực tế, bạn sẽ kiểm tra quyền của người dùng hiện tại
    // dựa trên featureId. Ví dụ:
    // const userPermissions = localStorage.getItem('userPermissions');
    // return userPermissions.includes(featureId);
    console.log(`Checking permissions for Feature ID: ${featureId}`);
    // Tạm thời trả về true để cho phép tất cả các API call
    return true; 
  }

  // Phương thức render() phải được các class con override.
  render() {
    return null;
  }
}

export default BasePage;