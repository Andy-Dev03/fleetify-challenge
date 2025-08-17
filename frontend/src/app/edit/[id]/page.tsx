"use client";

import FormData from "@/components/FormData";
import { Employees, DepartmentWithEmployees } from "@/types/dataTypes";
import { useParams, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";

export default function EditPage() {
  const searchParams = useSearchParams();
  const type = searchParams.get("type");

  const params = useParams();
  const id = params.id;

  const [initialData, setInitialData] = useState<
    Employees | DepartmentWithEmployees | null
  >(null);

  useEffect(() => {
    if (!id || !type) return;
    const fetchData = async () => {
      let url = "";
      if (type === "employees") {
        url = `http://localhost:8080/api/employee/${id}`;
      } else {
        url = `http://localhost:8080/api/departement/${id}`;
      }
      const res = await fetch(url, {
        method: "GET",
        cache: "no-store",
      });

      if (res.ok) {
        const data = await res.json();
        console.log(data.data);

        setInitialData(data.data);
      }
    };
    fetchData();
  }, [id, type]);

  return (
    <div className="pl-2 space-y-2">
      <h1 className="text-2xl font-bold">Edit Data</h1>
      <button
        className="mb-4 px-4 py-2 bg-gray-300 rounded hover:bg-gray-400 transition"
        onClick={() => {
          if (typeof window !== "undefined") {
            window.history.back();
          }
        }}
      >
        Back
      </button>

      <FormData
        formType={type === "employees" ? "employees" : "departments"}
        initialData={initialData}
        isEdit={true}
      />
    </div>
  );
}
