# CI/CD พื้นฐาน - เรียนรู้และใช้งานได้จริง

## CI/CD คืออะไร?

**CI (Continuous Integration)** = การรวมโค้ดอย่างต่อเนื่อง
- เมื่อคุณเขียนโค้ดเสร็จ push ขึ้น Git
- ระบบจะ**ทำการตรวจสอบอัตโนมัติ**ว่าโค้ดใช้งานได้หรือไม่
- รัน test, ตรวจสอบ syntax, build โปรเจค

**CD (Continuous Deployment/Delivery)** = การส่งมอบ/Deploy อย่างต่อเนื่อง
- เมื่อโค้ดผ่านการตรวจสอบแล้ว
- ระบบจะ**ส่งขึ้น server จริงอัตโนมัติ**
- ไม่ต้อง deploy ด้วยมือ

---

## ทำไมต้องใช้ CI/CD?

✅ **ลดข้อผิดพลาด** - ตรวจสอบโค้ดอัตโนมัติก่อน deploy
✅ **ประหยัดเวลา** - ไม่ต้อง build/test/deploy ด้วยมือ
✅ **Deploy บ่อยขึ้น** - อัพเดทฟีเจอร์ใหม่ได้เร็ว
✅ **ทำงานเป็นทีม** - รู้ทันทีว่าโค้ดใครทำพัง

---

## สิ่งที่ต้องเรียนรู้ (เรียงตามลำดับ)

### 1. Git พื้นฐาน (จำเป็นมาก!)
```bash
git add .
git commit -m "ข้อความอธิบายการแก้ไข"
git push origin main
```

**ทำไมต้องรู้:** CI/CD ทำงานเมื่อคุณ push โค้ดขึ้น Git

---

### 2. เครื่องมือ CI/CD ยอดนิยม

| เครื่องมือ | เหมาะกับ | ความยาก |
|-----------|---------|---------|
| **GitHub Actions** | โปรเจค GitHub | ⭐ ง่ายที่สุด |
| **GitLab CI/CD** | โปรเจค GitLab | ⭐⭐ ง่าย |
| **Jenkins** | องค์กรใหญ่ | ⭐⭐⭐⭐ ยาก |
| **CircleCI** | Startup | ⭐⭐ ง่าย |

**แนะนำเริ่มต้นที่:** GitHub Actions (ฟรี และใช้งานง่าย)

---

### 3. ไฟล์ Workflow พื้นฐาน

สร้างไฟล์: `.github/workflows/ci.yml`

```yaml
name: Go CI Pipeline

# เมื่อไหร่ให้ทำงาน
on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

# สิ่งที่ต้องทำ
jobs:
  build-and-test:
    runs-on: ubuntu-latest

    steps:
    # 1. ดึงโค้ดมา
    - uses: actions/checkout@v3

    # 2. ติดตั้ง Go
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    # 3. ติดตั้ง dependencies
    - name: Download dependencies
      run: go mod download

    # 4. รัน tests
    - name: Run tests
      run: go test -v ./...

    # 5. Build โปรเจค
    - name: Build
      run: go build -v ./...
```

---

### 4. Pipeline คืออะไร?

**Pipeline** = ขั้นตอนการทำงานอัตโนมัติ

```
โค้ด Push → ดึงโค้ด → ติดตั้ง → ทดสอบ → Build → Deploy
```

**ตัวอย่างจริง:**
1. คุณแก้ไขโค้ดและ `git push`
2. GitHub Actions เริ่มทำงาน
3. รัน `go mod download` ดาวน์โหลด dependencies
4. รัน `go test` ตรวจสอบโค้ด
5. รัน `go build` สร้างไฟล์ binary
6. ถ้าผ่านหมด → Deploy ขึ้น server

---

## ตัวอย่างโปรเจคจริง

### โปรเจค Go Web Server

**.github/workflows/deploy.yml**
```yaml
name: Deploy Go App

on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3

    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Download dependencies
      run: go mod download

    - name: Run tests
      run: go test -v ./...

    - name: Build
      run: go build -o app

    - name: Deploy to server
      run: |
        echo "Deploying to production..."
        # เพิ่มคำสั่ง deploy ของคุณที่นี่
```

---

### โปรเจค Go + Docker

**.github/workflows/docker.yml**
```yaml
name: Build and Push Docker

on: [push, pull_request]

jobs:
  docker:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3

    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Run tests
      run: go test -v ./...

    - name: Build Docker image
      run: docker build -t myapp:latest .

    - name: Push to registry
      run: |
        echo "Push image to registry..."
        # docker push myapp:latest
```

