# วิธี Deploy บนเครื่องตัวเอง

## เมื่อ push ไป develop branch
GitHub Actions จะ:
1. Build Docker image
2. Push ไป GitHub Container Registry (ghcr.io)
3. Tag: `develop` และ `develop-<commit-sha>`

## ดึง image และรันบนเครื่องตัวเอง

### 1. Login to GitHub Container Registry
```bash
# สร้าง Personal Access Token (PAT) ที่ https://github.com/settings/tokens
# เลือก scopes: read:packages

echo YOUR_GITHUB_TOKEN | docker login ghcr.io -u YOUR_GITHUB_USERNAME --password-stdin
```

### 2. Pull และ Run container
```bash
# Pull image ล่าสุด
docker pull ghcr.io/YOUR_GITHUB_USERNAME/learning-ci-cd:develop

# Run container
docker run -d \
  --name myapp \
  -p 8080:8080 \
  --restart unless-stopped \
  ghcr.io/YOUR_GITHUB_USERNAME/learning-ci-cd:develop
```

### 3. ทดสอบ
```bash
curl http://localhost:8080/health
curl http://localhost:8080/hello
```

### 4. Update เมื่อมี version ใหม่
```bash
# Stop และลบ container เก่า
docker stop myapp
docker rm myapp

# Pull version ใหม่และรัน
docker pull ghcr.io/YOUR_GITHUB_USERNAME/learning-ci-cd:develop
docker run -d --name myapp -p 8080:8080 --restart unless-stopped ghcr.io/YOUR_GITHUB_USERNAME/learning-ci-cd:develop
```

## Auto-update ด้วย Watchtower (Optional)

ถ้าอยากให้ auto-update เมื่อมี image ใหม่:

```bash
docker run -d \
  --name watchtower \
  -v /var/run/docker.sock:/var/run/docker.sock \
  containrrr/watchtower \
  --interval 300 \
  myapp
```

Watchtower จะตรวจสอบทุก 5 นาที และ update container อัตโนมัติ

## Docker Compose (แนะนำ)

สร้างไฟล์ `docker-compose.yml`:

```yaml
version: '3.8'

services:
  app:
    image: ghcr.io/YOUR_GITHUB_USERNAME/learning-ci-cd:develop
    container_name: myapp
    ports:
      - "8080:8080"
    restart: unless-stopped

  # Optional: Auto-update
  watchtower:
    image: containrrr/watchtower
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    command: --interval 300 myapp
```

รัน:
```bash
docker-compose up -d
```

Update:
```bash
docker-compose pull
docker-compose up -d
```
