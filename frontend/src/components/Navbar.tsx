"use client";
import React from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  CirclePlus,
  Database,
  CalendarClock,
  Rocket,
  Menu,
} from "lucide-react";

const navItems = [
  { href: "/", label: "All Data", icon: <Database className="mr-2 w-4" /> },
  {
    href: "/create",
    label: "Create Data",
    icon: <CirclePlus className="mr-2 w-4" />,
  },
  {
    href: "/absence",
    label: "Absense Employee",
    icon: <CalendarClock className="mr-2 w-4" />,
  },
];

export default function Navbar() {
  const pathname = usePathname();
  const [menuOpen, setMenuOpen] = React.useState(false);
  return (
    <nav className="bg-gray-900 text-white p-4 w-full md:h-screen md:flex md:flex-col">
      <div className="flex items-center justify-between md:block">
        <div className="flex items-center gap-2 mb-4 md:mb-4">
          <span className="text-xl font-bold">Fleetify Challenge</span>
        </div>
        <button
          className="md:hidden p-2 rounded focus:outline-none focus:ring-2 focus:ring-gray-400"
          onClick={() => setMenuOpen((prev) => !prev)}
          aria-label="Toggle menu"
        >
          <Menu className="w-6 h-6" />
        </button>
      </div>
      <ul
        className={`space-y-2 md:block ${
          menuOpen ? "block" : "hidden"
        } md:space-y-2 md:mt-0 mt-4`}
      >
        {navItems.map((item) => (
          <li
            key={item.href}
            className={`flex items-center rounded-lg px-2 py-2 transition-colors duration-200 ${
              pathname === item.href ? "bg-gray-700" : "hover:bg-gray-800"
            }`}
          >
            {item.icon}
            <Link href={item.href} className="w-full">
              {item.label}
            </Link>
          </li>
        ))}
      </ul>
    </nav>
  );
}
