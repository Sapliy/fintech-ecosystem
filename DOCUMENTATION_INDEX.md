# 📋 Sapliy Documentation Index

Welcome to Sapliy - the event-driven automation & policy platform for fintech and business flows.

> **Hybrid SaaS + Self-Hosted**: One codebase for small/medium companies (SaaS) and large enterprises (self-hosted). Open-source first, enterprise-ready.

---

## 📚 Documentation Overview

### For Quick Reference

- **[Quick Reference Guide](./QUICK_REFERENCE.md)** - Fast lookup for common tasks, APIs, and troubleshooting

### For Understanding the Platform

- **[Architecture](./ARCHITECTURE.md)** - Complete system design, deployment models, security, and technical decisions
- **[Business Model](./BUSINESS_MODEL.md)** - Pricing tiers, customer segments, acquisition strategy, and revenue projections

### For Enterprise Deployment

- **[Enterprise Guide](./ENTERPRISE_GUIDE.md)** - Complete guide for large organizations deploying self-hosted
  - Infrastructure planning (AWS, GCP, on-prem)
  - Security architecture & hardening
  - Deployment procedures & runbooks
  - Compliance (HIPAA, PCI-DSS, FedRAMP)
  - Disaster recovery & operational excellence

---

## 🚀 Getting Started

### **I'm a Developer** (Want to build with Sapliy)

