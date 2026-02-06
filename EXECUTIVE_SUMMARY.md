# Sapliy Executive Summary

> **Event-driven automation & policy platform for fintech and business flows**
>
> Hybrid SaaS + Self-Hosted | Open-Source First | Enterprise-Ready

---

## The Problem

Companies face critical challenges automating financial and business workflows:

1. **Fragmented Stack**
   - Webhook endpoints scattered across integrations
   - No central place to manage automation rules
   - Inconsistent audit trails & compliance tracking

2. **Safety Concerns**
   - No separation between test/production automation
   - Risk of breaking workflows during development
   - Manual approval workflows are error-prone

3. **Compliance Burden**
   - Difficult to prove who did what, when, why
   - Data sovereignty requirements (HIPAA, GDPR, FedRAMP)
   - Audit trails scattered across systems

4. **Scaling Infrastructure**
   - Webhook infrastructure scales linearly with integrations
   - Building multi-step flows requires custom code
   - No policy engine for complex business rules

---

## The Solution: Sapliy

**Sapliy is a platform that lets companies:**

1. ✅ **Emit events from their apps** (via SDK)
2. ✅ **Build no-code automation flows** (react to events)
3. ✅ **Control access with policies** (who can do what)
4. ✅ **Maintain audit trails** (compliance ready)
5. ✅ **Test safely** (test mode / live mode separation)
6. ✅ **Deploy anywhere** (SaaS or self-hosted)

### Example Use Case: Payment Processing

```
Event: "payment.initiated"
  ↓
Flow: "validate_and_process_payment"
  ├─ Action 1: Call payment gateway
  ├─ Action 2: Create ledger entry
  ├─ Action 3: Send confirmation email
  ├─ Action 4: Log to audit trail
  └─ Action 5: Webhook to customer app
  ↓
Result: Fully automated, audited, compliant workflow
```

---

## The Business Model

### 🌐 SaaS (Managed Service)

- **Customers**: Startups, SMBs ($0-$99/month)
- **Revenue**: Subscription + usage-based pricing
- **Gross Margin**: 75%
- **Customer Lifetime Value**: $5K - $50K

### 💼 Self-Hosted (Enterprise)

- **Customers**: Large enterprises ($2K-$500K+/year)
- **Revenue**: Annual license + support
- **Gross Margin**: 85%
- **Customer Lifetime Value**: $50K - $500K+

### 🔄 Cross-Selling

- Customers often start with SaaS (test mode)
- Graduate to self-hosted as compliance/scale requirements increase
- Upgrade path: $99/month (SaaS) → $50K+/year (Enterprise)

---

## Market Opportunity

### Total Addressable Market (TAM)

**SaaS Segment**: $10B+

- Developer tools & automation platforms
- Target: 100K+ startups & SMBs
- Revenue opportunity: $500M+ (if 5% penetration)

**Enterprise Segment**: $5B+

- Enterprise automation, compliance, data residency
- Target: 10K+ large enterprises
- Revenue opportunity: $2B+ (if 5% penetration)

**Total TAM**: $15B+

### Competitive Landscape

| Competitor   | SaaS | Self-Hosted | Fintech-Focus | Open-Source |
| ------------ | ---- | ----------- | ------------- | ----------- |
| **Zapier**   | ✅   | ❌          | ❌            | ❌          |
| **Make.com** | ✅   | ✅          | ❌            | ❌          |
| **n8n**      | ✅   | ✅          | ❌            | ✅          |
| **Stripe**   | ✅   | Native      | ✅            | ❌          |
| **Sapliy**   | ✅   | ✅          | ✅            | ✅          |

**Sapliy's advantage**: Only platform with hybrid SaaS+self-hosted, fintech-optimized, AND open-source.

---

## Go-to-Market Strategy

### Phase 1: Product Launch (Months 1-3)

- ✅ Launch SaaS platform (api.sapliy.io)
- ✅ Target early adopters (Y Combinator, Techstars alumni)
- ✅ Build Node.js SDK
- ✅ Achieve product-market fit

**Metrics**: 500 free users, 20 paying customers, $1K MRR

### Phase 2: Community & SDKs (Months 4-6)

- 📦 Open-source fintech-ecosystem
- 🐍 Python & Go SDKs
- 📚 Complete documentation
- 🤝 Build community (Discord, GitHub)

**Metrics**: 5K free users, 150 paying customers, $10K MRR

### Phase 3: Enterprise Launch (Months 7-9)

- 🏢 Self-hosted deployment option
- 📋 Compliance & security certifications
- 💼 Enterprise sales & support
- 🚀 Targeted enterprise outreach

**Metrics**: 3 enterprise customers, $40K ARR

### Phase 4: Scale & Monetization (Months 10-12)

- 💰 Premium features & integrations
- 📊 Advanced analytics & compliance
- 🎯 $50K+ ARR from enterprise
- 🌍 Multi-region deployment

**Metrics**: $500K+ combined SaaS + Enterprise revenue

---

## Financial Projections

### Year 1

