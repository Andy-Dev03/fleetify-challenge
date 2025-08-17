"use client";

import { useState, useEffect } from "react";
import TableData from "../components/TableData";
import {
  Department,
  Employees,
  AttendanceLog,
  TableType,
} from "../types/dataTypes";

// fetch helper
async function fetchData<T>(url: string): Promise<T> {
  const res = await fetch(url, { method: "GET", cache: "no-store" });
  if (!res.ok) throw new Error(`Failed to fetch: ${url}`);
  const data = await res.json();
  return data.data;
}

export default function Main() {
  // Condition to switch between table types
  const [tableType, setTableType] = useState<TableType>("departments");

  // State for each table type (departments, employees, attendanceLogs)
  const [departments, setDepartments] = useState<Department[]>([]);
  const [employees, setEmployees] = useState<Employees[]>([]);
  const [attendanceLogs, setAttendanceLogs] = useState<AttendanceLog[]>([]);

  // State for filters (attendance logs)
  const [filterDate, setFilterDate] = useState<string>("");
  const [filterDepartment, setFilterDepartment] = useState<string>("");

  const getDepartements = async () => {
    const data = await fetchData<Department[]>(
      "http://localhost:8080/api/departements"
    );
    setDepartments(data);
  };

  const getEmployees = async () => {
    const data = await fetchData<Employees[]>(
      "http://localhost:8080/api/employees"
    );
    setEmployees(data);
  };

  const getAttendanceLogs = async (date?: string, departmentId?: string) => {
    let url = "http://localhost:8080/api/attendance/logs";
    const params = [];
    if (date) params.push(`date=${date}`);
    if (departmentId) params.push(`department_id=${departmentId}`);
    if (params.length > 0) url += `?${params.join("&")}`;
    const data = await fetchData<AttendanceLog[]>(url);
    console.log(data);

    setAttendanceLogs(data);
  };

  useEffect(() => {
    getDepartements();
    getEmployees();
    getAttendanceLogs(filterDate, filterDepartment);
  }, [filterDate, filterDepartment]);

  // Button group for switching table type
  const ButtonGroup = () => (
    <div className="flex gap-2 mb-4">
      {[
        { label: "Departments", value: "departments" },
        { label: "Employees", value: "employees" },
        { label: "Attendance Logs", value: "attendanceLogs" },
      ].map((btn) => (
        <button
          key={btn.value}
          onClick={() => setTableType(btn.value as TableType)}
          className={`px-4 py-1 rounded-md transition ${
            tableType === btn.value
              ? "bg-gray-800 text-white"
              : "bg-gray-200 hover:bg-gray-300"
          }`}
        >
          {btn.label}
        </button>
      ))}
    </div>
  );

  // Filter section for attendance logs
  const AttendanceFilter = () => (
    <div className="mb-4 flex gap-4 items-center">
      <input
        type="date"
        value={filterDate}
        onChange={(e) => setFilterDate(e.target.value)}
        className="border px-2 py-1 rounded"
      />
      <select
        value={filterDepartment}
        onChange={(e) => setFilterDepartment(e.target.value)}
        className="border px-2 py-1 rounded"
      >
        <option value="">All Departments</option>
        {(departments ?? []).map((dept) => (
          <option key={dept.id} value={dept.id}>
            {dept.department_name}
          </option>
        ))}
      </select>
    </div>
  );

  // Table section
  const TableSection = () => {
    if (tableType === "departments") {
      const handleDeleteDepartment = async () => {
        await getDepartements();
        await getEmployees();
        await getAttendanceLogs(filterDate, filterDepartment);
      };
      return (
        <TableData
          tableType="departments"
          data={departments}
          onDeleteSuccess={handleDeleteDepartment}
        />
      );
    }
    if (tableType === "employees") {
      const handleDeleteEmployee = async () => {
        await getEmployees();
        await getAttendanceLogs(filterDate, filterDepartment);
      };
      return (
        <TableData
          tableType="employees"
          data={employees}
          onDeleteSuccess={handleDeleteEmployee}
        />
      );
    }
    if (tableType === "attendanceLogs") {
      return <TableData tableType="attendanceLogs" data={attendanceLogs} />;
    }
    return null;
  };

  return (
    <div className="pl-2 space-y-2">
      <h1 className="font-bold text-2xl">All of The Data</h1>
      <p className="text-gray-600">Choose the data you want to read</p>
      <ButtonGroup />
      {tableType === "attendanceLogs" && <AttendanceFilter />}
      <TableSection />
    </div>
  );
}