1. Start with [Quick Reference - Local Setup](./QUICK_REFERENCE.md#1-local-setup)
2. Read [Architecture - Core Mental Model](./ARCHITECTURE.md#core-mental-model)
3. Check out [Quick Reference - Common Endpoints](./QUICK_REFERENCE.md#api-reference-common-endpoints)
4. Explore the [SDKs](#-supported-sdks-and-packages)

### **I'm an Architect** (Need to understand system design)

1. Read [Architecture Overview](./ARCHITECTURE.md)
2. Review [Deployment Models](./ARCHITECTURE.md#deployment-models)
3. Check [Security Model](./ARCHITECTURE.md#security-model)
4. See [Performance & Scalability](./ARCHITECTURE.md#performance--scalability)

### **I'm an Enterprise Customer** (Planning self-hosted deployment)

1. Start with [Enterprise Guide - Overview](./ENTERPRISE_GUIDE.md#enterprise-overview)
2. Review [Pre-Deployment Planning](./ENTERPRISE_GUIDE.md#pre-deployment-planning)
3. Follow [Deployment Strategies](./ENTERPRISE_GUIDE.md#deployment-strategies)
4. Set up [Compliance & Audit](./ENTERPRISE_GUIDE.md#compliance--audit)

### **I'm a Product/Business Person** (Need to understand the business)

1. Read [Business Model Overview](./BUSINESS_MODEL.md#business-model-overview)
2. Review [Customer Segments](./BUSINESS_MODEL.md#customer-segments)
3. Check [Pricing & Revenue Strategy](./BUSINESS_MODEL.md#saas-pricing--growth-strategy)
4. See [Competitive Positioning](./BUSINESS_MODEL.md#competitive-positioning)

---

## 📦 Supported SDKs and Packages

| Language             | Package                     | Status | Repository                                                         |
| -------------------- | --------------------------- | ------ | ------------------------------------------------------------------ |
| **Node.js**          | `@sapliyio/fintech`         | ✅ GA  | [fintech-sdk-node](https://github.com/sapliy/fintech-sdk-node)     |
| **Python**           | `sapliyio-fintech`          | ✅ GA  | [fintech-sdk-python](https://github.com/sapliy/fintech-sdk-python) |
| **Go**               | `fintech-sdk-go`            | ✅ GA  | [fintech-sdk-go](https://github.com/sapliy/fintech-sdk-go)         |
| **React Components** | `@sapliyio/fintech-ui`      | ✅ GA  | [fintech-ui](https://github.com/sapliy/fintech-ui)                 |
| **Testing**          | `@sapliyio/fintech-testing` | ✅ GA  | [fintech-testing](https://github.com/sapliy/fintech-testing)       |
| **CLI**              | `@sapliyio/sapliy-cli`      | ✅ GA  | [sapliy-cli](https://github.com/sapliy/sapliy-cli)                 |

---

## 🏗️ Core Concepts (30-Second Overview)

### **Organization**

- Your top-level account
- Contains users, teams, policies
- Root of all access control

### **Zone** 🔑

- Isolated automation space
- Has test mode & live mode
- Contains flows, events, logs
- Think: "Stripe Account + Webhook Endpoint combined"

### **Event** ⚡

- Anything that happens in your app
- Emitted from SDK: `sapliy.emit('checkout.clicked', {...})`
- Triggers flows automatically
- No event = nothing happens

### **Flow** 🔄

- Automated response to events
- Listens to event type + zone
- Executes actions: webhooks, notifications, audit logs
- This is the core value proposition

**Example Flow:**

```
Event: "payment.completed"
  ↓
Flow: "send_confirmation_email"
  ↓
Actions: Send webhook to email service
         Create ledger entry
         Log audit trail
```

---

## 🌍 Deployment Models at a Glance

### **SaaS (Managed by Sapliy)**

```
api.sapliy.io
├─ Target: Startups, SMBs, fast-growing companies
├─ Setup time: Minutes
├─ Maintenance: Zero (we handle everything)
├─ Security: SOC 2, DDoS protection
├─ Scaling: Automatic
└─ Cost: $0-$99/month (free to Pro tier)
```

### **Self-Hosted (Your Infrastructure)**

```
fintech.yourcompany.com
├─ Target: Enterprises, regulated industries
├─ Setup time: 4-6 weeks
├─ Maintenance: Your team
├─ Security: Your VPC/on-prem, custom policies
├─ Scaling: Manual + Kubernetes
└─ Cost: $2K-$500K+/year (depending on scale)
```

**Key Point**: Both use the **same codebase**. The difference is where it runs.

---

## 🔐 Security Highlights

### Encryption

- ✅ AES-256 at rest
- ✅ TLS 1.3 in transit
- ✅ Field-level encryption for sensitive data
- ✅ Customer-managed keys (CMK) available

### Compliance

- ✅ SOC 2 Type II certified (SaaS)
- ✅ HIPAA-ready with BAA support
- ✅ PCI-DSS compliant (no card storage)
- ✅ GDPR/CCPA data handling
- ✅ Audit logs (immutable, tamper-proof)

### Access Control

- ✅ RBAC (Role-Based Access Control)
- ✅ ABAC (Attribute-Based Access Control)
- ✅ MFA (Multi-Factor Authentication)
- ✅ API key scoping & rotation
- ✅ IP whitelisting (Enterprise)

---

## 📊 Architecture Quick View

```
┌─────────────────────┐
│  SDKs / CLI / UI    │
│ (Node, Go, Python)  │
└──────────┬──────────┘
           │
           ↓
┌─────────────────────────────────────┐
│   API Gateway (8080)                │
│   Auth Service (8081)               │
│   Payments Service (8082)           │
│   Ledger Service (8083)             │
│   Zone Manager                      │
│   Flow Engine                       │
└──────────┬──────────────────────────┘
           │
    ┌──────┴───────┬─────────┬─────────┐
    ↓              ↓         ↓         ↓
┌────────┐   ┌──────────┐ ┌────────┐ ┌──────┐
│  PG DB │   │  Kafka   │ │ Redis  │ │RabMQ │
│ Audit  │   │ Events   │ │ Cache  │ │Queues│
│ Flows  │   │ Ledger   │ └────────┘ └──────┘
│ Events │   └──────────┘
└────────┘

[Webhooks / Notifications / Audit Logs]
```

---

## 🚢 Development Roadmap

### Phase 1: MVP / SaaS Launch ✅ (Months 1-3)

- Core backend, event ingestion, flow execution
- Node.js SDK
- Test/live mode support
- Basic SOC 2 compliance

### Phase 2: Open-Source & SDK Expansion 🔄 (Months 4-6)

- Publish as open-source (MIT license)
- Python & Go SDKs
- Testing toolkit
- Example applications

### Phase 3: Self-Hosted Option 🎯 (Months 7-9)

- Docker images & Kubernetes Helm charts
- AWS/GCP/Azure deployment guides
- HIPAA/PCI-DSS documentation
- Enterprise licensing & support

### Phase 4: Monetization & Advanced Features 💰 (Months 10-12)

- Usage-based pricing
- Premium plugins & integrations
- Advanced analytics & compliance reports
- AI-powered flow suggestions

---

## 💰 Pricing Summary

### SaaS Tiers

| Tier       | Price  | Events/Month | Zones | Live Mode |
| ---------- | ------ | ------------ | ----- | --------- |
| Free       | $0     | 1K           | 1     | ❌        |
| Starter    | $29    | 10K          | 3     | ✅        |
| Pro        | $99    | 100K         | ∞     | ✅        |
| Enterprise | Custom | ∞            | ∞     | ✅        |

### Self-Hosted Licenses (Annual)

| License    | Price  | Employees | Deployment    |
| ---------- | ------ | --------- | ------------- |
| Startup    | $1,999 | <50       | Single region |
| Growth     | $9,999 | <500      | Multi-AZ      |
| Enterprise | Custom | Unlimited | Multi-region  |

---

## 🔗 Key Links

### Development

- 📦 [GitHub Organization](https://github.com/sapliy)
- 🐛 [Report Issues](https://github.com/sapliy/fintech-ecosystem/issues)
- 🔧 [API Reference](./ARCHITECTURE.md#sdk--client-specifications)
- 📚 [SDKs & Packages](#-supported-sdks-and-packages)

### Community & Support

- 💬 [Discord Community](https://discord.gg/sapliy)
- 📧 [Enterprise Support](mailto:contact@sapliy.io)
- 🔒 [Security Issues](mailto:security@sapliy.io)
- 🌐 [Website](https://sapliy.io)

### Documentation

- 📖 [Full Architecture Docs](./ARCHITECTURE.md)
- 🏢 [Enterprise Deployment Guide](./ENTERPRISE_GUIDE.md)
- 💼 [Business Model & Pricing](./BUSINESS_MODEL.md)
- ⚡ [Quick Reference](./QUICK_REFERENCE.md)

---

## 🎯 Key Principles

1. **One Core Codebase** → Works seamlessly for both SaaS and self-hosted
2. **Hybrid-First Design** → Small users = SaaS, Big users = self-hosted
3. **Open-Source Foundation** → Build trust, get feedback, encourage contributions
4. **Safety-First Automation** → Test/live zones reduce risk, audit trails for compliance
5. **Extensibility Everywhere** → SDKs, UI, policies, connectors — all pluggable

---

## ❓ FAQ

### Is Sapliy open-source?

✅ Yes! The core fintech-ecosystem is MIT-licensed and open-source. SDKs and tools are also open-source.

### Can I self-host Sapliy?

✅ Yes! We provide Docker images, Kubernetes Helm charts, and full deployment guides for AWS, GCP, Azure, and on-prem.

### Do I have to choose between SaaS and self-hosted upfront?

❌ No! Start with SaaS (test mode), graduate to self-hosted as you scale. We have migration tools to help.

### What happens to my data if Sapliy goes down?

🔒 SaaS: Automatic failover across regions, 99.99% uptime SLA
🔒 Self-Hosted: Your infrastructure, your backup strategy

### How much does enterprise self-hosted cost?

💰 Starts at $2K/year (startups). Grows based on deployment size.
Get a custom quote: contact@sapliy.io

### Can I try before committing?

✅ Yes! Free SaaS tier (1K events/month, test mode only). No credit card required.

### What about compliance (HIPAA, PCI-DSS)?

✅ SaaS: SOC 2 certified, working towards HIPAA
✅ Self-Hosted: HIPAA-ready with BAA, PCI-DSS compliant

---

## 📝 License

MIT © [Sapliy](https://github.com/sapliy)

All code, documentation, and examples are open-source and free to use, modify, and distribute.

---

## 🤝 Contributing

We welcome contributions! See [CONTRIBUTING.md](https://github.com/sapliy/fintech-ecosystem/blob/main/CONTRIBUTING.md) for guidelines.

---

**Last Updated**: January 2024
**Version**: 1.0.0
