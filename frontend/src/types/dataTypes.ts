// ===================== Types =====================
export type Department = {
  id: number;
  department_name: string;
  max_clock_in_time: string;
  max_clock_out_time: string;
};

export type Employees = {
  id: number;
  employee_id: string;
  department_id: number;
  name: string;
  address: string;
  department: Department;
};

export type AttendanceLog = {
  id: number;
  employee_id: string;
  name?: string;
  attendance_id: string;
  date_attendance: string;
  attendance_type: number;
  description: string;
  department?: string;
  clock_in?: string;
  clock_out?: string;
};

export type FormData = {
  departmentName: string;
  maxClockIn: string;
  maxClockOut: string;
  DepartmentID: string;
  Name: string;
  Address: string;
  EmployeeID: string;
  [key: string]: string;
}

export type DepartmentWithEmployees = Department & {
  employees: Employees[];
};

export type TableType = "departments" | "employees" | "attendanceLogs";

