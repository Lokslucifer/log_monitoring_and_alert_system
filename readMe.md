# 🔍 Log Monitoring and Alert System

A distributed cloud-ready system that **collects**, **processes**, and **analyzes logs in real time**, and triggers **alerts** for critical events.  
Built with **Go**, **Kafka**, and **RabbitMQ**, designed for scalability and fault tolerance.

---

## 🧠 Overview

This project simulates a production-grade log monitoring pipeline:
- Collects logs from servers.
- Streams them via Kafka.
- Parses, stores, and analyzes logs.
- Sends alerts via RabbitMQ → Slack/Email.

---

## ⚙️ Architecture

    ┌──────────┐        ┌──────────┐        ┌────────────┐        ┌────────────┐
    │ Collector│──────▶│   Kafka  │──────▶│  Processor │──────▶│   RabbitMQ  │
    └──────────┘        └──────────┘        └────────────┘        └────────────┘
                                                                    │
                                                                    ▼
                                                             ┌────────────┐
                                                             │ Alert System│
                                                             └────────────┘


---

## 🧩 Modules

### 1. Collector
- Reads raw logs from files or streams.
- Handles file rotation and streaming.
- Pushes logs to **Kafka** for centralized processing.

### 2. Processor
- Consumes logs from **Kafka**.
- Parses and classifies logs (regex/JSON-based).
- Stores parsed logs in **PostgreSQL** or **SQLite**.
- Forwards critical logs to **RabbitMQ** for alerting.
- Gin Http server expose the logs via port 8090 for analysis purpose.

### 3. Alert System
- Listens to alerts from **RabbitMQ**.
- Sends notifications to **Slack**, **Email**, or any configured channel.
- Supports message filtering, rate limiting, and retries.
- Set SLACK_WEB_HOOK url in .env file in alert system

---



