"use client";

import { useRouter } from "next/navigation";
import { Edit } from "lucide-react";
import { useState } from "react";
import { toast } from "react-toastify";

export default function EditButton({
  id,
  type,
  onEdit,
}: {
  id: number;
  type?: "departments" | "employees";
  onEdit?: () => void;
}) {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);

  const handleEdit = () => {
    setIsLoading(true);
    try {
      router.push(`/edit/${id}?type=${type}`);
      if (onEdit) onEdit();
    } catch (error) {
      toast.error("Gagal mengarahkan ke halaman edit.");
      setIsLoading(false);
    }
  };

  return (
    <button
      onClick={handleEdit}
      disabled={isLoading}
      className={`px-2 py-1 rounded flex items-center gap-1 transition-colors duration-200 bg-blue-600 hover:bg-blue-700 text-white font-medium shadow`}
      style={{ fontSize: "0.95rem", minWidth: 0 }}
      title="Edit"
    >
      {isLoading ? (
        <div className="animate-spin rounded-full h-4 w-4 border-t-2 border-b-2 border-white"></div>
      ) : (
        <Edit className="w-4 h-4 mr-1" />
      )}
    </button>
  );
}
