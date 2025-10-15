# Backend Spatium

Selamat datang di repositori backend Spatium. Backend ini adalah inti dari aplikasi Spatium, yang menyediakan layanan API untuk mendukung fitur-fitur kesehatan mental.

Backend ini dirancang untuk menjadi server yang stabil, aman, dan efisien, memastikan pengalaman pengguna yang lancar dan andal.

---

## 🧠 Teknologi yang Digunakan

Spatium dibangun dengan fokus pada kecepatan, keamanan, dan efisiensi sumber daya, menggunakan ekosistem Golang yang ringan namun powerful.

### ⚙️ Backend

- **Go (Golang)**: Bahasa utama untuk membangun backend yang efisien, concurrency-friendly, dan mudah di-deploy.
- **Fiber**: Web framework modern berbasis fasthttp yang sangat cepat untuk membangun API RESTful.
- **GORM**: ORM (Object Relational Mapping) untuk mengelola relasi data dengan PostgreSQL secara elegan.
- **PostgreSQL**: Sistem basis data relasional yang andal dengan dukungan kuat untuk relasi antar entitas seperti users, posts, dan mood tags.
- **Go-Migrate**: Digunakan untuk manajemen versi skema database dan migrasi otomatis antar environment.
- **JWT (JSON Web Token)**: Menangani autentikasi anonim secara aman melalui token unik yang merepresentasikan identitas pengguna.

---

## Fitur API Utama

Backend ini menyediakan serangkaian endpoint API untuk mendukung fungsionalitas utama aplikasi:

- **Autentikasi & Pengguna**: Mengelola pendaftaran, login, dan profil pengguna.
- **Curhat Anonim**: Endpoint untuk membuat, melihat, dan mengelola curhatan anonim.
- **Mood & Refleksi**: API untuk melacak dan menyimpan data mood harian serta refleksi diri.
- **Komunitas**: Endpoint untuk interaksi komunitas, termasuk postingan, komentar, dan dukungan real-time.
- **AI Empatik**: Endpoint khusus untuk berinteraksi dengan model AI yang memberikan respons empatik.