---

## การ Deploy จริง (ตัวอย่างง่ายๆ)

### Deploy Go App ไปยัง Server ผ่าน SSH

```yaml
name: Deploy Go to Server

on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3

    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Build for Linux
      run: GOOS=linux GOARCH=amd64 go build -o app

    - name: Deploy to server
      uses: appleboy/scp-action@master
      with:
        host: ${{ secrets.SERVER_HOST }}
        username: ${{ secrets.SERVER_USER }}
        key: ${{ secrets.SSH_KEY }}
        source: "app"
        target: "/var/www/myapp"
```

---

## Secrets คืออะไร?

**Secrets** = ข้อมูลลับที่ไม่ควรเปิดเผย
- API Keys
- รหัสผ่าน Database
- Token ต่างๆ

**วิธีเพิ่ม Secrets:**
1. ไปที่ GitHub → Settings → Secrets and variables → Actions
2. คลิก "New repository secret"
3. ใส่ชื่อและค่า
4. ใช้งานใน workflow: `${{ secrets.YOUR_SECRET }}`

---

## เริ่มต้นอย่างไร? (5 ขั้นตอน)

### ขั้นตอนที่ 1: สร้างโปรเจคทดสอบ
```bash
mkdir my-ci-cd-test
cd my-ci-cd-test
go mod init github.com/yourusername/my-ci-cd-test
```

### ขั้นตอนที่ 2: สร้างไฟล์ทดสอบ
```go
// main.go
package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello CI/CD!")
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
```

```go
// main_test.go
package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/", nil)
    w := httptest.NewRecorder()

    handler(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }

    expected := "Hello CI/CD!"
    if w.Body.String() != expected {
        t.Errorf("Expected %s, got %s", expected, w.Body.String())
    }
}
```

### ขั้นตอนที่ 3: ทดสอบรันโปรเจค
```bash
# รันโปรเจค
go run main.go

# รัน tests
go test -v
```

### ขั้นตอนที่ 4: สร้าง GitHub Repo
```bash
git init
git add .
git commit -m "Initial commit"
git remote add origin YOUR_GITHUB_URL
git push -u origin main
```

### ขั้นตอนที่ 5: สร้าง Workflow
สร้างไฟล์ `.github/workflows/ci.yml` (ใช้โค้ดด้านบน)

---

## เคล็ดลับการเรียนรู้

✅ **เริ่มจากโปรเจคเล็ก** - อย่าเอาโปรเจคใหญ่มาทดสอบครั้งแรก
✅ **ทำทีละขั้นตอน** - เพิ่ม step ทีละอย่าง อย่าเพิ่มครั้งเดียวทั้งหมด
✅ **ดู Logs** - เมื่อมี error ดู logs ใน GitHub Actions
✅ **ลอกตัวอย่าง** - ดูโปรเจคคนอื่นใน GitHub

---

## ข้อผิดพลาดที่มือใหม่มักพบ

❌ **ไม่มี test** - สร้าง test ง่ายๆ สักอัน
❌ **ใส่ secrets ในโค้ด** - ใช้ GitHub Secrets เสมอ
❌ **Workflow ซับซ้อนเกินไป** - เริ่มจากง่ายๆ ก่อน
❌ **ไม่อ่าน error logs** - อ่าน logs ทุกครั้งที่มี error

---

## ทรัพยากรเพิ่มเติม

📚 **เอกสารภาษาไทย:**
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [GitLab CI/CD](https://docs.gitlab.com/ee/ci/)

🎥 **วิดีโอแนะนำ:**
- ค้นหา "GitHub Actions Tutorial" บน YouTube

💡 **ฝึกฝน:**
- สร้างโปรเจคง่ายๆ และเพิ่ม CI/CD
- ทดลอง push โค้ดและดู pipeline ทำงาน

---

## สรุป

**CI/CD ไม่ยากอย่างที่คิด!**

1. เขียนโค้ด → Push ไป Git
2. Workflow ทำงานอัตโนมัติ
3. ตรวจสอบ → Build → Deploy

**เริ่มวันนี้เลย:** สร้างไฟล์ `.github/workflows/ci.yml` ในโปรเจคของคุณ

---

*สร้างโดย CI/CD Learning Guide - เริ่มต้นง่าย ใช้งานได้จริง 🚀*
