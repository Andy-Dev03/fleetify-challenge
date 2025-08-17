import DeleteButton from "./ButtonDelete";
import EditButton from "./EditButton";
import { Department, Employees, AttendanceLog } from "../types/dataTypes";
import { useState, useEffect } from "react";

export default function TableData({
  tableType,
  data,
  onDeleteSuccess,
}: {
  tableType: "departments" | "employees" | "attendanceLogs";
  data: Department[] | Employees[] | AttendanceLog[];
  onDeleteSuccess?: () => void;
}) {
  const [tableData, setTableData] = useState(data);
  useEffect(() => {
    setTableData(data);
  }, [data]);
  const headers = {
    departments: [
      "No",
      "Department Name",
      "Max Check-In Time",
      "Max Check-Out Time",
      "Actions",
    ],
    employees: [
      "No",
      "Employee ID",
      "Employee Name",
      "Address",
      "Department Name",
      "Actions",
    ],
    attendanceLogs: [
      "No",
      "Attendance ID",
      "Employee Name",
      "Department Name",
      "Date Attendance",
      "Attendance Type",
      "Clock In",
      "Clock Out",
      "Description",
    ],
  };

  return (
    <div className="w-full overflow-x-auto rounded-lg shadow">
      <table className="min-w-full border border-gray-200 bg-white text-xs sm:text-sm">
        <thead>
          <tr className="bg-gray-100">
            {headers[tableType].map((header) => (
              <th
                key={header}
                className="p-2 text-center font-semibold border whitespace-nowrap"
              >
                {header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {!tableData || tableData.length === 0 ? (
            <tr>
              <td
                colSpan={headers[tableType].length}
                className="p-4 text-center text-gray-500 whitespace-nowrap"
              >
                No data available
              </td>
            </tr>
          ) : (
            <>
              {tableType === "departments" &&
                (tableData as Department[]).map((dept, i) => (
                  <tr key={dept.id} className="hover:bg-gray-50">
                    <td className="p-2 text-center border whitespace-nowrap">
                      {i + 1}
                    </td>
                    <td className="p-2 text-center border whitespace-nowrap">
                      {dept.department_name}
                    </td>
                    <td className="p-2 text-center border whitespace-nowrap">
                      {dept.max_clock_in_time}
                    </td>
                    <td className="p-2 text-center border whitespace-nowrap">
                      {dept.max_clock_out_time}
                    </td>
                    <td className="p-2 text-center border whitespace-nowrap">
                      <div className="flex justify-center gap-2 flex-wrap">
                        <EditButton type={tableType} id={dept.id} />
                        <DeleteButton
                          type={tableType}
                          id={dept.id}
                          onDelete={() => {
                            if (onDeleteSuccess) {
                              onDeleteSuccess();
                            }
                          }}
                        />
                      </div>
                    </td>
                  </tr>
                ))}
              {tableType === "employees" &&
                (tableData as Employees[]).map((emp, i) => (
                  <tr key={emp.id} className="hover:bg-gray-50">
                    <td className="p-2 text-center border">{i + 1}</td>
                    <td className="p-2 text-center border whitespace-nowrap">
                      {emp.employee_id}
                    </td>
                    <td className="p-2 text-center border">{emp.name}</td>
                    <td className="p-2 text-center border">{emp.address}</td>
                    <td className="p-2 text-center border whitespace-nowrap">
                      {emp.department.department_name}
                    </td>
                    <td className="p-2 text-center border whitespace-nowrap">
                      <div className="flex justify-center gap-2 flex-wrap">
                        <EditButton type={tableType} id={emp.id} />
                        <DeleteButton
                          type={tableType}
                          id={emp.id}
                          onDelete={() => {
                            if (onDeleteSuccess) {
                              onDeleteSuccess();
                            }
                          }}
                        />
                      </div>
                    </td>
                  </tr>
                ))}
              {tableType === "attendanceLogs" &&
                (tableData as AttendanceLog[]).map((log, i) => (
                  <tr key={i} className="hover:bg-gray-50">
                    <td className="p-2 text-center border">{i + 1}</td>
                    <td className="p-2 text-center border whitespace-nowrap">
                      {log.attendance_id}
                    </td>
                    <td className="p-2 text-center border">{log.name}</td>
                    <td className="p-2 text-center border">{log.department}</td>
                    <td className="p-2 text-center border whitespace-nowrap">
                      {log.date_attendance}
                    </td>
                    <td className="p-2 text-center border">
                      {log.attendance_type === 1
                        ? "In"
                        : log.attendance_type === 2
                        ? "Out"
                        : "-"}
                    </td>
                    <td className="p-2 text-center border">
                      {log.clock_in ?? "-"}
                    </td>
                    <td className="p-2 text-center border">
                      {log.clock_out ?? "-"}
                    </td>
                    <td className="p-2 text-center border">
                      {log.description}
                    </td>
                  </tr>
                ))}
            </>
          )}
        </tbody>
      </table>
    </div>
  );
}
