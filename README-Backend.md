# 📘 Attendance Management API (Golang)

## Tech Stack

- **Framework**: [Gin](https://gin-gonic.com/)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: MySQL (via `gorm`)
- **Bahasa**: Go

---

## 📂 Struktur Proyek

```
.
├── controllers/         # Handler API
│   ├── attendance_controller.go
│   ├── department_controller.go
│   └── employee_controller.go
├── models/              # ORM Models
│   ├── attendance.go
│   ├── attendance_history.go
│   ├── department.go
│   └── employee.go
├── routes.go            # Registrasi route API
├── config/              # Database connection (DB instance)
```

---

## 📑 Models

### `Employee`

```go
ID           uint
EmployeeID   string   // contoh: EMP-001
DepartmentID uint
Name         string
Address      string
```

### `Department`

```go
ID              uint
DepartmentName  string
MaxClockInTime  string   // format HH:MM:SS
MaxClockOutTime string   // format HH:MM:SS
Employees       []Employee
```

### `Attendance`

```go
ID           uint
EmployeeID   string
AttendanceID string   // contoh: ATT-001
ClockIn      time.Time
ClockOut     *time.Time
Employee     Employee
```

### `AttendanceHistory`

```go
ID             uint
EmployeeID     string
AttendanceID   string
DateAttendance time.Time
AttendanceType int    // 1=Clock In, 2=Clock Out
Description    string
```

---

## API Endpoints

### Employee

| Method | Endpoint            | Deskripsi                           |
| ------ | ------------------- | ----------------------------------- |
| GET    | `/api/employees`    | Ambil semua employee + department   |
| GET    | `/api/employee/:id` | Ambil detail employee               |
| POST   | `/api/employee`     | Tambah employee baru                |
| PATCH  | `/api/employee/:id` | Update data employee                |
| DELETE | `/api/employee/:id` | Hapus employee + attendance terkait |

### Department

| Method | Endpoint               | Deskripsi                                              |
| ------ | ---------------------- | ------------------------------------------------------ |
| GET    | `/api/departements`    | Ambil semua department + employees                     |
| GET    | `/api/departement/:id` | Detail department                                      |
| POST   | `/api/departement`     | Tambah department baru                                 |
| PATCH  | `/api/departement/:id` | Update department                                      |
| DELETE | `/api/departement/:id` | Hapus department + semua employee + attendance terkait |

### Attendance

| Method | Endpoint               | Deskripsi                                          |
| ------ | ---------------------- | -------------------------------------------------- |
| POST   | `/api/attendance`      | Clock In (absen masuk)                             |
| PUT    | `/api/attendance/:id`  | Clock Out (absen keluar)                           |
| GET    | `/api/attendance/logs` | Ambil log absensi (filter by tanggal / departemen) |

---

## 📝 Catatan Penting

- **EmployeeID** dan **AttendanceID** digenerate otomatis (`EMP-xxx`, `ATT-xxx`).
- AttendanceHistory menyimpan jejak setiap kali Clock In / Clock Out.
- Delete Department → semua employee & attendance terkait ikut terhapus.
- Jam maksimal Clock In/Out mengikuti konfigurasi department.

# API Documentation

## Endpoints :

### Employee

- `GET /api/employees`
- `GET /api/employee/:id`
- `POST /api/employee`
- `PATCH /api/employee/:id`
- `DELETE /api/employee/:id`

### Department

- `GET /api/departements`
- `GET /api/departement/:id`
- `POST /api/departement`
- `PATCH /api/departement/:id`
- `DELETE /api/departement/:id`

### Attendance

- `POST /api/attendance`
- `PUT /api/attendance/:id`
- `GET /api/attendance/logs`

---

## 1. GET /api/employees

**Response (200 - OK)**

```json
{
  "data": [
    {
      "id": 1,
      "employee_id": "EMP-001",
      "department_id": 1,
      "name": "John Doe",
      "address": "Jakarta",
      "created_at": "2025-08-17T08:00:00Z",
      "updated_at": "2025-08-17T08:00:00Z",
      "department": {
        "id": 1,
        "department_name": "IT",
        "max_clock_in_time": "09:00:00",
        "max_clock_out_time": "17:00:00",
        "created_at": "2025-08-17T08:00:00Z",
        "updated_at": "2025-08-17T08:00:00Z"
      }
    }
  ]
}
```

**Response (500 - Internal Server Error)**

```json
{
  "error": "Internal Server Error"
}
```

---

## 2. GET /api/employee/:id

**Description**  
Get employee detail by ID.

**Request Params**

```json
{
  "id": "integer (required)"
}
```

**Response (200 - OK)**

```json
{
  "data": {
    "id": 1,
    "employee_id": "EMP-001",
    "department_id": 1,
    "name": "John Doe",
    "address": "Jakarta",
    "created_at": "2025-08-17T08:00:00Z",
    "updated_at": "2025-08-17T08:00:00Z",
    "department": {
      "id": 1,
      "department_name": "IT",
      "max_clock_in_time": "09:00:00",
      "max_clock_out_time": "17:00:00",
      "created_at": "2025-08-17T08:00:00Z",
      "updated_at": "2025-08-17T08:00:00Z"
    }
  }
}
```

**Response (404 - Not Found)**

```json
{
  "error": "Employee not found"
}
```

---

## 3. POST /api/employee

**Description**  
Create a new employee.

**Request Body**

```json
{
  "department_id": 1,
  "name": "Jane Doe",
  "address": "Bandung"
}
```

**Response (200 - OK)**

```json
{
  "data": {
    "id": 2,
    "employee_id": "EMP-002",
    "department_id": 1,
    "name": "Jane Doe",
    "address": "Bandung",
    "created_at": "2025-08-17T08:00:00Z",
    "updated_at": "2025-08-17T08:00:00Z",
    "department": {
      "id": 2,
      "department_name": "FrontEnd Dev",
      "max_clock_in_time": "09:00:00",
      "max_clock_out_time": "17:00:00",
      "created_at": "2025-08-17T08:00:00Z",
      "updated_at": "2025-08-17T08:00:00Z"
    }
  }
}
```

**Response (400 - Bad Request)**

```json
{
  "error": "Department is required"

}
OR
{
  "error": "Name is required"
}
OR
{
  "error": "Address is required"
}
```

---

## 4. PATCH /api/employee/:id

**Description**  
Update employee by ID.

**Request Params**

```json
{
  "id": "integer (required)"
}
```

**Request Body**

```json
{
  "department_id": 1,
  "name": "Jane Smith",
  "address": "Surabaya"
}
```

**Response (200 - OK)**

```json
{
  "data": {
    "id": 2,
    "employee_id": "EMP-002",
    "department_id": 1,
    "name": "Jane Smith",
    "address": "Surabaya",
    "created_at": "2025-08-17T08:00:00Z",
    "updated_at": "2025-08-17T09:00:00Z",
    "department": {
      "id": 1,
      "department_name": "IT",
      "max_clock_in_time": "09:00:00",
      "max_clock_out_time": "17:00:00",
      "created_at": "2025-08-17T08:00:00Z",
      "updated_at": "2025-08-17T08:00:00Z"
    }
  }
}
```

**Response (400 - Bad Request)**

```json
{
  "error": "Department is required"

}
OR
{
  "error": "Name is required"
}
OR
{
  "error": "Address is required"
}
```

**Response (404 - Not Found)**

```json
{
  "error": "Employee not found"
}
```

---

## 5. DELETE /api/employee/:id

**Description**  
Delete employee and related attendance by ID.

**Response (200 - OK)**

```json
{
  "message": "Employee deleted successfully"
}
```

**Response (404 - Not Found)**

```json
{
  "error": "Employee not found"
}
```

---

## 6. GET /api/departements

**Description**  
Get all departments including employees.

**Response (200 - OK)**

```json
{
  "data": [
    {
      "id": 1,
      "department_name": "IT",
      "max_clock_in_time": "09:00:00",
      "max_clock_out_time": "17:00:00",
      "created_at": "2025-08-17T08:00:00Z",
      "updated_at": "2025-08-17T09:00:00Z",
      "employees": [
        {
          "id": 1,
          "employee_id": "EMP-001",
          "department_id": 1,
          "name": "John Doe",
          "address": "Jakarta",
          "created_at": "2025-08-17T08:00:00Z",
          "updated_at": "2025-08-17T09:00:00Z"
        }
      ]
    }
  ]
}
```

**Response (500 - Internal Server Error)**

```json
{
  "error": "Internal Server Error"
}
```

---

## 7. GET /api/departement/:id

**Description**  
Get department detail including employees.

**Response (200 - OK)**

```json
{
  "data": {
    "id": 1,
    "department_name": "IT",
    "max_clock_in_time": "09:00:00",
    "max_clock_out_time": "17:00:00",
    "created_at": "2025-08-17T08:00:00Z",
    "updated_at": "2025-08-17T09:00:00Z",
    "employees": [
      {
        "id": 1,
        "employee_id": "EMP-001",
        "department_id": 1,
        "name": "John Doe",
        "address": "Jakarta",
        "created_at": "2025-08-17T08:00:00Z",
        "updated_at": "2025-08-17T09:00:00Z"
      }
    ]
  }
}
```

**Response (404 - Not Found)**

```json
{
  "error": "Department not found"
}
```

---

## 8. POST /api/departement

**Description**  
Create a new department.

**Request Body (x-www-form-urlencoded)**

```json
{
  "department_name": 1,
  "max_clock_in_time": "09:00",
  "max_clock_out_time": "17:00"
}
```

**Response (200 - OK)**

```json
{
  "data": {
    "id": 1,
    "department_name": "IT",
    "max_clock_in_time": "09:00:00",
    "max_clock_out_time": "17:00:00",
    "created_at": "2025-08-17T08:00:00Z",
    "updated_at": "2025-08-17T08:00:00Z",
    "employees": []
  }
}
```

**Response (400 - Bad Request)**

```json
{
  "error": "Department Name is required"
}
OR
{
  "error": "Max Clock In Time is required"
}
OR
{
  "error": "Max Clock Out Time is required"
}
OR
{
    "error": "invalid format for max_clock_in_time, expected HH:mm"
}
OR
{
    "error": "invalid format for max_clock_out_time, expected HH:mm"
}

```

---

## 9. PATCH /api/departement/:id

**Description**  
Update department by ID.

**Request Body (x-www-form-urlencoded)**

```json
{
  "department_name": 1,
  "max_clock_in_time": "09:00",
  "max_clock_out_time": "17:00"
}
```

**Response (200 - OK)**

```json
{
  "data": {
    "id": 1,
    "department_name": "Finance",
    "max_clock_in_time": "08:30:00",
    "max_clock_out_time": "16:30:00",
    "created_at": "2025-08-17T08:00:00Z",
    "updated_at": "2025-08-17T08:00:00Z"
  }
}
```

**Response (400 - Bad Request)**

```json
{
  "error": "Department Name is required"
}
OR
{
  "error": "Max Clock In Time is required"
}
OR
{
  "error": "Max Clock Out Time is required"
}
OR
{
    "error": "invalid format for max_clock_in_time, expected HH:mm"
}
OR
{
    "error": "invalid format for max_clock_out_time, expected HH:mm"
}

```

**Response (404 - Not Found)**

```json
{
  "error": "Department not found"
}
```

---

## 10. DELETE /api/departement/:id

**Description**  
Delete department and related employees and attendance.

**Response (200 - OK)**

```json
{
  "message": "Department deleted successfully"
}
```

**Response (404 - Not Found)**

```json
{
  "error": "Department not found"
}
```

---

## 11. POST /api/attendance

**Description**  
Employee clock in.

**Request Body**

```json
{
  "employee_id": "EMP-001",
  "clock_in": "2025-08-17 08:55:00"
}
```

**Response (200 - OK)**

```json
{
  "data": {
    "id": 1,
    "employee_id": "EMP-001",
    "attendance_id": "ATT-001",
    "clock_in": "2025-08-17T08:55:00Z",
    "clock_out": null
  }
}
```

**Response (400 - Bad Request)**

```json
{
  "error": "Invalid input format"
}
OR
{
    "error": "Employee is required"
}
OR
{
    "error": "Clock In is required"
}
```

---

## 12. PUT /api/attendance/:id

**Description**  
Employee clock out.

**Request Params**

```json
{
  "attendance_id": "string(required)"
}
```

**Request Body**

```json
{
  "clock_out": "2025-08-17 17:05:00"
}
```

**Response (200 - OK)**

```json
{
  "data": {
    "id": 1,
    "employee_id": "EMP-001",
    "attendance_id": "ATT-001",
    "clock_in": "2025-08-17T08:55:00Z",
    "clock_out": "2025-08-17T17:05:00Z"
  }
}
```

**Response (400 - Bad Request)**

```json
{
  "error": "Invalid input format"
}
OR
{
    "error": "Employee is required"
}
OR
{
    "error": "Clock In is required"
}
```

**Response (404 - Not Found)**

```json
{
  "error": "Attendance not found"
}
```

---

## 13. GET /api/attendance/logs

**Description**  
Get attendance logs with optional filters (date, department).

**Request Query**

```
?date=2025-08-17&department_id=1
```

**Response (200 - OK)**

```json
{
  "data": [
    {
      "id": 1,
      "employee_id": "EMP-001",
      "attendance_id": "ATT-001",
      "name": "John Doe",
      "date_attendance": "2025-08-17 08:55:00",
      "attendance_type": 1,
      "description": "On Time (Check-in)",
      "department": "IT",
      "clock_in": "08:55:00",
      "clock_out": "17:05:00"
    }
  ]
}
```

---
