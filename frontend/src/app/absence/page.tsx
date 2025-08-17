"use client";
import { useRouter } from "next/navigation";
import { useEffect, useState, useCallback } from "react";
import { toast } from "react-toastify";
import { Employees } from "@/types/dataTypes";

export default function AbsensePage() {
  const router = useRouter();
  const [mode, setMode] = useState<"checkin" | "checkout">("checkin");
  const [employeeId, setEmployeeId] = useState("");
  const [clockIn, setClockIn] = useState("");
  const [clockOut, setClockOut] = useState("");
  const [attendanceId, setAttendanceId] = useState("");
  const [employees, setEmployees] = useState<Employees[]>([]);

  // Fetch employees
  const fetchEmployees = useCallback(async () => {
    try {
      const res = await fetch("http://localhost:8080/api/employees");
      const data = await res.json();
      setEmployees(Array.isArray(data.data) ? data.data : []);
    } catch {
      toast.error("Failed to fetch employees");
    }
  }, []);

  useEffect(() => {
    fetchEmployees();
  }, [fetchEmployees]);

  // Format date to MySQL DATETIME
  const formatDateTime = (dateStr: string) => {
    if (!dateStr) return "";
    return new Date(dateStr).toISOString().replace("T", " ").substring(0, 19);
  };

  // Check-in handler
  const handleCheckIn = async (e: React.FormEvent) => {
    e.preventDefault();
    // Check for empty employeeId or clockIn before request
    if (!employeeId) {
      toast.error("Employee cannot be empty!");
      return;
    }
    if (!clockIn) {
      toast.error("Clock In time cannot be empty!");
      return;
    }
    try {
      const res = await fetch("http://localhost:8080/api/attendance", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          employee_id: employeeId,
          clock_in: formatDateTime(clockIn),
        }),
      });
      const data = await res.json();
      if (res.ok) {
        toast.success(data.message || "Check-in successful!");
        setEmployeeId("");
        setClockIn("");
      } else if (
        res.status === 400 ||
        (data.message && data.message.toLowerCase().includes("bad request"))
      ) {
        toast.error("Bad request: Please check your input!");
      } else {
        toast.error(data.message || "Check-in failed!");
      }
    } catch {
      toast.error("Check-in failed!");
    }
  };

  // Check-out handler
  const handleCheckOut = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!attendanceId) {
      toast.error("Attendance ID cannot be empty!");
      return;
    }
    if (!clockOut) {
      toast.error("Clock Out time cannot be empty!");
      return;
    }
    try {
      const res = await fetch(
        `http://localhost:8080/api/attendance/${attendanceId}`,
        {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            clock_out: formatDateTime(clockOut),
          }),
        }
      );
      const data = await res.json();
      if (res.ok) {
        toast.success(data.message || "Check-out successful!");
        setAttendanceId("");
        setClockOut("");
      } else {
        if (
          res.status === 404 ||
          (data.message && data.message.toLowerCase().includes("not found"))
        ) {
          toast.error("Attendance ID is invalid");
        } else if (
          res.status === 400 ||
          (data.message && data.message.toLowerCase().includes("bad request"))
        ) {
          toast.error("Bad request: Please check your input!");
        } else {
          toast.error(data.message || "Check-out failed!");
        }
      }
    } catch {
      toast.error("Check-out failed!");
    }
  };

  return (
    <div className="pl-2 overflow-x-auto">
      <h1 className="text-2xl font-bold mb-4">Employee Attendance</h1>
      <div className="flex gap-4 mb-6">
        <button
          className={`px-4 py-2 rounded ${
            mode === "checkin" ? "bg-gray-800 text-white" : "bg-gray-200"
          }`}
          onClick={() => setMode("checkin")}
          type="button"
        >
          Check In
        </button>
        <button
          className={`px-4 py-2 rounded ${
            mode === "checkout" ? "bg-gray-800 text-white" : "bg-gray-200"
          }`}
          onClick={() => setMode("checkout")}
          type="button"
        >
          Check Out
        </button>
      </div>

      {mode === "checkin" && (
        <form
          onSubmit={handleCheckIn}
          className="mb-8 space-y-4 w-full max-w-md"
        >
          <div>
            <label className="block mb-1">Employee</label>
            <select
              value={employeeId}
              onChange={(e) => setEmployeeId(e.target.value)}
              className="border px-2 py-1 w-full sm:max-w-sm rounded"
            >
              <option value="" disabled>
                Select Employee
              </option>
              {employees.map((emp) => (
                <option key={emp.employee_id} value={emp.employee_id}>
                  {emp.name}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label className="block mb-1">Clock In (timestamp)</label>
            <input
              type="datetime-local"
              value={clockIn}
              onChange={(e) => setClockIn(e.target.value)}
              className="border px-2 py-1 w-full sm:max-w-sm rounded"
            />
          </div>
          <button
            type="submit"
            className="bg-gray-800 text-white px-4 py-2 rounded w-full sm:w-auto"
          >
            Check In
          </button>
        </form>
      )}

      {mode === "checkout" && (
        <form onSubmit={handleCheckOut} className="space-y-4 w-full max-w-md">
          <div>
            <label className="block mb-1">Attendance ID</label>
            <input
              type="text"
              value={attendanceId}
              onChange={(e) => setAttendanceId(e.target.value)}
              className="border px-2 py-1 w-full sm:max-w-sm rounded"
            />
          </div>
          <div>
            <label className="block mb-1">Clock Out (timestamp)</label>
            <input
              type="datetime-local"
              value={clockOut}
              onChange={(e) => setClockOut(e.target.value)}
              className="border px-2 py-1 w-full sm:max-w-sm rounded"
            />
          </div>
          <button
            type="submit"
            className="bg-gray-800 text-white px-4 py-2 rounded w-full sm:w-auto"
          >
            Check Out
          </button>
        </form>
      )}
    </div>
  );
}
