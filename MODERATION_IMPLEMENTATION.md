# Implementasi Keamanan & Moderasi Konten

## Ringkasan

Dokumen ini menjelaskan implementasi fitur keamanan dan moderasi konten untuk aplikasi Spatium, menangani concern kritis tentang keamanan pengguna dalam aplikasi kesehatan mental.

## Fitur yang Diimplementasikan

### 1. 🛡️ Automated Content Moderation

**Backend (Go):**
- **File:** `utils/moderation.go`
- **Fungsi Utama:**
  - `ModerateContent()` - Memeriksa konten untuk material berbahaya/toxic
  - `SanitizeContent()` - Mengganti kata-kata kasar dengan asterisk
  - `ContainsCrisisIndicators()` - Mendeteksi indikator krisis kesehatan mental
  - `GetCrisisSupportMessage()` - Mengembalikan pesan dukungan krisis

**Fitur Moderasi:**
1. **Word Filter (Local):**
   - Daftar kata-kata toxic dalam Bahasa Indonesia & Inggris
   - Daftar trigger words untuk indikator krisis (bunuh diri, self-harm)
   - Regex pattern untuk hate speech

2. **OpenAI Moderation API:**
   - Integrasi dengan OpenAI Moderation endpoint
   - Deteksi: hate, harassment, self-harm, violence, sexual content
   - Fallback jika API tidak tersedia

3. **Crisis Detection:**
   - Mendeteksi kata-kata terkait bunuh diri dan self-harm
   - Tidak memblokir, tapi menampilkan sumber bantuan
   - Pesan dukungan dengan nomor hotline (Into The Light, Yayasan Pulih, dll)

**Integrasi:**
- `services/post_service.go` - Moderasi sebelum post dipublikasikan
- `services/comment_service.go` - Moderasi sebelum komentar dipublikasikan

### 2. 🚨 Fitur Report

**Backend:**
- **Model:** `models/report.go`
- **Repository:** `repositories/report_repository.go`
- **Service:** `services/report_service.go`
- **Controller:** `controllers/report_controller.go`
- **Migration:** `database/migrations/000012_create_reports_table.up.sql`

**Endpoints:**
```
POST   /api/v1/reports              - Membuat laporan baru
GET    /api/v1/reports/me           - Melihat laporan saya
GET    /api/v1/reports/reasons      - Mendapatkan daftar alasan laporan
GET    /api/v1/admin/reports/pending - [Admin] Melihat laporan pending
PUT    /api/v1/admin/reports/:id/status - [Admin] Update status laporan
```

**Jenis Laporan (ReportType):**
- `post` - Melaporkan postingan
- `comment` - Melaporkan komentar
- `user` - Melaporkan pengguna

**Alasan Laporan (ReportReason):**
- `harassment` - Pelecehan
- `hate_speech` - Ujaran Kebencian
- `spam` - Spam
- `self_harm` - Menyakiti Diri Sendiri
- `violence` - Kekerasan
- `inappropriate` - Konten Tidak Pantas
- `other` - Lainnya

**Frontend (Flutter):**
- **Model:** `lib/features/timeline/data/models/report_model.dart`
- **Widget:** `lib/features/timeline/presentation/widget/report_dialog.dart`

### 3. 🚫 Fitur Block

**Backend:**
- **Model:** `models/block.go`
- **Repository:** `repositories/block_repository.go`
- **Service:** `services/block_service.go`
- **Controller:** `controllers/block_controller.go`
- **Migration:** `database/migrations/000013_create_blocks_table.up.sql`

**Endpoints:**
```
POST   /api/v1/users/block          - Memblokir pengguna
DELETE /api/v1/users/block/:user_id - Membuka blokir
GET    /api/v1/users/blocked        - Daftar pengguna yang diblokir
```

**Efek Block:**
1. Postingan pengguna yang diblokir tidak muncul di timeline
2. Komentar pengguna yang diblokir tidak terlihat
3. Pengguna yang memblokir tidak bisa dilihat oleh yang diblokir (mutual hiding)

**Frontend (Flutter):**
- **Model:** `lib/features/timeline/data/models/block_model.dart`
- **Data Source:** `lib/features/timeline/data/datasources/moderation_remote_data_source.dart`
- **Provider:** `lib/features/timeline/presentation/providers/moderation_provider.dart`
- **Widgets:**
  - `lib/features/timeline/presentation/widget/block_dialog.dart` - Dialog konfirmasi block
  - `lib/features/timeline/presentation/widget/post_actions_menu.dart` - Menu aksi post/comment

## Cara Penggunaan

### Frontend - Menambahkan Menu ke Post Card

