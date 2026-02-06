# Sapliy Enterprise Deployment & Security Guide

> Complete guide for enterprise customers deploying Sapliy self-hosted with maximum security, compliance, and reliability.

---

## Table of Contents

1. [Enterprise Overview](#enterprise-overview)
2. [Pre-Deployment Planning](#pre-deployment-planning)
3. [Security Architecture](#security-architecture)
4. [Deployment Strategies](#deployment-strategies)
5. [Operational Runbooks](#operational-runbooks)
6. [Compliance & Audit](#compliance--audit)
7. [Disaster Recovery](#disaster-recovery)
8. [Support & Escalation](#support--escalation)

---

## Enterprise Overview

### Who Should Use Enterprise Self-Hosted?

Sapliy Self-Hosted is ideal for organizations that:

- **Need data sovereignty**: Sensitive data must remain in-country or within your infrastructure
- **Face regulatory requirements**: HIPAA, PCI-DSS, FedRAMP, SOX compliance
- **Have scale requirements**: Processing 100K+ events/second
- **Require customization**: Custom integrations, policies, or deployment models
- **Value operational control**: Internal team manages infrastructure and upgrades
- **Have existing infrastructure**: AWS, GCP, Azure, or on-premise data centers

### Enterprise Features

#### Security & Compliance

✅ Multi-layer encryption (AES-256 at rest, TLS 1.3 in transit)
✅ HIPAA-compliant audit logging & BAA support
✅ PCI-DSS ready (no card data, tokenization focus)
✅ FedRAMP certified infrastructure support
✅ SOX-compliant audit trails
✅ GDPR/CCPA data handling & exportability
✅ Customizable data retention policies
✅ Field-level encryption for sensitive data

#### Operational

✅ Multi-region failover & disaster recovery
✅ 99.95% uptime SLA guarantee
✅ <15 minute RTO, <1 minute RPO
✅ Zero-downtime deployments & schema migrations
✅ Advanced monitoring & observability (OpenTelemetry)
✅ On-call support with 1-hour response SLA
✅ Quarterly penetration testing

#### Technical

✅ Kubernetes-native deployment (Helm charts)
✅ Configurable API endpoints (use your domain)
✅ Private VPC deployment options
✅ Network isolation & security groups
✅ Custom certificate management
✅ Secrets management integration (Vault, Secrets Manager)
✅ Database replication & failover automation

---

## Pre-Deployment Planning

### 1. Infrastructure Assessment

#### Cloud Infrastructure (AWS Example)

```yaml
Required AWS Resources:
├── VPC
│   ├── Private Subnets (3x AZs for databases)
│   ├── Private Subnets (3x AZs for application)
│   ├── Public Subnets (3x AZs for load balancers)
│   ├── NAT Gateways (1 per AZ for egress)
│   └── Internet Gateway (for load balancer access)
│
├── EKS Cluster
│   ├── 3+ nodes across AZs (minimum)
│   ├── Auto-scaling group (6-20 nodes)
│   ├── Instance type: m5.2xlarge or larger
│   └── Storage: EBS gp3 volumes for data
│
├── RDS PostgreSQL
│   ├── Multi-AZ deployment
│   ├── db.r5.2xlarge or larger (for 100K+ events/sec)
│   ├── Read replicas in each region
│   ├── Automated backups (daily, 30-day retention)
│   └── Encryption: AWS KMS managed keys
│
├── ElastiCache Redis
│   ├── Multi-AZ cluster mode enabled
│   ├── cache.r6g.xlarge or larger
│   ├── Automatic failover enabled
│   └── Encryption at rest & in-transit
│
├── MSK Kafka (or self-managed)
│   ├── 3 brokers minimum (1 per AZ)
│   ├── kafka.m5.2xlarge or larger
│   ├── Encryption: KMS managed keys
│   └── 7-day retention policy
│
├── Secrets Manager
│   ├── Database credentials
│   ├── API keys
│   ├── TLS certificates
│   └── Encryption: AWS KMS
│
└── S3 Buckets
    ├── Event archive (cold storage)
    ├── Backup storage (encrypted)
    ├── Log storage (immutable, versioned)
    └── Terraform state (encrypted, versioned)
```

#### On-Premise Infrastructure (Example)

```yaml
Required Hardware/VM Resources:
├── Load Balancers (3x)
│   └── Network load balancers for 100K+ events/sec
│
├── Kubernetes Cluster (6+ nodes minimum)
│   ├── Control plane: 3 nodes (4 CPU, 16 GB RAM each)
│   ├── Worker nodes: 6+ nodes (8+ CPU, 32+ GB RAM each)
│   └── Storage: 500 GB+ per node
│
├── PostgreSQL (HA setup)
│   ├── Primary + 2 replicas (minimum)
│   ├── 16+ CPU cores, 64+ GB RAM
│   ├── High-speed storage (NVMe SSDs)
│   └── Backup server
│
├── Redis Cluster (3 nodes minimum)
│   ├── 8+ CPU, 32+ GB RAM per node
│   └── Replication with automatic failover
│
├── Kafka Cluster (3+ brokers)
│   ├── 8+ CPU, 32+ GB RAM per broker
│   ├── High-speed storage for logs
│   └── Quorum-based setup
│
└── Backup Storage
    └── Dedicated NAS/SAN with encryption
```

### 2. Capacity Planning

#### Event Ingestion Scaling

```
Baseline Calculation:
- Target: 100K events/second
- Backend instances needed: 100K / 10K per instance = 10 instances
- With 2x redundancy for HA: 20 instances minimum
- Kubernetes nodes: 20 pods / 3 pods per node = 7 nodes
- Add 50% headroom: 10-12 nodes total

Database Scaling:
- 100K events/sec = 8.64B events/day
- Storage: 8.64B × 1KB per event = 8.64 TB/day
- Retention: 90 days = 777.6 TB storage
- Database size: 1.5 TB (with indexing overhead)
- Use RDS read replicas for analytics queries
- Archive to S3 after 30 days
```

#### Memory & CPU Requirements

```
Per Application Instance (for 10K events/sec):
- CPU: 4 cores (2 reserved, 2 burstable)
- Memory: 8 GB (6 GB reserved, 2 GB for JVM)
- Network: 1 Gbps throughput

Per Database Node (handling 100K events/sec):
- CPU: 16 cores (reserved)
- Memory: 64 GB (reserved)
- Storage: IOPS: 10,000+ random IOPS
- Throughput: 500 MB/s sustained
```

### 3. Network Design

#### High-Level Network Topology

```
Internet
   ↓
CDN/DDoS Protection (Cloudflare, AWS Shield)
   ↓
Load Balancer (Layer 7)
   ↓
┌─────────────────────────────────────────┐
│          Public Subnet (AZ1-3)          │
│  Load Balancer / NAT Gateway            │
└─────────────────────────────────────────┘
   ↓ (private endpoint)
┌─────────────────────────────────────────┐
│    Private Subnets (Application Tier)   │
│  Kubernetes Cluster (API Gateway, Auth) │
└─────────────────────────────────────────┘
   ↓ (private endpoint)
┌─────────────────────────────────────────┐
│    Private Subnets (Data Tier)          │
│  PostgreSQL, Redis, Kafka (encrypted)   │
└─────────────────────────────────────────┘
```

#### Security Group Rules

```
Load Balancer SG:
- Inbound: 443 (HTTPS) from 0.0.0.0/0
- Inbound: 80 (HTTP) from 0.0.0.0/0 (redirect to 443)
- Outbound: 8080 to App SG

Application SG:
- Inbound: 8080 from Load Balancer SG
- Inbound: 8081 from Load Balancer SG (Auth)
- Inbound: 8082 from Load Balancer SG (Payments)
- Inbound: 8083 from Load Balancer SG (Ledger)
- Outbound: 5432 to Database SG
- Outbound: 6379 to Redis SG
- Outbound: 9092 to Kafka SG

Database SG:
- Inbound: 5432 from App SG only
- Inbound: 5432 from Bastion SG (admin only)
- Outbound: none (stateful)

Redis SG:
- Inbound: 6379 from App SG only
- Outbound: none (stateful)

Kafka SG:
- Inbound: 9092 from App SG only
- Outbound: none (stateful)

Bastion SG (optional):
- Inbound: 22 from Corporate IP/VPN only
- Outbound: 22 to App SG
- Outbound: 5432 to Database SG
```

---

## Security Architecture

### 1. Encryption Strategy

#### Data at Rest

```yaml
PostgreSQL:
  - Table encryption: LUKS with AES-256
  - Backup encryption: AWS KMS or local key management
  - Key rotation: Every 90 days
  - Key storage: HSM (Hardware Security Module) recommended

Redis:
  - Encryption at rest: Enabled (platform-specific)
  - Backup encryption: AES-256 via AWS S3 encryption
  - Key rotation: Every 90 days

Kafka:
  - Message encryption: TLS 1.3 between brokers
  - At-rest: LUKS encryption of broker storage
  - Backup: Encrypted S3 archival

Application Secrets:
  - HashiCorp Vault or AWS Secrets Manager
  - Key storage: HSM or managed service
  - Rotation: Automatic (30-90 days)
```

#### Data in Transit

```yaml
Client to Load Balancer:
  - TLS 1.3 mandatory
  - HSTS header enforced (strict-transport-security)
  - Certificate: CA-signed (not self-signed)

Load Balancer to Application:
  - TLS 1.3 or mTLS for internal services
  - Separate certificates per service

Application to Database:
  - SSL/TLS with certificate verification
  - Connection pooling with SSL mode: require
  - Credentials: Rotated secrets

Application to Redis/Kafka:
  - TLS 1.3 with certificate verification
  - SASL/SCRAM authentication
```

### 2. Authentication & Authorization

#### Service-to-Service Authentication

```
├── mTLS (mutual TLS)
│   ├── Each service has unique certificate
│   ├── Certificate rotation: every 90 days
│   └── CRL (Certificate Revocation List) enabled
│
├── Service Accounts (Kubernetes)
│   ├── Each pod has unique service account
│   ├── RBAC policies for each account
│   └── Audit logging of all access
│
└── API Key Scoping
    ├── Zone-scoped keys (production vs. test)
    ├── Permission-scoped keys (events:emit, flows:read)
    ├── Expiration: 90 days default
    └── Rotation tracking: audit log of all rotations
```

#### User Authentication

```
├── OAuth 2.0 / OIDC (recommended)
│   ├── Integration with Okta, Azure AD, Auth0
│   ├── MFA enforcement: TOTP or FIDO2
│   ├── Session timeout: 30 minutes inactivity
│   └── Token refresh: 60-minute expiration
│
├── SAML 2.0 (Enterprise)
│   ├── Integration with corporate IdP
│   ├── Assertion encryption: AES-256
│   ├── Signed assertions only
│   └── Artifact resolution enabled
│
└── Basic Auth (Deprecated)
    └── Not recommended; use OAuth 2.0 instead
```

### 3. Audit & Logging

#### Immutable Audit Log Design

```sql
CREATE TABLE audit_logs_immutable (
  id UUID PRIMARY KEY,
  sequence_number BIGSERIAL UNIQUE,  -- For tamper detection
  timestamp TIMESTAMP NOT NULL,
  actor_id UUID NOT NULL,
  action VARCHAR NOT NULL,
  resource_type VARCHAR,
  resource_id UUID,
  changes JSONB,
  ip_address INET,
  user_agent VARCHAR,
  result VARCHAR (success|failure),
  error_message TEXT,
  signature VARCHAR,  -- HMAC-SHA256 for tamper detection
  signature_timestamp TIMESTAMP,
  CHECK (signature IS NOT NULL)  -- Enforce signature
) PARTITION BY RANGE (timestamp);

-- Weekly partitions for immutability
CREATE TABLE audit_logs_immutable_2024_w01 PARTITION OF audit_logs_immutable
  FOR VALUES FROM ('2024-01-01') TO ('2024-01-08');

-- Deny all deletes & updates on immutable logs
CREATE POLICY no_delete ON audit_logs_immutable
  USING (false);

CREATE POLICY no_update ON audit_logs_immutable
  USING (false);
```

#### Audit Events to Log

```
Authentication:
- Login success/failure
- API key creation/rotation/deletion
- MFA enabled/disabled
- Session timeout/logout

Authorization:
- Permission changes
- Role assignments/removals
- Access denied events
- Policy changes

Data Access:
- Event retrieval
- Flow access
- Zone access
- Configuration reads

Configuration Changes:
- Flow creation/update/deletion
- Zone configuration changes
- Integration updates
- Policy changes

Security Events:
- Failed authentication (3+ failures = alert)
- Privilege escalation attempts
- Certificate expiration warnings
- Secret rotation events
- Rate limit violations
```

---

## Deployment Strategies

### 1. Kubernetes Deployment (AWS EKS)

#### Pre-requisites

```bash
# AWS CLI configured with appropriate credentials
aws sts get-caller-identity

# kubectl & helm installed
kubectl version --client
helm version

# Create VPC & EKS Cluster
aws eks create-cluster \
  --name sapliy-prod \
  --version 1.27 \
  --region us-east-1 \
  --role-arn arn:aws:iam::ACCOUNT:role/eks-service-role \
  --resources-vpc-config subnetIds=subnet-1,subnet-2,subnet-3
```

#### Helm Chart Installation

```bash
# Add Sapliy Helm repository
helm repo add sapliy https://charts.sapliy.io
helm repo update

# Create namespace
kubectl create namespace sapliy-prod

# Create secrets for sensitive data
kubectl create secret generic sapliy-secrets \
  --from-literal=db_password=$(openssl rand -base64 32) \
  --from-literal=redis_password=$(openssl rand -base64 32) \
  --namespace sapliy-prod

# Deploy using Helm
helm install sapliy-release sapliy/sapliy \
  --namespace sapliy-prod \
  --values values-prod.yaml \
  --timeout 10m

# Verify deployment
kubectl get pods -n sapliy-prod
kubectl get svc -n sapliy-prod
```

#### values-prod.yaml

```yaml
# Global settings
global:
  environment: production
  tls_enabled: true
  log_level: info

# API Gateway settings
api_gateway:
  replicas: 10
  resources:
    requests:
      cpu: 4
      memory: 8Gi
    limits:
      cpu: 8
      memory: 16Gi
  pdb:
    minAvailable: 2

# Auth Service settings
auth:
  replicas: 3
  resources:
    requests:
      cpu: 2
      memory: 4Gi
    limits:
      cpu: 4
      memory: 8Gi

# Flow Engine settings
flow_engine:
  replicas: 10
  resources:
    requests:
      cpu: 4
      memory: 8Gi
    limits:
      cpu: 8
      memory: 16Gi

# Database configuration
postgresql:
  enabled: false # Use RDS instead
  externalHost: prod-db.c2t3yl1x4x5c.us-east-1.rds.amazonaws.com
  username: sapliy_admin
  password: ${DB_PASSWORD} # From secret
  database: sapliy_prod

# Redis configuration
redis:
  enabled: false # Use ElastiCache instead
  externalHost: prod-redis.xyz123.ng.0001.use1.cache.amazonaws.com
  port: 6379
  password: ${REDIS_PASSWORD} # From secret

# Kafka configuration
kafka:
  enabled: false # Use MSK instead
  brokers:
    - kafka1.prod.internal:9092
    - kafka2.prod.internal:9092
    - kafka3.prod.internal:9092

# Ingress configuration
ingress:
  enabled: true
  className: nginx
  hosts:
    - host: api.fintech.yourcompany.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: sapliy-tls
      hosts:
        - api.fintech.yourcompany.com

# Monitoring
monitoring:
  enabled: true
  prometheus:
    enabled: true
  alertmanager:
    enabled: true
    slackWebhook: https://hooks.slack.com/...

# Security settings
security:
  podSecurityPolicy: restricted
  networkPolicy:
    enabled: true
  rbac:
    enabled: true
  serviceAccountName: sapliy-sa
```

### 2. Database Initialization

#### PostgreSQL Setup (AWS RDS)

```bash
# Create RDS instance (high-availability setup)
aws rds create-db-instance \
  --db-instance-identifier sapliy-prod \
  --db-instance-class db.r5.2xlarge \
  --engine postgres \
  --engine-version 15.3 \
  --allocated-storage 500 \
  --storage-type gp3 \
  --storage-encrypted \
  --kms-key-id arn:aws:kms:us-east-1:ACCOUNT:key/KEY_ID \
  --multi-az \
  --backup-retention-period 30 \
  --enable-log-exports postgresql \
  --enable-performance-insights \
  --performance-insights-retention-period 7 \
  --enable-iam-database-authentication

# Create read replicas
aws rds create-db-instance-read-replica \
  --db-instance-identifier sapliy-prod-replica-1 \
  --source-db-instance-identifier sapliy-prod \
  --db-instance-class db.r5.2xlarge

# Initialize schema
psql -h sapliy-prod.c2t3yl1x4x5c.us-east-1.rds.amazonaws.com \
  -U sapliy_admin \
  -d sapliy_prod \
  -f /path/to/schema.sql

# Create indexes for performance
psql -h sapliy-prod.c2t3yl1x4x5c.us-east-1.rds.amazonaws.com \
  -U sapliy_admin \
  -d sapliy_prod \
  -f /path/to/indexes.sql
```

### 3. Secrets Management (AWS Secrets Manager)

```bash
# Create database credentials secret
aws secretsmanager create-secret \
  --name sapliy/prod/db \
  --description "Sapliy Production Database Credentials" \
  --secret-string '{
    "username": "sapliy_admin",
    "password": "'$(openssl rand -base64 32)'",
    "engine": "postgres",
    "host": "sapliy-prod.c2t3yl1x4x5c.us-east-1.rds.amazonaws.com",
    "port": 5432,
    "dbname": "sapliy_prod"
  }'

# Create API keys secret
aws secretsmanager create-secret \
  --name sapliy/prod/api_keys \
  --secret-string '{
    "sk_live_XXXXX": "secret_key_value",
    "pk_live_XXXXX": "publishable_key_value"
  }'

# Create TLS certificate secret (for mTLS)
aws secretsmanager create-secret \
  --name sapliy/prod/tls_cert \
  --secret-string file:///path/to/cert.pem
```

---

## Operational Runbooks

### 1. Incident Response

#### Database Connection Issues

```
1. Check database status
   aws rds describe-db-instances --db-instance-identifier sapliy-prod

2. Check application logs for connection errors
   kubectl logs -n sapliy-prod -l app=api-gateway --tail=100

3. Verify security group rules allow traffic
   aws ec2 describe-security-groups --group-ids sg-xxxxx

4. Check connection pool exhaustion
   SELECT count(*) FROM pg_stat_activity;

5. If primary is down, promote replica
   aws rds promote-read-replica --db-instance-identifier sapliy-prod-replica-1

6. Scale down traffic while recovering
   kubectl scale deployment api-gateway -n sapliy-prod --replicas=1
```

#### High Event Latency

```
1. Check Kafka lag
   kafka-consumer-groups.sh --bootstrap-server localhost:9092 \
     --group sapliy-flow-executor --describe

2. Check database query performance
   EXPLAIN ANALYZE SELECT * FROM events WHERE zone_id = $1 LIMIT 100;

3. Check Redis memory usage
   redis-cli INFO memory

4. Scale event processing
   kubectl scale deployment flow-engine -n sapliy-prod --replicas=15

5. Monitor using Prometheus
   - rate(events_processed_total[5m])
   - histogram_quantile(0.99, flow_execution_duration_seconds)
```

#### Webhook Delivery Failures

```
1. Check webhook service logs
   kubectl logs -n sapliy-prod -l app=webhook-dispatcher --tail=200

2. Verify retry queue size
   SELECT COUNT(*) FROM webhook_queue WHERE status = 'failed';

3. Check target endpoint availability
   curl -v https://target.example.com/webhook --max-time 5

4. Manual retry of failed webhooks
   INSERT INTO webhook_queue (event_id, status, retry_count)
   SELECT id, 'pending', 0 FROM webhooks WHERE status = 'failed'
   AND created_at > NOW() - INTERVAL '1 hour';
```

### 2. Backup & Recovery

#### Daily Backup Verification

```bash
#!/bin/bash
# Daily backup check script

# Check RDS automated backups
aws rds describe-db-instances \
  --db-instance-identifier sapliy-prod \
  --query 'DBInstances[0].LatestRestorableTime'

# Verify S3 backup bucket
aws s3 ls s3://sapliy-prod-backups/daily/ --recursive

# Restore to point-in-time (dry run)
echo "Latest restorable time: $(date)"

# Monthly restore test (non-production)
aws rds restore-db-instance-to-point-in-time \
  --source-db-instance-identifier sapliy-prod \
  --target-db-instance-identifier sapliy-prod-restore-test \
  --restore-time $(date -u +"%Y-%m-%dT%H:%M:%S.000Z") \
  --no-publicly-accessible
```

#### Point-in-Time Recovery Procedure

```bash
#!/bin/bash
# Recovery script (use with caution!)

# 1. Identify recovery point
RECOVERY_TIME="2024-01-15T14:30:00Z"

# 2. Restore to new instance
aws rds restore-db-instance-to-point-in-time \
  --source-db-instance-identifier sapliy-prod \
  --target-db-instance-identifier sapliy-prod-recovered \
  --restore-time $RECOVERY_TIME

# 3. Wait for recovery to complete
aws rds wait db-instance-available \
  --db-instance-identifier sapliy-prod-recovered

# 4. Test recovery instance
psql -h sapliy-prod-recovered.c2t3yl1x4x5c.us-east-1.rds.amazonaws.com \
  -U sapliy_admin \
  -d sapliy_prod \
  -c "SELECT COUNT(*) FROM events;"

# 5. If successful, swap DNS
aws route53 change-resource-record-sets \
  --hosted-zone-id ZXXXXX \
  --change-batch file:///path/to/route53_update.json
```

---

## Compliance & Audit

### 1. HIPAA Compliance Checklist

```
✅ Technical Safeguards:
  - [ ] Encryption at rest (AES-256)
  - [ ] Encryption in transit (TLS 1.3)
  - [ ] Access controls (RBAC, MFA)
  - [ ] Audit logging (immutable)
  - [ ] Data integrity checks (HMAC)
  - [ ] Transmission security

✅ Administrative Safeguards:
  - [ ] Workforce security (access controls)
  - [ ] Information access management (least-privilege)
  - [ ] Security training (annual)
  - [ ] Security incident procedures
  - [ ] Sanction policy (data breach response)
  - [ ] Workforce security (role-based access)

✅ Physical Safeguards:
  - [ ] Facility access controls (VPN, MFA)
  - [ ] Workstation security (encryption, MFA)
  - [ ] Media and device controls (encryption, disposal)
  - [ ] Facility security (firewalls, WAF)

✅ Documentation:
  - [ ] System security plan
  - [ ] Risk assessment (annual)
  - [ ] Security policies & procedures
  - [ ] Business associate agreements (BAA)
  - [ ] Breach notification procedures
  - [ ] Incident response plan
```

### 2. PCI-DSS Compliance Checklist

```
✅ Network Security:
  - [ ] Firewall configuration (deny by default)
  - [ ] No default passwords
  - [ ] Network segmentation (card data isolated)
  - [ ] IDS/IPS enabled
  - [ ] Regular network scans

✅ Data Protection:
  - [ ] No storage of full PAN (Primary Account Number)
  - [ ] Encryption of transmissions (TLS 1.3)
  - [ ] Tokenization instead of card storage
  - [ ] Secure key management
  - [ ] Limited data retention

✅ Access Control:
  - [ ] User identification (unique IDs)
  - [ ] Restricted access by role
  - [ ] Password requirements (min 8 chars, complexity)
  - [ ] Access logging (who, what, when)
  - [ ] Physical access controls

✅ Regular Testing:
  - [ ] Quarterly vulnerability scans
  - [ ] Annual penetration testing
  - [ ] Intrusion detection monitoring
  - [ ] Log reviews (monthly minimum)
```

### 3. Audit Log Export

```bash
#!/bin/bash
# Export audit logs for compliance review

# Export last 30 days of audit logs
psql -h sapliy-prod.c2t3yl1x4x5c.us-east-1.rds.amazonaws.com \
  -U sapliy_admin \
  -d sapliy_prod \
  -c "\COPY (
    SELECT
      id, actor_id, action, resource_type, resource_id,
      timestamp, ip_address, result, signature
    FROM audit_logs_immutable
    WHERE timestamp > NOW() - INTERVAL '30 days'
    ORDER BY timestamp DESC
  ) TO STDOUT WITH CSV HEADER;" | \
  gpg --encrypt --recipient compliance@yourcompany.com > \
  audit_logs_2024_01.csv.gpg

# Upload to secure storage
aws s3 cp audit_logs_2024_01.csv.gpg \
  s3://compliance-logs-bucket/sapliy/ \
  --sse aws:kms \
  --sse-kms-key-id arn:aws:kms:us-east-1:ACCOUNT:key/KEY_ID
```

---

## Disaster Recovery

### 1. Failover Procedures

#### Automatic Failover (RDS)

```
RDS Multi-AZ handles automatic failover:
1. Primary database detects failure
2. Automatic failover triggers (typically < 2 minutes)
3. Read replica is promoted to primary
4. Application connections retry and succeed
5. Manual monitoring to ensure stability

Configuration:
- Enable Multi-AZ: true
- Backup retention: 30 days
- Enhanced monitoring: enabled
- Performance insights: enabled
```

#### Manual Failover (Application)

```bash
#!/bin/bash
# Trigger manual failover to another region

# 1. Stop writes to primary region
kubectl scale deployment api-gateway -n sapliy-prod --replicas=0

# 2. Verify replica lag is minimal
psql -h sapliy-prod-replica-dr.c2t3yl1x4x5c.us-west-2.rds.amazonaws.com \
  -U sapliy_admin \
  -d sapliy_prod \
  -c "SELECT MAX(write_lag) FROM pg_stat_replication;"

# 3. Promote DR replica to primary
aws rds promote-read-replica \
  --db-instance-identifier sapliy-prod-replica-dr

# 4. Update DNS to point to DR region
aws route53 change-resource-record-sets \
  --hosted-zone-id ZXXXXX \
  --change-batch file:///path/to/route53_dr.json

# 5. Deploy application to DR region
kubectl apply -f sapliy-dr-deployment.yaml

# 6. Verify DR functionality
curl https://api.fintech.yourcompany.com/health
```

### 2. Recovery Time Objectives

| Component           | RTO     | RPO    | Notes                       |
| ------------------- | ------- | ------ | --------------------------- |
| **Database (RDS)**  | <2 min  | <1 min | Multi-AZ automatic failover |
| **Application**     | <5 min  | 0 min  | Stateless, quick scale-up   |
| **Redis Cache**     | <1 min  | 0 min  | Recreated on startup        |
| **Kafka Events**    | <5 min  | <1 min | Broker replication          |
| **Complete System** | <15 min | <1 min | All components failover     |

---

## Support & Escalation

### Sapliy Enterprise Support Tiers

#### Premium Support (24/7/365)

- **Response time**: 1 hour (all issues)
- **Resolution time**: 4 hours (critical), 8 hours (high)
- **Escalation**: Direct to engineering team
- **Included**: Quarterly business reviews, architecture guidance
- **Cost**: Contact sales

#### Standard Support (Business hours)

- **Response time**: 4 hours
- **Resolution time**: 24 hours (critical), 48 hours (high)
- **Escalation**: Support team only
- **Included**: Monthly check-ins
- **Cost**: Included with license

### Critical Issue Escalation Path

```
1. Customer reports issue via support portal
   ↓ (< 15 minutes)
2. Support engineer acknowledges & begins diagnosis
   ↓ (< 30 minutes)
3. If critical:
   - Page on-call engineer
   - Create incident room (Slack)
   - Begin war room bridge call
   ↓ (< 1 hour)
4. Engineering team joins & works towards resolution
   ↓ (< 4 hours for fix or workaround)
5. Post-incident review & RCA
   ↓ (< 24 hours)
6. Communication of findings & preventive measures
```

### SLA Guarantees

```
Availability SLA: 99.95% uptime (43.2 minutes downtime/month)
- Measured across all availability zones
- Excludes scheduled maintenance windows
- Excludes customer-caused outages

Latency SLA: p99 < 100ms for API endpoints
- Measured from load balancer to application
- Excludes external service latency (webhooks)

Support Response SLA: 1 hour for critical issues
- Measured from ticket creation
- Business hours: 24/7 (Premium support)
- 8x5 (Standard support)
```

---

## Conclusion

This guide provides comprehensive guidance for deploying and operating Sapliy in enterprise environments. For additional support or customization needs, contact your Sapliy account manager.

**Key Resources:**

- Enterprise Documentation: https://docs.sapliy.io/enterprise
- API Reference: https://docs.sapliy.io/api
- Support Portal: https://support.sapliy.io
- Security Updates: https://security.sapliy.io
