"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { toast } from "react-toastify";
import { Department } from "../types/dataTypes";
import type {
  FormData,
  Employees,
  DepartmentWithEmployees,
} from "../types/dataTypes";

type FormType = "departments" | "employees";

const DEPARTMENT_FIELDS = [
  { name: "departmentName", label: "Department Name", type: "text" },
  { name: "maxClockIn", label: "Max Clock In (Time)", type: "time" },
  { name: "maxClockOut", label: "Max Clock Out (Time)", type: "time" },
];

const EMPLOYEE_FIELDS = [
  { name: "DepartmentID", label: "Department", type: "select" },
  { name: "Name", label: "Name", type: "text" },
  { name: "Address", label: "Address", type: "text" },
];

export default function FormData({
  formType,
  initialData = null,
  isEdit = false,
}: {
  formType: FormType;
  initialData?: Employees | DepartmentWithEmployees | null;
  isEdit?: boolean;
}) {
  const router = useRouter();
  const [departments, setDepartments] = useState<Department[]>([]);
  const [formData, setFormData] = useState<FormData>({
    departmentName: "",
    maxClockIn: "",
    maxClockOut: "",
    DepartmentID: "",
    Name: "",
    Address: "",
    EmployeeID: "",
  });

  // Set form data when editing or resetting
  useEffect(() => {
    if (isEdit && initialData) {
      if (formType === "departments") {
        setFormData({
          departmentName:
            (initialData as Department | DepartmentWithEmployees)
              .department_name || "",
          // slice(0,5) → agar sesuai format input type="time"
          maxClockIn: (
            (initialData as Department | DepartmentWithEmployees)
              .max_clock_in_time || ""
          ).slice(0, 5),
          maxClockOut: (
            (initialData as Department | DepartmentWithEmployees)
              .max_clock_out_time || ""
          ).slice(0, 5),
          DepartmentID: "",
          Name: "",
          Address: "",
          EmployeeID: "",
        });
      } else if (formType === "employees") {
        setFormData({
          departmentName: "",
          maxClockIn: "",
          maxClockOut: "",
          DepartmentID:
            (initialData as Employees).department_id?.toString() || "",
          Name: (initialData as Employees).name || "",
          Address: (initialData as Employees).address || "",
          EmployeeID: (initialData as Employees).employee_id || "",
        });
      }
    } else {
      setFormData({
        departmentName: "",
        maxClockIn: "",
        maxClockOut: "",
        DepartmentID: "",
        Name: "",
        Address: "",
        EmployeeID: "",
      });
    }
  }, [isEdit, initialData, formType]);

  // Handle input changes
  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  // Fetch departments for employee form
  const fetchDepartments = async () => {
    try {
      const res = await fetch("http://localhost:8080/api/departements", {
        method: "GET",
        cache: "no-store",
      });
      if (!res.ok) throw new Error("Failed to fetch departments");
      const data = await res.json();
      setDepartments(data.data);
    } catch (error) {
      toast.error("Error fetching departments");
    }
  };

  // Fetch departments when formType is employees
  useEffect(() => {
    if (formType === "employees") {
      fetchDepartments();
    }
  }, [formType]);

  // Handle form submission
  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (formType === "departments") {
      if (!formData.departmentName) {
        toast.error("Department Name cannot be empty!");
        return;
      } else if (!formData.maxClockIn) {
        toast.error("Max Clock In cannot be empty!");
        return;
      } else if (!formData.maxClockOut) {
        toast.error("Max Clock Out cannot be empty!");
        return;
      }
    } else if (formType === "employees") {
      if (!formData.DepartmentID) {
        toast.error("Department cannot be empty!");
        return;
      } else if (!formData.Name) {
        toast.error("Name cannot be empty!");
        return;
      } else if (!formData.Address) {
        toast.error("Address cannot be empty!");
        return;
      }
    }
    try {
      let payload;
      let url;
      let method = isEdit ? "PATCH" : "POST";
      if (formType === "departments") {
        payload = {
          DepartmentName: formData.departmentName,
          // gunakan value yang lama (initialData) jika input kosong
          MaxClockInTimeStr:
            formData.maxClockIn ||
            (initialData as Department).max_clock_in_time.slice(0, 5),
          MaxClockOutTimeStr:
            formData.maxClockOut ||
            (initialData as Department).max_clock_out_time.slice(0, 5),
        };
        url =
          isEdit && initialData && "id" in initialData
            ? `http://localhost:8080/api/departement/${
                (initialData as Department | DepartmentWithEmployees).id
              }`
            : "http://localhost:8080/api/departement";
      } else if (formType === "employees") {
        payload = {
          DepartmentID: Number(formData.DepartmentID),
          Name: formData.Name,
          Address: formData.Address,
          EmployeeID: formData.EmployeeID,
        };
        url =
          isEdit && initialData && "id" in initialData
            ? `http://localhost:8080/api/employee/${
                (initialData as Employees).id
              }`
            : "http://localhost:8080/api/employee";
      }
      const res = await fetch(url!, {
        method,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      if (!res.ok)
        throw new Error(
          isEdit
            ? `Failed to update ${formType}`
            : `Failed to create ${formType}`
        );

      // Reset form when creating
      if (!isEdit) {
        setFormData({
          departmentName: "",
          maxClockIn: "",
          maxClockOut: "",
          DepartmentID: "",
          Name: "",
          Address: "",
          EmployeeID: "",
        });
      }

      toast.success(
        `${formType.charAt(0).toUpperCase() + formType.slice(1)} ${
          isEdit ? "updated" : "created"
        } successfully!`
      );
      isEdit ? router.push("/") : router.refresh();
    } catch (error) {
      toast.error(`Error creating ${formType}`);
      console.error(`Error creating ${formType}:`, error);
    }
  };

  return (
    <div className="border p-4 rounded-sm shadow-md max-w-md">
      <h1 className="mb-2 text-xl font-bold text-center">
        {isEdit
          ? `Edit ${
              formType.charAt(0).toUpperCase() + formType.slice(1).toLowerCase()
            }`
          : `Add ${
              formType.charAt(0).toUpperCase() + formType.slice(1).toLowerCase()
            }`}
      </h1>

      {isEdit && initialData && (
        <div className="mb-4 p-2 bg-gray-100 rounded text-gray-700">
          <span className="font-semibold">Editing:</span>{" "}
          {formType === "departments"
            ? (initialData as Department | DepartmentWithEmployees)
                .department_name
            : formType === "employees"
            ? (initialData as Employees).name
            : ""}
        </div>
      )}

      <p className="border-t border-t-black my-2 mb-4"></p>

      <form onSubmit={handleSubmit} className="space-y-4">
        {/* Department Form Fields */}
        {formType === "departments" && (
          <>
            {DEPARTMENT_FIELDS.map((field) => (
              <div key={field.name}>
                <label htmlFor={field.name} className="block">
                  {field.label}
                </label>
                <input
                  type={field.type}
                  className="border border-black px-2 py-1 rounded-lg w-full"
                  name={field.name}
                  onChange={handleChange}
                  value={formData[field.name] || ""}
                />
              </div>
            ))}
          </>
        )}

        {/* Employee Form Fields */}
        {formType === "employees" && (
          <>
            <div>
              <label htmlFor="DepartmentID">Department</label>
              <select
                name="DepartmentID"
                id="DepartmentID"
                className="border border-black px-2 py-1 rounded-lg w-full"
                value={formData.DepartmentID}
                onChange={handleChange}
              >
                <option value="" disabled>
                  Select Department
                </option>
                {(departments ?? []).map((dept) => (
                  <option key={dept.id} value={dept.id}>
                    {dept.department_name}
                  </option>
                ))}
              </select>
            </div>
            {EMPLOYEE_FIELDS.filter(
              (field) => field.name !== "DepartmentID"
            ).map((field) => (
              <div key={field.name}>
                <label htmlFor={field.name}>{field.label}</label>
                <input
                  type={field.type}
                  className="border border-black px-2 py-1 rounded-lg w-full"
                  name={field.name}
                  onChange={handleChange}
                  value={formData[field.name] || ""}
                />
              </div>
            ))}
          </>
        )}
        <div>
          <button
            type="submit"
            className="bg-gray-800 text-white px-4 py-2 rounded-md w-full hover:bg-gray-700 transition"
          >
            {isEdit
              ? "Save Changes"
              : `Create ${
                  formType.charAt(0).toUpperCase() + formType.slice(1)
                }`}
          </button>
        </div>
      </form>
    </div>
  );
}
