# 🧠 🎯 VISIÓN DEL PROYECTO

> **Sistema de Data Warehouse con ETL concurrente en Go y API REST para análisis estadístico universitario**

---

## 🧱 📦 ARQUITECTURA DEFINIDA (baseline)

```text
cmd/
 ├── api/        → servidor REST
 ├── etl/        → runner ETL

internal/
 ├── extract/
 ├── transform/
 ├── load/
 ├── pipeline/
 ├── repository/
 ├── domain/
 ├── service/

pkg/
 ├── db/
 ├── logger/
 ├── config/
```