```dart
import 'package:spatium/features/timeline/presentation/widget/post_actions_menu.dart';

// Dalam widget post card:
PostActionsMenu(
  postId: post.publicId,
  postContent: post.content,
  postOwnerUserId: post.userId,  // jika tersedia
  postOwnerAlias: post.alias,    // jika tersedia
  isOwner: post.isOwner,
  onDelete: () => deletePost(post.publicId),
  onReported: () => refreshPosts(),
  onBlocked: () => refreshPosts(),
)
```

### Frontend - Menampilkan Halaman Blocked Users

```dart
import 'package:spatium/features/timeline/presentation/widget/block_dialog.dart';

Navigator.push(
  context,
  MaterialPageRoute(builder: (context) => const BlockedUsersPage()),
);
```

## Database Migrations

Jalankan migrations untuk membuat tabel baru:

```bash
# Menggunakan golang-migrate
migrate -path database/migrations -database "postgres://..." up
```

Atau tables akan otomatis di-create via GORM AutoMigrate saat server start.

## Konfigurasi

### Environment Variables

Pastikan `.env` memiliki:
```
OPENAI_API_KEY=your_openai_api_key  # Untuk OpenAI Moderation API
```

## Flow Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                    USER CREATES POST/COMMENT                     │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                    CONTENT MODERATION CHECK                      │
│  1. Check toxic words (local list)                              │
│  2. Check hate speech patterns (regex)                          │
│  3. Check OpenAI Moderation API (if available)                  │
└─────────────────────────────────────────────────────────────────┘
                                │
            ┌───────────────────┼───────────────────┐
            │                   │                   │
            ▼                   ▼                   ▼
      ┌──────────┐       ┌──────────┐       ┌──────────────┐
      │ BLOCKED  │       │  PASSED  │       │CRISIS DETECT │
      │  ❌      │       │  ✅      │       │  ⚠️          │
      └──────────┘       └──────────┘       └──────────────┘
            │                   │                   │
            ▼                   ▼                   ▼
      Return Error        Save Post         Save Post +
      with Reason         to Database       Show Support
                                           Resources
```

## Keamanan Tambahan yang Disarankan

1. **Rate Limiting:** Batasi jumlah report per user per hari
2. **Admin Dashboard:** Buat dashboard untuk moderator meninjau laporan
3. **Email Notification:** Kirim email ke admin untuk laporan critical
4. **Auto-ban:** Implementasi auto-ban untuk user dengan banyak laporan valid
5. **Content Hash:** Simpan hash konten yang diblokir untuk mencegah re-posting

## File yang Dibuat/Dimodifikasi

### Backend (Go):
- ✅ `utils/moderation.go` (NEW)
- ✅ `models/report.go` (NEW)
- ✅ `models/block.go` (NEW)
- ✅ `repositories/report_repository.go` (NEW)
- ✅ `repositories/block_repository.go` (NEW)
- ✅ `repositories/user_repository.go` (MODIFIED - added GetByID, GetByPublicID)
- ✅ `repositories/post_repository.go` (MODIFIED - added GetAllExcluding)
- ✅ `repositories/comment_repository.go` (MODIFIED - added GetCommentsByPostIDExcluding)
- ✅ `services/moderation_service.go` (NEW)
- ✅ `services/report_service.go` (NEW)
- ✅ `services/block_service.go` (NEW)
- ✅ `services/post_service.go` (MODIFIED - integrated moderation)
- ✅ `services/comment_service.go` (MODIFIED - integrated moderation)
- ✅ `controllers/report_controller.go` (NEW)
- ✅ `controllers/block_controller.go` (NEW)
- ✅ `controllers/post_controller.go` (MODIFIED - updated for block filtering)
- ✅ `controllers/comment_controller.go` (MODIFIED - updated for block filtering)
- ✅ `routes/route.go` (MODIFIED - added new endpoints)
- ✅ `server/server.go` (MODIFIED - wiring new components)
- ✅ `database/migrations/000012_create_reports_table.up.sql` (NEW)
- ✅ `database/migrations/000012_create_reports_table.down.sql` (NEW)
- ✅ `database/migrations/000013_create_blocks_table.up.sql` (NEW)
- ✅ `database/migrations/000013_create_blocks_table.down.sql` (NEW)

### Frontend (Flutter):
- ✅ `lib/core/constants/api_constants.dart` (MODIFIED - added endpoints)
- ✅ `lib/features/timeline/data/models/report_model.dart` (NEW)
- ✅ `lib/features/timeline/data/models/block_model.dart` (NEW)
- ✅ `lib/features/timeline/data/datasources/moderation_remote_data_source.dart` (NEW)
- ✅ `lib/features/timeline/presentation/providers/moderation_provider.dart` (NEW)
- ✅ `lib/features/timeline/presentation/widget/report_dialog.dart` (NEW)
- ✅ `lib/features/timeline/presentation/widget/block_dialog.dart` (NEW)
- ✅ `lib/features/timeline/presentation/widget/post_actions_menu.dart` (NEW)
