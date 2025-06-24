# UNVSDI - Dependency Injection Library for Go

[![Go Report Card](https://goreportcard.com/badge/github.com/your-username/unvsdi)](https://goreportcard.com/report/github.com/your-username/unvsdi)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Một thư viện Dependency Injection (DI) linh hoạt cho Go, hỗ trợ **Singleton**, **Scoped**, và **Transient** lifecycle, với khả năng tự động resolve dependency và phát hiện circular dependency.

## 📌 Features
- ✅ Hỗ trợ Generic (`TOwner`): Dependency biết ai sở hữu nó (ví dụ: `Singleton[App, Logger]`).
- ✅ Tích hợp **embedded struct**.
- ✅ Tự động phát hiện **circular dependency**.
- ✅ Lifecycle rõ ràng: Singleton/Scoped/Transient.
- ✅ Dễ dàng mocking cho unit test.

---

## 🚀 Cài đặt
```bash
go get github.com/your-username/unvsdi