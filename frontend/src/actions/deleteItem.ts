"use server";

export async function deleteItem(type: string, id: number) {
  let endpoint = "";
  if (type === "departments") {
    endpoint = `http://localhost:8080/api/departement/${id}`;
  } else if (type === "employees") {
    endpoint = `http://localhost:8080/api/employee/${id}`;
  } else {
    throw new Error("Unknown delete type");
  }
  
  const res = await fetch(endpoint, {
    method: "DELETE",
    cache: "no-store",
  });
  if (!res.ok) {
    const errorData = await res.json();
    throw new Error(errorData.message || "Failed to delete item");
  }
}
