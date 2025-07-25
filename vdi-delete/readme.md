# UNVSDI - Dependency Injection Library for Go / DI库 / DIライブラリ / Thư viện DI

[![Go Report Card](https://goreportcard.com/badge/github.com/your-username/unvsdi)](https://goreportcard.com/report/github.com/your-username/unvsdi)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

<details>
<summary><strong>English</strong></summary>

A flexible Dependency Injection (DI) library for Go, supporting **Singleton**, **Scoped**, and **Transient** lifecycles, with automatic dependency resolution and circular dependency detection.

## 📌 Features
- ✅ Generic `TOwner` support: Dependencies know their owner (e.g., `Singleton[App, Logger]`).
- ✅ Embedded struct integration.
- ✅ Automatic circular dependency detection.
- ✅ Clear lifecycle: Singleton/Scoped/Transient.
- ✅ Easy mocking for unit tests.

</details>

<details>
<summary><strong>中文</strong></summary>

一个灵活的Go依赖注入(DI)库，支持**Singleton**、**Scoped**和**Transient**生命周期，具有自动依赖解析和循环依赖检测功能。

## 📌 功能
- ✅ 支持泛型`TOwner`：依赖知道其所有者(例如`Singleton[App, Logger]`)
- ✅ 嵌入式结构体支持
- ✅ 自动检测循环依赖
- ✅ 明确的生命周期：Singleton/Scoped/Transient
- ✅ 便于单元测试mock

</details>

<details>
<summary><strong>日本語</strong></summary>

Go用の柔軟な依存性注入(DI)ライブラリ。**Singleton**、**Scoped**、**Transient**のライフサイクルをサポートし、自動的な依存関係解決と循環依存検出機能を備えています。

## 📌 特徴
- ✅ ジェネリック`TOwner`サポート：依存関係が所有者を認識(例:`Singleton[App, Logger]`)
- ✅ 埋め込み構造体の統合
- ✅ 循環依存の自動検出
- ✅ 明確なライフサイクル：Singleton/Scoped/Transient
- ✅ 単体テストのモック作成が容易

</details>

<details>
<summary><strong>Tiếng Việt</strong></summary>

Thư viện Dependency Injection (DI) linh hoạt cho Go, hỗ trợ vòng đời **Singleton**, **Scoped** và **Transient**, với khả năng tự động resolve dependency và phát hiện circular dependency.

## 📌 Tính năng
- ✅ Hỗ trợ Generic (`TOwner`): Dependency biết ai sở hữu nó (ví dụ: `Singleton[App, Logger]`)
- ✅ Tích hợp embedded struct
- ✅ Tự động phát hiện circular dependency
- ✅ Vòng đời rõ ràng: Singleton/Scoped/Transient
- ✅ Dễ dàng mocking cho unit test

</details>

---

## 🚀 Installation / 安装 / インストール / Cài đặt
```bash
go get github.com/nttlong/igorm/unvs-di