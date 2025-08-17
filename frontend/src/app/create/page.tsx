"use client";

import FormData from "@/components/FormData";
import { useState } from "react";

type FormType = "departments" | "employees";
const formTypes: FormType[] = ["departments", "employees"];
const formTypeLabels = {
  departments: "Departments",
  employees: "Employees",
};

export default function CreateData() {
  const [formType, setFormType] = useState<FormType>("departments");
  return (
    <div className="pl-2 space-y-2">
      <h1 className="text-2xl font-bold">Create Data</h1>
      <p className="text-gray-600">Choose the data you want to make</p>
      <div className="flex gap-2 mb-4">
        {formTypes.map((type) => (
          <button
            key={type}
            onClick={() => setFormType(type)}
            className={`px-4 py-1 rounded-md transition ${
              formType === type
                ? "bg-gray-800 text-white"
                : "bg-gray-200 hover:bg-gray-300"
            }`}
          >
            {formTypeLabels[type]}
          </button>
        ))}
      </div>
      <FormData formType={formType} />
    </div>
  );
}
