# Backend Spatium

Selamat datang di repositori backend Spatium. Backend ini adalah inti dari aplikasi Spatium, yang menyediakan layanan API untuk mendukung fitur-fitur kesehatan mental.

Backend ini dirancang untuk menjadi server yang stabil, aman, dan efisien, memastikan pengalaman pengguna yang lancar dan andal.

---

## Teknologi yang Digunakan

Backend Spatium dibangun menggunakan ekosistem JavaScript, dengan fokus pada performa dan skalabilitas:

- **Node.js & Express.js**: Menggunakan Node.js sebagai runtime dan Express.js sebagai framework web untuk membangun API RESTful yang cepat dan terstruktur.
- **MongoDB & Mongoose**: Mengandalkan MongoDB sebagai database NoSQL yang fleksibel, dikelola dengan Mongoose untuk memodelkan data secara efektif.
- **JWT (JSON Web Tokens)**: Untuk mengelola autentikasi pengguna secara aman, memastikan hanya pengguna yang sah yang bisa mengakses data mereka.
- **Socket.IO**: Untuk mengaktifkan komunikasi real-time, mendukung fitur komunitas di mana pengguna bisa berinteraksi secara instan.
- **TensorFlow.js**: Integrasi untuk fungsionalitas AI, seperti memberikan respons empatik berdasarkan pemahaman emosi pengguna.

---

## Fitur API Utama

Backend ini menyediakan serangkaian endpoint API untuk mendukung fungsionalitas utama aplikasi:

- **Autentikasi & Pengguna**: Mengelola pendaftaran, login, dan profil pengguna.
- **Curhat Anonim**: Endpoint untuk membuat, melihat, dan mengelola curhatan anonim.
- **Mood & Refleksi**: API untuk melacak dan menyimpan data mood harian serta refleksi diri.
- **Komunitas**: Endpoint untuk interaksi komunitas, termasuk postingan, komentar, dan dukungan real-time.
- **AI Empatik**: Endpoint khusus untuk berinteraksi dengan model AI yang memberikan respons empatik.