| Metric                  | Q1  | Q2  | Q3   | Q4    |
| ----------------------- | --- | --- | ---- | ----- |
| Free Users              | 500 | 2K  | 5K   | 10K   |
| Paying Customers (SaaS) | 20  | 75  | 250  | 500   |
| MRR (SaaS)              | $1K | $4K | $16K | $38K  |
| Enterprise Customers    | 0   | 0   | 1    | 3     |
| ARR (Enterprise)        | $0  | $0  | $10K | $100K |
| Total Revenue           | $1K | $4K | $26K | $138K |

**Year 1 Total Revenue**: ~$170K

### Year 2-3

| Year   | SaaS ARR | Enterprise ARR | Total ARR | Growth |
| ------ | -------- | -------------- | --------- | ------ |
| Year 1 | $450K    | $100K          | $550K     | —      |
| Year 2 | $1.5M    | $500K          | $2M       | 3.6x   |
| Year 3 | $4M      | $1.5M          | $5.5M     | 2.75x  |

---

## Key Success Metrics

### SaaS Metrics

- **Free-to-paid conversion**: 10-15% (target)
- **CAC (Customer Acquisition Cost)**: $150
- **CAC Payback Period**: 2-3 months
- **LTV (Lifetime Value)**: $5K - $50K
- **Churn Rate**: <5% monthly
- **NRR (Net Revenue Retention)**: >120%

### Enterprise Metrics

- **Sales Cycle**: 120-180 days
- **CAC**: $10K - $20K
- **CAC Payback Period**: 12-18 months
- **LTV**: $50K - $500K+
- **Churn Rate**: <2% annually
- **NRR**: >130%

### Product Metrics

- **Event ingestion**: 10K+ events/sec
- **Flow execution latency**: <100ms p99
- **Webhook delivery success**: 99.99%
- **API uptime**: 99.95% (self-hosted), 99.99% (SaaS)

---

## Funding Needs

### Series A: $3-5M

- **Use of Funds**:
  - Engineering (50%): Expand team, advanced features
  - Sales & Marketing (30%): Enterprise sales, growth marketing
  - Operations (10%): Compliance, security, infrastructure
  - General (10%): Legal, HR, overhead
- **Milestones**:
  - $2M ARR (end of Year 2)
  - 30+ Enterprise customers
  - 5K+ SaaS customers
  - SOC 2 + HIPAA certification

---

## Team & Hiring

### Current Team

- **CEO**: Product vision, strategy, fundraising
- **CTO**: Architecture, engineering leadership
- **VP Product**: Product roadmap, customer discovery

### Hiring Plan (Next 12 Months)

- 2x Full-stack engineers (SaaS/API)
- 1x Infrastructure/DevOps engineer
- 1x Security engineer
- 1x Sales development rep (SDR)
- 1x Customer success manager (CSM)
- 1x Marketing manager

---

## Risk Mitigation

### Technical Risks

| Risk                          | Mitigation                                                           |
| ----------------------------- | -------------------------------------------------------------------- |
| Scalability (10K+ events/sec) | Built on proven stack (Kafka, PostgreSQL); load testing in progress  |
| Security/Compliance           | SOC 2 roadmap; HIPAA-ready architecture; security audits quarterly   |
| Vendor lock-in                | Open-source foundation; easy migration paths; configurable endpoints |

### Business Risks

| Risk            | Mitigation                                                           |
| --------------- | -------------------------------------------------------------------- |
| Competition     | Differentiation: hybrid, fintech-focused, open-source                |
| Churn           | Product quality; free tier to reduce switching costs; strong support |
| Market adoption | Community-first; open-source advantage; network effects              |

---

## Conclusion

Sapliy addresses a $15B+ market opportunity with a unique hybrid business model:

1. **🌐 SaaS** for startups & SMBs (high volume, recurring revenue)
2. **💼 Self-Hosted** for enterprises (high value, sticky, long contracts)
3. **💻 Open-Source** foundation (community trust, differentiation)

**By Year 3**, we project:

- **$5.5M ARR** in combined revenue
- **5K+ SaaS customers**
- **20+ Enterprise customers**
- **Industry leadership** in event-driven fintech automation

**Next Steps**:

1. Complete Series A fundraising ($3-5M)
2. Launch self-hosted enterprise option (Q2 2024)
3. Achieve $100K MRR in SaaS (Q3 2024)
4. Close first $500K+ enterprise deal (Q4 2024)

---

## Appendix: Key Resources

### Documentation

- [Full Architecture](./ARCHITECTURE.md) — Technical design & system overview
- [Enterprise Guide](./ENTERPRISE_GUIDE.md) — Self-hosted deployment & operations
- [Business Model](./BUSINESS_MODEL.md) — Detailed pricing, segments, and GTM
- [Quick Reference](./QUICK_REFERENCE.md) — Developer quick-start guide

### Links

- **Website**: https://sapliy.io
- **GitHub**: https://github.com/sapliy
- **Discord**: https://discord.gg/sapliy
- **Contact**: contact@sapliy.io

---

**Prepared**: January 2024  
**Document Version**: 1.0  
**Confidentiality**: Company Confidential
