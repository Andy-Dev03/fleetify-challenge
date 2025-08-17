"use client";
import { Trash2 } from "lucide-react";
import { toast } from "react-toastify";
import { deleteItem } from "../actions/deleteItem";
import { useState } from "react";

type ButtonDeleteProps = {
  type: string;
  id: number;
  onDelete?: () => void;
};

export default function DeleteButton({
  type,
  id,
  onDelete,
}: ButtonDeleteProps) {
  const [loading, setLoading] = useState(false);

  const handleDelete = async () => {
    setLoading(true);
    try {
      await deleteItem(type, id);
      toast.success("Berhasil menghapus data!");
      if (onDelete) onDelete();
    } catch (error: any) {
      toast.error(error?.message || "Gagal menghapus data!");
    } finally {
      setLoading(false);
    }
  };

  return (
    <button
      type="button"
      onClick={handleDelete}
      className={`px-2 py-1 rounded flex items-center gap-1 transition-colors duration-200 bg-red-600 hover:bg-red-700 text-white font-medium shadow`}
      style={{ fontSize: "0.95rem", minWidth: 0 }}
      title="Delete"
      disabled={loading}
    >
      {loading ? (
        <div className="animate-spin rounded-full h-4 w-4 border-t-2 border-b-2 border-white"></div>
      ) : (
        <Trash2 className="w-4 h-4 mr-1" />
      )}
    </button>
  );
}
